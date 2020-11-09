package fields

import (
  "errors"

  "encoding/json"
  "database/sql/driver"
  
  "github.com/google/uuid"
)

type Owner struct {
  ID              uuid.UUID  `json:"id,omitempty"`
  Login           string    `json:"login"`
  EMail           string    `json:"email"`
  Avatar          string    `json:"avatar"`
  DisplayName     string    `json:"displayname"`
}

func (a Owner) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *Owner) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
