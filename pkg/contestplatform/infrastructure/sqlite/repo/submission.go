package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	domainmodel "contest-platform/pkg/contestplatform/domain/model"
)

func NewSubmissionRepository(db *sql.DB) domainmodel.SubmissionRepository {
	return &submissionRepository{db: db}
}

type submissionRepository struct {
	db *sql.DB
}

func (repo *submissionRepository) NextID() domainmodel.SubmissionID {
	return domainmodel.SubmissionID(newID("submission"))
}

func (repo *submissionRepository) Find(id domainmodel.SubmissionID) (domainmodel.Submission, error) {
	row := repo.db.QueryRow(`
		SELECT id, problem_id, language, source_code, verdict, test_results_json, created_at_unix
		FROM submissions
		WHERE id = ?
	`, string(id))

	var (
		submissionID  string
		problemID     string
		language      string
		sourceCode    string
		verdict       string
		resultsJSON   string
		createdAtUnix int64
	)
	if err := row.Scan(&submissionID, &problemID, &language, &sourceCode, &verdict, &resultsJSON, &createdAtUnix); err != nil {
		return nil, fmt.Errorf("find submission: %w", err)
	}

	var testResults []domainmodel.TestResult
	if resultsJSON != "" {
		if err := json.Unmarshal([]byte(resultsJSON), &testResults); err != nil {
			return nil, fmt.Errorf("unmarshal submission results: %w", err)
		}
	}

	return domainmodel.SubmissionFromSnapshot(domainmodel.SubmissionSnapshot{
		ID:          domainmodel.SubmissionID(submissionID),
		ProblemID:   domainmodel.ProblemID(problemID),
		Language:    domainmodel.Language(language),
		SourceCode:  sourceCode,
		Verdict:     domainmodel.Verdict(verdict),
		TestResults: testResults,
		CreatedAt:   time.Unix(createdAtUnix, 0).UTC(),
	}), nil
}

func (repo *submissionRepository) FindLatest(problemID domainmodel.ProblemID) (domainmodel.Submission, error) {
	row := repo.db.QueryRow(`
		SELECT id, problem_id, language, source_code, verdict, test_results_json, created_at_unix
		FROM submissions
		WHERE problem_id = ?
		ORDER BY created_at_unix DESC, id DESC
		LIMIT 1
	`, string(problemID))

	var (
		submissionID  string
		foundProblem  string
		language      string
		sourceCode    string
		verdict       string
		resultsJSON   string
		createdAtUnix int64
	)
	if err := row.Scan(&submissionID, &foundProblem, &language, &sourceCode, &verdict, &resultsJSON, &createdAtUnix); err != nil {
		return nil, fmt.Errorf("find latest submission: %w", err)
	}

	var testResults []domainmodel.TestResult
	if resultsJSON != "" {
		if err := json.Unmarshal([]byte(resultsJSON), &testResults); err != nil {
			return nil, fmt.Errorf("unmarshal latest submission results: %w", err)
		}
	}

	return domainmodel.SubmissionFromSnapshot(domainmodel.SubmissionSnapshot{
		ID:          domainmodel.SubmissionID(submissionID),
		ProblemID:   domainmodel.ProblemID(foundProblem),
		Language:    domainmodel.Language(language),
		SourceCode:  sourceCode,
		Verdict:     domainmodel.Verdict(verdict),
		TestResults: testResults,
		CreatedAt:   time.Unix(createdAtUnix, 0).UTC(),
	}), nil
}

func (repo *submissionRepository) Store(submission domainmodel.Submission) error {
	snapshot := domainmodel.SnapshotSubmission(submission)
	resultsJSON, err := json.Marshal(snapshot.TestResults)
	if err != nil {
		return fmt.Errorf("marshal submission results: %w", err)
	}

	_, err = repo.db.Exec(`
		INSERT INTO submissions (id, problem_id, language, source_code, verdict, test_results_json, created_at_unix)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			problem_id = excluded.problem_id,
			language = excluded.language,
			source_code = excluded.source_code,
			verdict = excluded.verdict,
			test_results_json = excluded.test_results_json,
			created_at_unix = excluded.created_at_unix
	`,
		string(snapshot.ID),
		string(snapshot.ProblemID),
		string(snapshot.Language),
		snapshot.SourceCode,
		string(snapshot.Verdict),
		string(resultsJSON),
		snapshot.CreatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("store submission: %w", err)
	}

	return nil
}
