package models

import (
  "reflect"

  "github.com/golang/glog"
  
  "github.com/jinzhu/gorm"

  "github.com/Lunkov/lib-ref"
)

func (db *DBConn) DBAutoMigrate(connectStr string) bool {
  dbm, err := gorm.Open("postgres", connectStr)
  if err != nil {
    glog.Errorf("ERR: MODELS: failed to connect database: %v\n", err)
    return false
  }
  if glog.V(9) {
    dbm.LogMode(true)
  }
  defer dbm.Close()
  
  createUserACL := false
  createGroupACL := false
  
  for _, model := range db.mods {
    class := db.getBaseByName(model.BaseClass)

    _, ok := ref.RunMethodIfExists(db.getBaseByName(model.BaseClass), "DBMigrate", db, model.CODE)

    if !ok && class != nil {
      cl := reflect.New(db.getBaseByNameR(model.BaseClass))
      dbm.Table(model.CODE).AutoMigrate(cl.Elem().Interface())
      
      if model.UseUserTablePermissions {
        dbm.Table(model.CODE + "_acl_user").AutoMigrate(PermissionUser{})
        createUserACL = true
      }
      if model.UseGroupTablePermissions {
        dbm.Table(model.CODE + "_acl_group").AutoMigrate(PermissionGroup{})
        createGroupACL = true
      }
    }
  }
  if createUserACL {
    dbm.Table(tblUserACL).AutoMigrate(UserACL{})
  }
  if createGroupACL {
    dbm.Table(tblGroupACL).AutoMigrate(GroupACL{})
  }
  return true
}
