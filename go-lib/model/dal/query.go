package dal

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// HasID interface provides a method od set an object ID
type HasID interface {
	SetID(string)
}

// Edge is an interface to convert simple Edge into a database representation
type Edge interface {
	// ToEdgeDO accepts a full entity ID that involves the collection name
	// it returns a DO (database object) that should be inserted directly into DB
	ToEdgeDO(fullEntityID string) interface{}
}

// DBQueryOne - is a helper function to query one object from DB into `dest`.
// @dest must be a pointer value
func dbQueryOne(ctx context.Context, dest interface{}, query string, bindVars map[string]interface{}, db driver.Database) (*driver.DocumentMeta, errstack.E) {
	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, errstack.WrapAsInf(err, "Can not execute the query")
	}
	defer errstack.CallAndLog(logger, cursor.Close)
	if !cursor.HasMore() {
		return nil, NewNotFound("")
	}
	meta, err := cursor.ReadDocument(ctx, dest)
	return &meta, errstack.WrapAsDomainF(err, "Can't parse result (%s) into %t type", query, dest)
}

// DBQueryOne - is a helper function to query one object from DB into `dest`.
// Doesn't fail if output value is a derived result and doesn't have an ID (i.e. boolean).
// @dest must be a pointer value
func DBQueryOne(ctx context.Context, dest interface{}, query string, bindVars map[string]interface{}, db driver.Database) errstack.E {
	_, err := dbQueryOne(ctx, dest, query, bindVars, db)
	return err
}

// DBQueryFirst - queries first entry of the returned list. Fails if element is not a collection entry. Query has to return a list.
// @dest must be a pointer value
func DBQueryFirst(ctx context.Context, dest interface{}, query string, bindVars map[string]interface{}, db driver.Database) errstack.E {
	meta, err := dbQueryOne(ctx, dest, query, bindVars, db)
	if err != nil {
		return err
	}
	// The key is empty when the output of the query is calculated. E.g. a boolean.
	if meta.Key == "" {
		return errstack.NewReqF("Not a collection entry")
	}
	return nil
}

// DBExec - is a helper function to execute a query which doesn't return anything
// (eg remove query).
func DBExec(ctx context.Context, query string, bindVars map[string]interface{}, db driver.Database) errstack.E {
	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return errstack.WrapAsInf(err, "Can not execute the query")
	}
	errstack.Log(logger, cursor.Close())
	return nil
}

// DBGetOneFromColl is a helper function to get one object from DB into `dest` using collection
// and document id. @dest must be a pointer value
func DBGetOneFromColl(ctx context.Context, dest interface{}, id string, collName dbconst.Col, db driver.Database) errstack.E {
	if id == "" {
		return model.ErrNoID
	}
	col, err := GetCollByName(ctx, db, collName)
	if err != nil {
		return model.ErrDbCollection(err, collName)
	}
	_, err = col.ReadDocument(ctx, id, dest)
	return MaybeWrapAsNotFound(err)
}

func trimAbsoluteID(absoluteID string) string {
	trimmed := strings.Split(absoluteID, "/")
	if len(trimmed) == 2 {
		return trimmed[1]
	}
	return absoluteID
}

// DBQueryMany executes a query and reads elements into the provided `dest`. It must be
// a pointer (not nil) to a Slice of elements compatible with expected data.
func DBQueryMany(ctx context.Context, dest interface{}, query string, bindVars map[string]interface{}, db driver.Database) errstack.E {
	destv := reflect.ValueOf(dest)
	if destv.Kind() != reflect.Ptr || destv.Elem().Kind() != reflect.Slice {
		return errstack.NewDomain("Destination must be a pointer to a slice")
	}
	if dest == nil {
		return errstack.NewDomain("Destination must be a non-nil pointer")
	}
	slicev := destv.Elem()
	etype := slicev.Type().Elem()

	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		return errstack.WrapAsInf(err, "Can't execute query")
	}
	defer errstack.CallAndLog(logger, cursor.Close)
	for cursor.HasMore() {
		elem := reflect.New(etype)
		if _, err = cursor.ReadDocument(ctx, elem.Interface()); err != nil {
			return errstack.WrapAsDomain(err, "Can't read value")
		}
		slicev = reflect.Append(slicev, elem.Elem())
	}
	destv.Elem().Set(slicev)
	return nil
}

// insertHasID is a helper function which will insert an obj into collection and assign
// a newly created ID to it. `obj` must be a pointer.
func insertHasID(ctx context.Context, coll dbconst.Col, obj HasID, db driver.Database) (driver.DocumentMeta, errstack.E) {
	meta, err := InsertAny(ctx, coll, obj, db)
	if err == nil {
		obj.SetID(meta.Key)
	}
	return meta, err
}

// InsertAny inserts an object into collection and doesn't assign ID after insertion.
// `obj` must be a pointer.
func InsertAny(ctx context.Context, coll dbconst.Col, obj interface{}, db driver.Database) (driver.DocumentMeta, errstack.E) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return driver.DocumentMeta{}, errstack.NewDomainF("Expects pointer type, got %T", obj)
	}
	col, err := GetCollByName(ctx, db, coll)
	if err != nil {
		return driver.DocumentMeta{}, model.ErrDbCollection(err, coll)
	}
	meta, err := col.CreateDocument(ctx, obj)
	return meta, errstack.WrapAsInf(err, "Can't insert document into '"+string(coll)+"' collection")
}

func existsQuery(filterQuery string) string {
	return fmt.Sprintf("RETURN FIRST(%s LIMIT 1 RETURN 1) != null", filterQuery)
}

// InsertIntoGraph inserts entity and associated edge.
func InsertIntoGraph(ctx context.Context, db driver.Database, docColl, edgeColl dbconst.Col, entity HasID, edge Edge) (driver.DocumentMeta, errstack.E) {
	if entity == nil {
		return driver.DocumentMeta{}, errstack.NewDomainF("Can't insert nil object [%s]", edgeColl)
	}
	dm, err := insertHasID(ctx, docColl, entity, db)
	if err != nil {
		return dm, errstack.WrapAsInfF(err, "Can't insert into coll [%s]", docColl)
	}
	edgeDO := edge.ToEdgeDO(dm.ID.String())
	_, err = InsertAny(ctx, edgeColl, &edgeDO, db)
	return dm, errstack.WrapAsInfF(err, "Can't insert into edges; documents used in the edge may not exist [%s]", edgeColl)
}

// DeleteByID deletes entry by its ID
func DeleteByID(ctx context.Context, db driver.Database, colletion dbconst.Col, id string) errstack.E {
	col, err := GetCollByName(ctx, db, colletion)
	if err != nil {
		return model.ErrDbCollection(err, colletion)
	}
	_, err = col.RemoveDocument(ctx, id)
	return errstack.WrapAsInf(err, fmt.Sprintf("can't delete '%s' from '%s'", id, colletion))
}
