package decimal

import (
	"database/sql/driver"
)

// Value implements the driver.Valuer interface for database serialization.
func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (d *Decimal) Scan(value interface{}) error {
	dec, err := NewFromInterface(value)
	if err != nil {
		return err
	}
	d.nat = dec.nat
	return nil
}
