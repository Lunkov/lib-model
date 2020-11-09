package models

import (
  "time"
  "github.com/google/uuid"
  
  "github.com/Lunkov/lib-model/fields"
)

////////////////////////////////
// Document
///////////////////////////////

type Biomaterial struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`

  // Каталожный Номер
  CatalogNum     string        `db:"catalog_num"         json:"catalog_num"           yaml:"catalog_num"`
  // Происхождение
  Origin         string        `db:"origin"       json:"origin"         yaml:"origin"`
 
  // Дата поступления
  AdmissionAt      time.Time     `db:"admission_at;default: now()"  json:"admission_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
 
  Name           string        `db:"name"         json:"name"           yaml:"name"`
  Description    string        `db:"description"  json:"description"    yaml:"description"`

  Owner    fields.RelashionShip    `db:"owner"        json:"owner,ommitempty"                gorm:"type:jsonb;"`

  Authors  fields.RelashionShips   `db:"authors"         json:"authors,ommitempty"               gorm:"type:jsonb;"`
}
