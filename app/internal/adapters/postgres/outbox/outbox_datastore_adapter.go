package outboxdatastore

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	postgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql"
	"github.com/sanctumlabs/curtz/app/internal/ports"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database/postgres"
)

// leaseKey is the advisory lock key of the relay lease, derived from a fixed name so every instance computes the same one.
var leaseKey = func() int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte("curtz.outbox.relay"))
	return int64(h.Sum64())
}()

var _ ports.OutboxDatastore = (*Adapter)(nil)

// Adapter implements ports.OutboxDatastore on Postgres.
type Adapter struct {
	queries *postgresql.Queries
	// leaseConfig opens the dedicated connection the lease lives on. It is parsed through the pool's parser so the
	// pool_* parameters of the connection string are removed; a plain connection would send them to the server.
	leaseConfig *pgx.ConnConfig
	timeout     time.Duration
}

// NewAdapter builds the adapter. connString is the same connection string the pool was built from; timeout bounds every
// statement, because the relay runs some of them on a context that does not carry a deadline.
func NewAdapter(client database.PostgresDatabaseClient, connString string, timeout time.Duration) (*Adapter, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse the connection string for the outbox lease: %w", err)
	}
	// The lease connection is housekeeping and must not create a span per statement.
	poolConfig.ConnConfig.Tracer = nil

	return &Adapter{
		queries:     postgresql.New(client.GetDB()),
		leaseConfig: poolConfig.ConnConfig,
		timeout:     timeout,
	}, nil
}

func (a *Adapter) bounded(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, a.timeout)
}

// Acquire takes the relay lease: a session-level advisory lock on a dedicated connection. The lock disappears with
// the connection, so a crashed instance, or a failover to a new primary, releases it without anyone cleaning up.
func (a *Adapter) Acquire(ctx context.Context) (ports.Lease, error) {
	ctx, cancel := a.bounded(ctx)
	defer cancel()

	conn, err := pgx.ConnectConfig(ctx, a.leaseConfig)
	if err != nil {
		return nil, fmt.Errorf("open the outbox lease connection: %w", err)
	}
	var acquired bool
	if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", leaseKey).Scan(&acquired); err != nil {
		_ = conn.Close(context.WithoutCancel(ctx))
		return nil, fmt.Errorf("take the outbox lease: %w", err)
	}
	if !acquired {
		_ = conn.Close(context.WithoutCancel(ctx))
		return nil, ports.ErrNotLeader
	}
	return &lease{conn: conn, timeout: a.timeout}, nil
}

type lease struct {
	conn    *pgx.Conn
	timeout time.Duration
}

// Alive runs a trivial statement on the lease's own connection; it fails once the connection, and with it the lock, is gone.
func (l *lease) Alive(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()
	var one int
	if err := l.conn.QueryRow(ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("the outbox lease connection is not usable: %w", err)
	}
	return nil
}

// Release closes the connection, which releases the lock.
func (l *lease) Release() {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()
	if err := l.conn.Close(ctx); err != nil {
		slog.Warn("closing the outbox lease connection", "error", err)
	}
}

// Claim returns up to perDestination of the oldest unsent, unparked events of every destination, in creation order.
func (a *Adapter) Claim(ctx context.Context, perDestination int) ([]ports.OutboxEvent, error) {
	ctx, cancel := a.bounded(ctx)
	defer cancel()

	rows, err := a.queries.QueryClaimOutboxEvents(ctx, int32(perDestination))
	if err != nil {
		return nil, fmt.Errorf("claim outbox events: %w", err)
	}
	events := make([]ports.OutboxEvent, 0, len(rows))
	for _, row := range rows {
		id, idErr := postgres.UUIDToString(row.ID)
		if idErr != nil {
			return nil, fmt.Errorf("outbox event has an invalid id: %w", idErr)
		}
		events = append(events, ports.OutboxEvent{
			ID:           id,
			Destination:  row.Destination,
			PartitionKey: row.PartitionKey.String,
			EventType:    row.EventType,
			Headers:      parseHeaders(id, row.Headers),
			Payload:      row.Payload,
			Attempts:     int(row.Attempts),
			CreatedAt:    row.CreatedAt.Time,
		})
	}
	return events, nil
}

// parseHeaders reads the stored headers JSON. A row whose headers are not an object of strings must not block every
// event behind it, so it is published without headers and the problem is logged.
func parseHeaders(id string, raw []byte) map[string]string {
	headers := map[string]string{}
	if err := json.Unmarshal(raw, &headers); err != nil {
		slog.Warn("outbox event has unreadable headers, publishing it without them", "event_id", id, "error", err)
		return map[string]string{}
	}
	return headers
}

// MarkSent records that the events were acknowledged by the broker.
func (a *Adapter) MarkSent(ctx context.Context, ids []string) error {
	uuids := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		parsed, err := postgres.StringToUUID(id)
		if err != nil {
			return fmt.Errorf("invalid outbox event id %q: %w", id, err)
		}
		uuids = append(uuids, parsed)
	}
	ctx, cancel := a.bounded(ctx)
	defer cancel()

	if _, err := a.queries.QueryMarkOutboxEventsSent(ctx, uuids); err != nil {
		return fmt.Errorf("mark outbox events as sent: %w", err)
	}
	return nil
}

// RecordRejection counts a permanent rejection of the event by the broker and returns how many it has had.
func (a *Adapter) RecordRejection(ctx context.Context, id, reason string) (int, error) {
	uuid, err := postgres.StringToUUID(id)
	if err != nil {
		return 0, fmt.Errorf("invalid outbox event id %q: %w", id, err)
	}
	ctx, cancel := a.bounded(ctx)
	defer cancel()

	attempts, err := a.queries.QueryRecordOutboxEventRejection(ctx, postgresql.QueryRecordOutboxEventRejectionParams{
		ID: uuid, ErrorMessage: pgtype.Text{String: reason, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("record the rejection of an outbox event: %w", err)
	}
	return int(attempts), nil
}

// Park stops the relay from retrying the event; the row stays, with its reason.
func (a *Adapter) Park(ctx context.Context, id, reason string) error {
	uuid, err := postgres.StringToUUID(id)
	if err != nil {
		return fmt.Errorf("invalid outbox event id %q: %w", id, err)
	}
	ctx, cancel := a.bounded(ctx)
	defer cancel()

	if err := a.queries.QueryParkOutboxEvent(ctx, postgresql.QueryParkOutboxEventParams{
		ID: uuid, ErrorMessage: pgtype.Text{String: reason, Valid: true},
	}); err != nil {
		return fmt.Errorf("park an outbox event: %w", err)
	}
	return nil
}

// Purge deletes up to limit events sent before olderThan and returns how many it deleted.
func (a *Adapter) Purge(ctx context.Context, olderThan time.Time, limit int) (int, error) {
	ctx, cancel := a.bounded(ctx)
	defer cancel()

	deleted, err := a.queries.QueryPurgeSentOutboxEvents(ctx, postgresql.QueryPurgeSentOutboxEventsParams{
		OlderThan: pgtype.Timestamptz{Time: olderThan, Valid: true},
		LimitBy:   int32(limit),
	})
	if err != nil {
		return 0, fmt.Errorf("purge sent outbox events: %w", err)
	}
	return int(deleted), nil
}

// Backlog reports the undelivered and parked events.
func (a *Adapter) Backlog(ctx context.Context) (ports.Backlog, error) {
	ctx, cancel := a.bounded(ctx)
	defer cancel()

	row, err := a.queries.QueryOutboxBacklog(ctx)
	if err != nil {
		return ports.Backlog{}, fmt.Errorf("read the outbox backlog: %w", err)
	}
	backlog := ports.Backlog{Unsent: row.Unsent, Parked: row.Parked}
	if row.OldestUnsent.Valid {
		backlog.OldestUnsent = row.OldestUnsent.Time
	}
	return backlog, nil
}
