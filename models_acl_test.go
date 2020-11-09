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
  
  
  assert.Equal(t, dbCreate, aclCalcCRUD("c"))
  assert.Equal(t, dbDelete, aclCalcCRUD("d"))
  assert.Equal(t, dbCreate | dbRead | dbUpdate | dbDelete, aclCalcCRUD("crud"))
  
  c := aclCalcCRUD("ocrud")
  assert.Equal(t, "ocrud", c.String())
  c = aclCalcCRUD("cru")
  assert.Equal(t, "cru", c.String())
  
  var res bool
  
  userAdmin := base.User{EMail: "admin@", Group: "user_crm", Groups: []string{"admin", "user_crm"}}
  userUser  := base.User{EMail: "user@",  Group: "user", Groups: []string{"admin_system", "user_crm"}}
  userOther := base.User{EMail: "guest@", Group: "user_crm", Groups: []string{"admin_system", "user_crm"}}
  
  modelOrg := ModelInfo{CODE: "org", Permissions: map[string]ModelCRUD{"admin": ModelCRUD{CRUD: "crud"}}}
  aclRecalcCRUD(&modelOrg.Permissions)

  res = aclCheck(&modelOrg, nil, dbCreate)
  assert.Equal(t, false, res)

  res = aclCheck(&modelOrg, &userAdmin, dbCreate)
  assert.Equal(t, true, res)

  res = aclCheck(&modelOrg, &userUser, dbCreate)
  assert.Equal(t, false, res)

  res = aclCheck(&modelOrg, &userOther, dbCreate)
  assert.Equal(t, false, res)


  modelMsg := ModelInfo{CODE: "messages", Permissions: map[string]ModelCRUD{"user": ModelCRUD{CRUD: "rud"}}}
  aclRecalcCRUD(&modelMsg.Permissions)

  res = aclCheck(&modelMsg, &userAdmin, dbCreate)
  assert.Equal(t, false, res)

  res = aclCheck(&modelMsg, &userUser, dbCreate)
  assert.Equal(t, false, res)
  res = aclCheck(&modelMsg, &userUser, dbRead)
  assert.Equal(t, true, res)
  res = aclCheck(&modelMsg, &userUser, dbUpdate)
  assert.Equal(t, true, res)
  res = aclCheck(&modelMsg, &userUser, dbDelete)
  assert.Equal(t, true, res)

  res = aclCheck(&modelMsg, &userOther, dbCreate)
  assert.Equal(t, false, res)
  
  
}
