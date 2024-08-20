package vo

import "encoding/json"

type ValueObject[T comparable] interface {
	Value() T
	Validate() error
	IsEqual(other T) bool
}

type baseValueObject[T comparable] struct {
	value T
}

func (b *baseValueObject[T]) Value() T {
	return b.value
}

func (b *baseValueObject[T]) IsEqual(other T) bool {
	return b.value == other
}

// Marshal serializes the value object into JSON
func (b *baseValueObject[T]) Marshal() ([]byte, error) {
	return json.Marshal(b.value)
}

// Unmarshal deserializes the JSON data into the value object
func (b *baseValueObject[T]) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &b.value)
}
