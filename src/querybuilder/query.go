package querybuilder

import (
	"fmt"
	"strings"
)

const PREFIXES = `PREFIX lr: <http://linkrec.example.org/schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX list: <http://jena.hpl.hp.com/ARQ/list#>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX skos: <http://www.w3.org/2004/02/skos/core#>`

type Query struct {
	sel     *Select
	where   *Where
	groupBy *GroupBy
}

type Select struct {
	fields       []string
	groupConcats []GroupConcat
}

type GroupConcat struct {
	sep      string
	field    string
	as       string
	distinct bool
}

type Where struct {
	clauses         []SubWhere
	optionalClauses []OptionalClause
	extractions     []WhereExtraction
	filters         []Filter
	binds           []Bind
	subQueries      []string
}

type SubWhere struct {
	subj    *WhereSubject
	clauses []WhereClause
}

type WhereExtraction struct {
	field     string
	attribute string
	binding   string
}

type WhereClause struct {
	field   string
	binding string
}

type WhereSubject struct {
	binding string
	obj     string
}

type GroupBy struct {
	fields []string
}

type Filter struct {
	field            string
	value            []string
	op               FilterOp
	opWithPrevFilter FilterType
}

type Bind struct {
	field string
	alias string
}

type FilterType string

const (
	AND FilterType = "&&"
	OR             = "||"
)

type FilterOp string

const (
	LE FilterOp = "<="
	LT          = "<"
	GE          = ">="
	GT          = ">"
	EQ          = "="
	IN          = "IN"
)

type OptionalClause struct {
	clauses []WhereExtraction
}

func QueryBuilder() *Query {
	return &Query{
		sel:     &Select{},
		where:   &Where{},
		groupBy: &GroupBy{},
	}
}

func (q *Query) Select(fields []string) *Query {
	for _, element := range fields {
		q.sel.fields = append(q.sel.fields, element)
	}
	return q
}

func (q *Query) GroupConcat(field string, sep string, as string, distinct bool) *Query {
	groupConcat := GroupConcat{}
	groupConcat.field = field
	groupConcat.sep = sep
	groupConcat.as = as
	groupConcat.distinct = distinct
	q.sel.groupConcats = append(q.sel.groupConcats, groupConcat)
	return q
}

func (q *Query) WhereSubject(binding string, obj string) *Query {
	var sub = SubWhere{
		subj: &WhereSubject{binding: binding, obj: obj},
	}
	q.where.clauses = append(q.where.clauses, sub)
	return q
}

func (q *Query) Where(field string, binding string) *Query {
	if len(q.where.clauses) == 0 {
		q.where.clauses = append(q.where.clauses, SubWhere{
			subj: &WhereSubject{binding: "user", obj: "User"},
		})
	}

	var c = WhereClause{
		field:   field,
		binding: binding,
	}
	var curr = &q.where.clauses[len(q.where.clauses)-1]
	curr.clauses = append(curr.clauses, c)

	return q
}

func (q *Query) WhereExtraction(field string, attr string, binding string) *Query {
	var ex = WhereExtraction{
		field:     field,
		binding:   binding,
		attribute: attr,
	}
	q.where.extractions = append(q.where.extractions, ex)
	return q
}

func (q *Query) WhereSubQuery(sub string) *Query {
	q.where.subQueries = append(q.where.subQueries, sub)
	return q
}

func (q *Query) GroupBy(fields []string) *Query {
	q.groupBy.fields = fields
	return q
}

func (q *Query) Bind(field string, alias string) *Query {
	q.where.binds = append(q.where.binds, Bind{field, alias})
	return q
}

func (q *Query) Filter(field string, value []string, op FilterOp) *Query {
	q.where.filters = append(
		q.where.filters,
		Filter{
			field: field,
			value: value,
			op:    op,
		},
	)
	return q
}

func (q *Query) AndFilter(field string, value []string, op FilterOp) *Query {
	q.where.filters = append(
		q.where.filters,
		Filter{
			field:            field,
			value:            value,
			op:               op,
			opWithPrevFilter: AND,
		},
	)
	return q
}

func (q *Query) OrFilter(field string, value []string, op FilterOp) *Query {
	q.where.filters = append(
		q.where.filters,
		Filter{
			field:            field,
			value:            value,
			op:               op,
			opWithPrevFilter: OR,
		},
	)
	return q
}

func (q *Query) NewOptional(obj string, attr string, binding string) *Query {
	we := WhereExtraction{field: obj, attribute: attr, binding: binding}
	q.where.optionalClauses = append(q.where.optionalClauses, OptionalClause{
		clauses: []WhereExtraction{we},
	})
	return q
}

func (q *Query) AddOptionalTriple(field string, attr string, bind string) *Query {
	clauses := q.where.optionalClauses
	currClause := &clauses[len(clauses)-1]
	currClause.clauses = append(currClause.clauses, WhereExtraction{
		field:     field,
		attribute: attr,
		binding:   bind,
	})
	return q
}

func (q *Query) Build() string {
	var output = PREFIXES + "\n\n"
	output += buildSelect(*q.sel) + "\n"
	output += buildWhere(*q.where) + "\n"
	if len(q.groupBy.fields) > 0 {
		output += buildGroupBy(*q.groupBy) + "\n"
	}
	return output
}

func (q *Query) BuildSubQuery() string {
	var output = "{\n"
	output += buildSelect(*q.sel) + "\n"
	output += buildWhere(*q.where) + "\n"
	output += buildGroupBy(*q.groupBy) + "\n"
	return output + "}"
}

func buildSelect(sel Select) string {
	var output = "SELECT"
	for _, el := range sel.fields {
		output += " ?" + el
	}
	for _, groupConcat := range sel.groupConcats {
		if groupConcat.distinct {
			output += fmt.Sprintf(" (GROUP_CONCAT(DISTINCT ?%s; separator=\"%s\") AS ?%s)", groupConcat.field, groupConcat.sep, groupConcat.as)
		} else {
			output += fmt.Sprintf(" (GROUP_CONCAT(?%s; separator=\"%s\") AS ?%s)", groupConcat.field, groupConcat.sep, groupConcat.as)
		}
	}
	return output
}

func buildWhere(wh Where) string {
	var output = "WHERE {\n"
	for _, sub := range wh.clauses {
		var subwhere = fmt.Sprintf("?%s a lr:%s ;\n", sub.subj.binding, sub.subj.obj)
		for idx, cl := range sub.clauses {
			subwhere += fmt.Sprintf("lr:%s ?%s", cl.field, cl.binding)
			if (idx) == len(sub.clauses)-1 {
				subwhere += " ."
			} else {
				subwhere += " ;"
			}
			subwhere += "\n"
		}
		output += subwhere
	}
	for _, ex := range wh.extractions {
		output += fmt.Sprintf("?%s lr:%s ?%s .\n", ex.field, ex.attribute, ex.binding)
	}
	for _, sub := range wh.subQueries {
		output += sub + "\n"
	}

	for _, opt := range wh.optionalClauses {
		output += "OPTIONAL {\n"
		for _, cl := range opt.clauses {
			output += fmt.Sprintf("?%s %s ?%s .\n", cl.field, cl.attribute, cl.binding)
		}
		output += "}\n"
	}

	if len(wh.binds) > 0 {
		output += buildBinds(wh.binds) + "\n"
	}
	if len(wh.filters) > 0 {
		output += buildFilter(wh.filters) + "\n"
	}
	output += "}"
	return output
}

func buildGroupBy(gb GroupBy) string {
	var output = "GROUP BY"
	for _, el := range gb.fields {
		output += " ?" + el
	}
	return output
}

func buildFilter(filters []Filter) string {
	var output = ""
	for _, fil := range filters {
		var filter = "FILTER("
		if fil.op == IN {
			filter += fmt.Sprintf("?%s IN (%s)", fil.field, strings.Join(fil.value, ", "))
		} else {
			if len(fil.opWithPrevFilter) == 0 {
				if fil.value[0] == "\"en\"" {
					filter += fmt.Sprintf("LANG(?%s) %s %s", fil.field, fil.op, fil.value[0])
				} else {
					filter += fmt.Sprintf("?%s %s %s", fil.field, fil.op, fil.value[0])
				}
			} else {
				filter += fmt.Sprintf(" %s ?%s %s %s", fil.opWithPrevFilter, fil.field, fil.op, fil.value[0])
			}
		}
		output += filter + ")\n"
	}
	return output
}

func buildBinds(binds []Bind) string {
	var output = ""
	for _, bin := range binds {
		output += fmt.Sprintf("BIND(?%s AS ?%s)", bin.field, bin.alias)
	}
	return output
}

func main() {
	q :=
		QueryBuilder().Select([]string{"id", "title", "description", "location", "postedById", "startDate", "endDate", "status", "degreeType", "degreeField"}).
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
			Where("requiredSkill", "skill").
			WhereExtraction("postedBy", "Id", "postedById")
	q.Filter("name", []string{"name"}, EQ)
	// q.Filter("requiredEducation", []string{string(*requiredEducation)}, EQ)
	q.Filter("status", []string{"true"}, EQ)
	qs := q.GroupBy([]string{"id", "title", "description", "location", "postedById", "startDate", "endDate", "status", "degreeType", "degreeField"}).Build()

	fmt.Println(qs)

}
