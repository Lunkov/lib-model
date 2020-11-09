package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
)

type BankAccount struct {
  BIK                    string     `json:"bik"                         yaml:"bik"`
  BankName               string     `json:"bank_name"                   yaml:"bank_name"`
  CorrespondentAccount   string     `json:"correspondent_account"       yaml:"correspondent_account"`
  Account                string     `json:"account"                     yaml:"account"`
  Currency_CODE          string     `json:"currency_code"               yaml:"currency_code"`
  Default                bool       `json:"default_account"             yaml:"default_account"`
}

type BankAccounts []BankAccount

func (a BankAccounts) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *BankAccounts) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
