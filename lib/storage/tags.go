package storage

import (
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
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

// UnmarshalText takes a tag from syslog and a timestamp and
// sets the key appropriately. The syslog format is:
// docker/v0/<build-uuid>/<container-id>/<container-name>
func (t *TagV0) UnmarshalText(v []byte) error {
	parts := strings.Split(string(v), "/")
	if len(parts) != 5 {
		return errors.Wrapf(ErrorKeyInvalidNumberOfTagPartsV0, "tag %s has %d part(s)", string(v), len(parts))
	}

	if parts[0] != "docker" || parts[1] != "v0" {
		return errors.Wrapf(ErrorKeyInvalidPrefixV0, "tag %s has prefix %s/%s", parts[0], parts[1])
	}

	var err error
	t.BuildID, err = uuid.FromString(parts[2])
	if err != nil {
		return errors.Wrapf(err, "error: parsing V0 clio key uuid %s", parts[2])
	}

	t.ContainerID = parts[3]
	t.ContainerName = parts[4]
	return nil
}
