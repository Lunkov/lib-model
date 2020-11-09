package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
)

type CompanyContact struct {
  Name        string     `json:"name,omitempty"`
  Tel         string     `json:"tel,omitempty"`
  Mobile      string     `json:"mobile,omitempty"`
  Url         string     `json:"url,omitempty"`
  Email       string     `json:"email,omitempty"`
}

func (a CompanyContact) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *CompanyContact) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
