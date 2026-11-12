package outbox

import (
	"context"
	"fmt"
)

const deleteOutboxEventQuery = `
DELETE FROM outbox_events
WHERE event_id = $1
`

func (r *PostgresOutboxRepository) DeleteOutboxEvent(ctx context.Context, eventID int64) error {
	if _, err := r.pool.Exec(ctx, deleteOutboxEventQuery, eventID); err != nil {
		return fmt.Errorf("delete outbox event: %w", err)
	}
	return nil
}
