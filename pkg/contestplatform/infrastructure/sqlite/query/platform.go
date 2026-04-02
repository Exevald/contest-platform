package query

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	appquery "contest-platform/pkg/contestplatform/app/query"
)

func NewPlatformQueryService(db *sql.DB) appquery.PlatformQueryService {
	return &platformQueryService{db: db}
}

type platformQueryService struct {
	db *sql.DB
}

func (s *platformQueryService) ListProblems() ([]appquery.ProblemListItem, error) {
	rows, err := s.db.Query(`
		SELECT id, title
		FROM problems
		ORDER BY title, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list problems view: %w", err)
	}
	defer rows.Close()

	items := make([]appquery.ProblemListItem, 0)
	for rows.Next() {
		var item appquery.ProblemListItem
		if err = rows.Scan(&item.ID, &item.Title); err != nil {
			return nil, fmt.Errorf("scan problems view: %w", err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (s *platformQueryService) GetProblemDescription(problemID string) (string, error) {
	row := s.db.QueryRow(`
		SELECT description
		FROM problems
		WHERE id = ?
	`, problemID)

	var description string
	if err := row.Scan(&description); err != nil {
		return "", fmt.Errorf("get problem description: %w", err)
	}

	return description, nil
}

func (s *platformQueryService) GetSubmissionStatus(submissionID string) (*appquery.SubmissionView, error) {
	row := s.db.QueryRow(`
		SELECT id, problem_id, language, verdict, compilation_output, created_at_unix
		FROM submissions
		WHERE id = ?
	`, submissionID)

	view, err := scanSubmissionView(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &view, nil
}

func (s *platformQueryService) ListSubmissionHistory(problemID string) ([]appquery.SubmissionView, error) {
	rows, err := s.db.Query(`
		SELECT id, problem_id, language, verdict, compilation_output, created_at_unix
		FROM submissions
		WHERE problem_id = ?
		ORDER BY created_at_unix DESC, id DESC
	`, problemID)
	if err != nil {
		return nil, fmt.Errorf("list submission history: %w", err)
	}
	defer rows.Close()

	items := make([]appquery.SubmissionView, 0)
	for rows.Next() {
		item, scanErr := scanSubmissionView(rows.Scan)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func scanSubmissionView(scan func(dest ...any) error) (appquery.SubmissionView, error) {
	var (
		id                string
		problemID         string
		language          string
		verdict           string
		compilationOutput string
		createdAtUnix     int64
	)
	if err := scan(&id, &problemID, &language, &verdict, &compilationOutput, &createdAtUnix); err != nil {
		return appquery.SubmissionView{}, err
	}

	return appquery.SubmissionView{
		ID:                id,
		ProblemID:         problemID,
		Language:          language,
		Verdict:           verdict,
		CompilationOutput: compilationOutput,
		CreatedAt:         time.Unix(createdAtUnix, 0).UTC(),
	}, nil
}
