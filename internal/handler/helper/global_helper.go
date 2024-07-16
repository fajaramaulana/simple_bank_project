package helper

import "github.com/google/uuid"

func ConvertStringToUUID(s string) (uuid.UUID, error) {
	uuidRes, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuidRes, nil
}
