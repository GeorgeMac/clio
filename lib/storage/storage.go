package storage

import "time"

var (
	// ErrorKeyInvalidNumberOfTagPartsV0 is returned when a tag is incorrectly
	// formed because it doesn't contain enough of the correct parts.
	ErrorKeyInvalidNumberOfTagPartsV0 = ErrorKeyInvalidNumberOfTagParts(5)

	// ErrorKeyInvalidPrefixV0 is returned when a tag is prefix incorrectly
	ErrorKeyInvalidPrefixV0 = ErrorKeyInvalidPrefix("docker/v0")

	// Ensure SysLogLineV0 implements Line interface
	_ Line = &SysLogLineV0{}
)

// Storage is something which you can Put LogLine
// types in to.
type Storage interface {
	Put(Line) error
}

// Line is an interface that describes the main unit of
// storage in the storage package. A line represents a single
// row for a log event.
type Line interface {
	Tag() Tag
	Payload() []byte
	CreatedAt() time.Time
}

// SysLogLineV0 represents a single line from a syslog event.
type SysLogLineV0 struct {
	// embeded V0 tag
	TagV0

	// event time
	Time time.Time
	// line payload
	Data []byte
}

// Tag returns the emdedded TagV0 converted to a plain old Tag type.
func (s *SysLogLineV0) Tag() Tag { return Tag(s.TagV0) }

// Payload returns the underlying data from the Syslog payload.
func (s *SysLogLineV0) Payload() []byte { return s.Data }

// CreatedAt returns the time of the the event.
func (s *SysLogLineV0) CreatedAt() time.Time { return s.Time }
