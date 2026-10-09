package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// sqlc puts "-- name: Query :one" on the first line of every statement, which the library would turn into the span name
// "--" for every query.
func TestSpanName_UsesTheSqlcQueryName(t *testing.T) {
	cases := map[string]struct {
		statement string
		want      string
	}{
		"sqlc one":      {"-- name: QueryCreateUser :one\nINSERT INTO users (id) VALUES ($1)", "QueryCreateUser"},
		"sqlc many":     {"-- name: QueryOutboxEventsUnSent :many\nSELECT * FROM outbox_events", "QueryOutboxEventsUnSent"},
		"sqlc exec":     {"-- name: QueryDeleteUser :exec\nDELETE FROM users WHERE id = $1", "QueryDeleteUser"},
		"plain select":  {"select 1", "SELECT"},
		"begin":         {"begin", "BEGIN"},
		"leading space": {"  \n  update users set x = 1", "UPDATE"},
		"other comment": {"-- hello\nSELECT 1", "SELECT"},
		"empty":         {"", "query"},
		"name only":     {"-- name:", "query"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, spanName(tc.statement))
		})
	}
}
