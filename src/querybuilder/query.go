package querybuilder

import (
	"fmt"
)

const PREFIXES = `PREFIX lr: <http://linkrec.example.org/schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX list: <http://jena.hpl.hp.com/ARQ/list#>`

type Query struct {
	sel     *Select
	where   *Where
	groupBy *GroupBy
}

type Select struct {
	fields      []string
	groupConcat *GroupConcat
}

type GroupConcat struct {
	sep   string
	field string
	as    string
}

type Where struct {
	clauses []SubWhere
	filters []Filter
}

type SubWhere struct {
	subj    *WhereSubject
	clauses []WhereClause
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
	value            string
	op               FilterOp
	opWithPrevFilter FilterType
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
)

func QueryBuilder() *Query {
	return &Query{
		sel: &Select{
			groupConcat: &GroupConcat{},
		},
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

func (q *Query) GroupConcat(field string, sep string, as string) *Query {
	q.sel.groupConcat.field = field
	q.sel.groupConcat.sep = sep
	q.sel.groupConcat.as = as
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
	var c = WhereClause{
		field:   field,
		binding: binding,
	}
	var curr = &q.where.clauses[len(q.where.clauses)-1]
	curr.clauses = append(curr.clauses, c)
	return q
}

func (q *Query) GroupBy(fields []string) *Query {
	q.groupBy.fields = fields
	return q
}

func (q *Query) Filter(field string, value string, op FilterOp) *Query {
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

func (q *Query) AndFilter(field string, value string, op FilterOp) *Query {
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

func (q *Query) OrFilter(field string, value string, op FilterOp) *Query {
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

func (q *Query) Build() string {
	var output = PREFIXES + "\n\n"
	output += buildSelect(*q.sel) + "\n"
	output += buildWhere(*q.where) + "\n"
	output += buildGroupBy(*q.groupBy) + "\n"
	return output
}

func buildSelect(sel Select) string {
	var output = "SELECT"
	for _, el := range sel.fields {
		output += " ?" + el
	}
	var concat = sel.groupConcat
	if concat != nil {
		output += fmt.Sprintf(" (GROUP_CONCAT(?%s; separator=\"%s\") AS ?%s)", concat.field, concat.sep, concat.as)
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
	output += buildFilter(wh.filters) + "\n"
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
	if len(filters) == 0 {
		return ""
	}
	var output = "FILTER("
	for _, fil := range filters {
		if fil.opWithPrevFilter == "" {
			output += fmt.Sprintf("?%s %s \"%s\"", fil.field, fil.op, fil.value)
		} else {
			output += fmt.Sprintf(" %s ?%s %s \"%s\"", fil.opWithPrevFilter, fil.field, fil.op, fil.value)
		}
	}
	return output + ")"
}

func main() {
	var s = QueryBuilder().
		Select([]string{"userId", "userName"}).
		GroupConcat("skill", ", ", "skills").
		WhereSubject("user", "User").
		Where("Id", "userId").
		Where("hasName", "userName").
		Where("hasSkill", "skillList").
		Filter("userId", "1", EQ).
		OrFilter("userId", "6", GT).
		GroupBy([]string{"userName", "userId"}).
		Build()
	fmt.Println(s)
}
