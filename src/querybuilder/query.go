package main

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
