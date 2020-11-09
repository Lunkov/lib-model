package fields

import (
  "errors"
  
  "time"
  
  "encoding/json"
  "database/sql/driver"
  
  "github.com/google/uuid"
)

type WorkGroup struct {
  ID            uuid.UUID  `json:"id,omitempty"`
  Name          string     `json:"name,omitempty"`
  Fields        string     `json:"fields"`
  FieldsRead    string     `json:"fields_read"`
  FieldsUpdate  string     `json:"fields_update"`
  ExpiredAt     time.Time  `json:"expired_at"`
}

type WorkGroups []WorkGroup

func (a WorkGroup) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *WorkGroup) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}

func (a WorkGroups) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *WorkGroups) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
