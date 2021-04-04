package models

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/lib-auth/base"
)


/////////////////////////
// TESTS
/////////////////////////
func TestCheckACL(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	// flag.Set("v", "9")
	flag.Parse()
  
  dbconn := New()
  
  assert.Equal(t, dbCreate, dbconn.aclCalcCRUD("c"))
  assert.Equal(t, dbDelete, dbconn.aclCalcCRUD("d"))
  assert.Equal(t, dbCreate | dbRead | dbUpdate | dbDelete, dbconn.aclCalcCRUD("crud"))
  
  c := dbconn.aclCalcCRUD("ocrud")
  assert.Equal(t, "ocrud", c.String())
  c = dbconn.aclCalcCRUD("cru")
  assert.Equal(t, "cru", c.String())
  
  var res bool
  
  userAdmin := base.User{EMail: "admin@", Group: "user_crm", Groups: []string{"admin", "user_crm"}}
  userUser  := base.User{EMail: "user@",  Group: "user", Groups: []string{"admin_system", "user_crm"}}
  userOther := base.User{EMail: "guest@", Group: "user_crm", Groups: []string{"admin_system", "user_crm"}}
  
  modelOrg := ModelInfo{CODE: "org", Permissions: map[string]ModelCRUD{"admin": ModelCRUD{CRUD: "crud"}}}
  dbconn.aclRecalcCRUD(&modelOrg.Permissions)

  res = dbconn.aclCheck(&modelOrg, nil, dbCreate)
  assert.Equal(t, false, res)

  res = dbconn.aclCheck(&modelOrg, &userAdmin, dbCreate)
  assert.Equal(t, true, res)

  res = dbconn.aclCheck(&modelOrg, &userUser, dbCreate)
  assert.Equal(t, false, res)

  res = dbconn.aclCheck(&modelOrg, &userOther, dbCreate)
  assert.Equal(t, false, res)


  modelMsg := ModelInfo{CODE: "messages", Permissions: map[string]ModelCRUD{"user": ModelCRUD{CRUD: "rud"}}}
  dbconn.aclRecalcCRUD(&modelMsg.Permissions)

  res = dbconn.aclCheck(&modelMsg, &userAdmin, dbCreate)
  assert.Equal(t, false, res)

  res = dbconn.aclCheck(&modelMsg, &userUser, dbCreate)
  assert.Equal(t, false, res)
  res = dbconn.aclCheck(&modelMsg, &userUser, dbRead)
  assert.Equal(t, true, res)
  res = dbconn.aclCheck(&modelMsg, &userUser, dbUpdate)
  assert.Equal(t, true, res)
  res = dbconn.aclCheck(&modelMsg, &userUser, dbDelete)
  assert.Equal(t, true, res)

  res = dbconn.aclCheck(&modelMsg, &userOther, dbCreate)
  assert.Equal(t, false, res)
  
  
}
