package server

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	c "github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/server/graph"
	"github.com/angelorc/sinfonia-go/server/graph/generated"
	w3t "github.com/angelorc/sinfonia-go/server/web3token"
	"github.com/angelorc/sinfonia-go/utility"
	"github.com/labstack/echo"
)

var token string
var playgroundPassword string
var submissionToken string

// Get header value and add to gql resolvers
func getHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		headers := ctx.Request().Header
		token = utility.GetHeaderString("Authorization", headers)
		playgroundPassword = utility.GetHeaderString("Playground-Password", headers)

		return next(ctx)
	}
}

func InitGraphql(e *echo.Echo) {
	// Resolvers && Directives
	resolver := graph.Resolver{Token: &token}
	config := generated.Config{Resolvers: &resolver}
	config.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		web3token, err := w3t.NewWeb3TokenFromBearer(*resolver.Token)
		if err != nil {
			return nil, fmt.Errorf("access denied")
		}

		if web3token.GetDomain() != "test.com" {
			return nil, fmt.Errorf("access denied")
		}

		return next(ctx)
	}

	e.Use(getHeaders)

	// new custom handler based on gqlgen version 0.11.3
	queryHandler := handler.New(generated.NewExecutableSchema(config))

	// queryHandler.Use(&debug.Tracer{})
	queryHandler.AddTransport(transport.POST{})
	queryHandler.AddTransport(transport.MultipartForm{})
	queryHandler.SetQueryCache(lru.New(1000))
	queryHandler.Use(extension.AutomaticPersistedQuery{Cache: lru.New(100)})
	queryHandler.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		rc := graphql.GetOperationContext(ctx)
		if playgroundPassword == c.GetSecret("GRAPHQL_PLAYGROUND_PASS") {
			rc.DisableIntrospection = false
		} else {
			rc.DisableIntrospection = true
		}
		return next(ctx)
	})

	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL Playground", c.GetSecret("GRAPHQL_ENDPOINT"))))
	//e.POST("/query", echo.WrapHandler(dataloader.DataLoaderMiddleware(queryHandler)))
	e.POST("/query", echo.WrapHandler(queryHandler))
}
