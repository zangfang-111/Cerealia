The `model` package represents the Domain modeled using GraphQL.

The DB specific (if needed) functions and structures should be placed in /go-lib/dal* packages

## Package organization

+ go-lib/model          ← business domain goes here (represented by GraphQL)
+ go-lib/model/daladb   ← specific arangodb dal functions and model will go here
+ go-lib/resolver       ← default resolvers go here
+ go-lib/resolvermem    ← in memory resolvers (eg for tests) will go here
