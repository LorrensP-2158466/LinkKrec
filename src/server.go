package main

import (
	"LinkKrec/graph"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/knakk/rdf"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func loadRDFData(filename string) ([]rdf.Triple, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	triples, err := rdf.NewTripleDecoder(file, rdf.Turtle).DecodeAll()
	if err != nil {
		return nil, fmt.Errorf("error decoding triples: %v", err)
	}

	return triples, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Load RDF data
	triples, _ := loadRDFData("/rdf/load_rdf_data.ttl")

	for _, triple := range triples {
		fmt.Println(triple.Subj, triple.Pred, triple.Obj)
	}

	fmt.Println("Starting server on port " + port)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
