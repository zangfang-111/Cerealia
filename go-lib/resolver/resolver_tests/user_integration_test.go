package resolvertests

import (
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/resolver/testutil"
	. "github.com/robert-zaremba/checkers"
	. "gopkg.in/check.v1"
)

func (s *TradeIntegrationSuite) TestLogin(c *C) {
	// positive
	_, err := testutil.Login(s.noopResolver, testutil.SampleUser1)
	c.Assert(err, IsNil)

	// negative
	_, err = testutil.Login(s.noopResolver, model.UserLoginInput{
		Email:    "aa@aa.aa",
		Password: "bad password",
	})
	c.Assert(err, Not(ErrorContains), "(?i).*pass.*")
}

func (s *TradeIntegrationSuite) TestUserSignup(c *C) {
	mr := s.noopResolver.Mutation()
	// test for correct new user

	_, err := mr.UserSignup(testctx, &testutil.SampleNewUser)
	c.Check(err, IsNil, Comment("Failed to insert user"))

	// test for new user with email which already exists.
	nu1 := testutil.SampleNewUser
	nu1.FirstName = "AnotherJohn"
	_, err = mr.UserSignup(testctx, &nu1)
	c.Check(err, NotNil, Comment("Expected error for new user signup with existing email"))

	// delete new user and organization from db
	err = dal.DeleteUser(testctx, s.db, nu1.Email)
	c.Check(err, IsNil, Comment("Failed to delete new user from db"))
}

func (s *TradeIntegrationSuite) TestUserSignupNagative(c *C) {
	mr := s.noopResolver.Mutation()

	// test for wrong new user without firstname
	nu := testutil.SampleNewUser
	nu.FirstName = ""
	_, err := mr.UserSignup(testctx, &nu)
	c.Check(err, NotNil, Comment("Expected error for new user without name"))

	// test for wrong new user without company.
	nu1 := nu
	nu1.FirstName = "John"
	nu1.OrgID = ""
	_, err = mr.UserSignup(testctx, &nu1)
	c.Check(err, NotNil, Comment("Expected error for new user without company"))
}

func (s *TradeIntegrationSuite) TestCreateNewOrganization(c *C) {
	mr := s.noopResolver.Mutation()
	// test for correct new user

	org, err := mr.OrganizationCreate(testctx, testutil.SampleNewOrganization)
	c.Assert(err, IsNil, Comment("Failed to insert organization"))
	c.Check(org.ID, Not(IsEmpty))

	// negative test for new organization with name which already exists.
	newOrg1 := testutil.SampleNewOrganization
	newOrg1.Address = "testAddress1"
	newOrg1.Email = "test.org1@gmail.com"
	_, err = mr.OrganizationCreate(testctx, newOrg1)
	c.Check(err, NotNil)

	// negative test for new organization with address which already exists.
	newOrg2 := testutil.SampleNewOrganization
	newOrg2.Name = "testCompany2"
	newOrg1.Email = "test.org2@gmail.com"
	_, err = mr.OrganizationCreate(testctx, newOrg2)
	c.Check(err, NotNil)
	// delete new organization from db
	err = dal.DeleteOrg(testctx, s.db, org.ID)
	c.Check(err, IsNil, Comment("Failed to delete new org from db"))
}
