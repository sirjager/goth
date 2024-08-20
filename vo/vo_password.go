package vo

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	baseValueObject[string]
}

type HashedPassword struct {
	baseValueObject[string]
}

const (
	PasswordMinLength = 8
	PasswordMaxLength = 200
)

func NewPassword(value string) (*Password, error) {
	password := &Password{baseValueObject[string]{value}}
	if err := password.Validate(); err != nil {
		return nil, err
	}
	return password, nil
}

func (v *Password) Validate() error {
	if len(v.value) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(v.value) > 255 {
		return errors.New("password must be at a atmost 255 characters long")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range v.value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	var missingRequirements []string
	if !hasUpper {
		missingRequirements = append(missingRequirements, "an uppercase letter")
	}
	if !hasLower {
		missingRequirements = append(missingRequirements, "a lowercase letter")
	}
	if !hasNumber {
		missingRequirements = append(missingRequirements, "a number")
	}
	if !hasSpecial {
		missingRequirements = append(missingRequirements, "a special character")
	}

	if len(missingRequirements) > 0 {
		errorMessage := fmt.Sprintf(
			"password is missing %s",
			strings.Join(missingRequirements, ", "),
		)
		return errors.New(errorMessage)
	}

	return nil
}

func (v *Password) HashPassword(cost ...int) (*HashedPassword, error) {
	_cost := bcrypt.DefaultCost + 2
	_minCost := bcrypt.MinCost + 4
	if len(cost) == 1 {
		if cost[0] < (_minCost) {
			return nil, fmt.Errorf("cost too low; must be at least %d", _minCost)
		}
		if cost[0] > bcrypt.MaxCost {
			return nil, fmt.Errorf("cost too high; must be at most %d", bcrypt.MaxCost)
		}
		_cost = cost[0]
	}
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(v.value), _cost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash: %w", err)
	}
	return &HashedPassword{baseValueObject[string]{string(hashedValue)}}, nil
}

func (v *Password) VerifyPassword(hashedValue string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(v.value))
}

func (v *HashedPassword) VerifyPassword(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(v.value), []byte(plainPassword))
}

func MustParsePassword(value string) *Password {
	password, err := NewPassword(value)
	if err != nil {
		panic(err)
	}
	return password
}

// MustParseHashedPassword will be used to assert to HashedPassword from database user object
func MustParseHashedPassword(hashedValue string) *HashedPassword {
	// name MustParseHashedPassword is used for maintain name consistency, like MustParseEmail, MustParseID etc
	return &HashedPassword{baseValueObject[string]{hashedValue}}
}
