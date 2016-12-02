package storage

import "github.com/pkg/errors"

// ErrorKeyInvalidPrefix returns an error for a given
// desired prefix. It is used in a response to a tag that does
// not start with the desired prefix.
func ErrorKeyInvalidPrefix(prefix string) error {
	return errors.Errorf("invalid: should have prefix %s", prefix)
}

// ErrorKeyInvalidNumberOfTagParts returns an error for a given
// desired number of parts to a tag. It is used in a response to a tag
// provided with an unexpected number of parts after a split on the separator.
func ErrorKeyInvalidNumberOfTagParts(expected int) error {
	return errors.Errorf("invalid: should have %d tag parts", expected)
}
