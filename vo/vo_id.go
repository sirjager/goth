package vo

import (
	"github.com/google/uuid"
)

type ID struct {
	value string
}

// NewIDFrom returns a new ID from a string or an error
func NewIDFrom(value string) (*ID, error) {
	id := &ID{value: value}
	if err := id.Validate(); err != nil {
		return nil, err
	}
	return id, nil
}

func NewID() (*ID, error) {
	// generate id
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &ID{id.String()}, nil
}

func GenerateTestID() *ID {
	id, err := NewID()
	if err != nil {
		panic(err)
	}
	return id
}

func (v *ID) Value() string {
	return v.value
}

func (v *ID) IsEqual(other *ID) bool {
	return v.value == other.value
}

// MustParseID returns ID if valid else panics
func MustParseID(value string) *ID {
	id := &ID{value: value}
	if err := id.Validate(); err != nil {
		panic(err)
	}
	return id
}

func (v *ID) Validate() error {
	id, err := uuid.Parse(v.value)
	if err != nil {
		return err
	}
	v.value = id.String()
	return nil
}
