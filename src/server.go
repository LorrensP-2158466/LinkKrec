package main

import (
	"LinkKrec/graph"
	"LinkKrec/graph/loaders"
	"LinkKrec/graph/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

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

func GetQueryRepo(c *gin.Context) *sparql.Repo {
	val, exists := c.Get(string(util.QueryRepoKey))
	if !exists {
		panic("Something went horribly wrong with assigning the repo to the context")
	}
	return val.(*sparql.Repo)
}
func GetUpdateRepo(c *gin.Context) *sparql.Repo {
	val, exists := c.Get(string(util.UpdateRepoKey))
	if !exists {
		panic("Something went horribly wrong with assigning the repo to the context")
	}
	return val.(*sparql.Repo)
}

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

func signInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func loginGothUser(c *gin.Context, goth_user goth.User) (*util.UserSessionInfo, error) {
	query_repo := GetQueryRepo(c)

	// check if the user already exists
	res, err := query_repo.Query(fmt.Sprintf(`
			
    PREFIX lr: <http://linkRec.org/ontology/>
    PREFIX schema: <http://schema.org/>
    PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

    SELECT ?id ?email ?name ?accountCompleted
    WHERE{
    	?user a lr:User;
    		lr:Id ?id ;
    		lr:hasName ?name;
    		lr:hasEmail ?email;
    		lr:isEmplyer false ;
    		lr:isProfileComplete ?accountCompleted .

    	FILTER(?email = "%s")
    }
		`, goth_user.Email))

	if err != nil {
		return nil, err
	}
	// if not
	if len(res.Solutions()) == 0 {
		update_repo := GetUpdateRepo(c)
		uuid := uuid.New().String()
		insertQuery := fmt.Sprintf(`
	    PREFIX lr: <http://linkRec.org/ontology/>

	    INSERT DATA {
	        lr:User%s a lr:User ;
	        	lr:Id "%s" ;
	            lr:hasName "%s" ;
	            lr:hasEmail "%s" ;
	            lr:isEmployer false ;
	            lr:isProfileComplete false .
	    }`, uuid, uuid, goth_user.Name, goth_user.Email)

		err := update_repo.Update(insertQuery)
		if err != nil {
			return nil, err
		}
		return &util.UserSessionInfo{
			IsComplete: false,
			IsUser:     true,
			Email:      goth_user.Email,
			Id:         uuid}, nil
	} else {
		// user exists create session info and return
		sessInfo, err := util.MapPrimitiveBindingsToStruct[util.UserSessionInfo](res.Solutions()[0])
		if err != nil {
			return nil, err
		}
		sessInfo.Cookie = goth_user.AccessToken
		sessInfo.IsUser = true

		return &sessInfo, nil
	}
}

func callbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	goth_user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session, err := store.Get(c.Request, "user-session")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Values["access_token"] = goth_user.AccessToken

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	sessionInfo, err := loginGothUser(c, goth_user)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session.Values[util.SessionInfoKey] = &sessionInfo

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

		if _, ok := session.Values[util.SessionInfoKey]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Can't get session info", "authorization_urls": []string{"/auth/google"}})
			return
		}

		c.Next()
	}
}

func ginCtxToRawCtx(gqlHandler *handler.Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Create a new context with the Gin context information
		ctx := context.WithValue(c.Request.Context(), util.QueryRepoKey, c.Value(string(util.QueryRepoKey)))
		ctx = context.WithValue(ctx, util.UpdateRepoKey, c.Value(string(util.UpdateRepoKey)))
		ctx = context.WithValue(ctx, loaders.LoadersKey, c.Value(string(loaders.LoadersKey)))
		ctx = context.WithValue(ctx, util.SessionInfoKey, c.Value(string(util.SessionInfoKey)))

		// Create a new request with the modified context
		req := c.Request.WithContext(ctx)

		// Pass the modified request to the GraphQL handler
		gqlHandler.ServeHTTP(c.Writer, req)
	}
}

func setupRouter(repo *sparql.Repo, updateRepo *sparql.Repo) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set(string(util.QueryRepoKey), repo)
		c.Set(string(util.UpdateRepoKey), updateRepo)
		c.Next()
	})
	r.Use(loaders.Middleware(repo))
	r.Use(util.GinContextToContextMiddleware())

	r.GET("/auth/:provider", signInWithProvider)
	r.GET("/auth/:provider/callback", callbackHandler)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repo: repo}}))

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))
		protected.GET("/is_authorized", func(c *gin.Context) { c.String(http.StatusAccepted, "AUTHORIZED") })
		protected.GET("/graphql", gin.WrapH(srv))
		protected.POST("/graphql", gin.WrapH(srv))
		protected.GET("/test_sess_info", func(c *gin.Context) {
			session, _ := store.Get(c.Request, "user-session")
			c.JSON(http.StatusOK, session.Values[util.SessionInfoKey])
		})
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

	baseUrl := "http://localhost:3030/link_krec/"
	queryEndpoint := baseUrl + "query"
	mutateEndpoint := baseUrl + "update"

	// Connect to the SPARQL endpoints
	repo, err := sparql.NewRepo(queryEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect to the SPARQL endpoint: %v", err)
	}
	updateRepo, err := sparql.NewRepo(mutateEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect to SPARQL endpoint: %v", err)
	}

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallbackURL),
	)
	init_store()

	r := setupRouter(repo, updateRepo)
	log.Printf("Server running at http://localhost:%s/", port)
	log.Fatal(r.Run(":" + port))
}
