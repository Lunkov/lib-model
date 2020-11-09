package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
  
  "github.com/google/uuid"
)

type RelashionShip struct {
  ID          uuid.UUID  `json:"id,omitempty"`
  Name        string     `json:"name,omitempty"`
}

type RelashionShips []RelashionShip

func (a RelashionShip) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *RelashionShip) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}

func (a RelashionShips) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *RelashionShips) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
