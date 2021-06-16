//TODO: go:generate gorunpkg github.com/99designs/gqlgen

package resolver

import (
	"bitbucket.org/cerealia/apps/go-lib/stellar"
	"bitbucket.org/cerealia/apps/go-lib/stellar/txsource"

	"bitbucket.org/cerealia/apps/go-lib/gql"
	driver "github.com/arangodb/go-driver"
	"github.com/robert-zaremba/log15"
)

var logger = log15.Root()

// Resolver is an extension to gql.ResolverRoot which provides access to struct variables.
type Resolver interface {
	gql.ResolverRoot
	DB() driver.Database
}

// resolver represents the base resolver for all resolvers
type resolver struct {
	db             driver.Database
	stellarDriver  *stellar.Driver
	txSourceDriver txsource.Driver

	approveReqRes       gql.ApproveReqResolver
	docRes              gql.DocResolver
	mutationRes         gql.MutationResolver
	queryRes            gql.QueryResolver
	tradeRes            gql.TradeResolver
	tradeStageAddReqRes gql.TradeStageAddReqResolver
	tradeStageDocRes    gql.TradeStageDocResolver
	userRes             gql.UserResolver
	accessApprovalRes   gql.AccessApprovalResolver
	tradeOfferRes       gql.TradeOfferResolver
	moderatorRes        gql.StageModeratorResolver
	notificationRes     gql.NotificationResolver
}

// NewResolver initialize a new instance of resolver
func NewResolver(db driver.Database, stellarDriver *stellar.Driver, txSourceDriver txsource.Driver) Resolver {
	r := new(resolver)
	r.db = db
	r.txSourceDriver = txSourceDriver
	r.stellarDriver = stellarDriver
	r.approveReqRes = approveReqResolver{r}
	r.docRes = docResolver{r}
	r.mutationRes = mutationResolver{r}
	r.queryRes = queryResolver{r}
	r.tradeRes = tradeResolver{r}
	r.tradeStageAddReqRes = tradeStageAddReqResolver{r}
	r.tradeStageDocRes = tradeStageDocResolver{r}
	r.userRes = userResolver{r}
	r.accessApprovalRes = accessApprovalResolver{r}
	r.tradeOfferRes = tradeOfferRes{r}
	r.moderatorRes = stageModeratorResolver{r}
	r.notificationRes = notificationResolver{r}
	return r
}

// DB returns DB driver
func (r *resolver) DB() driver.Database {
	return r.db
}

// Mutation gets a mutation resolver
func (r *resolver) Mutation() gql.MutationResolver {
	return r.mutationRes
}

// Query gets a query resolver
func (r *resolver) Query() gql.QueryResolver {
	return r.queryRes
}

func (r *resolver) ApproveReq() gql.ApproveReqResolver {
	return approveReqResolver{r}
}

func (r *resolver) Doc() gql.DocResolver {
	return r.docRes
}

// Trade returns a trade resolver
func (r *resolver) Trade() gql.TradeResolver {
	return r.tradeRes
}

// TradeStageAddReq implements the resolver interface
func (r *resolver) TradeStageAddReq() gql.TradeStageAddReqResolver {
	return r.tradeStageAddReqRes
}

// TradeStageDoc returns tradeStageDoc resolver
func (r *resolver) TradeStageDoc() gql.TradeStageDocResolver {
	return r.tradeStageDocRes
}

// User returns user resolver
func (r *resolver) User() gql.UserResolver {
	return r.userRes
}

// AccessApproval implements the joinning user approval
func (r *resolver) AccessApproval() gql.AccessApprovalResolver {
	return r.accessApprovalRes
}

// TradeOffer implements the resolver interface
func (r *resolver) TradeOffer() gql.TradeOfferResolver {
	return r.tradeOfferRes
}

// StageModerator implements the trade stage moderator interface
func (r *resolver) StageModerator() gql.StageModeratorResolver {
	return r.moderatorRes
}

// Notification implements the notification interface
func (r *resolver) Notification() gql.NotificationResolver {
	return r.notificationRes
}
