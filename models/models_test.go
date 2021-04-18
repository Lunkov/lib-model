package models

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/lib-ref"
  "github.com/Lunkov/lib-model/fields"
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
  org1.Bank = make(fields.BankAccounts, 0, 0)
  
  o1_need := map[string]interface{}{"address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127282", "name":"OOO `Org`"}
  o1 := ref.ConvertToMap(&org1)
  assert.Equal(t, o1_need, o1)

  org2 := Organization{}
  
  ref.ConvertFromMap(&org2, &o1)
  assert.Equal(t, org1, org2)
}
