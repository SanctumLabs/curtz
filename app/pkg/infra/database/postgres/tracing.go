package postgres

import "strings"

// spanName names the span of a SQL statement. sqlc starts every statement with "-- name: Query :one", which the library
// would turn into the span name "--"; the query name is what a person looking at a trace wants, and it is as low in
// cardinality as the statement set. Any other statement is named after its first word (SELECT, BEGIN, ...), as the library does.
func spanName(statement string) string {
	statement = strings.TrimSpace(statement)
	if rest, found := strings.CutPrefix(statement, "-- name:"); found {
		if fields := strings.Fields(rest); len(fields) > 0 {
			return fields[0]
		}
		return "query"
	}
	if fields := strings.Fields(statement); len(fields) > 0 && !strings.HasPrefix(fields[0], "--") {
		return strings.ToUpper(fields[0])
	}
	if _, afterComment, found := strings.Cut(statement, "\n"); found && strings.HasPrefix(statement, "--") {
		return spanName(afterComment)
	}
	return "query"
}
