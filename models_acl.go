package models

import (
  "strings"
  "time"
  "github.com/google/uuid"
  "github.com/golang/glog"

  "github.com/Lunkov/lib-auth/base"
)

type TypeActionDB uint

const (
   dbUndef     TypeActionDB = 0
   dbCreate    TypeActionDB = 1
   dbRead      TypeActionDB = 2
   dbUpdate    TypeActionDB = 4
   dbDelete    TypeActionDB = 8
   dbOwner     TypeActionDB = 16
   MAX         TypeActionDB = 32
)

type UserACL struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  
  EMail           string        `db:"email"         json:"email"           yaml:"email"`
  Login           string        `db:"login"         json:"login"           yaml:"login"`
}

type GroupACL struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  
  Name           string        `db:"name"         json:"name"           yaml:"name"`
  Description    string        `db:"description"  json:"description"    yaml:"description"`
}

type PermissionGroup struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`

  Group_ID       uuid.UUID     `db:"group_id"                   json:"group_id"      yaml:"group_id"         gorm:"column:group_id;type:uuid"`
  Object_ID      uuid.UUID     `db:"object_id"                  json:"object_id"     yaml:"object_id"        gorm:"column:object_id;type:uuid"`
  
  IsOwner        bool          `db:"is_owner"                   json:"is_owner"      yaml:"is_owner"         gorm:"column:is_owner"`
  CanCreate      bool          `db:"can_create"                 json:"can_create"    yaml:"can_create"       gorm:"column:can_create"`
  CanRead        bool          `db:"can_read"                   json:"can_read"      yaml:"can_read"         gorm:"column:can_read"`
  CanUpdate      bool          `db:"can_update"                 json:"can_update"    yaml:"can_update"       gorm:"column:can_update"`
  CanDelete      bool          `db:"can_delete"                 json:"can_delete"    yaml:"can_delete"       gorm:"column:can_delete"`
}

type PermissionUser struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`

  User_ID        uuid.UUID     `db:"user_id"                    json:"user_id"       yaml:"user_id"          gorm:"column:user_id;type:uuid"`
  Object_ID      uuid.UUID     `db:"object_id"                  json:"object_id"     yaml:"object_id"        gorm:"column:object_id;type:uuid"`
  
  IsOwner        bool          `db:"is_owner"                   json:"is_owner"      yaml:"is_owner"         gorm:"column:is_owner"`
  CanCreate      bool          `db:"can_create"                 json:"can_create"    yaml:"can_create"       gorm:"column:can_create"`
  CanRead        bool          `db:"can_read"                   json:"can_read"      yaml:"can_read"         gorm:"column:can_read"`
  CanUpdate      bool          `db:"can_update"                 json:"can_update"    yaml:"can_update"       gorm:"column:can_update"`
  CanDelete      bool          `db:"can_delete"                 json:"can_delete"    yaml:"can_delete"       gorm:"column:can_delete"`
}

var cacheGroups = make(map[string]uuid.UUID)

func aclCheck(model *ModelInfo, user *base.User, action TypeActionDB) bool {
  if user == nil && len(model.Permissions) < 1 {
    return true
  }
  if user == nil && len(model.Permissions) > 0 {
    glog.Errorf("ERR: MODELS ACL: User(NULL) BUT len(model.Permissions) = %d", len(model.Permissions))
    return false
  }
  crud, ok := findGroupCRUD(model, user)
  if !ok {
    glog.Errorf("ERR: MODELS ACL: USER(%s). Access denied. MODEL '%s' CRUD not found", user.EMail, model.CODE)
    return false
  }
  if (crud.CRUDm & action) == 0 {
    glog.Errorf("ERR: MODELS ACL: USER(%s). Action denied.", user.EMail)
    return false
  }
  return true
}

func aclDBCheckID(model *ModelInfo, user *base.User, object_id uuid.UUID, action TypeActionDB) bool {
  return true
}

func aclSetOwner(modelName string, user *base.User, data *map[string]interface{}) bool {
  if user == nil {
    glog.Errorf("ERR: Access denied. Create|Update (UseOwnerField): User(NULL) => aclSetOwner(%s)", modelName)
    return false
  }
  (*data)["owner.id"] = user.ID
  (*data)["owner.login"] = user.Login
  (*data)["owner.email"] = user.EMail
  (*data)["owner.displayname"] = user.DisplayName
  (*data)["owner.avatar"] = user.Avatar
  return true
}

func aclIsOwner(modelName string, user *base.User, data *map[string]interface{}) bool {
  if user == nil {
    glog.Errorf("ERR: Access denied. Update (UseOwnerField): User(NULL) => aclIsOwner(%s)", modelName)
    return false
  }
  id, ok := (*data)["owner.id"]
  if !ok {
    glog.Errorf("ERR: Access denied. Update (UseOwnerField) aclIsOwner(%s): Not Found Owner ID.", modelName)
    return false
  }
  return id == user.ID
}

func findGroupCRUD(model *ModelInfo, user *base.User) (*ModelCRUD, bool) {
  crud, ok := model.Permissions[user.Group]
  if ok {
    return &crud, true
  }
  for _, group := range user.Groups {
    crud, ok := model.Permissions[group]
    if ok {
      return &crud, true
    }
  }
  crud, ok = model.Permissions["other"]
  if ok {
    return &crud, true
  }
  return nil, false
}

func aclRecalcCRUD(permissions *map[string]ModelCRUD) {
  for group, perm := range (*permissions) {
    perm.CRUDm = aclCalcCRUD(perm.CRUD)
    (*permissions)[group] = perm
  }
}

func aclCalcCRUD(crud string) TypeActionDB {
  crudLower := strings.ToLower(crud)
  c := dbUndef
  if strings.Contains(crudLower, "o") {
    c |= dbOwner
  }
  if strings.Contains(crudLower, "c") {
    c |= dbCreate
  }
  if strings.Contains(crudLower, "r") {
    c |= dbRead
  }
  if strings.Contains(crudLower, "u") {
    c |= dbUpdate
  }
  if strings.Contains(crudLower, "d") {
    c |= dbDelete
  }
  return c
}

func (crud *TypeActionDB) String() string {
  res := ""
  if ((*crud) & dbOwner) != 0 {
    res += "o"
  }
  if ((*crud) & dbCreate) != 0 {
    res += "c"
  }
  if ((*crud) & dbRead) != 0 {
    res += "r"
  }
  if ((*crud) & dbUpdate) != 0 {
    res += "u"
  }
  if ((*crud) & dbDelete) != 0 {
    res += "d"
  }
  return res
}

func aclFilterFields(model *ModelInfo, user *base.User, action TypeActionDB) (string, bool) {
  return "", false
}
