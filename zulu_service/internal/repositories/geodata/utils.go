package geodata

import "database/sql"

func getPointerIfValid(value sql.NullFloat64) *float64 {
	if value.Valid {
		return &value.Float64
	}
	return nil
}
