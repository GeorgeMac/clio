package storage

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrorKeyInvalidNumberOfTagPartsV0 is returned when a tag is incorrectly
	// formed because it doesn't contain enough of the correct parts.
	ErrorKeyInvalidNumberOfTagPartsV0 = ErrorKeyInvalidNumberOfTagParts(4)

	// ErrorKeyInvalidPrefixV0 is returned when a tag is prefix incorrectly
	ErrorKeyInvalidPrefixV0 = ErrorKeyInvalidPrefix("v0")
)

// Tag is a struct which represent the parts
// of a syslog tag, used to identify where logs
// came from.
type Tag struct {
	BuildID       uuid.UUID
	ContainerID   string
	ContainerName string
}

// TagV0 decorates Tag with an UnmarshalText method
// for parsing the tag from a slice of bytes.
// It implements the encoding.TestMarshaller interface.
type TagV0 Tag

// String returns a string representation of the Tag
func (t *TagV0) String() string {
	return fmt.Sprintf("V0 Tag BUILD_ID[%s] CONTAINER_ID[%s] CONTAINER_NAME[%s]", t.BuildID.String(), t.ContainerID, t.ContainerName)
}

// GoString returns a string representation of the Tag
func (t *TagV0) GoString() string {
	return t.String()
}

// UnmarshalText takes a tag from syslog and a timestamp and
// sets the key appropriately. The syslog format is:
// v0/<build-uuid>/<container-id>/<container-name>
func (t *TagV0) UnmarshalText(v []byte) error {
	parts := strings.Split(string(v), "/")
	if len(parts) != 4 {
		return errors.Wrapf(ErrorKeyInvalidNumberOfTagPartsV0, "tag %s has %d part(s)", string(v), len(parts))
	}

	if parts[0] != "v0" {
		return errors.Wrapf(ErrorKeyInvalidPrefixV0, "tag %s has prefix %s", parts[0])
	}

	var err error
	t.BuildID, err = uuid.FromString(parts[1])
	if err != nil {
		return errors.Wrapf(err, "error: parsing V0 clio key uuid %s", parts[2])
	}

	t.ContainerID = parts[2]
	t.ContainerName = parts[3]
	return nil
}
