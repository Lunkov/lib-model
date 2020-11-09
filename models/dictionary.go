package models

import (
  "time"
  "github.com/google/uuid"
)

////////////////////////////////
// Dictionary
///////////////////////////////
type FlatDictionary struct {
  ID             uuid.UUID     `db:"id"           json:"id"             yaml:"id"              gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at"   json:"created_at"     yaml:"created_at"      sql:"default: null"`
  UpdatedAt      time.Time     `db:"updated_at"   json:"updated_at"     yaml:"updated_at"      sql:"default: null"`
  DeletedAt     *time.Time     `db:"deleted_at"   json:"deleted_at"     yaml:"deleted_at"      sql:"default: null"`
  
  CODE           string        `db:"code"         json:"code"           yaml:"code"`
  Name           string        `db:"name"         json:"name"           yaml:"name"`
  Description    string        `db:"description"  json:"description"    yaml:"description"`
  
  Symbol         string        `db:"symbol"       json:"symbol"         yaml:"symbol"`
  Image          string        `db:"image"        json:"image"          yaml:"image"`
}

type Dictionary struct {
  ID             uuid.UUID     `db:"id"           json:"id"             yaml:"id"              gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at"   json:"created_at"     yaml:"created_at"      sql:"default: null"`
  UpdatedAt      time.Time     `db:"updated_at"   json:"updated_at"     yaml:"updated_at"      sql:"default: null"`
  DeletedAt     *time.Time     `db:"deleted_at"   json:"deleted_at"     yaml:"deleted_at"      sql:"default: null"`
  
  PARENT_ID      uuid.UUID     `db:"parent_id"    json:"parent_id"      yaml:"parent_id"       gorm:"column:parent_id;type:uuid;"  sql:"default: null"`
  
  CODE_ID        int           `db:"code_id"      json:"code_id"        yaml:"code_id"`
  CODE           string        `db:"code"         json:"code"           yaml:"code"`
  PARENT_CODE    string        `db:"parent_code"  json:"parent_code"    yaml:"parent_code"`
  Name           string        `db:"name"         json:"name"           yaml:"name"`
  Description    string        `db:"description"  json:"description"    yaml:"description"`
  
  Symbol         string        `db:"symbol"       json:"symbol"         yaml:"symbol"`
  Image          string        `db:"image"        json:"image"          yaml:"image"`
}
