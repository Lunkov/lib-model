package fields

import (
  "errors"
  "encoding/json"
  "database/sql/driver"
)

type File struct {
  Name                   string     `json:"name"                        yaml:"name"`
  Description            string     `json:"description"                 yaml:"description"`
  FileName               string     `json:"file_name"                   yaml:"file_name"`
  FileType               string     `json:"file_type"                   yaml:"file_type"`
}

type Files []File

func (a Files) Value() (driver.Value, error) {
  return json.Marshal(a)
}

func (a *Files) Scan(value interface{}) error {
  b, ok := value.([]byte)
  if !ok {
      return errors.New("type assertion to []byte failed")
  }

  return json.Unmarshal(b, &a)
}
