package main

import (
	"LinkKrec/graph"
	"LinkKrec/graph/loaders"
	"LinkKrec/graph/util"
	"LinkKrec/usersession"
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

const (
	port          = "8080"
	QueryRepoKey  = "queryrepo"
	UpdateRepoKey = "updateRepo"
)

var store *sessions.CookieStore

func GetQueryRepo(c *gin.Context) *sparql.Repo {
	val, exists := c.Get(string(QueryRepoKey))
	if !exists {
		panic("Something went horribly wrong with assigning the repo to the context")
	}
	return val.(*sparql.Repo)
}
func GetUpdateRepo(c *gin.Context) *sparql.Repo {
	val, exists := c.Get(string(UpdateRepoKey))
	if !exists {
		panic("Something went horribly wrong with assigning the update repo to the context")
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

func loginGothUser(c *gin.Context, goth_user goth.User) (*usersession.UserSessionInfo, error) {
	fmt.Println("LOGGING GOTH USER IN: ", goth_user)
	query_repo := GetQueryRepo(c)

	// check if the user already exists
	res, err := query_repo.Query(fmt.Sprintf(`			
	    PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX foaf: <http://xmlns.com/foaf/0.1/> 

		SELECT 
			?id ?email ?name ?ProfileCompleted
			GROUP_CONCAT(DISTINCT ?companyId; seperator", ") AS ?companyIds)
		WHERE{
			?user a lr:User;
				lr:Id ?id ;
				foaf:name ?name;
				foaf:mbox ?email;
				lr:isProfileComplete ?ProfileCompleted .

			OPTIONAL{
				?user a lr:User ;
					lr:hasCompany ?company .
				?company a lr:Company;
					lr:Id ?companyId .
			}

			FILTER(?email = "%s")
		}
		GROUP BY ?id ?email ?name ?ProfileCompleted
		`, goth_user.Email))

	if err != nil {
		return nil, err
	}
	// if not
	if len(res.Solutions()) == 0 {
		update_repo := GetUpdateRepo(c)
		uuid := uuid.New().String()
		insertQuery := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX foaf: <http://xmlns.com/foaf/0.1/> 

	    INSERT DATA {
	        lr:User%s a lr:User ;
	        	lr:Id "%s" ;
	            foaf:name "%s" ;
	            foaf:mbox "%s" ;
	            lr:isProfileComplete false .
	    }
	    `, uuid, uuid, goth_user.FirstName+" "+goth_user.LastName, goth_user.Email)

		err := update_repo.Update(insertQuery)
		fmt.Println("MADE UPDATE")
		if err != nil {
			return nil, err
		}
		return &usersession.UserSessionInfo{
			IsComplete: false,
			Email:      goth_user.Email,
			Id:         uuid,
			Cookie:     goth_user.AccessToken,
			CompanyIds: []string{},
		}, nil
	} else {
		// user exists create session info and return
		sessInfo, err := util.MapPrimitiveBindingsToStruct[usersession.UserSessionInfo](res.Solutions()[0])
		fmt.Println(sessInfo)
		if err != nil {
			return nil, err
		}
		sessInfo.Cookie = goth_user.AccessToken

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

	sessionInfo, err := loginGothUser(c, goth_user)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session.Values[usersession.SessionInfoKey] = sessionInfo

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

		if val, ok := session.Values[usersession.SessionInfoKey].(*usersession.UserSessionInfo); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Can't get session info: %T", session.Values[usersession.SessionInfoKey]), "authorization_urls": []string{"/auth/google"}})
		} else {
			c.Set(string(usersession.SessionInfoKey), val)
		}

		c.Next()
	}
}

func ginCtxToRawCtx(gqlHandler *handler.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// these 2 are used by internals/implementations of graphql which we don't want to modify
		keys := []string{loaders.LoadersKey, usersession.SessionInfoKey}

		ctx := c.Request.Context()
		for _, key := range keys {
			val, exists := c.Get(string(key))
			if !exists || val == nil {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Man wtf man: %s", key))
			}
			ctx = context.WithValue(ctx, key, val)
		}
		// put gin inside of the ctx so we can access them inside the resolvers
		ctx = context.WithValue(ctx, "ginCtx", c)

		// Create a new request with the modified context
		req := c.Request.WithContext(ctx)

		// Pass the modified request to the GraphQL handler
		gqlHandler.ServeHTTP(c.Writer, req)
	}
}

func setupRouter(repo *sparql.Repo, updateRepo *sparql.Repo) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set(string(QueryRepoKey), repo)
		c.Set(string(UpdateRepoKey), updateRepo)
		c.Next()
	})
	r.Use(loaders.Middleware(repo))
	// Middleware to save session after each request

	r.GET("/auth/:provider", signInWithProvider)
	r.GET("/auth/:provider/callback", callbackHandler)

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					Repo:       repo,
					UpdateRepo: updateRepo,
					Store:      store},
			},
		))

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	protected.Use(func(c *gin.Context) {
		c.Next() // Process request first

		// After processing the request, save the session
		session, _ := store.Get(c.Request, "user-session")
		if err := session.Save(c.Request, c.Writer); err != nil {
			// Log or handle session save error if needed
			log.Printf("Error saving session: %v", err)
		}
	})
	{
		protected.GET("/playground", gin.WrapH(playground.Handler("GraphQL playground", "/graphql")))
		protected.GET("/is_authorized", func(c *gin.Context) { c.String(http.StatusAccepted, "AUTHORIZED") })
		protected.GET("/graphql", ginCtxToRawCtx(srv))
		protected.POST("/graphql", ginCtxToRawCtx(srv))
		protected.GET("/test_sess_info", func(c *gin.Context) {
			session, _ := store.Get(c.Request, "user-session")
			c.JSON(http.StatusOK, session.Values[usersession.SessionInfoKey])
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
		log.Fatalf("Failed to connect to SPARQL endpoint: %v", err)
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
