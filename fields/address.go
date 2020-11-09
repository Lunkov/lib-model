package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
)

type Address struct {
  Country        string     `json:"country"          yaml:"country"`
  Index          string     `json:"index"            yaml:"index"`
  City           string     `json:"city"             yaml:"city"`
  Street         string     `json:"street"           yaml:"street"`
  House          string     `json:"house"            yaml:"house"`
  Room           string     `json:"room"             yaml:"room"`
}

func (a Address) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *Address) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
