package vo

import (
	"github.com/google/uuid"
)

type ID struct {
	baseValueObject[uuid.UUID]
}

func NewID() (*ID, error) {
	_uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	id := &ID{baseValueObject[uuid.UUID]{_uuid}}
	if err = id.Validate(); err != nil {
		return nil, err
	}
	return id, nil
}

func NewIDFrom(value string) (*ID, error) {
	_uuid, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}
	id := &ID{baseValueObject[uuid.UUID]{_uuid}}
	if err = id.Validate(); err != nil {
		return nil, err
	}
	return id, nil
}

// MustParseID returns ID if valid else panics
func MustParseID(value string) *ID {
	id := &ID{baseValueObject[uuid.UUID]{uuid.MustParse(value)}}
	if err := id.Validate(); err != nil {
		panic(err)
	}
	return id
}

func (v *ID) Validate() error {
	// currently not needed
	return nil
}
