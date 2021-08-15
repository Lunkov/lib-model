package models

import (
  "reflect"
  "github.com/golang/glog"
)

func (db *DBConn) BaseAdd(name string, class reflect.Type) {
  db.typeRegistry[name] = class
}

func (db *DBConn) BaseCount() int {
  return len(db.typeRegistry)
}

func (db *DBConn) getBaseByName(name string) interface{} {
  item, ok := db.typeRegistry[name]
  if !ok {
    glog.Errorf("ERR: BaseModel(%s) not found", name)
    return nil
  }
  return reflect.ValueOf(item).Interface()
}

func (db *DBConn) getBaseByNameR(name string) reflect.Type {
  item, ok := db.typeRegistry[name]
  if !ok {
    glog.Errorf("ERR: BaseModel(%s) not found", name)
    return nil
  }
  return item
}

