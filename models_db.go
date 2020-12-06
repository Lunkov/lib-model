package models

import (
  "strings"
  "reflect"

  "encoding/json"
  
  "github.com/google/uuid"
  "github.com/golang/glog"
  
  "github.com/jinzhu/gorm"

  "github.com/Lunkov/lib-auth/base"
  "github.com/Lunkov/lib-map"
  "github.com/Lunkov/lib-reflect"
)

const tblUserACL = "acl_user"
const tblGroupACL = "acl_group"


func DBAutoMigrate(connectStr string) bool {
  db, err := gorm.Open("postgres", connectStr)
  if err != nil {
    glog.Errorf("ERR: MODELS: failed to connect database: %v\n", err)
    return false
  }
  if glog.V(9) {
    db.LogMode(true)
  }
  defer db.Close()
  
  createUserACL := false
  createGroupACL := false
  
  for _, model := range mods {
    class := getBaseByName(model.BaseClass)

    _, ok := ref.RunMethodIfExists(model.CODE, getBaseByName(model.BaseClass), "DBMigrate", db, model.CODE)

    if !ok && class != nil {
      cl := reflect.New(getBaseByNameR(model.BaseClass))
      db.Table(model.CODE).AutoMigrate(cl.Elem().Interface())
      
      if model.UseUserTablePermissions {
        db.Table(model.CODE + "_acl_user").AutoMigrate(PermissionUser{})
        createUserACL = true
      }
      if model.UseGroupTablePermissions {
        db.Table(model.CODE + "_acl_group").AutoMigrate(PermissionGroup{})
        createGroupACL = true
      }
    }
  }
  if createUserACL {
    db.Table(tblUserACL).AutoMigrate(UserACL{})
  }
  if createGroupACL {
    db.Table(tblGroupACL).AutoMigrate(GroupACL{})
  }
  return true
}

func DBTableGet(modelName string, user *base.User, fields []string, order []string, offset int, limit int) ([]byte, int, bool) {
  m, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return []byte(""), 0, false
  }

  if !aclCheck(&m, user, dbRead) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Read: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Read: User(NULL) => Model(%s)", modelName)
    }
    return []byte(""), 0, false
  }
  
  return dbTableGet(&m, modelName, "", fields, order, offset, limit)
}

func dbTableGet(m *ModelInfo, modelName string, where string, fields []string, order []string, offset int, limit int) ([]byte, int, bool) {

  count := 0
  if m.UseDeletedAt {
    if where != "" {
      dbHandleRead.Table(m.CODE).Select("count(id)").Where("deleted_at IS NULL AND " + where).Count(&count)
    } else {
      dbHandleRead.Table(m.CODE).Select("count(id)").Where("deleted_at IS NULL").Count(&count)
    }
  } else {
    dbHandleRead.Table(m.CODE).Select("count(id)").Where(where).Count(&count)
  }

  var rowres []map[string]interface{}
  var jsonRes []byte
  
  sql1 := dbHandleRead.Table(m.CODE)
  
  if len(fields) > 0 && fields[0] != "" {
    sql1 = sql1.Select(fields)
  }

  str_order := strings.Join(order[:], ",")
  if glog.V(9) {
    glog.Errorf("DBG: Model(%s) Order: '%s'\n", modelName, str_order)
  }
  if str_order != "" {
    sql1 = sql1.Order(str_order)
  }
  sql1 = sql1.Offset(offset)
  if limit > 0 {
    sql1 = sql1.Limit(limit)
  }
  if m.UseDeletedAt {
    sql1 = sql1.Where("deleted_at IS NULL")
  }
  if glog.V(9) {
    glog.Errorf("DBG: Model(%s) WHERE: '%s'", modelName, where)
  }
  if where != "" {
    sql1 = sql1.Where(where)
  }
  
  rows, err := sql1.Rows()
  if err != nil {
    glog.Errorf("ERR: Model(%s) SQL err: %v\n", modelName, err)
    return []byte(""), count, false
  }
  defer rows.Close()

  for rows.Next() {
    row := GetClass(modelName)
    sql1.ScanRows(rows, row)
    
    rowMap := maps.ConvertToMap(row)
    rowres = append(rowres, rowMap)
  }
  jsonRes, _ = json.Marshal(rowres)
  return jsonRes, count, true
}

func DBInsert(modelName string, user *base.User, data *map[string]interface{}) bool {
  model, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return false
  }

  if !aclCheck(&model, user, dbCreate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Create: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Create: User(NULL) => Model(%s)", modelName)
    }
    return false
  }

  iv := GetClass(modelName)
  if model.UseOwnerField {
    if !aclSetOwner(modelName, user, data) {
      return false
    }
  }
  maps.ConvertFromMap(iv, data)
    
  sql1 := dbHandleWrite.Table(model.CODE)
  tr := sql1.Begin()
  err := tr.Create(iv).Error
  if err != nil {
    glog.Errorf("ERR: Model(%s) DB INSERT: %v", modelName, err)
    tr.Rollback()
    return false
  } else {
    tr.Commit()
    if glog.V(9) {
      glog.Errorf("DBG: Model(%s) Insert(data=%v)", modelName, iv)
    }
  }

  go sendNatsMsg(&model, dbCreate, iv)
  
  return true
}

func DBUpdate(modelName string, user *base.User, data *map[string]interface{}) bool {
  model, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return false
  }
  if !aclCheck(&model, user, dbUpdate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Update: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Update: User(NULL) => Model(%s)", modelName)
    }
    return false
  }

  sqlRead := dbHandleRead.Table(model.CODE)

  var err error
  err = nil

  readIV := GetClass(modelName)
  uid, oku := maps.GetMapFieldUUID(data, "ID")
  if !oku {
    glog.Errorf("ERR: Model(%s) not found ID\n", modelName)
    return false
  }
  if !aclDBCheckID(&model, user, uid, dbUpdate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Update: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Update: User(NULL) => Model(%s)", modelName)
    }
    return false
  }
  ok = maps.SetFieldUUID(readIV, uid, "ID")
  if !ok {
    glog.Errorf("ERR: Model(%s) Don`t SET ID<%v>\n", modelName, uid)
    return false
  }
  err = sqlRead.Where("id = ?", uid).First(readIV).Error
  if err != nil {
    glog.Errorf("ERR: Item(%s => ID<%v>) not found\n", modelName, uid)
    return false
  }
  iv := GetClass(modelName)
  if model.UseOwnerField {
    if !aclSetOwner(modelName, user, data) {
      return false
    }
  }
  if model.SimpleUpdateMode {
    maps.ConvertFromMap(iv, data)
  } else {
    dataRead := maps.ConvertToMap(readIV)
    maps.UnionMaps(&dataRead, data)
    maps.ConvertFromMap(iv, &dataRead)
  }

  /// Update
  sqlWrite := dbHandleWrite.Table(model.CODE)
  tr := sqlWrite.Begin()
  err = tr.Save(iv).Error
  if err != nil {
    glog.Errorf("ERR: Model(%s) DB INSERT: %v", modelName, err)
    tr.Rollback()
  } else {
    tr.Commit()
  }
  go sendNatsMsg(&model, dbUpdate, iv)
  
  return true
}


func DBGetItemByID(modelName string, user *base.User, id string) ([]byte, bool) {
  model, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return nil, false
  }
  if !aclCheck(&model, user, dbRead) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Read: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Read: User(NULL) => Model(%s)", modelName)
    }
    return nil, false
  }

  uid, err := uuid.Parse(id)
  if err != nil {
    glog.Errorf("ERR: Model(%s). ID<%s> is not UUID. %v", modelName, id, err)
    return nil, false
  }
  
  sql1 := dbHandleRead.Table(model.CODE)

  iv := GetClass(modelName)
  sql1 = sql1.Where("id = ?", uid)
  if model.UseOwnerField {
    sql1 = sql1.Where("owner->>'id' = ? OR work_groups @> '[{\"name\": \"?\"}]'", user.ID, user.Group)
  }
  err = sql1.First(iv).Error

  if err != nil {
    glog.Errorf("ERR: Model(%s) UUID<%v> not found. %v", modelName, id, err)
    return nil, false
  }
  
  rowMap := maps.ConvertToMap(iv)
  jsonRes, _ := json.Marshal(rowMap)
  
  return jsonRes, true
}

func DBGetItemByCODE(modelName string, user *base.User, code string) ([]byte, bool) {
  model, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return nil, false
  }
  if !aclCheck(&model, user, dbRead) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Read: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Read: User(NULL) => Model(%s)", modelName)
    }
    return nil, false
  }

  sql1 := dbHandleRead.Table(model.CODE)

  iv := GetClass(modelName)
  
  var err error
  
  sql1 = sql1.Where("code = ?", code)
  if model.UseOwnerField {
    sql1 = sql1.Where("owner->>'id' = ? OR work_groups @> '[{\"name\": ?}]'", user.ID, user.Group)
  }

  err = sql1.First(iv).Error
  
  if err != nil {
    glog.Errorf("ERR: Model(%s) CODE<%v> not found\n", modelName, code)
    return nil, false
  }
  
  rowMap := maps.ConvertToMap(iv)
  jsonRes, _ := json.Marshal(rowMap)
  
  return jsonRes, true
}

  
func DBUpdateItemBy(modelName string, user *base.User, search string, param string, data *map[string]interface{}) bool {
  m, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found", modelName)
    return false
  }
  if !aclCheck(&m, user, dbUpdate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Update: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Update: User(NULL) => Model(%s)", modelName)
    }
    return false
  }

  iv := GetClass(modelName)
  maps.ConvertFromMap(iv, data)

  sql1 := dbHandleWrite.Table(m.CODE)
  x := iv
  err := sql1.Where(search, param).First(x).Error
  
  tr := sql1.Begin()
  if err != nil {
    tr.Create(iv)
  } else {
    tr.Save(iv)
  }
  tr.Commit()
  return true
}

func DBDeleteItemByID(modelName string, user *base.User, id string) bool {
  model, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return false
  }
  if !aclCheck(&model, user, dbDelete) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Delete: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Delete: User(NULL) => Model(%s)", modelName)
    }
    return false
  }

  uid, err := uuid.Parse(id)
  if err != nil {
    glog.Errorf("ERR: Model(%s). ID<%s> is not UUID. %v", modelName, id, err)
    return false
  }
  
  sql1 := dbHandleWrite.Table(model.CODE)

  iv := GetClass(modelName)
  maps.SetFieldUUID(iv, uid, "ID")
  err = sql1.Delete(iv).Error
  // db.Unscoped().Delete(&order)
  //// DELETE FROM orders WHERE id=10;
  
  if err != nil {
    glog.Errorf("ERR: Model(%s) UUID<%v> not found\n", modelName, id)
    return false
  }
  go sendNatsMsg(&model, dbUpdate, iv)
  return true
}
