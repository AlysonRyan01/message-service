package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alysonryan/message-service/models"
	"github.com/lib/pq"
)

type MessageLogRepository struct {
	db *sql.DB
}

func NewMessageLogRepository(db *sql.DB) *MessageLogRepository {
	return &MessageLogRepository{
		db: db,
	}
}

func (r *MessageLogRepository) Create(ctx context.Context, m *models.Message) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var messageLogID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO message_log (message_log_status_id)
		VALUES ($1)
		RETURNING message_log_id`,
		models.StatusPending,
	).Scan(&messageLogID)
	if err != nil {
		return "", fmt.Errorf("failed to create message log: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO message (message_log_id, message_type_id, service_order_id, phone, payload)
		VALUES ($1, $2, $3, $4, $5)`,
		messageLogID, m.MessageTypeID, m.ServiceOrderID, m.Phone, m.Payload,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create message: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	return messageLogID, nil
}

func (r *MessageLogRepository) IncrementAttempts(ctx context.Context, messageLogID string, lastError string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE message_log
		SET attempts = attempts + 1,
		    last_error = $1
		WHERE message_log_id = $2`,
		lastError, messageLogID,
	)
	if err != nil {
		return fmt.Errorf("failed to increment attempts for message log %s: %w", messageLogID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read rows affected for message log %s: %w", messageLogID, err)
	}
	if rows == 0 {
		return fmt.Errorf("message log %s: %w", messageLogID, "message log not found")
	}

	return nil
}

func (r *MessageLogRepository) MarkSent(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE message_log
		SET message_log_status_id = $1,
		    sent_at = NOW()
		WHERE message_log_id = $2`,
		models.StatusSent, id,
	)
	return err
}

func (r *MessageLogRepository) MarkFailed(ctx context.Context, id string, errMsg string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE message_log
		SET message_log_status_id = $1,
		    attempts = attempts + 1,
		    last_error = $2
		WHERE message_log_id = $3`,
		models.StatusFailed, errMsg, id,
	)
	return err
}

func (r *MessageLogRepository) Delete(ctx context.Context, messageLogID string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM message_log WHERE message_log_id = $1`,
		messageLogID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete message log %s: %w", messageLogID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected for message log %s: %w", messageLogID, err)
	}

	if rows == 0 {
		return fmt.Errorf("message log %s: %w", messageLogID, "message log not found")
	}

	return nil
}

func (r *MessageLogRepository) FetchPending(ctx context.Context) ([]models.PendingMessage, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		SELECT l.message_log_id, m.message_type_id, m.service_order_id,
		       m.phone, m.payload, l.attempts
		FROM message_log l
		JOIN message m ON m.message_log_id = l.message_log_id
		WHERE l.message_log_status_id = $1
		ORDER BY l.created_at
		LIMIT 10
		FOR UPDATE OF 1 SKIP LOCKED`,
		models.StatusPending,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending messages: %w", err)
	}
	defer rows.Close()

	var messages []models.PendingMessage
	var ids []string

	for rows.Next() {
		var msg models.PendingMessage

		if err := rows.Scan(
			&msg.MessageLogID,
			&msg.MessageTypeID,
			&msg.ServiceOrderID,
			&msg.Phone,
			&msg.Payload,
			&msg.Attempts,
		); err != nil {
			return nil, fmt.Errorf("failed to scan pending message: %w", err)
		}

		messages = append(messages, msg)
		ids = append(ids, msg.MessageLogID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate pending messages: %w", err)
	}
	rows.Close()

	_, err = tx.ExecContext(ctx, `
		UPDATE message_log
		SET message_log_status_id = $1
		WHERE message_log_id = ANY($2)`,
		models.StatusProcessing, pq.Array(ids),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to mark as processing: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}


	return messages, nil
}

func (r *MessageLogRepository) Requeue(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE message_log
		SET message_log_status_id = $1, attempts = 0, last_error = NULL
		WHERE message_log_status_id = $2
		AND message_log_id = $3`,
		models.StatusPending, models.StatusFailed, id,
	)
	if err != nil {
		return fmt.Errorf("failed to requeue messages: %w", err)
	}
	return nil
}