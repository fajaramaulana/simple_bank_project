package validate

import "github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"

var currency = []string{"USD", "IDR", "EUR"}

func ValidateUsername(value string) error {
	// min 6
	if err := helper.ValidateMin(value, 6); err != nil {
		return err
	}

	// max 15
	if err := helper.ValidateMax(value, 15); err != nil {
		return err
	}

	// only alphanumeric
	if err := helper.ValidateAlphanum(value); err != nil {
		return err
	}

	return nil
}

func ValidatePassword(value string) error {
	// min 8
	if err := helper.ValidateMin(value, 8); err != nil {
		return err
	}

	// only alphanumeric
	if err := helper.ValidateAlphanum(value); err != nil {
		return err
	}
	return nil
}

func ValidateCurrency(value string) error {
	// min 3
	if err := helper.ValidateSlice(value, currency); err != nil {
		return err
	}

	return nil
}

func ValidateFullName(value string) error {
	// min 3
	if err := helper.ValidateMin(value, 3); err != nil {
		return err
	}

	// max 50
	if err := helper.ValidateMax(value, 50); err != nil {
		return err
	}

	// only alphabethic with space
	if err := helper.ValidateAlphaSpace(value); err != nil {
		return err
	}

	return nil
}
