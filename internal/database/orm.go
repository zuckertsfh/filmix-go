package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Operator string

const (
	OpEq   Operator = "="
	OpNe   Operator = "!="
	OpGt   Operator = ">"
	OpGte  Operator = ">="
	OpLt   Operator = "<"
	OpLte  Operator = "<="
	OpLike Operator = "LIKE"
	OpIn   Operator = "IN"
)

// --- Condition System ---

type Condition interface{}

type FieldCond struct {
	Field string
	Op    Operator
	Value any
}

type GroupCond struct {
	Op   string // AND, OR, NOT
	Args []Condition
}

func Cond(field string, cond FieldCond) FieldCond {
	cond.Field = field
	return cond
}

func Eq(v any) FieldCond   { return FieldCond{Op: OpEq, Value: v} }
func Ne(v any) FieldCond   { return FieldCond{Op: OpNe, Value: v} }
func Gt(v any) FieldCond   { return FieldCond{Op: OpGt, Value: v} }
func Gte(v any) FieldCond  { return FieldCond{Op: OpGte, Value: v} }
func Lt(v any) FieldCond   { return FieldCond{Op: OpLt, Value: v} }
func Lte(v any) FieldCond  { return FieldCond{Op: OpLte, Value: v} }
func Like(v any) FieldCond { return FieldCond{Op: OpLike, Value: v} }
func In(v any) FieldCond   { return FieldCond{Op: OpIn, Value: v} }

func And(conds ...Condition) GroupCond { return GroupCond{Op: "AND", Args: conds} }
func Or(conds ...Condition) GroupCond  { return GroupCond{Op: "OR", Args: conds} }
func Not(cond Condition) GroupCond     { return GroupCond{Op: "NOT", Args: []Condition{cond}} }

// --- Joins ---

type Join struct {
	Table    string
	On       string
	JoinType string
}

// --- Base Model ---

type BaseModel[T any] struct {
	table   string
	db      *sql.DB
	joins   []Join
	selects []string
	where   Condition // single root condition
	order   string
	limit   int
	offset  int
}

func NewModel[T any](db *sql.DB, table string) *BaseModel[T] {
	return &BaseModel[T]{db: db, table: table}
}

func (m *BaseModel[T]) Select(columns ...string) *BaseModel[T] {
	m.selects = columns
	return m
}

func (m *BaseModel[T]) Join(table, on, joinType string) *BaseModel[T] {
	m.joins = append(m.joins, Join{Table: table, On: on, JoinType: joinType})
	return m
}

func (m *BaseModel[T]) Where(cond Condition) *BaseModel[T] {
	m.where = cond
	return m
}

func (m *BaseModel[T]) OrderBy(field, direction string) *BaseModel[T] {
	m.order = fmt.Sprintf("%s %s", field, strings.ToUpper(direction))
	return m
}

func (m *BaseModel[T]) Limit(n int) *BaseModel[T] {
	m.limit = n
	return m
}

func (m *BaseModel[T]) Offset(n int) *BaseModel[T] {
	m.offset = n
	return m
}

// --- Build SQL ---

func (m *BaseModel[T]) buildQuery(single bool) (string, []any) {
	cols := "*"
	if len(m.selects) > 0 {
		cols = strings.Join(m.selects, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", cols, m.table)
	args := []any{}

	// Joins
	for _, j := range m.joins {
		query += fmt.Sprintf(" %s JOIN %s ON %s", j.JoinType, j.Table, j.On)
	}

	// Where
	if m.where != nil {
		whereSQL, whereArgs := buildCondition(m.where)
		query += " WHERE " + whereSQL
		args = append(args, whereArgs...)
	}

	// Order
	if m.order != "" {
		query += " ORDER BY " + m.order
	}

	// Limit & Offset
	if single {
		query += " LIMIT 1"
	} else {
		if m.limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", m.limit)
		}
		if m.offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", m.offset)
		}
	}

	return query, args
}

func buildCondition(c Condition) (string, []any) {
	switch cond := c.(type) {
	case FieldCond:
		if cond.Op == OpIn {
			values := cond.Value.([]any)
			ph := make([]string, len(values))
			args := make([]any, len(values))
			for i, v := range values {
				ph[i] = "?"
				args[i] = v
			}
			return fmt.Sprintf("%s IN (%s)", cond.Field, strings.Join(ph, ",")), args
		}
		return fmt.Sprintf("%s %s ?", cond.Field, cond.Op), []any{cond.Value}

	case GroupCond:
		if cond.Op == "NOT" {
			sqlPart, args := buildCondition(cond.Args[0])
			return fmt.Sprintf("NOT (%s)", sqlPart), args
		}
		var parts []string
		var args []any
		for _, sub := range cond.Args {
			sqlPart, subArgs := buildCondition(sub)
			parts = append(parts, sqlPart)
			args = append(args, subArgs...)
		}
		return "(" + strings.Join(parts, " "+cond.Op+" ") + ")", args
	}
	return "", nil
}

// --- Execute queries ---

func (m *BaseModel[T]) FindMany() ([]T, error) {
	query, args := m.buildQuery(false)
	fmt.Println("SQL:", query)
	fmt.Println("Args:", args)

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var results []T

	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		var t T
		v := reflect.ValueOf(&t).Elem()
		for i, col := range cols {
			field := v.FieldByNameFunc(func(n string) bool {
				return strings.EqualFold(n, col)
			})
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(values[i]))
			}
		}
		results = append(results, t)
	}

	return results, nil
}

func (m *BaseModel[T]) FindOne() (*T, error) {
	query, args := m.buildQuery(true)
	fmt.Println("SQL:", query)
	fmt.Println("Args:", args)

	row := m.db.QueryRow(query, args...)

	// fetch columns
	tmp, _ := m.db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 1", m.table))
	defer tmp.Close()
	cols, _ := tmp.Columns()

	values := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range values {
		ptrs[i] = &values[i]
	}

	if err := row.Scan(ptrs...); err != nil {
		return nil, err
	}

	var t T
	v := reflect.ValueOf(&t).Elem()
	for i, col := range cols {
		field := v.FieldByNameFunc(func(n string) bool {
			return strings.EqualFold(n, col)
		})
		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(values[i]))
		}
	}

	return &t, nil
}

func (m *BaseModel[T]) Update(entity T) (sql.Result, error) {
	setParts := []string{}
	args := []any{}

	// Use reflection on the struct
	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)

	// If pointer, deref
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}

		// Get db column name from tag, fallback to field name
		col := fieldType.Tag.Get("db")
		if col == "" {
			col = strings.ToLower(fieldType.Name)
		}

		// Add to SET clause
		setParts = append(setParts, fmt.Sprintf("%s = ?", col))
		args = append(args, fieldVal.Interface())
	}

	query := fmt.Sprintf("UPDATE %s SET %s", m.table, strings.Join(setParts, ", "))

	// Apply WHERE
	if m.where != nil {
		whereSQL, whereArgs := buildCondition(m.where)
		if whereSQL != "" {
			query += " WHERE " + whereSQL
			args = append(args, whereArgs...)
		}
	}

	fmt.Println("SQL:", query, "ARGS:", args)
	return m.db.Exec(query, args...)
}

func (m *BaseModel[T]) Delete() (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s", m.table)
	args := []any{}

	if m.where != nil {
		whereSQL, whereArgs := buildCondition(m.where)
		if whereSQL != "" {
			query += " WHERE " + whereSQL
			args = append(args, whereArgs...)
		}
	}

	fmt.Println("SQL:", query, "ARGS:", args)
	return m.db.Exec(query, args...)
}

func (m *BaseModel[T]) SoftDelete() (sql.Result, error) {
	now := time.Now()

	query := fmt.Sprintf("UPDATE %s SET deleted_at = ?", m.table)
	args := []any{now}

	if m.where != nil {
		whereSQL, whereArgs := buildCondition(m.where)
		if whereSQL != "" {
			query += " WHERE " + whereSQL
			args = append(args, whereArgs...)
		}
	}

	return m.db.Exec(query, args...)
}
