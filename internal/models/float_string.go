package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type FloatString float64

func (f *FloatString) UnmarshalJSON(b []byte) error {
	var asFloat float64
	var asString string

	// Intentar parsear como número primero
	if err := json.Unmarshal(b, &asFloat); err == nil {
		*f = FloatString(asFloat)
		return nil
	}

	// Si no, intentar parsear como string que contiene un número
	if err := json.Unmarshal(b, &asString); err != nil {
		return fmt.Errorf("FloatString: no se puede interpretar como float ni como string: %w", err)
	}

	parsed, err := strconv.ParseFloat(asString, 64)
	if err != nil {
		return fmt.Errorf("FloatString: error convirtiendo string a float: %w", err)
	}

	*f = FloatString(parsed)
	return nil
}
