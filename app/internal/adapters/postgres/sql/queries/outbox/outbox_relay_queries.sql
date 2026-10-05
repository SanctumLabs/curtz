-- name: QueryClaimOutboxEvents :many
-- The relay's batch: per destination, the oldest unsent, unparked, undeleted events in creation order, so a destination
-- whose topic is unavailable cannot starve the others. The creation order is (created_at, id); id is a UUIDv7 and keeps
-- the events of one transaction in the order they were recorded.
SELECT
  ranked.id,
  ranked.partition_key,
  ranked.destination,
  ranked.event_type,
  ranked.headers,
  ranked.payload,
  ranked.attempts,
  ranked.created_at
FROM (
  SELECT
    oe.*,
    row_number() OVER (PARTITION BY oe.destination ORDER BY oe.created_at, oe.id) AS position
  FROM outbox_events oe
  WHERE oe.sent_time IS NULL
    AND oe.parked_at IS NULL
    AND oe.deleted_at IS NULL
) ranked
WHERE ranked.position <= sqlc.arg(per_destination)::int
ORDER BY ranked.created_at, ranked.id;

-- name: QueryMarkOutboxEventsSent :execrows
UPDATE outbox_events
SET
  sent_time = now(),
  updated_at = now()
WHERE id = ANY(sqlc.arg(ids)::uuid[])
  AND sent_time IS NULL;

-- name: QueryRecordOutboxEventRejection :one
UPDATE outbox_events
SET
  attempts = attempts + 1,
  error_message = sqlc.arg(error_message),
  updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING attempts;

-- name: QueryParkOutboxEvent :exec
UPDATE outbox_events
SET
  parked_at = now(),
  error_message = sqlc.arg(error_message),
  updated_at = now()
WHERE id = sqlc.arg(id);

-- name: QueryPurgeSentOutboxEvents :execrows
DELETE FROM outbox_events
WHERE id IN (
  SELECT oe.id
  FROM outbox_events oe
  WHERE oe.sent_time IS NOT NULL
    AND oe.sent_time < sqlc.arg(older_than)::timestamptz
  ORDER BY oe.sent_time
  LIMIT sqlc.arg(limit_by)::int
);

-- name: QueryOutboxBacklog :one
SELECT
  count(*) FILTER (WHERE parked_at IS NULL)::bigint AS unsent,
  (min(created_at) FILTER (WHERE parked_at IS NULL))::timestamptz AS oldest_unsent,
  count(*) FILTER (WHERE parked_at IS NOT NULL)::bigint AS parked
FROM outbox_events
WHERE sent_time IS NULL
  AND deleted_at IS NULL;
