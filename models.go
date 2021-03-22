package models

import (
  "strings"
  "time"
  "os"
  "io/ioutil"
  "path/filepath"
  "reflect"

  "github.com/golang/glog"
  "github.com/jinzhu/gorm"

  _ "github.com/lib/pq"

  "gopkg.in/yaml.v2"

  "github.com/nats-io/nats.go"

  "github.com/Lunkov/lib-cache"
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/lib-map"
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

type CacheInfo struct {
  Mode          string      `yaml:"mode"`
  Expiry_time   int         `yaml:"expiry_time"`
  Redis         RedisInfo   `yaml:"redis"`
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
  CODE          string                     `json:"code"                         yaml:"code"`
  Order         int                        `json:"order"                        yaml:"order"`

  Title         string                     `json:"title"                        yaml:"title"`
  Columns       map[string]InfoColumn      `json:"columns"                      yaml:"columns"`
  
  BaseClass     string                     `json:"base_model"                   yaml:"base_model"`
  RefClass      reflect.Type               `json:"-"                            yaml:"-"`
  
  CacheConf     CacheInfo                  `json:"cache"                        yaml:"cache"`
  Cache         *cache.Cache               `json:"-"                            yaml:"-"`

  MultiLanguage bool                       `json:"multi_language"               yaml:"multi_language"`
  
  EventsStr     string                     `json:"events"                       yaml:"events"`
  EventsMask    TypeActionDB               `json:"events_mask"                  yaml:"events_mask"`
  
  SimpleUpdateMode      bool               `json:"simple_update_mode"           yaml:"simple_update_mode"`
  
  UseDeletedAt               bool          `json:"use_deleted_at"               yaml:"use_deleted_at"`
  UseGroupTablePermissions   bool          `json:"use_group_table_permissions"  yaml:"use_group_table_permissions"`
  UseUserTablePermissions    bool          `json:"use_user_table_permissions"   yaml:"use_user_table_permissions"`
  
  UseOwnerField              bool          `json:"use_owner_field"   yaml:"use_owner_field"`
  
  
  Permissions           map[string]ModelCRUD  `json:"acl_crud"                     yaml:"acl_crud"`
}

var dbHandleWrite     *gorm.DB
var dbHandleRead      *gorm.DB

var ncNatsMsg    *nats.Conn
var ecNatsMsg    *nats.EncodedConn
var configPath    string

var mods = make(map[string]ModelInfo)

func Init(connectStrWrite string, connectStrRead string, confPath string) bool {
  var err error
  
  configPath = confPath

  dbHandleWrite, err = gorm.Open("postgres", connectStrWrite)
  if err != nil {
    glog.Errorf("ERR: MODELS: failed to connect database (write): %v\n", err)
    return false
  }
  // Get generic database object sql.DB to use its functions
  sqlDB := dbHandleWrite.DB()
  // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
  sqlDB.SetMaxIdleConns(10)
  // SetMaxOpenConns sets the maximum number of open connections to the database.
  sqlDB.SetMaxOpenConns(100)
  // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
  sqlDB.SetConnMaxLifetime(time.Hour)
  
  dbHandleRead, err = gorm.Open("postgres", connectStrRead)
  if err != nil {
    glog.Errorf("ERR: MODELS: failed to connect database (read): %v\n", err)
    return false
  }
  // Get generic database object sql.DB to use its functions
  sqlDB = dbHandleRead.DB()
  // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
  sqlDB.SetMaxIdleConns(10)
  // SetMaxOpenConns sets the maximum number of open connections to the database.
  sqlDB.SetMaxOpenConns(100)
  // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
  sqlDB.SetConnMaxLifetime(time.Hour)


  if glog.V(9) {
    dbHandleWrite.LogMode(true)
    dbHandleRead.LogMode(true)
  }
  env.LoadFromFilesDB(dbHandleWrite, configPath + "/models", "", loadYAML)
  
  return true
}

func LoadData() {
  for _, model := range mods {
    glog.Infof("LOG: Init(%s)\n", model.CODE)
    dbTr := dbHandleWrite.Table(model.CODE).Begin()
    loadDataClass(dbTr, model.CODE, configPath + "/data", model.CODE)
    dbTr.Commit()
  }
}

func GetClass(modelName string) interface{} {
  model, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found", modelName)
    return nil
  }
  ref := getBaseByNameR(model.BaseClass)
  if !ok {
    glog.Errorf("ERR: BaseClass(%s) not found", model.BaseClass)
    return nil
  }
  return reflect.New(ref).Interface()
}

func loadDataClass(dbHandle *gorm.DB, modelCODE string, configPath string, tableName string) {
  filepath.Walk(configPath + "/" + tableName, func(filename string, f os.FileInfo, err error) error {
    if f != nil && f.IsDir() == false {
      if glog.V(2) {
        glog.Infof("FILE: %s\n", filename)
      }
      var err error
      yamlFile, err := ioutil.ReadFile(filename)
      if err != nil {
        glog.Errorf("ERR: ReadFile.yamlFile(%s)  #%v ", filename, err)
      } else {
        loadFileClass(dbHandle, modelCODE, filename, yamlFile)
      }
    }
    return nil
  })
}


func loadFileClass(dbHandle *gorm.DB, modelCODE string, filename string, yamlFile []byte) int {
  var err error
  var mapTmp = make(map[string]map[string]interface{})

  err = yaml.Unmarshal(yamlFile, mapTmp)
  if err != nil {
    glog.Errorf("ERR: ORG: yamlFile(%s): YAML: %v", filename, err)
  }
  if len(mapTmp) > 0 && dbHandle != nil {
    for _, data := range mapTmp {
      iv := GetClass(modelCODE)
      if iv == nil {
        continue
      }
      maps.ConvertFromMap(iv, &data)
      
      valTmp := iv

      code, ok := maps.GetFieldString(iv, "CODE")
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


func loadYAML(dbHandle *gorm.DB, filename string, yamlFile []byte) int {
  var err error
  var mapTmp = make(map[string]ModelInfo)

  err = yaml.Unmarshal(yamlFile, mapTmp)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  if(len(mapTmp) > 0) {
    for key, model := range mapTmp {
      rclass := getBaseByNameR(model.BaseClass)
      if rclass == nil {
        continue
      }
      model.RefClass = rclass
      model.EventsMask = aclCalcCRUD(model.EventsStr)
      aclRecalcCRUD(&model.Permissions)
      mods[key] = model
    }
  }

  return len(mapTmp)
}

func GetDBHandleRead() *gorm.DB {
  return dbHandleRead
}

func GetDBHandleWrite() *gorm.DB {
  return dbHandleWrite
}

func RunMethod(modelMethod string, args ...interface{}) bool {
  if glog.V(9) {
    glog.Infof("DBG: MODEL: RunMethod(%v)\n", modelMethod)
  }
  
  arM := strings.Split(modelMethod, ".")
  if len(arM) != 2 {
    glog.Errorf("ERR: RunMethod(%s) not found\n", modelMethod)
    return false
  }
  model, ok := mods[arM[0]]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", arM[0])
    return false
  }

  ref.RunMethodIfExists(getBaseByName(model.BaseClass), arM[1], args...)

  return true
}


func Close() {
  // baseClose()
  if dbHandleRead != nil {
    dbHandleRead.Close()
    dbHandleRead = nil
  }
  if dbHandleWrite != nil {
    dbHandleWrite.Close()
    dbHandleWrite = nil
  }
}
