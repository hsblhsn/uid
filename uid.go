// Package uid sourced from: https://github.com/ent/contrib/blob/master/entgql/internal/todopulid/ent/schema/pulid/pulid.go
package uid

import (
	"crypto/rand"
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"

	"github.com/oklog/ulid/v2"
)

// ID implements a uid - a prefixed ULID.
type ID string

const (
	PrefixLen = 5
	ULIDLen   = 26
	IDLen     = PrefixLen + ULIDLen
)

// newULID returns a new ULID for time.Now() using the default entropy source.
func newULID() ulid.ULID {
	return ulid.MustNew(ulid.Now(), rand.Reader)
}

// MustNew returns a new uid for time.Now() given a prefix. This uses the default entropy source.
func MustNew(prefix string) ID {
	if len(prefix) != PrefixLen {
		panic("uid: uid prefix must have to be 5 characters long")
	}
	return ID(prefix + newULID().String())
}

func (u ID) Prefix() string {
	if !u.IsValid() {
		return ""
	}
	return string(u)[:PrefixLen]
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (u *ID) UnmarshalGQL(v interface{}) error {
	return u.Scan(v)
}

// MarshalGQL implements the graphql.Marshaler interface.
func (u ID) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(u.String()))
}

// Scan implements the Scanner interface.
func (u *ID) Scan(src interface{}) error {
	if src == nil {
		src = ""
	}
	switch v := src.(type) {
	case string:
		*u = ID(v)
	case []byte:
		*u = ID(v)
	case ID:
		*u = v
	default:
		return fmt.Errorf("uid: cannot scan type %T into uid.ID", src)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (u *ID) Value() (driver.Value, error) {
	return string(*u), nil
}

// String implements fmt.Stringer interface.
func (u ID) String() string {
	return string(u)
}

// IsValid checks the length of an ID.
// And returns true if the ID length is valid.
func (u *ID) IsValid() bool {
	return len(*u) == IDLen
}