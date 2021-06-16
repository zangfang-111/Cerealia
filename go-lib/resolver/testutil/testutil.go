// Package testutil is meant for usage in integration tests only
// It should not be used in other application's places
package testutil

import (
	"context"
	"fmt"
	"net/http"

	"bitbucket.org/cerealia/apps/go-lib/gql"
	"bitbucket.org/cerealia/apps/go-lib/middleware"
	"bitbucket.org/cerealia/apps/go-lib/model"
	"bitbucket.org/cerealia/apps/go-lib/model/dal"
	"bitbucket.org/cerealia/apps/go-lib/model/dbconst"
	"bitbucket.org/cerealia/apps/go-lib/resolver"
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/secretkey"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txvalidation"
	driver "github.com/arangodb/go-driver"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/errstack"
	"github.com/stellar/go/keypair"
)

// SampleUser1 is a sample user to use in tests
var SampleUser1 = model.UserLoginInput{
	Email:    "kk@kk.kk",
	Password: "birthday",
}

// SampleUser2 is a sample user to use in tests
var SampleUser2 = model.UserLoginInput{
	Email:    "bb@bb.bb",
	Password: "birthday",
}

// SampleUser3 is a sample user to use in tests
var SampleUser3 = model.UserLoginInput{
	Email:    "ss@ss.ss",
	Password: "birthday",
}

// SampleUserModerator is a user that has moderator and trader capabilities
var SampleUserModerator = model.UserLoginInput{
	Email:    "jj@jj.jj",
	Password: "birthday",
}

// SampleNewUser is a sample new user to use in tests
var SampleNewUser = model.NewUserInput{
	FirstName: "John",
	LastName:  "Doe",
	Email:     "test.john@gmail.com",
	Password:  "JohnDoe12345",
	PublicKey: "GD3EPS4EBOK6ZELDEN466I6EU4LW7TK6UL6INZRC6OKLZKYXXESS64VE",
	OrgID:     "430823",
	OrgRole:   "manager",
}

// SampleNewOrganization is a sample new organization to use in tests
var SampleNewOrganization = model.OrgInput{
	Name:      "testCompany",
	Address:   "testAddress",
	Email:     "test.m",
	Telephone: "1poi23456789123",
}

const (
	// SampleUser1Seed is a key seed for the SampleUser1
	SampleUser1Seed = "SCVZ2D7TXSH7RNIYJOQSJOBTJSWE7AMDOHSVBQJA3OZABTJQ5MC47XID"
	// SampleUser2Seed is a key seed for the SampleUser2
	SampleUser2Seed = "SDOM5RGWXOYFKXLOL4OKUQ7EAKO4KVCTLEZCQCESIAVQOQFMPRWOKAKY"
	// SampleUserModeratorSeed is a key seed for trader-moderator
	SampleUserModeratorSeed = "SBIGQPMZASCJLA3PWOAWD5JWLQYQ6OHNXXZR34SGW4CTSBVPSWTCAYAQ"
	// SampleMasterSecret is a sample secret for whole application
	SampleMasterSecret = "SACCWZ55B6YVURGOCXIGKR77KHTPJ7SI2IU6P4EFZYA3ZLOORKOET2NB"
)

// Credentials is a user object that contains his context too
type Credentials struct {
	*model.User
	Ctx context.Context
}

// loginErr creates context for a logged in user
func loginErr(r resolver.Resolver, loginInput model.UserLoginInput) (*Credentials, error) {
	ctx := context.Background()
	authUser, err := r.Mutation().UserLogin(ctx, loginInput)
	if err != nil {
		return nil, err
	}
	authHandler := middleware.WithAuth(r.DB())
	req := &http.Request{
		Header: map[string][]string{
			"Authorization": []string{"Bearer " + authUser.Token},
		}}
	req = req.WithContext(ctx)
	rctx := routing.Context{Request: req}
	err = authHandler(&rctx)
	if err != nil {
		return nil, err
	}
	user, err := dal.GetUserByEmail(ctx, r.DB(), loginInput.Email)
	if err != nil {
		return nil, err
	}
	return &Credentials{
		user,
		rctx.Request.Context(),
	}, err
}

// Login creates context for a logged in user
func Login(r resolver.Resolver, loginInput model.UserLoginInput) (*Credentials, errstack.E) {
	creds, err := loginErr(r, loginInput)
	return creds, errstack.WrapAsInf(err)
}

// LoginTwo logs in with two users at once
func LoginTwo(r resolver.Resolver, firstUser, secondUser model.UserLoginInput) (*Credentials, *Credentials, error) {
	u1, err := Login(r, firstUser)
	if err != nil {
		return nil, nil, errstack.WrapAsInfF(err, "Can't login as first user")
	}
	u2, err := Login(r, secondUser)
	return u1, u2, errstack.WrapAsInfF(err, "Can't login as second user")
}

// MakeTradeInput creates a new trade input
func MakeTradeInput(name, buyerID, sellerID string, description *string) model.NewTradeInput {
	return model.NewTradeInput{
		TemplateID:  "1471516",
		Name:        "trade-test",
		BuyerID:     buyerID,
		SellerID:    sellerID,
		Description: description,
	}
}

// GetTrade finds a trade in DB
func GetTrade(userContext context.Context, resolver gql.ResolverRoot, tradeID string) (*model.Trade, error) {
	return resolver.Query().Trade(userContext, tradeID)
}

// ApproveDoc approves a document
func ApproveDoc(userContext context.Context, resolver gql.ResolverRoot, id model.TradeStageDocPath, signedTx string) (*model.TradeStageDoc, error) {
	return resolver.Mutation().TradeStageDocApprove(userContext, id, signedTx)
}

// CloseStage attempts to close a stage
func CloseStage(userContext context.Context, resolver gql.ResolverRoot, id model.TradeStagePath, signedTx, reason string) (*model.ApproveReq, error) {
	return resolver.Mutation().TradeStageCloseReq(userContext, id, signedTx, reason)
}

// CloseStageReqReject rejects the close request of a stage
func CloseStageReqReject(userContext context.Context, resolver gql.ResolverRoot, id model.TradeStagePath, signedTx, reason string) (*int, error) {
	return resolver.Mutation().TradeStageCloseReqReject(userContext, id, signedTx, reason)
}

// SignTx signs a transaction with given secrets
func SignTx(driver stellar.Driver, rawTx string, seeds ...string) (string, error) {
	txEnvelope, err := txvalidation.ReadEnvelopeBuilder(rawTx)
	if err != nil {
		return "", err
	}
	for _, seed := range seeds {
		parsed, err := secretkey.Parse(seed)
		if err != nil {
			return "", err
		}
		txEnvelope, err = driver.SignEnvelope(txEnvelope, *parsed)
		if err != nil {
			return "", err
		}
	}
	return txEnvelope.Base64()
}

// InsertPoolSourceAccounts creates two testnet pool accounts
func InsertPoolSourceAccounts(ctx context.Context, db driver.Database) errstack.E {
	secrets := []string{
		// These accounts are prepared to be used on testnet too
		"SBM5J676BY3G72ARP366XL5IHXTWE5LPLYY7C3VNI3KLV2D6H3KO5PEF",
		"SCOBNHIPYLBY2UV5VSEXLLB4EM6ZRX5Q26TRZR2BK2Q52C6IJ2O7TMUL",
	}
	for _, secret := range secrets {
		parsed, err := secretkey.Parse(secret)
		if err != nil {
			return err
		}
		_, err = upsertNewAcc(ctx, db, *parsed)
		if err != nil {
			return err
		}
	}
	return nil
}

func upsertNewAcc(ctx context.Context, db driver.Database, kp keypair.Full) (*model.TXSourceAcc, errstack.E) {
	lock := model.TXSourceAcc{}
	query := fmt.Sprintf(`
UPSERT {_key: @pubKey}
INSERT {_key: @pubKey, sKey: @sKey, type: "pool"}
REPLACE {_key: @pubKey, sKey: @sKey, type: "pool", lockExpiresAt: null, lockUserID: null, lockTradeID: null, lockUnlockedAt: null}
IN %s
RETURN NEW
`, dbconst.ColTxSourceAccs)
	bindVars := map[string]interface{}{
		"pubKey": kp.Address(),
		"sKey":   kp.Seed(),
	}
	err := dal.DBQueryFirst(ctx, &lock, query, bindVars, db)
	return &lock, errstack.WrapAsInfF(err, "Couldn't acquire a lock for an account '%s'. Source account pool is exhausted.", kp.Address())
}

// DeleteAllPoolSourceAccounts deletes all pool accounts from the DB
func DeleteAllPoolSourceAccounts(ctx context.Context, db driver.Database) {
	output := []model.TXSourceAcc{}
	query := fmt.Sprintf(`
for i in %s
filter i.type == "%s"
remove i._key in %s
`, dbconst.ColTxSourceAccs, model.TxSourceAccTypePool, dbconst.ColTxSourceAccs)
	_ = dal.DBQueryMany(ctx, &output, query, map[string]interface{}{}, db)
}

// UnlockSourceAccs unlocks all locks (including trade locks)
func UnlockSourceAccs(ctx context.Context, db driver.Database) {
	output := []model.TXSourceAcc{}
	query := fmt.Sprintf(`
for i in %s
UPDATE {_key: i._key, lockTradeID:null, lockUserID:null, lockExpiresAt:null, lockUnlockedAt:null} IN %s
return i
`, dbconst.ColTxSourceAccs, dbconst.ColTxSourceAccs)
	_ = dal.DBQueryMany(ctx, &output, query, map[string]interface{}{}, db)
}

// CleanSourceAccs performs cleaning operations for source accounts
func CleanSourceAccs(ctx context.Context, db driver.Database) errstack.E {
	DeleteAllPoolSourceAccounts(ctx, db)
	err := InsertPoolSourceAccounts(ctx, db)
	if err != nil {
		return err
	}
	UnlockSourceAccs(ctx, db)
	return nil
}

// CreateTrade logs in with several users and creates a noop trade
func CreateTrade(r resolver.Resolver, desc *string) (buyer, seller, third, moderator *Credentials, trade *model.Trade, err errstack.E) {
	buyer, err = Login(r, SampleUser1)
	if err != nil {
		return
	}
	seller, err = Login(r, SampleUser2)
	if err != nil {
		return
	}
	third, err = Login(r, SampleUser3)
	if err != nil {
		return
	}
	moderator, err = Login(r, SampleUserModerator)
	if err != nil {
		return
	}
	trade, errt := r.Mutation().TradeCreate(buyer.Ctx, MakeTradeInput("test-trade", buyer.ID, seller.ID, desc))
	err = errstack.WrapAsInf(errt, "Can't create trade")
	return
}
