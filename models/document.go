package models

import (
  "time"
  "github.com/google/uuid"
  
  "github.com/Lunkov/lib-model/fields"
)

////////////////////////////////
// Document
///////////////////////////////

type Document struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  
  Number       string        `db:"num_doc"            json:"num_doc"                               gorm:"column:num_doc;type:varchar(64)"`
  StartDT      time.Time     `db:"start_time"         json:"start_time"        yaml:"start_time"   sql:"default: null"`
  FinishDT     time.Time     `db:"finish_time"        json:"finish_time"       yaml:"finish_time"  sql:"default: null"`

  FromOrg      fields.RelashionShip    `db:"from_org"        json:"from_org,ommitempty"                gorm:"type:jsonb;"`
  ToOrgs       fields.RelashionShips   `db:"to_orgs"         json:"to_orgs,ommitempty"                 gorm:"type:jsonb;"`
  
  Files        fields.RelashionShips   `db:"files"         json:"files,ommitempty"                 gorm:"type:jsonb;"`
}
