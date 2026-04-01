package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	domainmodel "contest-platform/pkg/contestplatform/domain/model"
)

func NewProblemRepository(db *sql.DB) domainmodel.ProblemRepository {
	return &problemRepository{db: db}
}

type problemRepository struct {
	db *sql.DB
}

func (repo *problemRepository) NextID() domainmodel.ProblemID {
	return domainmodel.ProblemID(newID("problem"))
}

func (repo *problemRepository) List() ([]domainmodel.Problem, error) {
	rows, err := repo.db.Query(`
		SELECT id, title, description, time_limit_ns, memory_limit_bytes, test_cases_json
		FROM problems
		ORDER BY title, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list problems: %w", err)
	}
	defer rows.Close()

	problems := make([]domainmodel.Problem, 0)
	for rows.Next() {
		problem, scanErr := scanProblem(rows.Scan)
		if scanErr != nil {
			return nil, scanErr
		}
		problems = append(problems, problem)
	}

	return problems, rows.Err()
}

func (repo *problemRepository) Find(id domainmodel.ProblemID) (domainmodel.Problem, error) {
	row := repo.db.QueryRow(`
		SELECT id, title, description, time_limit_ns, memory_limit_bytes, test_cases_json
		FROM problems
		WHERE id = ?
	`, string(id))

	return scanProblem(row.Scan)
}

func (repo *problemRepository) Store(problem domainmodel.Problem) error {
	snapshot := domainmodel.SnapshotProblem(problem)
	testCasesJSON, err := json.Marshal(snapshot.TestCases)
	if err != nil {
		return fmt.Errorf("marshal problem test cases: %w", err)
	}

	_, err = repo.db.Exec(`
		INSERT INTO problems (id, title, description, time_limit_ns, memory_limit_bytes, test_cases_json)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			description = excluded.description,
			time_limit_ns = excluded.time_limit_ns,
			memory_limit_bytes = excluded.memory_limit_bytes,
			test_cases_json = excluded.test_cases_json
	`,
		string(snapshot.ID),
		string(snapshot.Title),
		snapshot.Description,
		snapshot.Constraints.TimeLimit.Nanoseconds(),
		snapshot.Constraints.MemoryLimit,
		string(testCasesJSON),
	)
	if err != nil {
		return fmt.Errorf("store problem: %w", err)
	}

	return nil
}

func scanProblem(scan func(dest ...any) error) (domainmodel.Problem, error) {
	var (
		id           string
		title        string
		description  string
		timeLimitNS  int64
		memoryLimit  uint64
		testCasesRaw string
	)
	if err := scan(&id, &title, &description, &timeLimitNS, &memoryLimit, &testCasesRaw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("scan problem: %w", err)
	}

	var testCases []domainmodel.TestCase
	if testCasesRaw != "" {
		if err := json.Unmarshal([]byte(testCasesRaw), &testCases); err != nil {
			return nil, fmt.Errorf("unmarshal problem test cases: %w", err)
		}
	}

	return domainmodel.ProblemFromSnapshot(domainmodel.ProblemSnapshot{
		ID:          domainmodel.ProblemID(id),
		Title:       domainmodel.Title(title),
		Description: description,
		Constraints: domainmodel.Constraints{
			TimeLimit:   time.Duration(timeLimitNS),
			MemoryLimit: memoryLimit,
		},
		TestCases: testCases,
	})
}
