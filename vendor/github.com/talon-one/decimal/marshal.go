package decimal

import "bytes"

// UnmarshalText implements the encoding.TextMarshaler interface for serialization
func (d Decimal) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for deserialization
func (d *Decimal) UnmarshalText(buf []byte) error {
	tmp, err := NewFromString(string(buf))
	if err != nil {
		return err
	}
	d.native = tmp.native
	return nil
}

// MarshalJSON implements the json.Marshaler interface for serialization
func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.MarshalText()
}

// UnmarshalText implements the json.Unmarshaler interface for deserialization
func (d *Decimal) UnmarshalJSON(buf []byte) error {
	return d.UnmarshalText(bytes.Trim(bytes.TrimSpace(buf), `"`))
}
