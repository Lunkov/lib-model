package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
)

type CompanyPerson struct {
  Position    string     `json:"position,omitempty"         yaml:"position"`
  Legal       string     `json:"legal,omitempty"            yaml:"legal"`
  FirstName   string     `json:"first_name,omitempty"       yaml:"first_name"`
  LastName    string     `json:"last_name,omitempty"        yaml:"last_name"`
  MiddleName  string     `json:"middle_name,omitempty"      yaml:"middle_name"`
  Tel         string     `json:"tel,omitempty"              yaml:"tel"`
  Mobile      string     `json:"mobile,omitempty"           yaml:"mobile"`
  Email       string     `json:"email,omitempty"            yaml:"email"`
  Photo       string     `json:"photo,omitempty"            yaml:"photo"`
}

func (a CompanyPerson) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *CompanyPerson) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }
  return json.Unmarshal(b, &a)
}
