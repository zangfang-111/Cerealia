package dal

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"

	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

func (s *DalSuite) TestDBQueryOne(c *C) {
	var exists bool
	query := existsQuery(`FOR d IN users FILTER @userEmail IN d.emails`)
	bindVars := map[string]interface{}{
		"userEmail": "ss@ss.ss",
	}
	errs := DBQueryOne(testctx, &exists, query, bindVars, s.db)
	c.Check(errs, IsNil, Comment("Failed to get user by email"))
	c.Check(exists, IsTrue, Comment("the query value should be true"))

	var trade model.Trade
	query = `FOR d IN trades FILTER d.buyer.userID==@buyerUserID return d`
	bindVars = map[string]interface{}{
		"buyerUserID": "2",
	}
	errs = DBQueryOne(testctx, &trade, query, bindVars, s.db)
	c.Check(errs, IsNil, Comment("Failed to get trade by buyerUserID"))

	// negative test
	query = `FOR d IN trades FILTER d.buyer.userID==@buyerUserID return dd` // wrong query
	bindVars = map[string]interface{}{
		"buyerUserID": "2",
	}
	errs = DBQueryOne(testctx, &trade, query, bindVars, s.db)
	c.Check(errs, NotNil)
}

func (s *DalSuite) TestDBGetOneFromColl(c *C) {
	var user model.User
	errs := DBGetOneFromColl(testctx, &user, "1", dbconst.ColUsers, s.db)
	c.Check(errs, IsNil, Comment("Failed to get user by id"))
	errs = DBGetOneFromColl(testctx, &user, "2", dbconst.ColUsers, s.db)
	c.Check(errs, IsNil, Comment("Failed to get user by id"))

	// negative test
	errs = DBGetOneFromColl(testctx, &user, "100", dbconst.ColUsers, s.db) // wrong user id
	c.Check(errs, NotNil)
}

func (s *DalSuite) TestDBQueryMany(c *C) {
	var users []model.User
	query := "for d in users return d"
	errs := DBQueryMany(testctx, &users, query, nil, s.db)
	c.Check(errs, IsNil, Comment("Failed to get all users"))

	var trades []model.Trade
	query = "for d in trades Filter d.buyer.userID==@buyerUserID return d"
	bindVars := map[string]interface{}{
		"buyerUserID": "2",
	}
	errs = DBQueryMany(testctx, &trades, query, bindVars, s.db)
	c.Check(errs, IsNil, Comment("Failed to get trades filtered by buyerUserID"))

	// negative test
	query = `FOR d IN trades FILTER d.buyer.userID==@buyerUserID return dd` // wrong query
	bindVars = map[string]interface{}{
		"buyerUserID": "2",
	}
	errs = DBQueryMany(testctx, &trades, query, bindVars, s.db)
	c.Check(errs, NotNil)
}

func (s *DalSuite) TestInsertUser(c *C) {
	user := model.User{
		FirstName: "John",
		LastName:  "Doe",
		Emails:    []string{"test.userinsert@gmail.com"},
	}
	meta, errs := InsertAny(testctx, dbconst.ColUsers, &user, s.db)
	c.Check(errs, IsNil, Comment("Failed to insert new user"))
	c.Check(meta, NotNil)
	// delete inserted new user
	userCol, err := s.db.Collection(testctx, string(dbconst.ColUsers))
	c.Check(err, IsNil, Comment("Failed to connect to users collection"))
	_, err = userCol.RemoveDocument(testctx, meta.Key)
	c.Check(err, IsNil, Comment("Failed to delete new user inserted"))
}

func (s *DalSuite) TestTrimAbsoluteID(c *C) {
	id := trimAbsoluteID("hello-world/123")
	c.Check(id, Equals, "123")
	id2 := trimAbsoluteID("sadsf/aaaaaaadfsaf")
	c.Check(id2, Equals, "aaaaaaadfsaf")
	id3 := trimAbsoluteID("hello-world-123")
	c.Check(id3, Equals, "hello-world-123")
	id4 := trimAbsoluteID("hello/world/123")
	c.Check(id4, Equals, "hello/world/123")
	id5 := trimAbsoluteID("")
	c.Check(id5, Equals, "")
}
