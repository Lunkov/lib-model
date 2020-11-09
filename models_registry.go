package models

import (
  "reflect"
  "github.com/golang/glog"
)

var typeRegistry = make(map[string]reflect.Type)

func BaseAdd(name string, class reflect.Type) {
  typeRegistry[name] = class
}

func BaseCount() int {
  return len(typeRegistry)
}

func getBaseByName(name string) interface{} {
  item, ok := typeRegistry[name]
  if !ok {
    glog.Errorf("ERR: BaseModel(%s) not found\n", name)
    return nil
  }
	return reflect.ValueOf(item).Interface()
}

func getBaseByNameR(name string) reflect.Type {
  item, ok := typeRegistry[name]
  if !ok {
    glog.Errorf("ERR: BaseModel(%s) not found\n", name)
    return nil
  }
	return item
}

