package helper

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func MapToSlice(m map[string]string, separator string) []string {
	var result []string
	for key, value := range m {
		result = append(result, fmt.Sprintf("%s%s%s", key, separator, value))
	}
	return result
}

func NumericToBigInt(n pgtype.Numeric) *big.Int {
	result := new(big.Int)
	result.SetBytes(n.Int.Bytes())
	exp := int64(n.Exp)
	if exp > 0 {
		result.Mul(result, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(exp), nil))
	} else if exp < 0 {
		result.Div(result, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(-exp), nil))
	}
	return result
}

func bigIntToFloat64(n *big.Int) float64 {
	floatVal, _ := new(big.Float).SetInt(n).Float64()
	return floatVal
}

func NumerictoFloat64(n pgtype.Numeric) float64 {
	return bigIntToFloat64(NumericToBigInt(n))
}

func NumericToString(n pgtype.Numeric) string {
	return NumericToBigInt(n).String()
}
