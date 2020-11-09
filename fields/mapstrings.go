package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
)


type MapStrings map[string]string

func (p MapStrings) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *MapStrings) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]string)
	if !ok {
		return errors.New("Type assertion .(map[string]string) failed.")
	}

	return nil
}
