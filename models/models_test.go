package models

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/lib-maps"
)


/////////////////////////
// TESTS
/////////////////////////
func TestOrganization(t *testing.T) {
  var org1 Organization
  org1.Name = "OOO `Org`"
  org1.AddressLegal.Country = "Russia"
  org1.AddressLegal.Index   = "127282"
  org1.AddressLegal.City    = "Moscow"
  org1.Bank = make(FBankAccounts, 0, 0)
  
  o1_need := map[string]interface{}{"address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "name":"OOO `Org`"}
  o1 := org1.ConvertToMap()
  assert.Equal(t, o1_need, o1)

  org2 := Organization{}
  
  maps.ConvertFromMap(&org2, &o1)
  assert.Equal(t, org1, org2)
}
