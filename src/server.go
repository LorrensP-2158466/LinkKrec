package main

import (
	"LinkKrec/graph"
	"LinkKrec/graph/loaders"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/knakk/sparql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	endpointURL := "http://localhost:3030/link_krec/sparql"

	// Connect to the SPARQL endpoint
	repo, err := sparql.NewRepo(endpointURL)
	if err != nil {
		log.Fatalf("Failed to connect to the SPARQL endpoint: %v", err)
	}

	fmt.Println("Starting server on port " + port)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repo: repo}}))

	injected_srv := loaders.Middleware(repo, srv)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", injected_srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
