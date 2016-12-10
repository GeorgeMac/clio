package storage

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	// Ensure SysLogLineV0 implements Line interface
	_ Line = &SysLogLineV0{}
)

// Line is an interface that describes the main unit of
// storage in the storage package
// A line represents a single row for a log event not
// yet persisted
type Line interface {
	Tag() Tag
	Payload() []byte
	CreatedAt() time.Time
}

type PersistedLine struct {
	BuildID       uuid.UUID
	CreatedAt     time.Time
	InsertedAt    time.Time
	ContainerID   string
	ContainerName string
	Payload       []byte
}

// SysLogLineV0 represents a single line from a syslog event
type SysLogLineV0 struct {
	// embeded V0 tag
	TagV0

	// event time
	Time time.Time
	// line payload
	Data []byte
}

// Tag returns the emdedded TagV0 converted to a plain old Tag type
func (s *SysLogLineV0) Tag() Tag { return Tag(s.TagV0) }

// Payload returns the underlying data from the Syslog payload
func (s *SysLogLineV0) Payload() []byte { return s.Data }

// CreatedAt returns the time of the the event
func (s *SysLogLineV0) CreatedAt() time.Time { return s.Time }

// String returns a string representation of the Syslog line
func (s *SysLogLineV0) String() string {
	return fmt.Sprintf("(%s) timestamp %v content %v", s.TagV0, s.Time, string(s.Data))
}

// GoString returns a string representation of the Syslog line
func (s *SysLogLineV0) GoString() string {
	return s.String()
}

// UnmarshalLogParts takes a map[sting]interface{} (syslog logParts) and
// unmarshalls the parts tag, content and timestamp
func (s *SysLogLineV0) UnmarshalLogParts(parts map[string]interface{}) error {
	// parse tag
	tag, err := retrieveBytes(parts, "tag")
	if err != nil {
		return err
	}

	// unmarshal in to the V0 tag
	if err := s.TagV0.UnmarshalText(tag); err != nil {
		return err
	}

	// parse content
	content, err := retrieveBytes(parts, "content")
	if err != nil {
		return err
	}

	s.Data = content

	// parse timestamp
	timestamp, err := retrieveTime(parts, "timestamp")
	if err != nil {
		return err
	}

	s.Time = timestamp

	return nil
}

func retrieveBytes(parts map[string]interface{}, key string) ([]byte, error) {
	if value, ok := parts[key]; !ok {
		return nil, fmt.Errorf("missing key %s", key)
	} else if valueString := value.(string); !ok {
		return nil, fmt.Errorf("unexpected value type (%T) %v", value, value)
	} else {
		return []byte(valueString), nil
	}

	return nil, nil
}

func retrieveTime(parts map[string]interface{}, key string) (time.Time, error) {
	if value, ok := parts[key]; !ok {
		return time.Time{}, fmt.Errorf("missing key %s", key)
	} else if timestamp := value.(time.Time); !ok {
		return time.Time{}, fmt.Errorf("unexpected value type (%T) %v", value, value)
	} else {
		return timestamp, nil
	}

	return time.Time{}, nil
}
