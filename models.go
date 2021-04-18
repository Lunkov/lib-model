package models

import (
  "strings"
  "time"
  "os"
  "io/ioutil"
  "path/filepath"
  "reflect"
  
  "github.com/google/uuid"

  "github.com/golang/glog"
  "github.com/jinzhu/gorm"

  _ "github.com/lib/pq"

  "gopkg.in/yaml.v2"

  "github.com/nats-io/nats.go"

  "github.com/Lunkov/lib-cache"
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-ref"
)

type InfoColumn struct {
  Title                 string  `json:"title"                yaml:"title"`
  Type                  string  `json:"type"                 yaml:"type"`
}

type RedisInfo struct {
  Url               string   `yaml:"url"`
  Max_connections   int      `yaml:"max_connections"`
}

type ModelCRUD struct {
  CRUD          string                     `json:"crud"                         yaml:"crud"`
  CRUDm         TypeActionDB               `json:"crudm"                        yaml:"crudm"`
  Fields        string                     `json:"fields"                       yaml:"fields"`
  FieldsCreate  string                     `json:"fields_create"                yaml:"fields_create"`
  FieldsRead    string                     `json:"fields_read"                  yaml:"fields_read"`
  FieldsUpdate  string                     `json:"fields_update"                yaml:"fields_update"`
  FieldsDelete  string                     `json:"fields_delete"                yaml:"fields_delete"`
}

type ModelInfo struct {
  CODE          string                              `json:"code"                         yaml:"code"`
  Order         int                                 `json:"order"                        yaml:"order"`

  Title         string                              `json:"title"                        yaml:"title"`
  Columns                    map[string]InfoColumn  `json:"columns"                      yaml:"columns"`
  
  BaseClass                  string                 `json:"base_model"                   yaml:"base_model"`
  RefClass                   reflect.Type           `json:"-"                            yaml:"-"`
  
  CacheConf                  cache.CacheConfig      `json:"cache"                        yaml:"cache"`
  Cache                      cache.ICache           `json:"-"                            yaml:"-"`

  MultiLanguage              bool                   `json:"multi_language"               yaml:"multi_language"`
  
  EventsStr                  string                 `json:"events"                       yaml:"events"`
  EventsMask                 TypeActionDB           `json:"events_mask"                  yaml:"events_mask"`
  
  SimpleUpdateMode           bool                   `json:"simple_update_mode"           yaml:"simple_update_mode"`
  
  UseDeletedAt               bool                   `json:"use_deleted_at"               yaml:"use_deleted_at"`
  UseGroupTablePermissions   bool                   `json:"use_group_table_permissions"  yaml:"use_group_table_permissions"`
  UseUserTablePermissions    bool                   `json:"use_user_table_permissions"   yaml:"use_user_table_permissions"`
   
  UseOwnerField              bool                   `json:"use_owner_field"   yaml:"use_owner_field"`
  
  
  Permissions                map[string]ModelCRUD   `json:"acl_crud"                     yaml:"acl_crud"`
}

type DBConn struct {
  configPath        string
  
  HandleWrite      *gorm.DB
  HandleRead       *gorm.DB

  ncNatsMsg        *nats.Conn
  ecNatsMsg        *nats.EncodedConn

  mods              map[string]ModelInfo
  
  cacheGroups       map[string]uuid.UUID
  typeRegistry      map[string]reflect.Type
}


func New() *DBConn {
  return &DBConn{ mods: make(map[string]ModelInfo),
                  cacheGroups: make(map[string]uuid.UUID),
                  typeRegistry: make(map[string]reflect.Type),
                  }
}

func (db *DBConn) Init(connectStrWrite string, connectStrRead string, confPath string) bool {
  var err error
  
  db.configPath = confPath

  if len(connectStrWrite) > 10 {
    db.HandleWrite, err = gorm.Open("postgres", connectStrWrite)
    if err != nil {
      glog.Errorf("ERR: MODELS: failed to connect database (write): %v\n", err)
      return false
    }
    // Get generic database object sql.DB to use its functions
    sqlDB := db.HandleWrite.DB()
    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
    sqlDB.SetMaxIdleConns(10)
    // SetMaxOpenConns sets the maximum number of open connections to the database.
    sqlDB.SetMaxOpenConns(100)
    // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
    sqlDB.SetConnMaxLifetime(time.Hour)
  }
  
  if len(connectStrRead) > 10 {
    db.HandleRead, err = gorm.Open("postgres", connectStrRead)
    if err != nil {
      glog.Errorf("ERR: MODELS: failed to connect database (read): %v\n", err)
      return false
    }
    // Get generic database object sql.DB to use its functions
    sqlDB := db.HandleRead.DB()
    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
    sqlDB.SetMaxIdleConns(10)
    // SetMaxOpenConns sets the maximum number of open connections to the database.
    sqlDB.SetMaxOpenConns(100)
    // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
    sqlDB.SetConnMaxLifetime(time.Hour)
  }


  if glog.V(9) {
    db.HandleWrite.LogMode(true)
    db.HandleRead.LogMode(true)
  }
  env.LoadFromFilesDB(db.HandleWrite, db.configPath + "/models", "", db.loadYAML)
  
  return true
}

func (db *DBConn) LoadData() {
  for _, model := range db.mods {
    glog.Infof("LOG: Init(%s)", model.CODE)
    dbTr := db.HandleWrite.Table(model.CODE).Begin()
    db.loadDataClass(dbTr, model.CODE, db.configPath + "/data", model.CODE)
    dbTr.Commit()
  }
}

func (db *DBConn) GetClass(modelName string) interface{} {
  model, ok := db.mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found", modelName)
    return nil
  }
  ref := db.getBaseByNameR(model.BaseClass)
  if !ok {
    glog.Errorf("ERR: BaseClass(%s) not found", model.BaseClass)
    return nil
  }
  return reflect.New(ref).Interface()
}

func (db *DBConn) loadDataClass(dbHandle *gorm.DB, modelCODE string, configPath string, tableName string) {
  filepath.Walk(configPath + "/" + tableName, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false {
      if glog.V(2) {
        glog.Infof("FILE: %s", filename)
      }
      var err error
      yamlFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.yamlFile(%s)  #%v ", filename, err)
      } else {
        db.loadFileClass(dbHandle, modelCODE, filename, yamlFile)
      }
    }
    return nil
  })
}


func (db *DBConn) loadFileClass(dbHandle *gorm.DB, modelCODE string, filename string, yamlFile []byte) int {
  var err error
  var mapTmp = make(map[string]map[string]interface{})

  err = yaml.Unmarshal(yamlFile, mapTmp)
  if err != nil {
    glog.Errorf("ERR: ORG: yamlFile(%s): YAML: %v", filename, err)
  }
  if len(mapTmp) > 0 && dbHandle != nil {
    for _, data := range mapTmp {
      iv := db.GetClass(modelCODE)
      if iv == nil {
        continue
      }
      ref.ConvertFromMap(iv, &data)
      
      valTmp := iv

      code, ok := ref.GetFieldString(iv, "CODE")
      if ok {
        if errW := dbHandle.First(valTmp, "code = ?", code).Error; errW != nil {
          err = dbHandle.Create(iv).Error
        } else {
          err = dbHandle.Save(iv).Error
        }
        if err != nil {
          glog.Errorf("ERR: ORG: yaml.Marshal: yamlFile(%s): %s: %v (value = %v)", filename, code, err, data)
        }
      }
    }
  }
  return len(mapTmp)
}


func (db *DBConn) loadYAML(dbHandle *gorm.DB, filename string, yamlFile []byte) int {
  var err error
  var mapTmp = make(map[string]ModelInfo)

  err = yaml.Unmarshal(yamlFile, mapTmp)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  if(len(mapTmp) > 0) {
    for key, model := range mapTmp {
      rclass := db.getBaseByNameR(model.BaseClass)
      if rclass == nil {
        continue
      }
      model.RefClass = rclass
      model.EventsMask = db.aclCalcCRUD(model.EventsStr)
      db.aclRecalcCRUD(&model.Permissions)
      db.mods[key] = model
    }
  }

  return len(mapTmp)
}

func (db *DBConn) GetDBHandleRead() *gorm.DB {
  return db.HandleRead
}

func (db *DBConn) GetDBHandleWrite() *gorm.DB {
  return db.HandleWrite
}

func (db *DBConn) RunMethod(modelMethod string, args ...interface{}) bool {
  if glog.V(9) {
    glog.Infof("DBG: MODEL: RunMethod(%v)\n", modelMethod)
  }
  
  arM := strings.Split(modelMethod, ".")
  if len(arM) != 2 {
    glog.Errorf("ERR: RunMethod(%s) not found\n", modelMethod)
    return false
  }
  model, ok := db.mods[arM[0]]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", arM[0])
    return false
  }

  ref.RunMethodIfExists(db.getBaseByName(model.BaseClass), arM[1], args...)

  return true
}


func (db *DBConn) Close() {
  // baseClose()
  if db.HandleRead != nil {
    db.HandleRead.Close()
    db.HandleRead = nil
  }
  if db.HandleWrite != nil {
    db.HandleWrite.Close()
    db.HandleWrite = nil
  }
}
