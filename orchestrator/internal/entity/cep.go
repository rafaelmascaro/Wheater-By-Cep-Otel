package entity

import (
	"errors"
	"regexp"
)

var ErrInvalidZipcode = errors.New("invalid zipcode")

type CEP string

func NewCEP(cep string) (CEP, error) {
	c := CEP(cep)
	err := c.IsValid()
	if err != nil {
		return "", err
	}
	return c, nil
}

func (c CEP) IsValid() error {
	regex := regexp.MustCompile(`^\d{8}$`)
	if !regex.MatchString(string(c)) {
		return ErrInvalidZipcode
	}
	return nil
}
