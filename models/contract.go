package models

import (
  "time"
  "github.com/google/uuid"
  
  "github.com/Lunkov/lib-model/fields"
)

////////////////////////////////
// Contract
///////////////////////////////

type Contract struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`

  
  Number       string        `db:"num_doc"            json:"num_doc"                               gorm:"column:num_doc;type:varchar(64)"`
  StartDT      time.Time     `db:"start_time"         json:"start_time"        yaml:"start_time"   sql:"default: null"`
  FinishDT     time.Time     `db:"finish_time"        json:"finish_time"       yaml:"finish_time"  sql:"default: null"`

  Customer     fields.RelashionShip   `db:"customer"           json:"customer,ommitempty"                 gorm:"type:jsonb;"`
  Executor     fields.RelashionShip   `db:"executor"           json:"executor,ommitempty"                 gorm:"type:jsonb;"`
  
  Partners     fields.RelashionShips  `db:"partners"           json:"partners,ommitempty"                 gorm:"type:jsonb;"`
}
