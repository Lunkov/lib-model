package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
)

type CompanyRegister struct {
  INN         string     `json:"inn,omitempty"       yaml:"inn"`
  OGRN        string     `json:"ogrn,omitempty"      yaml:"ogrn"`
}

func (a CompanyRegister) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *CompanyRegister) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
