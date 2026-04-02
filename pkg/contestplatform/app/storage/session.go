package storage

type ParticipantSession struct {
	ParticipantCode string
	Theme           string
}

type SessionStorage interface {
	Load() (*ParticipantSession, error)
	Save(session ParticipantSession) error
}
