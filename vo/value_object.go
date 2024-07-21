package vo

type ValueObject interface {
	Value() string
	Validate() error
}
