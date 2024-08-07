package vo

import (
	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID
}

// NewIDFrom returns a new ID from a string or an error
func NewIDFrom(value string) (*ID, error) {
	id := &ID{value: uuid.MustParse(value)}
	if err := id.Validate(); err != nil {
		return nil, err
	}
	return id, nil
}

func NewID() (*ID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &ID{id}, nil
}

func (v *ID) Value() uuid.UUID {
	return v.value
}

func (v *ID) IsEqual(other *ID) bool {
	return v.value.String() == other.value.String()
}

// MustParseID returns ID if valid else panics
func MustParseID(value string) *ID {
	id := &ID{value: uuid.MustParse(value)}
	return id
}

func (v *ID) Validate() error {
	return nil
}
