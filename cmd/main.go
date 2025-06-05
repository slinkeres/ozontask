package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slinkeres/ozontask/graph"
	"github.com/slinkeres/ozontask/internal/consts"
	"github.com/slinkeres/ozontask/database"
	"github.com/slinkeres/ozontask/internal/gateway"
	im "github.com/slinkeres/ozontask/internal/gateway/in-memory"
	"github.com/slinkeres/ozontask/internal/gateway/postgres"
	"github.com/slinkeres/ozontask/internal/server/graphql"
	"github.com/slinkeres/ozontask/internal/service"
	lg "github.com/slinkeres/ozontask/internal/logger"
)


func main() {
	logger := lg.InitLogger()
	logger.Info.Print("Executing InitLogger")

	envFile := ".env"
	if len(os.Args) >= 2 {
		envFile = os.Args[1]
	}

	logger.Info.Print("Executing InitConfig")
	logger.Info.Printf("Reading %s \n", envFile)
	if err := godotenv.Load(envFile); err != nil { 
		logger.Info.Print("Cannot load env file")
		logger.Err.Fatalf(err.Error())
	}

	logger.Info.Print("Connecting to Postgres")

	postgresDb, err := database.NewPostgresDB()

	if err != nil {
		logger.Info.Print("Failed to connect to database")
		logger.Err.Fatalf(err.Error())
	}

	var gateways *gateway.Gateways

	logger.Info.Print("Creating Gateways.")
	logger.Info.Print("USE_IN_MEMORY = ", os.Getenv("USE_IN_MEMORY"))

	if os.Getenv("USE_IN_MEMORY") == "true" {
		posts := im.NewPostsInMemory(consts.PostsPullSize)
		comments := im.NewCommentsInMemory(consts.CommentsPullSize)
		gateways = gateway.NewGateways(posts, comments)
	} else {
		posts := postgres.NewPostsPostgres(postgresDb)
		comments := postgres.NewCommentsPostgres(postgresDb)
		gateways = gateway.NewGateways(posts, comments)
	}

	logger.Info.Print("Creating Services")
	services := service.NewServices(gateways, logger)

	logger.Info.Print("Creating graphql server")
	port := os.Getenv("PORT")
	if port == ""{
		logger.Info.Print("Set default port")
		port = "8080"
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graphql.Resolver{
		PostsService:      services.Posts,
		CommentsService:   services.Comments,
		CommentsObservers: graphql.NewCommentsObserver(),
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logger.Info.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	logger.Err.Fatal(http.ListenAndServe(":"+port, nil))

}