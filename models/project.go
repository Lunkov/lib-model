package models

import (
  "time"
  "github.com/google/uuid"
  
  "github.com/Lunkov/lib-model/fields"
)

////////////////////////////////
// Project
///////////////////////////////

type Project struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`

  CODE           string    `db:"code"         json:"code"          yaml:"code"             gorm:"type:varchar(96);default: null"`

  Name           string    `db:"name"         json:"name"          yaml:"name"             sql:"column:name"        gorm:"column:name;type:varchar(256)"`
  Description    string    `db:"description"  json:"description"   yaml:"description"      gorm:"default: null"`

  StartAt      time.Time     `db:"start_at;default: now()"   json:"start_at"     sql:"default: now()"    gorm:"type:timestamp with time zone"`
  FinishAt     time.Time     `db:"finish_at;default: now()"  json:"finish_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  
  UrlLogo        string      `db:"url_logo"     json:"url_logo"      yaml:"url_logo"`

  Owner          fields.Owner           `db:"owner"              json:"owner,ommitempty"                 gorm:"type:jsonb;"`

  WorkGroups     fields.WorkGroups      `db:"work_groups"        json:"work_groups,ommitempty"               gorm:"type:jsonb;"`

  ContractDoc    fields.RelashionShip   `db:"contract"           json:"contract,ommitempty"                 gorm:"type:jsonb;"`
  Customer       fields.RelashionShip   `db:"customer"           json:"customer,ommitempty"                 gorm:"type:jsonb;"`
  Executor       fields.RelashionShip   `db:"executor"           json:"executor,ommitempty"                 gorm:"type:jsonb;"`
  
  BudgetEstimate float32    `db:"budget_estimate"     json:"budget_estimate"      yaml:"budget_estimate"`
  BudgetNow      float32    `db:"budget_now"          json:"budget_now"           yaml:"budget_now"`
}
