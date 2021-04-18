package fields

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/lib-ref"
)


/////////////////////////
// TESTS
/////////////////////////
func TestAddress(t *testing.T) {
  var adr Address
  adr.Country = "Russia"
  adr.Index   = "127282"
  adr.City    = "Moscow"
  
  m_need := map[string]interface{}{"city":"Moscow", "country":"Russia", "index":"127282"}
  m := ref.ConvertToMap(adr)
  assert.Equal(t, m_need, m)

}
