package models

import (
  "time"
  "github.com/google/uuid"
  
  "github.com/Lunkov/lib-model/fields"
)

////////////////////////////////
// Organization
///////////////////////////////

type Organization struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`

  CODE           string    `db:"code"         json:"code"          yaml:"code"             gorm:"type:varchar(96);default: null"`

  Name           string    `db:"name"         json:"name"          yaml:"name"             sql:"column:name"        gorm:"column:name;type:varchar(256)"`
  Description    string    `db:"description"  json:"description"   yaml:"description"      gorm:"default: null"`
  
  UrlLogo        string    `db:"url_logo"     json:"url_logo"      yaml:"url_logo"`
  UrlIcon        string    `db:"url_icon"     json:"url_icon"      yaml:"url_icon"`
  UrlMain        string    `db:"url_main"     json:"url_main"      yaml:"url_main"`
  
  Support          fields.CompanyContact   `db:"support"            json:"support"                        yaml:"support"             gorm:"column:support;type:jsonb;"`
  
  CEO              fields.CompanyPerson    `db:"ceo"                json:"ceo,omitempty"                  yaml:"ceo"                 gorm:"column:ceo;type:jsonb;"`
  ChiefAccountant  fields.CompanyPerson    `db:"chief_accountant"   json:"chief_accountant,omitempty"     yaml:"chief_accountant"    gorm:"column:chief_accountant;type:jsonb;"`
  Signer           fields.CompanyPerson    `db:"signer"             json:"signer,omitempty"               yaml:"signer"              gorm:"column:signer;type:jsonb;"`

  AddressLegal      fields.Address         `db:"address_legal"             json:"address_legal,omitempty"                 gorm:"type:jsonb;"`
  AddressBilling    fields.Address         `db:"address_billing"           json:"address_billing,omitempty"               gorm:"type:jsonb;"`
  AddressShipping   fields.Address         `db:"address_shipping"          json:"address_shipping,omitempty"              gorm:"type:jsonb;"`
  
  Register       fields.CompanyRegister   `db:"register"       json:"register,ommitempty"             gorm:"type:jsonb;"`
  Bank           fields.BankAccounts      `db:"bank"           json:"bank,ommitempty"                 gorm:"type:jsonb;"`
}
