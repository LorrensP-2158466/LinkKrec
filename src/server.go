package main

import (
	"LinkKrec/graph"
	"LinkKrec/graph/loaders"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/knakk/sparql"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const port = "8080"

var store *sessions.CookieStore

func init_store() {
	randomKey := securecookie.GenerateRandomKey(32)
	store = sessions.NewCookieStore([]byte(randomKey))
	gothic.Store = store

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
	}
}


	baseUrl := "http://localhost:3030/link_krec/"
	queryEndpoint := baseUrl + "query"
	mutateEndpoint := baseUrl + "update"

	// Connect to the SPARQL endpoints
	repo, err := sparql.NewRepo(queryEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect to the SPARQL endpoint: %v", err)
	}
	updateRepo, err := sparql.NewRepo(mutateEndpoint)

func signInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func callbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session, err := store.Get(c.Request, "user-session")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}


	session.Values["access_token"] = user.AccessToken
	if err := session.Save(c.Request, c.Writer); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/is_authorized")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "user-session")
		if err != nil {
			session.Options.MaxAge = -1
			_ = session.Save(c.Request, c.Writer)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "authorization_urls": []string{"/auth/google"}})
			return
		}

		if _, ok := session.Values["access_token"]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "authorization_urls": []string{"/auth/google"}})
			return
		}

		c.Next()
	}
}

func setupRouter(repo *sparql.Repo) *gin.Engine {
	r := gin.Default()


	r.GET("/auth/:provider", signInWithProvider)
	r.GET("/auth/:provider/callback", callbackHandler)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repo: repo}}))
	injected_srv := loaders.Middleware(repo, srv)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))
		protected.GET("/is_authorized", func(c *gin.Context) { c.String(http.StatusAccepted, "AUTHORIZED") })
		protected.GET("/graphql", gin.WrapH(injected_srv))
	}

	return r
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Missing required environment variables")
	}

	endpointURL := "http://localhost:3030/link_krec/sparql"
	repo, err := sparql.NewRepo(endpointURL)
	if err != nil {
		log.Fatalf("Failed to connect to SPARQL endpoint: %v", err)
	}

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
	)
	init_store()

	r := setupRouter(repo)
	log.Printf("Server running at http://localhost:%s/", port)
	log.Fatal(r.Run(":" + port))
}
