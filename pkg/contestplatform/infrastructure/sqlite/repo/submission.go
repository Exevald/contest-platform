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
		SELECT 
		    id, 
		    problem_id,
		    participant_code,
		    language,
		    source_code,
		    verdict,
		    test_results_json,
		    compilation_output,
		    created_at_unix
		FROM submissions
		WHERE id = ?
	`, string(id))

	var (
		submissionID      string
		problemID         string
		participantCode   string
		language          string
		sourceCode        string
		verdict           string
		resultsJSON       string
		compilationOutput string
		createdAtUnix     int64
	)
	if err := row.Scan(
		&submissionID,
		&problemID,
		&participantCode,
		&language,
		&sourceCode,
		&verdict,
		&resultsJSON,
		&compilationOutput,
		&createdAtUnix,
	); err != nil {
		return nil, fmt.Errorf("find submission: %w", err)
	}

	var testResults []domainmodel.TestResult
	if resultsJSON != "" {
		if err := json.Unmarshal([]byte(resultsJSON), &testResults); err != nil {
			return nil, fmt.Errorf("unmarshal submission results: %w", err)
		}
	}

	return domainmodel.SubmissionFromSnapshot(domainmodel.SubmissionSnapshot{
		ID:                domainmodel.SubmissionID(submissionID),
		ProblemID:         domainmodel.ProblemID(problemID),
		ParticipantCode:   participantCode,
		Language:          domainmodel.Language(language),
		SourceCode:        sourceCode,
		Verdict:           domainmodel.Verdict(verdict),
		TestResults:       testResults,
		CompilationOutput: compilationOutput,
		CreatedAt:         time.Unix(createdAtUnix, 0).UTC(),
	}), nil
}

func (repo *submissionRepository) FindLatest(problemID domainmodel.ProblemID) (domainmodel.Submission, error) {
	row := repo.db.QueryRow(`
		SELECT 
		    id,
		    problem_id,
		    participant_code,
		    language,
		    source_code,
		    verdict,
		    test_results_json,
		    compilation_output,
		    created_at_unix
		FROM submissions
		WHERE problem_id = ?
		ORDER BY created_at_unix DESC, id DESC
		LIMIT 1
	`, string(problemID))

	var (
		submissionID      string
		foundProblem      string
		participantCode   string
		language          string
		sourceCode        string
		verdict           string
		resultsJSON       string
		compilationOutput string
		createdAtUnix     int64
	)
	if err := row.Scan(&submissionID, &foundProblem, &participantCode, &language, &sourceCode, &verdict, &resultsJSON, &compilationOutput, &createdAtUnix); err != nil {
		return nil, fmt.Errorf("find latest submission: %w", err)
	}

	var testResults []domainmodel.TestResult
	if resultsJSON != "" {
		if err := json.Unmarshal([]byte(resultsJSON), &testResults); err != nil {
			return nil, fmt.Errorf("unmarshal latest submission results: %w", err)
		}
	}

	return domainmodel.SubmissionFromSnapshot(domainmodel.SubmissionSnapshot{
		ID:                domainmodel.SubmissionID(submissionID),
		ProblemID:         domainmodel.ProblemID(foundProblem),
		ParticipantCode:   participantCode,
		Language:          domainmodel.Language(language),
		SourceCode:        sourceCode,
		Verdict:           domainmodel.Verdict(verdict),
		TestResults:       testResults,
		CompilationOutput: compilationOutput,
		CreatedAt:         time.Unix(createdAtUnix, 0).UTC(),
	}), nil
}

func (repo *submissionRepository) ListByProblem(problemID domainmodel.ProblemID) ([]domainmodel.Submission, error) {
	rows, err := repo.db.Query(`
		SELECT
		    id,
		    problem_id,
		    participant_code,
		    language,
		    source_code,
		    verdict,
		    test_results_json,
		    compilation_output,
		    created_at_unix
		FROM submissions
		WHERE problem_id = ?
		ORDER BY created_at_unix DESC, id DESC
	`, string(problemID))
	if err != nil {
		return nil, fmt.Errorf("list submissions: %w", err)
	}
	defer rows.Close()

	submissions := make([]domainmodel.Submission, 0)
	for rows.Next() {
		var (
			submissionID      string
			foundProblem      string
			participantCode   string
			language          string
			sourceCode        string
			verdict           string
			resultsJSON       string
			compilationOutput string
			createdAtUnix     int64
		)
		if err = rows.Scan(
			&submissionID,
			&foundProblem,
			&participantCode,
			&language,
			&sourceCode,
			&verdict,
			&resultsJSON,
			&compilationOutput,
			&createdAtUnix,
		); err != nil {
			return nil, fmt.Errorf("scan submissions: %w", err)
		}

		var testResults []domainmodel.TestResult
		if resultsJSON != "" {
			if err = json.Unmarshal([]byte(resultsJSON), &testResults); err != nil {
				return nil, fmt.Errorf("unmarshal submission history results: %w", err)
			}
		}

		submissions = append(submissions, domainmodel.SubmissionFromSnapshot(domainmodel.SubmissionSnapshot{
			ID:                domainmodel.SubmissionID(submissionID),
			ProblemID:         domainmodel.ProblemID(foundProblem),
			ParticipantCode:   participantCode,
			Language:          domainmodel.Language(language),
			SourceCode:        sourceCode,
			Verdict:           domainmodel.Verdict(verdict),
			TestResults:       testResults,
			CompilationOutput: compilationOutput,
			CreatedAt:         time.Unix(createdAtUnix, 0).UTC(),
		}))
	}

	return submissions, rows.Err()
}

func (repo *submissionRepository) Store(submission domainmodel.Submission) error {
	snapshot := domainmodel.SnapshotSubmission(submission)
	resultsJSON, err := json.Marshal(snapshot.TestResults)
	if err != nil {
		return fmt.Errorf("marshal submission results: %w", err)
	}

	_, err = repo.db.Exec(`
		INSERT INTO
		    submissions 
		(id, problem_id, participant_code, language, source_code, verdict, test_results_json, compilation_output, created_at_unix)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			problem_id = excluded.problem_id,
			participant_code = excluded.participant_code,
			language = excluded.language,
			source_code = excluded.source_code,
			verdict = excluded.verdict,
			test_results_json = excluded.test_results_json,
			compilation_output = excluded.compilation_output,
			created_at_unix = excluded.created_at_unix
	`,
		string(snapshot.ID),
		string(snapshot.ProblemID),
		snapshot.ParticipantCode,
		string(snapshot.Language),
		snapshot.SourceCode,
		string(snapshot.Verdict),
		string(resultsJSON),
		snapshot.CompilationOutput,
		snapshot.CreatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("store submission: %w", err)
	}

	return nil
}
