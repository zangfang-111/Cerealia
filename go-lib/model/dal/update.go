package dal

import (
	"context"

	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/errstack"
)

// UpdateDoc updates document by it's key
func UpdateDoc(ctx context.Context, db driver.Database, collection dbconst.Col, key string, obj interface{}) (driver.DocumentMeta, errstack.E) {
	if obj == nil {
		return driver.DocumentMeta{}, errstack.NewDomain("Can't update nil object")
	}
	col, err := GetCollByName(ctx, db, collection)
	if err != nil {
		return driver.DocumentMeta{}, model.ErrDbCollection(err, collection)
	}
	dm, err := col.UpdateDocument(ctx, key, obj)
	return dm, errstack.WrapAsInf(err, "Can't update "+string(collection)+" object")
}

func replaceDoc(ctx context.Context, db driver.Database, collection dbconst.Col, key string, obj interface{}) (driver.DocumentMeta, errstack.E) {
	if obj == nil {
		return driver.DocumentMeta{}, errstack.NewDomain("Can't replace nil object")
	}
	col, err := GetCollByName(ctx, db, collection)
	if err != nil {
		return driver.DocumentMeta{}, model.ErrDbCollection(err, collection)
	}
	dm, err := col.ReplaceDocument(ctx, key, obj)
	return dm, errstack.WrapAsInf(err, "Can't replace "+string(collection)+" object")
}

func deleteDoc(ctx context.Context, db driver.Database, collection dbconst.Col, key string) errstack.E {
	col, err := GetCollByName(ctx, db, collection)
	if err != nil {
		return model.ErrDbCollection(err, collection)
	}
	_, err = col.RemoveDocument(ctx, key)
	return errstack.WrapAsInf(err, "Can't remove object from "+string(collection))
}
