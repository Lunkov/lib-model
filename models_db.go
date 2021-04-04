package models

import (
  "strings"
  "reflect"

  "encoding/json"
  
  "github.com/google/uuid"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-auth/base"
  "github.com/Lunkov/lib-ref"
)

const tblUserACL = "acl_user"
const tblGroupACL = "acl_group"


func (db *DBConn) DBTableGet(modelName string, user *base.User, fields []string, order []string, offset int, limit int) ([]byte, int, bool) {
  m, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return []byte(""), 0, false
  }

  if !db.aclCheck(&m, user, dbRead) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Read: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Read: User(NULL) => Model(%s)", modelName)
    }
    return []byte(""), 0, false
  }
  
  return db.dbTableGet(&m, modelName, "", fields, order, offset, limit)
}

func (db *DBConn) dbTableGet(m *ModelInfo, modelName string, where string, fields []string, order []string, offset int, limit int) ([]byte, int, bool) {

  count := 0
  if m.UseDeletedAt {
    if where != "" {
      db.HandleRead.Table(m.CODE).Select("count(id)").Where("deleted_at IS NULL AND " + where).Count(&count)
    } else {
      db.HandleRead.Table(m.CODE).Select("count(id)").Where("deleted_at IS NULL").Count(&count)
    }
  } else {
    db.HandleRead.Table(m.CODE).Select("count(id)").Where(where).Count(&count)
  }

  var rowres []map[string]interface{}
  var jsonRes []byte
  
  sql1 := db.HandleRead.Table(m.CODE)
  
  if len(fields) > 0 && fields[0] != "" {
    sql1 = sql1.Select(fields)
  }
  if len(order) > 0 {
    t_orders := make([]string, 0)
    r := strings.NewReplacer("desc", "", "acs", "") 
    p := reflect.New(db.getBaseByNameR(m.BaseClass))
    for _, item := range order {
      if ref.FieldExists(p, r.Replace(item)) {
        t_orders = append(t_orders, item)
      }
    }
    str_order := strings.Join(t_orders[:], ".")
    if glog.V(9) {
      glog.Errorf("DBG: Model(%s) Order: '%s'\n", modelName, str_order)
    }
    if str_order != "" {
      sql1 = sql1.Order(str_order)
    }
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
    row := db.GetClass(modelName)
    sql1.ScanRows(rows, row)
    
    rowMap := ref.ConvertToMap(row)
    rowres = append(rowres, rowMap)
  }
  jsonRes, _ = json.Marshal(rowres)
  return jsonRes, count, true
}

func (db *DBConn) DBInsert(modelName string, user *base.User, data *map[string]interface{}) bool {
  model, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return false
  }

  if !db.aclCheck(&model, user, dbCreate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Create: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Create: User(NULL) => Model(%s)", modelName)
    }
    return false
  }

  iv := db.GetClass(modelName)
  if model.UseOwnerField {
    if !db.aclSetOwner(modelName, user, data) {
      return false
    }
  }
  ref.ConvertFromMap(iv, data)
    
  sql1 := db.HandleWrite.Table(model.CODE)
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

  go db.sendNatsMsg(&model, dbCreate, iv)
  
  return true
}

func (db *DBConn) DBUpdate(modelName string, user *base.User, data *map[string]interface{}) bool {
  model, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return false
  }
  if !db.aclCheck(&model, user, dbUpdate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Update: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Update: User(NULL) => Model(%s)", modelName)
    }
    return false
  }

  sqlRead := db.HandleRead.Table(model.CODE)

  var err error
  err = nil

  readIV := db.GetClass(modelName)
  uid, oku := ref.GetMapFieldUUID(data, "ID")
  if !oku {
    glog.Errorf("ERR: Model(%s) not found ID\n", modelName)
    return false
  }
  if !db.aclDBCheckID(&model, user, uid, dbUpdate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Update: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Update: User(NULL) => Model(%s)", modelName)
    }
    return false
  }
  ok = ref.SetFieldUUID(readIV, uid, "ID")
  if !ok {
    glog.Errorf("ERR: Model(%s) Don`t SET ID<%v>\n", modelName, uid)
    return false
  }
  err = sqlRead.Where("id = ?", uid).First(readIV).Error
  if err != nil {
    glog.Errorf("ERR: Item(%s => ID<%v>) not found\n", modelName, uid)
    return false
  }
  iv := db.GetClass(modelName)
  if model.UseOwnerField {
    if !db.aclSetOwner(modelName, user, data) {
      return false
    }
  }
  if model.SimpleUpdateMode {
    ref.ConvertFromMap(iv, data)
  } else {
    dataRead := ref.ConvertToMap(readIV)
    ref.UnionMapsLK(&dataRead, data)
    ref.ConvertFromMap(iv, &dataRead)
  }

  /// Update
  sqlWrite := db.HandleWrite.Table(model.CODE)
  tr := sqlWrite.Begin()
  err = tr.Save(iv).Error
  if err != nil {
    glog.Errorf("ERR: Model(%s) DB INSERT: %v", modelName, err)
    tr.Rollback()
  } else {
    tr.Commit()
  }
  go db.sendNatsMsg(&model, dbUpdate, iv)
  
  return true
}


func (db *DBConn) DBGetItemByID(modelName string, user *base.User, id string) ([]byte, bool) {
  model, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return nil, false
  }
  if !db.aclCheck(&model, user, dbRead) {
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
  
  sql1 := db.HandleRead.Table(model.CODE)

  iv := db.GetClass(modelName)
  sql1 = sql1.Where("id = ?", uid)
  if model.UseOwnerField {
    sql1 = sql1.Where("owner->>'id' = ? OR work_groups @> '[{\"name\": \"?\"}]'", user.ID, user.Group)
  }
  err = sql1.First(iv).Error

  if err != nil {
    glog.Errorf("ERR: Model(%s) UUID<%v> not found. %v", modelName, id, err)
    return nil, false
  }
  
  rowMap := ref.ConvertToMap(iv)
  jsonRes, _ := json.Marshal(rowMap)
  
  return jsonRes, true
}

func (db *DBConn) DBGetItemByCODE(modelName string, user *base.User, code string) ([]byte, bool) {
  model, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return nil, false
  }
  if !db.aclCheck(&model, user, dbRead) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Read: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Read: User(NULL) => Model(%s)", modelName)
    }
    return nil, false
  }

  sql1 := db.HandleRead.Table(model.CODE)

  iv := db.GetClass(modelName)
  
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
  
  rowMap := ref.ConvertToMap(iv)
  jsonRes, _ := json.Marshal(rowMap)
  
  return jsonRes, true
}

  
func (db *DBConn) DBUpdateItemBy(modelName string, user *base.User, search string, param string, data *map[string]interface{}) bool {
  m, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found", modelName)
    return false
  }
  if !db.aclCheck(&m, user, dbUpdate) {
    if user != nil {
      glog.Errorf("ERR: Access denied. Update: User(%s) => Model(%s)", user.EMail, modelName)
    } else {
      glog.Errorf("ERR: Access denied. Update: User(NULL) => Model(%s)", modelName)
    }
    return false
  }

  iv := db.GetClass(modelName)
  ref.ConvertFromMap(iv, data)

  sql1 := db.HandleWrite.Table(m.CODE)
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

func (db *DBConn) DBDeleteItemByID(modelName string, user *base.User, id string) bool {
  model, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return false
  }
  if !db.aclCheck(&model, user, dbDelete) {
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
  
  sql1 := db.HandleWrite.Table(model.CODE)

  iv := db.GetClass(modelName)
  ref.SetFieldUUID(iv, uid, "ID")
  err = sql1.Delete(iv).Error
  // db.Unscoped().Delete(&order)
  //// DELETE FROM orders WHERE id=10;
  
  if err != nil {
    glog.Errorf("ERR: Model(%s) UUID<%v> not found\n", modelName, id)
    return false
  }
  go db.sendNatsMsg(&model, dbUpdate, iv)
  return true
}
