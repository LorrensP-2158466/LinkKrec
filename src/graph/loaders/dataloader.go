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
	UserLoader            *dataloadgen.Loader[string, *model.User]
	VacancyLoader         *dataloadgen.Loader[string, *model.Vacancy]
	EmployerLoader        *dataloadgen.Loader[string, *model.Employer]
	EducationEntryLoader  *dataloadgen.Loader[string, *model.EducationEntry]
	ExperienceEntryLoader *dataloadgen.Loader[string, *model.ExperienceEntry]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *sparql.Repo) *Loaders {
	// define the data loader
	ur := &DataBase{Repo: conn}
	return &Loaders{
		UserLoader:            dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
		VacancyLoader:         dataloadgen.NewLoader(ur.getVacancies, dataloadgen.WithWait(time.Millisecond)),
		EmployerLoader:        dataloadgen.NewLoader(ur.getEmployers, dataloadgen.WithWait(time.Millisecond)),
		EducationEntryLoader:  dataloadgen.NewLoader(ur.getEducationEntries, dataloadgen.WithWait(time.Millisecond)),
		ExperienceEntryLoader: dataloadgen.NewLoader(ur.getExperienceEntries, dataloadgen.WithWait(time.Millisecond)),
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
