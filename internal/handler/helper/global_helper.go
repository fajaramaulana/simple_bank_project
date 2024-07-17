package helper

import (
	"strings"

	"github.com/google/uuid"
)

func ConvertStringToUUID(s string) (uuid.UUID, error) {
	uuidRes, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuidRes, nil
}

func Implode(currencies map[string]bool, separator string) string {
	keys := make([]string, 0, len(currencies))
	for key := range currencies {
		keys = append(keys, key)
	}
	return strings.Join(keys, separator)
}
