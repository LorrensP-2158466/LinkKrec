package loaders

import (
	"LinkKrec/graph/model"
	"LinkKrec/graph/util"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/knakk/sparql"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	LoadersKey = "dataloaders"
)

type DataBase struct {
	Repo *sparql.Repo
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader              *dataloadgen.Loader[string, *model.User]
	VacancyLoader           *dataloadgen.Loader[string, *model.Vacancy]
	EmployerLoader          *dataloadgen.Loader[string, *model.Employer]
	EducationEntryLoader    *dataloadgen.Loader[string, *model.EducationEntry]
	ExperienceEntryLoader   *dataloadgen.Loader[string, *model.ExperienceEntry]
	ConnectionRequestLoader *dataloadgen.Loader[string, *model.ConnectionRequest]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *sparql.Repo) *Loaders {
	// define the data loader
	ur := &DataBase{Repo: conn}
	return &Loaders{
		UserLoader:              dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
		VacancyLoader:           dataloadgen.NewLoader(ur.getVacancies, dataloadgen.WithWait(time.Millisecond)),
		EmployerLoader:          dataloadgen.NewLoader(ur.getEmployers, dataloadgen.WithWait(time.Millisecond)),
		EducationEntryLoader:    dataloadgen.NewLoader(ur.getEducationEntries, dataloadgen.WithWait(time.Millisecond)),
		ExperienceEntryLoader:   dataloadgen.NewLoader(ur.getExperienceEntries, dataloadgen.WithWait(time.Millisecond)),
		ConnectionRequestLoader: dataloadgen.NewLoader(ur.getConnectionRequests, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(conn *sparql.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		loader := NewLoaders(conn)
		c.Set(string(LoadersKey), loader)
		c.Next()
	}
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	fmt.Println("Request", ctx.Value(util.QueryRepoKey))
	fmt.Println("util", ctx.Value(util.UpdateRepoKey))
	fmt.Println("loaders", ctx.Value(LoadersKey))
	fmt.Println("util", ctx.Value(util.SessionInfoKey))
	return ctx.Value(LoadersKey).(*Loaders)
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

func GetEmployer(ctx context.Context, employerID string) (*model.Employer, error) {
	loaders := For(ctx)
	return loaders.EmployerLoader.Load(ctx, employerID)
}

func GetEmployers(ctx context.Context, employerIDs []string) ([]*model.Employer, error) {
	loaders := For(ctx)
	return loaders.EmployerLoader.LoadAll(ctx, employerIDs)
}

func GetEducationEntry(ctx context.Context, educationEntryID string) (*model.EducationEntry, error) {
	loaders := For(ctx)
	return loaders.EducationEntryLoader.Load(ctx, educationEntryID)
}

func GetEducationEntries(ctx context.Context, educationEntryIDs []string) ([]*model.EducationEntry, error) {
	loaders := For(ctx)
	return loaders.EducationEntryLoader.LoadAll(ctx, educationEntryIDs)
}

func GetExperienceEntry(ctx context.Context, experienceEntryID string) (*model.ExperienceEntry, error) {
	loaders := For(ctx)
	return loaders.ExperienceEntryLoader.Load(ctx, experienceEntryID)
}

func GetExperienceEntries(ctx context.Context, experienceEntryIDs []string) ([]*model.ExperienceEntry, error) {
	loaders := For(ctx)
	return loaders.ExperienceEntryLoader.LoadAll(ctx, experienceEntryIDs)
}

func GetConnectionRequest(ctx context.Context, connectionRequestID string) (*model.ConnectionRequest, error) {
	loaders := For(ctx)
	return loaders.ConnectionRequestLoader.Load(ctx, connectionRequestID)
}

func getConnectionRequests(ctx context.Context, connectionRequestIDs []string) ([]*model.ConnectionRequest, error) {
	loaders := For(ctx)
	return loaders.ConnectionRequestLoader.LoadAll(ctx, connectionRequestIDs)
}
