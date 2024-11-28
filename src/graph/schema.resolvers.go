package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	graph_model "LinkKrec/graph/model"
	query_builder "LinkKrec/querybuilder"
	"context"
	"fmt"
	"strings"
)

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input graph_model.RegisterUserInput) (*graph_model.User, error) {
	panic(fmt.Errorf("not implemented: RegisterUser - registerUser"))
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input graph_model.UpdateUserInput) (*graph_model.User, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
}

// UpdateUserProfile is the resolver for the updateUserProfile field.
func (r *mutationResolver) UpdateUserProfile(ctx context.Context, id string, input graph_model.UpdateProfileInput) (*graph_model.User, error) {
	panic(fmt.Errorf("not implemented: UpdateUserProfile - updateUserProfile"))
}

// ManageConnection is the resolver for the manageConnection field.
func (r *mutationResolver) ManageConnection(ctx context.Context, userID string, connectedUserID string, action string) (*graph_model.AskedConnection, error) {
	panic(fmt.Errorf("not implemented: ManageConnection - manageConnection"))
}

// NotifyProfileVisit is the resolver for the notifyProfileVisit field.
func (r *mutationResolver) NotifyProfileVisit(ctx context.Context, visitorID string, visitedUserID string) (*graph_model.Notification, error) {
	panic(fmt.Errorf("not implemented: NotifyProfileVisit - notifyProfileVisit"))
}

// CreateVacancy is the resolver for the createVacancy field.
func (r *mutationResolver) CreateVacancy(ctx context.Context, employerID string, input graph_model.CreateVacancyInput) (*graph_model.Vacancy, error) {
	panic(fmt.Errorf("not implemented: CreateVacancy - createVacancy"))
}

// UpdateVacancy is the resolver for the updateVacancy field.
func (r *mutationResolver) UpdateVacancy(ctx context.Context, id string, input graph_model.CreateVacancyInput) (*graph_model.Vacancy, error) {
	panic(fmt.Errorf("not implemented: UpdateVacancy - updateVacancy"))
}

// DeleteVacancy is the resolver for the deleteVacancy field.
func (r *mutationResolver) DeleteVacancy(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteVacancy - deleteVacancy"))
}

// UpdateUserLookingForOpportunities is the resolver for the updateUserLookingForOpportunities field.
func (r *mutationResolver) UpdateUserLookingForOpportunities(ctx context.Context, userID string, looking bool) (*graph_model.User, error) {
	panic(fmt.Errorf("not implemented: UpdateUserLookingForOpportunities - updateUserLookingForOpportunities"))
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id string) (*graph_model.User, error) {
	var q = query_builder.QueryBuilder().
		Select([]string{"name", "id"}).
		GroupConcat("skill", ", ", "skills").
		WhereSubject("user", "User").
		Where("Id", "id").
		Where("hasName", "name").
		Where("hasSkill", "skill").
		Filter("id", "\"1\"", query_builder.EQ).
		GroupBy([]string{"id", "name"}).
		Build()

	res, err := r.Repo.Query(q)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result, err := MapStringBindingsToStruct[graph_model.User](res.Bindings())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context, name *string, location *string, isEmployer *bool, skills []*string, lookingForOpportunities *bool) ([]*graph_model.User, error) {
	fmt.Println("GetUsers")

	var skillSubQuery = query_builder.QueryBuilder().
		WhereSubject("user", "User").
		Select([]string{"user"}).
		GroupConcat("skill", ", ", "skills").
		Where("hasSkill", "skill").
		GroupBy([]string{"user"}).
		BuildSubQuery()

	var q = query_builder.
		QueryBuilder().Select([]string{"userId", "userName", "email", "isEmployer", "location", "lookingForOppurtunities", "skills", "connections"}).
		WhereSubject("user", "User").
		Where("Id", "userId").
		Where("hasName", "userName").
		Where("hasEmail", "email").
		Where("isEmployer", "isEmployer").
		Where("hasLocation", "location").
		Where("isLookingForOpportunities", "lookingForOpportunities").
		WhereSubQuery(skillSubQuery).
		//GroupBy([]string{"userId", "userName"}).
		Build()
	fmt.Println(q)

	res, err := r.Repo.Query(
		q)
	if err != nil {
		return nil, err
	}

	var users []*graph_model.User

	for i := 0; i < len(res.Results.Bindings); i++ {
		// Map the first result row to the User model
		row := res.Results.Bindings[i]

		// Extract values
		userId := row["userId"].Value
		userName := row["userName"].Value
		userSkills := row["skills"].Value

		// Split the skills string into a slice
		skillsList := strings.Split(userSkills, ", ")

		// Convert skillsList to []*string
		var skillsPtrList []*string
		for _, skill := range skillsList {
			skillCopy := skill
			skillsPtrList = append(skillsPtrList, &skillCopy)
		}

		// Create a graph_model.User instance
		user := &graph_model.User{
			ID:     userId,
			Name:   userName,
			Skills: skillsPtrList,
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserConnections is the resolver for the getUserConnections field.
func (r *queryResolver) GetUserConnections(ctx context.Context, userID string, skill *string, location *string, depth *int) ([]*graph_model.User, error) {
	panic(fmt.Errorf("not implemented: GetUserConnections - getUserConnections"))
}

// GetVacancies is the resolver for the getVacancies field.
func (r *queryResolver) GetVacancies(ctx context.Context, search *string, location *string, requiredSkills []*string, minEducation *graph_model.DegreeType, isActive *bool) ([]*graph_model.Vacancy, error) {
	panic(fmt.Errorf("not implemented: GetVacancies - getVacancies"))
}

// GetVacancy is the resolver for the getVacancy field.
func (r *queryResolver) GetVacancy(ctx context.Context, id string) (*graph_model.Vacancy, error) {
	panic(fmt.Errorf("not implemented: GetVacancy - getVacancy"))
}

// GetEmployers is the resolver for the getEmployers field.
func (r *queryResolver) GetEmployers(ctx context.Context, name *string, location *string) ([]*graph_model.Employer, error) {
	panic(fmt.Errorf("not implemented: GetEmployers - getEmployers"))
}

// GetEmployer is the resolver for the getEmployer field.
func (r *queryResolver) GetEmployer(ctx context.Context, id string) (*graph_model.Employer, error) {
	panic(fmt.Errorf("not implemented: GetEmployer - getEmployer"))
}

// GetNotifications is the resolver for the getNotifications field.
func (r *queryResolver) GetNotifications(ctx context.Context, userID string, since *string) ([]*graph_model.Notification, error) {
	panic(fmt.Errorf("not implemented: GetNotifications - getNotifications"))
}

// GetConnectionRequests is the resolver for the getConnectionRequests field.
func (r *queryResolver) GetConnectionRequests(ctx context.Context, userID string, status *bool) ([]*graph_model.AskedConnection, error) {
	panic(fmt.Errorf("not implemented: GetConnectionRequests - getConnectionRequests"))
}

// NewConnectionRequest is the resolver for the newConnectionRequest field.
func (r *subscriptionResolver) NewConnectionRequest(ctx context.Context, forUserID string) (<-chan *graph_model.AskedConnection, error) {
	panic(fmt.Errorf("not implemented: NewConnectionRequest - newConnectionRequest"))
}

// ConnectionRequestStatusUpdate is the resolver for the connectionRequestStatusUpdate field.
func (r *subscriptionResolver) ConnectionRequestStatusUpdate(ctx context.Context, forUserID string) (<-chan *graph_model.AskedConnection, error) {
	panic(fmt.Errorf("not implemented: ConnectionRequestStatusUpdate - connectionRequestStatusUpdate"))
}

// NewMatchingVacancy is the resolver for the newMatchingVacancy field.
func (r *subscriptionResolver) NewMatchingVacancy(ctx context.Context, userID string) (<-chan *graph_model.Vacancy, error) {
	panic(fmt.Errorf("not implemented: NewMatchingVacancy - newMatchingVacancy"))
}

// NewNotification is the resolver for the newNotification field.
func (r *subscriptionResolver) NewNotification(ctx context.Context, forUserID string) (<-chan *graph_model.Notification, error) {
	panic(fmt.Errorf("not implemented: NewNotification - newNotification"))
}

// Connections is the resolver for the connections field.
func (r *userResolver) Connections(ctx context.Context, obj *graph_model.User) ([]*graph_model.User, error) {
	panic(fmt.Errorf("not implemented: Connections - connections"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
