package models

import (
  "strings"
  "strconv"

  "net/http"
  
  "github.com/golang/glog"
  "github.com/Lunkov/lib-auth/base"
)

func parserWhere(in string, connect string) string {
  where := ""
  if len(in) < 2 {
    return where
  }
  in = in[1:len(in)-1]
  arP := strings.Split(in, ",")
  for _, p := range arP {
    item := strings.Split(p, ".")
    if len(item) == 3 {
      if len(where) > 0 {
        where += " " + connect + " "
      }
      where += item[0] + " " + operators(item[1]) + " " + item[2]
    }
  }
  return where
}
  
func parserAND(in string) string {
  return parserWhere(in, "AND")
}

func parserOR(in string) string {
  return parserWhere(in, "OR")
}

func operators(op string) string {
  switch (op) {
  case "eq":
     return "="
  case "gt":
     return ">"
  case "gte":
     return ">="
  case "lt":
     return "<"
  case "lte":
     return "<="
  case "neq":
     return "!="
  case "is":
     return "IS"
  case "not":
     return "IS NOT"
  }
  return ""
}

func (db *DBConn) HTTPTableGet(r *http.Request, modelName string, user *base.User) ([]byte, int, bool) {
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

  r.ParseForm()
  offset, _ := strconv.Atoi(r.Form.Get("offset"))
  limit, _ := strconv.Atoi(r.Form.Get("limit"))
  fields := r.Form.Get("select")
  ar_fields := strings.Split(fields, ",")
  order := strings.ReplaceAll(r.Form.Get("order"), ".", " ")
  ar_order := strings.Split(order, ",")
  
  where := parserAND(r.Form.Get("and"))
  
  return db.dbTableGet(&m, modelName, where, ar_fields, ar_order, offset, limit)
}
