package loaders

import (
	"LinkKrec/graph/model"
	"context"
	"net/http"
	"time"

	"github.com/knakk/sparql"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type DataBase struct {
	Repo *sparql.Repo
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader              *dataloadgen.Loader[string, *model.User]
	VacancyLoader           *dataloadgen.Loader[string, *model.Vacancy]
	CompanyLoader           *dataloadgen.Loader[string, *model.Company]
	EducationEntryLoader    *dataloadgen.Loader[string, *model.EducationEntry]
	ConnectionRequestLoader *dataloadgen.Loader[string, *model.ConnectionRequest]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *sparql.Repo) *Loaders {
	// define the data loader
	ur := &DataBase{Repo: conn}
	return &Loaders{
		UserLoader:              dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
		VacancyLoader:           dataloadgen.NewLoader(ur.getVacancies, dataloadgen.WithWait(time.Millisecond)),
		CompanyLoader:           dataloadgen.NewLoader(ur.getCompanies, dataloadgen.WithWait(time.Millisecond)),
		EducationEntryLoader:    dataloadgen.NewLoader(ur.getEducationEntries, dataloadgen.WithWait(time.Millisecond)),
		ConnectionRequestLoader: dataloadgen.NewLoader(ur.getConnectionRequests, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(conn *sparql.Repo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(conn)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetUser returns single user by id efficiently
func GetUser(ctx context.Context, userID string) (*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)
}

// GetUsers returns many users by ids efficiently
func GetUsers(ctx context.Context, userIDs []string) ([]*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}

// GetVacancy returns single vacancy by id efficiently
func GetVacancy(ctx context.Context, vacancyID string) (*model.Vacancy, error) {
	loaders := For(ctx)
	return loaders.VacancyLoader.Load(ctx, vacancyID)
}

// GetVacancies returns many vacancies by ids efficiently
func GetVacancies(ctx context.Context, vacancyIDs []string) ([]*model.Vacancy, error) {
	loaders := For(ctx)
	return loaders.VacancyLoader.LoadAll(ctx, vacancyIDs)
}

func GetCompany(ctx context.Context, companyID string) (*model.Company, error) {
	loaders := For(ctx)
	return loaders.CompanyLoader.Load(ctx, companyID)
}

func GetCompanies(ctx context.Context, companyIDs []string) ([]*model.Company, error) {
	loaders := For(ctx)
	return loaders.CompanyLoader.LoadAll(ctx, companyIDs)
}

func GetEducationEntry(ctx context.Context, educationEntryID string) (*model.EducationEntry, error) {
	loaders := For(ctx)
	return loaders.EducationEntryLoader.Load(ctx, educationEntryID)
}

func GetEducationEntries(ctx context.Context, educationEntryIDs []string) ([]*model.EducationEntry, error) {
	loaders := For(ctx)
	return loaders.EducationEntryLoader.LoadAll(ctx, educationEntryIDs)
}

func GetConnectionRequest(ctx context.Context, connectionRequestID string) (*model.ConnectionRequest, error) {
	loaders := For(ctx)
	return loaders.ConnectionRequestLoader.Load(ctx, connectionRequestID)
}

func getConnectionRequests(ctx context.Context, connectionRequestIDs []string) ([]*model.ConnectionRequest, error) {
	loaders := For(ctx)
	return loaders.ConnectionRequestLoader.LoadAll(ctx, connectionRequestIDs)
}
