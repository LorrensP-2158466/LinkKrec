package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"LinkKrec/graph/loaders"
	"LinkKrec/graph/model"
	"LinkKrec/graph/util"
	query_builder "LinkKrec/querybuilder"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/knakk/rdf"
)

// Location is the resolver for the location field.
func (r *companyResolver) Location(ctx context.Context, obj *model.Company) (*model.Location, error) {
	return loaders.GetLocation(ctx, obj.Location.ID)
}

// Vacancies is the resolver for the vacancies field.
func (r *companyResolver) Vacancies(ctx context.Context, obj *model.Company) ([]*model.Vacancy, error) {
	ids := util.Map(obj.Vacancies, func(v *model.Vacancy) string {
		return v.ID
	})
	return loaders.GetVacancies(ctx, ids)
}

// Employees is the resolver for the employees field.
func (r *companyResolver) Employees(ctx context.Context, obj *model.Company) ([]*model.User, error) {
	ids := util.Map(obj.Employees, func(u *model.User) string {
		return u.ID
	})
	return loaders.GetUsers(ctx, ids)
}

// FromUser is the resolver for the fromUser field.
func (r *connectionRequestResolver) FromUser(ctx context.Context, obj *model.ConnectionRequest) (*model.User, error) {
	return loaders.GetUser(ctx, obj.FromUser.ID)
}

// ConnectedToUser is the resolver for the connectedToUser field.
func (r *connectionRequestResolver) ConnectedToUser(ctx context.Context, obj *model.ConnectionRequest) (*model.User, error) {
	return loaders.GetUser(ctx, obj.ConnectedToUser.ID)
}

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented: RegisterUser - registerUser"))
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
}

// UpdateUserProfile is the resolver for the updateUserProfile field.
func (r *mutationResolver) UpdateUserProfile(ctx context.Context, id string, input model.UpdateProfileInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented: UpdateUserProfile - updateUserProfile"))
}

// AddConnectionRequest is the resolver for the addConnectionRequest field.
func (r *mutationResolver) AddConnectionRequest(ctx context.Context, fromUserID string, connectedToUserID string) (*model.ConnectionRequest, error) {
	requestID := uuid.New().String()

	q := fmt.Sprintf(`
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX lr: <http://linkrec.example.org/schema#>

		INSERT {
		  lr:connectionRequest%s a lr:ConnectionRequest ;
		      lr:Id "%s" ;
		      lr:fromUser ?fromUser ;
		      lr:connectedToUser ?toUser ;
		      lr:requestStatus false .
		}
		WHERE {
		  ?fromUser a lr:User ;
		            lr:Id "%s" .
		  ?toUser a lr:User ;
		          lr:Id "%s" .
		}
		`, requestID, requestID, fromUserID, connectedToUserID)

	err := r.UpdateRepo.Update(q)
	if err != nil {
		return nil, err
	}
	fmt.Println("err:", err)

	// If the query was successful, return the updated user
	return loaders.GetConnectionRequest(ctx, requestID)
}

// SetConnectionRequestStatusFalse is the resolver for the setConnectionRequestStatusFalse field.
func (r *mutationResolver) SetConnectionRequestStatusFalse(ctx context.Context, id string) (*model.ConnectionRequest, error) {
	q := fmt.Sprintf(`
        PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
        PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
        PREFIX lr: <http://linkrec.example.org/schema#>
        PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

        DELETE {
            ?connectionRequest lr:status ?status .
        }
        INSERT {
            ?connectionRequest lr:status false .
        }
        WHERE {
            ?connectionRequest a lr:ConnectionRequest ;
                               lr:Id "%s" ;
                               lr:status ?status .
        }
    `, id)

	err := r.UpdateRepo.Update(q)
	if err != nil {
		return nil, err
	}

	return loaders.GetConnectionRequest(ctx, id)
}

// NotifyProfileVisit is the resolver for the notifyProfileVisit field.
func (r *mutationResolver) NotifyProfileVisit(ctx context.Context, visitorID string, visitedUserID string) (*model.Notification, error) {
	notificationID := uuid.New().String()

	visitedByUser, err := loaders.GetUser(ctx, visitorID)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

		INSERT DATA {
		  lr:notification%s a lr:Notification ;
		      lr:Id "%s" ;
		      lr:notificationTitle "Profile Visit" ;
		      lr:notificationMessage "Your profile has been visited by %s" ;
		      lr:forUser lr:User%s ;
		      lr:notificationCreatedAt "%s"^^xsd:dateTime .
		}
	`, notificationID, notificationID, visitorID, visitedByUser.Name, time.Now())

	updErr := r.UpdateRepo.Update(q)
	if updErr != nil {
		return nil, updErr
	}

	return loaders.GetNotification(ctx, notificationID)
}

// CreateVacancy is the resolver for the createVacancy field.
func (r *mutationResolver) CreateVacancy(ctx context.Context, companyID string, input model.CreateVacancyInput) (*model.Vacancy, error) {
	vacancyID := uuid.New().String()

	// TO DO: add location shit so it isn't a string but an object made with gisco.
	// Address is in input.Location

	var skillQueries string
	if input.RequiredSkills != nil {
		skillQueries = ""
		for _, skill := range input.RequiredSkills {
			skillQueries += fmt.Sprintf(" ;\nlr:requiredSkill esco_skill:%s", *skill)
		}
		skillQueries += " .\n"
	} else {
		skillQueries = ".\n"
	}

	q := fmt.Sprintf(`
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

		INSERT DATA {
		  lr:vacancy%s a lr:Vacancy ;
		      lr:Id "%s" ;
		      lr:vacancyTitle "%s" ;
			  lr:vacancyDescription "%s" ;
			  lr:vacancyLocation "%s" ;
			  lr:postedBy lr:Company%s ;
			  lr:vacancyStartDate "%s"^^xsd:date ;
			  lr:vacancyEndDate "%s"^^xsd:date ;
			  lr:vacancyStatus %t;
		      lr:requiredDegreeType lr:%s ;
			  lr:requiredDegreeField lr:%s ;
			  lr:requiredExperienceDuration %d %s
		}
		`, vacancyID, vacancyID, input.Title, input.Description, input.Location, companyID, input.StartDate, input.EndDate, input.Status, input.RequiredDegreeType, input.RequiredDegreeField, input.RequiredExperienceDuration, skillQueries)

	err := r.UpdateRepo.Update(q)
	if err != nil {
		return nil, err
	}

	// If the query was successful, return the updated user
	return loaders.GetVacancy(ctx, vacancyID)
}

// UpdateVacancy is the resolver for the updateVacancy field.
func (r *mutationResolver) UpdateVacancy(ctx context.Context, id string, input model.UpdateVacancyInput) (*model.Vacancy, error) {
	var deleteParts, insertParts string

	// Conditionally add fields to the DELETE/INSERT sections
	if input.Title != nil {
		deleteParts += "?vacancy lr:vacancyTitle ?oldTitle .\n"
		insertParts += fmt.Sprintf("?vacancy lr:vacancyTitle \"%s\" .\n", *input.Title)
	}
	if input.Description != nil {
		deleteParts += "?vacancy lr:vacancyDescription ?oldDescription .\n"
		insertParts += fmt.Sprintf("?vacancy lr:vacancyDescription \"%s\" .\n", *input.Description)
	}
	if input.Location != nil {
		deleteParts += "?vacancy lr:vacancyLocation ?oldLocation .\n"
		insertParts += fmt.Sprintf("?vacancy lr:vacancyLocation \"%s\" .\n", *input.Location)
	}
	if input.StartDate != nil {
		deleteParts += "?vacancy lr:vacancyStartDate ?oldStartDate .\n"
		insertParts += fmt.Sprintf("?vacancy lr:vacancyStartDate \"%s\"^^xsd:date .\n", *input.StartDate)
	}
	if input.EndDate != nil {
		deleteParts += "?vacancy lr:vacancyEndDate ?oldEndDate .\n"
		insertParts += fmt.Sprintf("?vacancy lr:vacancyEndDate \"%s\"^^xsd:date .\n", *input.EndDate)
	}
	if input.RequiredDegreeType != nil {
		deleteParts += "?vacancy lr:requiredDegreeType ?oldDegreeType .\n"
		insertParts += fmt.Sprintf("?vacancy lr:requiredDegreeType lr:%s .\n", *input.RequiredDegreeType)
	}
	if input.RequiredDegreeField != nil {
		deleteParts += "?vacancy lr:requiredDegreeField ?oldDegreeField .\n"
		insertParts += fmt.Sprintf("?vacancy lr:requiredDegreeField lr:%s .\n", *input.RequiredDegreeField)
	}
	if input.RequiredExperienceDuration != nil {
		deleteParts += "?vacancy lr:requiredExperienceDuration ?oldExperience .\n"
		insertParts += fmt.Sprintf("?vacancy lr:requiredExperienceDuration %d .\n", *input.RequiredExperienceDuration)
	}
	// check if the status is provided
	if input.Status != nil {
		deleteParts += "?vacancy lr:vacancyStatus ?oldStatus .\n"
		insertParts += fmt.Sprintf("?vacancy lr:vacancyStatus %t .\n", *input.Status)
	}
	if len(input.RequiredSkills) != 0 {
		deleteParts += "?vacancy lr:requiredSkill ?oldSkill .\n"
		for _, skill := range input.RequiredSkills {
			insertParts += fmt.Sprintf("?vacancy lr:requiredSkill esco_skill:%s .\n", *skill)
		}
	}

	// If no fields are provided, return an error
	if deleteParts == "" && insertParts == "" {
		return nil, fmt.Errorf("no fields provided for update")
	}

	// Construct the full SPARQL query
	q := fmt.Sprintf(`
        PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
        PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
        PREFIX lr: <http://linkrec.example.org/schema#>
        PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

        DELETE {
            %s
        }
        INSERT {
            %s
        }
        WHERE {
            ?vacancy a lr:Vacancy ;
                     lr:Id "%s" .
            %s
        }
    `, deleteParts, insertParts, id, deleteParts)
	fmt.Println(q)

	err := r.UpdateRepo.Update(q)
	if err != nil {
		return nil, err
	}

	return loaders.GetVacancy(ctx, id)
}

// DeleteVacancy is the resolver for the deleteVacancy field.
func (r *mutationResolver) DeleteVacancy(ctx context.Context, id string) (*bool, error) {
	q := fmt.Sprintf(`
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>
		
		DELETE {
          ?vacancy ?p ?o .
          ?location ?lp ?lo .
        }
        WHERE {
          ?vacancy a lr:Vacancy ;
                   lr:Id "%s" ;
                   ?p ?o .
          OPTIONAL {
            ?vacancy lr:vacancyLocation ?location .
            ?location ?lp ?lo .
          }
        }
	`, id)

	err := r.UpdateRepo.Update(q)
	if err != nil {
		succes := false
		return &succes, err
	}

	success := true
	return &success, nil
}

// UpdateUserLookingForOpportunities is the resolver for the updateUserLookingForOpportunities field.
func (r *mutationResolver) UpdateUserLookingForOpportunities(ctx context.Context, userID string, looking bool) (*model.User, error) {
	// Convert the `looking` boolean to a string representation
	lookingStr := strconv.FormatBool(looking)

	// SPARQL update query to change the isLookingForOpportunities value
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>

        DELETE {
            ?user lr:isLookingForOpportunities ?currentValue .
        }
        INSERT {
            ?user lr:isLookingForOpportunities "%s" .
        }
        WHERE {
            ?user a lr:User ;
            lr:Id "%s" ;
            lr:isLookingForOpportunities ?currentValue .
        }
    `, lookingStr, userID)

	err := r.UpdateRepo.Update(q)
	if err != nil {
		return nil, err
	}
	fmt.Println(err)

	// If the query was successful, return the updated user
	return loaders.GetUser(ctx, userID)
}

// ForUser is the resolver for the forUser field.
func (r *notificationResolver) ForUser(ctx context.Context, obj *model.Notification) (*model.User, error) {
	return loaders.GetUser(ctx, obj.ForUser.ID)
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	return loaders.GetUser(ctx, id)
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context, name *string, location *string, skills []*string, lookingForOpportunities *bool) ([]*model.User, error) {
	q := query_builder.
		QueryBuilder().
		Select([]string{"id", "name", "email", "locationId", "lookingForOpportunities"}).
		GroupConcat("skill", ", ", "skills", true).
		GroupConcat("connectionName", ", ", "connections", true).
		GroupConcat("educationEntry", ", ", "educations", true).
		GroupConcat("companyId", ", ", "companies", true).
		WhereSubject("user", "User").
		Where("Id", "id").
		Where("hasName", "name").
		Where("hasEmail", "email").
		Where("isLookingForOpportunities", "isLookingForOpportunities").
		Bind("isLookingForOpportunities", "lookingForOpportunities").
		NewOptional("user", "lr:hasLocation", "location").
		AddOptionalTriple("location", "lr:Id", "locationId").
		NewOptional("user", "lr:hasSkill", "escoSkill").
		AddOptionalTriple("escoSkill", "skos:prefLabel", "skill").
		NewOptional("user", "lr:hasConnection", "connection").
		AddOptionalTriple("connection", "lr:Id", "connectionName").
		NewOptional("user", "lr:hasEducation", "education").
		AddOptionalTriple("education", "lr:Id", "educationEntry").
		NewOptional("user", "lr:hasCompany", "company").
		AddOptionalTriple("company", "lr:Id", "companyId").
		Filter("skill", []string{"\"en\""}, query_builder.EQ)
	if name != nil {
		q.Filter("name", []string{*name}, query_builder.EQ)
	}
	if location != nil {
		q.Filter("location", []string{*location}, query_builder.EQ)
	}
	if len(skills) > 0 {
		convSkills := util.Map(skills, func(s *string) string {
			return fmt.Sprintf("\"%s\"", *s)
		})
		q.Filter("skill", convSkills, query_builder.IN)
	}
	if lookingForOpportunities != nil {
		q.Filter("isLookingForOpportunities", []string{strconv.FormatBool(*lookingForOpportunities)}, query_builder.EQ)
	}
	qs := q.GroupBy([]string{"id", "name", "email", "locationId", "lookingForOpportunities"}).Build()
	res, err := r.Repo.Query(qs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	users := make([]*model.User, 0)
	for _, user := range res.Solutions() {
		obj, err := util.MapRdfUserToGQL(user)
		if err != nil {
			return nil, err
		}
		users = append(users, obj)
	}
	return users, nil
}

// GetVacancies is the resolver for the getVacancies field.
func (r *queryResolver) GetVacancies(ctx context.Context, title *string, location *string, requiredEducation *model.DegreeType, status *bool) ([]*model.Vacancy, error) {
	q := query_builder.
		QueryBuilder().Select([]string{"id", "title", "description", "location", "postedById", "startDate", "endDate", "status", "degreeType", "degreeField", "experienceDuration"}).
		GroupConcat("skill", ", ", "skills", true).
		WhereSubject("vacancy", "Vacancy").
		Where("Id", "id").
		Where("vacancyTitle", "title").
		Where("vacancyDescription", "description").
		Where("vacancyLocation", "location").
		Where("postedBy", "postedBy").
		Where("vacancyStartDate", "startDate").
		Where("vacancyEndDate", "endDate").
		Where("vacancyStatus", "status").
		Where("requiredDegreeType", "degreeType").
		Where("requiredDegreeField", "degreeField").
		Where("requiredExperienceDuration", "experienceDuration").
		Where("requiredSkill", "skill").
		WhereExtraction("postedBy", "Id", "postedById")
	if title != nil {
		q.Filter("name", []string{*title}, query_builder.EQ)
	}
	if location != nil {
		q.Filter("location", []string{*location}, query_builder.EQ)
	}
	if requiredEducation != nil {
		q.Filter("requiredEducation", []string{string(*requiredEducation)}, query_builder.EQ)
	}
	if status != nil {
		q.Filter("status", []string{strconv.FormatBool(*status)}, query_builder.EQ)
	}
	qs := q.GroupBy([]string{"id", "title", "description", "location", "postedById", "startDate", "endDate", "status", "degreeType", "degreeField", "experienceDuration"}).Build()

	fmt.Println(qs)
	res, err := r.Repo.Query(qs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	vacancies := make([]*model.Vacancy, 0)
	for _, user := range res.Solutions() {
		obj, err := util.MapRdfVacancyToGQL(user)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, obj)
	}
	return vacancies, nil
}

// GetVacancy is the resolver for the getVacancy field.
func (r *queryResolver) GetVacancy(ctx context.Context, id string) (*model.Vacancy, error) {
	return loaders.GetVacancy(ctx, id)
}

// GetCompanies is the resolver for the getCompanies field.
func (r *queryResolver) GetCompanies(ctx context.Context, name *string, location *string) ([]*model.Company, error) {
	q := query_builder.
		QueryBuilder().Select([]string{"id", "name", "email", "locationId"}).
		GroupConcat("vacancyId", ", ", "vacancies", true).
		GroupConcat("employeeId", ", ", "employees", true).
		WhereSubject("company", "Company").
		Where("Id", "id").
		Where("companyName", "name").
		Where("companyEmail", "email").
		NewOptional("company", "lr:companyLocation", "location").
		AddOptionalTriple("location", "lr:Id", "locationId").
		NewOptional("company", "lr:hasVacancy", "vacancy").
		AddOptionalTriple("vacancy", "lr:Id", "vacancyId").
		NewOptional("company", "lr:hasEmployee", "employee").
		AddOptionalTriple("employee", "lr:Id", "employeeId")
	if name != nil {
		q.Filter("name", []string{*name}, query_builder.EQ)
	}
	if location != nil {
		q.Filter("location", []string{*location}, query_builder.EQ)
	}
	qs := q.GroupBy([]string{"id", "name", "email", "locationId"}).Build()

	res, err := r.Repo.Query(qs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	companies := make([]*model.Company, 0)
	for _, company := range res.Solutions() {
		obj, err := util.MapRdfCompanyToGQL(company)
		if err != nil {
			return nil, err
		}
		companies = append(companies, obj)
	}
	return companies, nil
}

// GetCompany is the resolver for the getCompany field.
func (r *queryResolver) GetCompany(ctx context.Context, id string) (*model.Company, error) {
	return loaders.GetCompany(ctx, id)
}

// GetNotifications is the resolver for the getNotifications field.
func (r *queryResolver) GetNotifications(ctx context.Context, userID string) ([]*model.Notification, error) {
	q := query_builder.
		QueryBuilder().Select([]string{"id", "title", "message", "forUserId", "createdAt"}).
		WhereSubject("notification", "Notification").
		Where("Id", "id").
		Where("notificationTitle", "title").
		Where("notificationMessage", "message").
		Where("forUser", "forUser").
		NewOptional("notification", "lr:notificationCreatedAt", "createdAt").
		WhereExtraction("forUser", "Id", "forUserId")
	if userID != "" {
		quotedUserID := fmt.Sprintf("\"%s\"", userID)
		q.Filter("forUserId", []string{quotedUserID}, query_builder.EQ)
	}
	qs := q.GroupBy([]string{"id", "title", "message", "forUserId", "createdAt"}).Build()

	res, err := r.Repo.Query(qs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	notifications := make([]*model.Notification, 0)
	for _, notification := range res.Solutions() {
		obj, err := util.MapRdfNotificationToGQL(notification)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, obj)
	}
	return notifications, nil
}

// GetConnectionRequests is the resolver for the getConnectionRequests field.
func (r *queryResolver) GetConnectionRequests(ctx context.Context, userID string, status *bool) ([]*model.ConnectionRequest, error) {
	q := query_builder.
		QueryBuilder().Select([]string{"id", "fromUserId", "connectedToUserId", "status"}).
		WhereSubject("connectionRequest", "ConnectionRequest").
		Where("Id", "id").
		Where("fromUser", "user").
		Where("connectedToUser", "connectedTo").
		Where("requestStatus", "status").
		WhereExtraction("user", "Id", "fromUserId").
		WhereExtraction("connectedTo", "Id", "connectedToUserId")
	if userID != "" {
		quotedUserID := fmt.Sprintf("\"%s\"", userID)
		q.Filter("connectedToUserId", []string{quotedUserID}, query_builder.EQ)
	}
	if status != nil && userID != "" {
		q.AndFilter("status", []string{strconv.FormatBool(*status)}, query_builder.EQ)
	} else if status != nil {
		q.Filter("status", []string{strconv.FormatBool(*status)}, query_builder.EQ)
	}
	qs := q.GroupBy([]string{"id", "fromUserId", "connectedToUserId", "status"}).Build()

	res, err := r.Repo.Query(qs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	connectionRequests := make([]*model.ConnectionRequest, 0)
	for _, connectionRequest := range res.Solutions() {
		obj, err := util.MapRdfConnectionRequestToGQL(connectionRequest)
		if err != nil {
			return nil, err
		}
		connectionRequests = append(connectionRequests, obj)
	}
	return connectionRequests, nil
}

// MatchVacancyToUsers is the resolver for the matchVacancyToUsers field.
func (r *queryResolver) MatchVacancyToUsers(ctx context.Context, vacancyID string, maxDist float64) ([]*model.User, error) {
	q := fmt.Sprintf(`
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		PREFIX list: <http://jena.hpl.hp.com/ARQ/list#>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX lfn: <http://www.dotnetrdf.org/leviathan#>

		SELECT DISTINCT ?userId ?distanceInKm
		WHERE {
		# constants
			VALUES (?pi ?earthRadius ) { ( 3.1415926535 6378.137 ) }
		
			?vacancy lr:Id "%s" ;
				lr:vacancyLocation ?vacancyLoc ;
				OPTIONAL { ?vacancy lr:requiredSkill ?requiredSkill }
				OPTIONAL { ?vacancy lr:requiredDegreeType ?requiredDegreeType }
				OPTIONAL { ?vacancy lr:requiredDegreeField ?requiredDegreeField }


			
			?vacancyLoc lr:longitude ?long2 ;
				lr:latitude ?lat2 .
				
			?user lr:Id ?userId ;
				lr:hasEmail ?userEmail ;
				lr:hasName ?userName ;
				lr:hasLocation ?userLoc;
				lr:hasSkill ?requiredSkill;  # Direct match instead of FILTER
				lr:hasEducation ?education ;
				lr:isProfileComplete true ;
				lr:isLookingForOpportunities true .
				
			?userLoc lr:longitude ?long1 ;
				lr:latitude ?lat1 .
				
			?education lr:degreeType ?userDegreeType;
				lr:degreeField ?userDegreeField.

			?userDegreeType rdfs:subClassOf* ?requiredDegreeType .
			?userDegreeField rdfs:subClassOf* ?requiredDegreeField .
			
			# haversine
			BIND(?earthRadius * 2 * lfn:sin-1(lfn:sqrt(
				lfn:pow(lfn:sin((?lat2 - ?lat1) * ?pi / 360), 2) +
				lfn:cos(?lat1 * ?pi / 180) * lfn:cos(?lat2 * ?pi / 180) *
				lfn:pow(lfn:sin((?long2 - ?long1) * ?pi / 360), 2)
			)) AS ?distanceInKm)
		
			FILTER(?distanceInKm <= %f)
		}
	`, vacancyID, maxDist)

	res, _ := r.Repo.Query(q)

	userIds := res.Bindings()["userId"]
	var users = make([]*model.User, 0)
	if userIds != nil {
		ids := util.Map(userIds, func(u rdf.Term) string {
			return u.String()
		})
		users, _ = loaders.GetUsers(ctx, ids)
	}
	return users, nil
}

// MatchUserToVacancies is the resolver for the matchUserToVacancies field.
func (r *queryResolver) MatchUserToVacancies(ctx context.Context, userID string, maxDist float64) ([]*model.Vacancy, error) {
	q := fmt.Sprintf(`
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		PREFIX list: <http://jena.hpl.hp.com/ARQ/list#>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX lfn: <http://www.dotnetrdf.org/leviathan#>

		SELECT DISTINCT ?vacancyId ?matchedSkill ?distanceInKm WHERE {
			# constants
			VALUES (?pi ?earthRadius) {
				( 3.1415926535 6378.137 )
			}

			?user lr:Id "%s" ;
				lr:hasLocation ?userLoc ;
				lr:hasSkill ?userSkill ;
				lr:hasEducation ?education .

			?userLoc lr:longitude ?long1 ;
					lr:latitude ?lat1 .

			?education lr:degreeType ?userDegreeType ;
					lr:degreeField ?userDegreeField .

			# Vacancy matching
			?vacancy lr:Id ?vacancyId ;
				lr:vacancyLocation ?vacancyLoc ;
			OPTIONAL { ?vacancy lr:requiredSkill ?userSkill }
			OPTIONAL { ?vacancy lr:requiredDegreeType ?requiredDegreeType }
			OPTIONAL { ?vacancy lr:requiredDegreeField ?requiredDegreeField }

			?vacancyLoc lr:longitude ?long2 ;
						lr:latitude ?lat2 .

			?requiredDegreeType rdfs:subClassOf* ?userDegreeType .
			?requiredDegreeField rdfs:subClassOf* ?userDegreeField .

			BIND(?earthRadius * 2 * lfn:sin-1(lfn:sqrt(
				lfn:pow(lfn:sin((?lat2 - ?lat1) * ?pi / 360), 2) +
				lfn:cos(?lat1 * ?pi / 180) * lfn:cos(?lat2 * ?pi / 180) *
				lfn:pow(lfn:sin((?long2 - ?long1) * ?pi / 360), 2)
			)) AS ?distanceInKm)

			FILTER(?distanceInKm <= %f)
		}
	`, userID, maxDist)

	res, _ := r.Repo.Query(q)

	vacancyIds := res.Bindings()["vacancyId"]
	var vacancies = make([]*model.Vacancy, 0)
	if vacancyIds != nil {
		ids := util.Map(vacancyIds, func(u rdf.Term) string {
			return u.String()
		})
		vacancies, _ = loaders.GetVacancies(ctx, ids)
	}
	return vacancies, nil
}

// NewConnectionRequest is the resolver for the newConnectionRequest field.
func (r *subscriptionResolver) NewConnectionRequest(ctx context.Context, forUserID string) (<-chan *model.ConnectionRequest, error) {
	panic(fmt.Errorf("not implemented: NewConnectionRequest - newConnectionRequest"))
}

// ConnectionRequestStatusUpdate is the resolver for the connectionRequestStatusUpdate field.
func (r *subscriptionResolver) ConnectionRequestStatusUpdate(ctx context.Context, forUserID string) (<-chan *model.ConnectionRequest, error) {
	panic(fmt.Errorf("not implemented: ConnectionRequestStatusUpdate - connectionRequestStatusUpdate"))
}

// NewMatchingVacancy is the resolver for the newMatchingVacancy field.
func (r *subscriptionResolver) NewMatchingVacancy(ctx context.Context, userID string) (<-chan *model.Vacancy, error) {
	panic(fmt.Errorf("not implemented: NewMatchingVacancy - newMatchingVacancy"))
}

// NewNotification is the resolver for the newNotification field.
func (r *subscriptionResolver) NewNotification(ctx context.Context, forUserID string) (<-chan *model.Notification, error) {
	panic(fmt.Errorf("not implemented: NewNotification - newNotification"))
}

// Location is the resolver for the location field.
func (r *userResolver) Location(ctx context.Context, obj *model.User) (*model.Location, error) {
	return loaders.GetLocation(ctx, obj.Location.ID)
}

// Connections is the resolver for the connections field.
func (r *userResolver) Connections(ctx context.Context, obj *model.User) ([]*model.User, error) {
	ids := util.Map(obj.Connections, func(u *model.User) string {
		return u.ID
	})
	return loaders.GetUsers(ctx, ids)
}

// Education is the resolver for the education field.
func (r *userResolver) Education(ctx context.Context, obj *model.User) ([]*model.EducationEntry, error) {
	ids := util.Map(obj.Education, func(e *model.EducationEntry) string {
		return e.ID
	})
	return loaders.GetEducationEntries(ctx, ids)
}

// Companies is the resolver for the companies field.
func (r *userResolver) Companies(ctx context.Context, obj *model.User) ([]*model.Company, error) {
	ids := util.Map(obj.Companies, func(c *model.Company) string {
		return c.ID
	})
	return loaders.GetCompanies(ctx, ids)
}

// PostedBy is the resolver for the postedBy field.
func (r *vacancyResolver) PostedBy(ctx context.Context, obj *model.Vacancy) (*model.Company, error) {
	return loaders.GetCompany(ctx, obj.PostedBy.ID)
}

// Company returns CompanyResolver implementation.
func (r *Resolver) Company() CompanyResolver { return &companyResolver{r} }

// ConnectionRequest returns ConnectionRequestResolver implementation.
func (r *Resolver) ConnectionRequest() ConnectionRequestResolver {
	return &connectionRequestResolver{r}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Notification returns NotificationResolver implementation.
func (r *Resolver) Notification() NotificationResolver { return &notificationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

// Vacancy returns VacancyResolver implementation.
func (r *Resolver) Vacancy() VacancyResolver { return &vacancyResolver{r} }

type companyResolver struct{ *Resolver }
type connectionRequestResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type notificationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type vacancyResolver struct{ *Resolver }
