package storage

import (
	"database/sql"
	"errors"
	"fmt"

	appmodel "contest-platform/pkg/contestplatform/app/storage"
)

const singletonSessionID = "current"

func NewSessionStorage(db *sql.DB) appmodel.SessionStorage {
	return &sessionStorage{db: db}
}

type sessionStorage struct {
	db *sql.DB
}

func (storage *sessionStorage) Load() (*appmodel.ParticipantSession, error) {
	row := storage.db.QueryRow(`
		SELECT participant_code, theme
		FROM participant_session
		WHERE id = ?
	`, singletonSessionID)

	var session appmodel.ParticipantSession
	if err := row.Scan(&session.ParticipantCode, &session.Theme); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("load participant session: %w", err)
	}

	return &session, nil
}

func (storage *sessionStorage) Save(session appmodel.ParticipantSession) error {
	_, err := storage.db.Exec(`
		INSERT INTO participant_session (id, participant_code, theme)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			participant_code = excluded.participant_code,
			theme = excluded.theme
	`, singletonSessionID, session.ParticipantCode, session.Theme)
	if err != nil {
		return fmt.Errorf("save participant session: %w", err)
	}

	return nil
}
