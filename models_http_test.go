package models

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
)


/////////////////////////
// HTTP TESTS
/////////////////////////
func TestHTTPWhere(t *testing.T) {
  flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()

  in := "(grade.gte.90,student.is.true,age.gt.14)"
  o1_need := "grade >= 90 AND student IS true AND age > 14"
  o1 := parserAND(in)
  assert.Equal(t, o1_need, o1)

  in = "(grade.gte.90,student.is.true,age.gt.14)"
  o1_need = "grade >= 90 OR student IS true OR age > 14"
  o1 = parserOR(in)
  assert.Equal(t, o1_need, o1)

  in = ""
  o1_need = ""
  o1 = parserOR(in)
  assert.Equal(t, o1_need, o1)

}
