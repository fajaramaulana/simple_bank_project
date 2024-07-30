package validate

import "github.com/fajaramaulana/simple_bank_project/internal/grpcapi/helper"

var currency = []string{"USD", "IDR", "EUR"}

func ValidateUsername(value string) error {
	// required
	if err := helper.ValidateRequired(value); err != nil {
		return err
	}

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

func ValidateUsernameNoRequired(value string) error {
	if value != "" {
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
	return nil
}

func ValidatePassword(value string) error {
	// required
	if err := helper.ValidateRequired(value); err != nil {
		return err
	}

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

func ValidatePasswordNotRequired(value string) error {
	if value != "" {
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
	return nil
}

func ValidateCurrency(value string) error {
	if err := helper.ValidateSlice(value, currency); err != nil {
		return err
	}

	return nil
}

func ValidateCurrencyNotRequired(value string) error {
	if value != "" {
		if err := helper.ValidateSlice(value, currency); err != nil {
			return err
		}
		return nil
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

func ValidateFullNameNotRequired(value string) error {
	if value != "" {
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
	return nil
}

func ValidateUserUUID(value string) error {
	// required
	if err := helper.ValidateRequired(value); err != nil {
		return err
	}

	// UUID
	if err := helper.ValidateUUID(value); err != nil {
		return err
	}

	return nil
}
