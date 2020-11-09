package models

import (
  "time"
  "github.com/google/uuid"
  "github.com/lib/pq"
)

////////////////////////////////
// Person
///////////////////////////////

type Person struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`

  Login         string          `json:"login"         db:"login"`
  EMail         string          `json:"email"         db:"email"                                gorm:"column:email;unique;not null"`
  DisplayName   string          `json:"display_name"  db:"display_name"`
  Avatar        string          `json:"avatar"        db:"avatar"`
  Group         string          `json:"group"`
  Groups        pq.StringArray  `json:"groups"        sql:"column:groups;type:varchar(64)[]"    gorm:"column:groups;type:varchar(64)[]"` // 
  TimeLogin     time.Time       `json:"-"`
  AuthCode      string          `json:"-"`
  Disable       bool            `json:"disable"       db:"disable"`
}
