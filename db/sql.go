package db

import (
	"fmt"
	"strings"

	"crawshaw.io/sqlite"
)

// bindFunc is a function that binds one or more values to a sqlite.Stmt. The
// starting param index is passed into the bindFunc when sql.Prep is called.
// The bindFunc must return the number of binds it called on the Stmt so that
// the param index can be advanced.
type bindFunc func(*sqlite.Stmt, int) int

// sql is a SQL builder for more complex queries. It allows for adding binds to
// a sqlite.Stmt before it is prepared. As SQL is appended to the query,
// bindFuncs are queued for later when sql.Prep() is called. Do not copy a
// non-zero sql.
type sql struct {
	sql   strings.Builder
	binds []bindFunc
}

// Appends a trailing ';' to the SQL and calls conn.Prep. Finally all bindFuncs
// are called and the stmt is returned ready for its first Stmt.Step() call.
func (sql *sql) Prep(conn *sqlite.Conn) *sqlite.Stmt {
	sql.sql.WriteString(`;`)
	log.Debug(sql.sql.String())
	stmt := conn.Prep(sql.sql.String())
	param := 1
	for _, bind := range sql.binds {
		param += bind(stmt, param)
	}
	return stmt
}

// Append str to the SQL and append the binds.
func (sql *sql) Append(str string, binds ...bindFunc) {
	sql.sql.WriteString(str)
	sql.binds = append(sql.binds, binds...)
}

// Append n comma separated params placeholders (e.g. "?, ?, ... , ?") to the
// SQL and append the binds.
func (sql *sql) Bind(n int, binds ...bindFunc) {
	str := strings.TrimRight(strings.Repeat("?, ", n), ", ")
	sql.Append(str, binds...)
}

const MaxLimit = 600

// Append "LIMIT ?, ?" to the SQL and the appropriate page and limit binds.
func (sql *sql) Paginate(page, limit uint64) {
	if limit == 0 || limit > MaxLimit {
		limit = MaxLimit
	}
	sql.Append(` LIMIT ?, ?`, func(s *sqlite.Stmt, p int) int {
		s.BindInt64(p, int64((page-1)*limit))
		s.BindInt64(p+1, int64(limit))
		return 2
	})
}

// Append "ORDER BY "id" ASC or DESC". No binds are added.
func (sql *sql) Order(order string) {
	switch strings.ToLower(order) {
	case "asc", "":
		sql.Append(` ORDER BY "id" ASC`)
	case "desc":
		sql.Append(` ORDER BY "id" DESC`)
	default:
		panic(fmt.Errorf("invalid order: %v", order))
	}
}

// Combines Order and Paginate in one call.
func (sql *sql) OrderPaginate(order string, page, limit uint64) {
	sql.Order(order)
	sql.Paginate(page, limit)
}