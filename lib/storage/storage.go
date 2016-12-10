package storage

import uuid "github.com/satori/go.uuid"

// Writer is something which you can Write Line
// types in to
type Writer interface {
	Write(Line) error
}

// Reader is something which returns a read-only
// channel of lines for a given UUID.
type Reader interface {
	Read(uuid.UUID) (<-chan PersistedLine, error)
}
