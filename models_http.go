package models

import (
  "strings"
  "strconv"
  "fmt"
  "io"
  "errors"

  "net/http"
  "github.com/golang/gddo/httputil/header"

  "encoding/json"
  
  "github.com/golang/glog"
  "github.com/Lunkov/lib-auth/base"
)

type malformedRequest struct {
  status int
  msg    string
}

func (mr *malformedRequest) Error() string {
  return mr.msg
}

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

func HTTPTableGet(r *http.Request, modelName string, user *base.User) ([]byte, int, bool) {
  m, ok := mods[modelName]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", modelName)
    return []byte(""), 0, false
  }

  if !aclCheck(&m, user, dbRead) {
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
  
  return dbTableGet(&m, modelName, where, ar_fields, ar_order, offset, limit)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
  if r.Header.Get("Content-Type") != "" {
    value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
    if value != "application/json" {
      msg := "Content-Type header is not application/json"
      return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
    }
  }

  r.Body = http.MaxBytesReader(w, r.Body, 1048576)

  dec := json.NewDecoder(r.Body)
  // TODO
  // dec.DisallowUnknownFields()

  err := dec.Decode(&dst)
  if err != nil {
    var syntaxError *json.SyntaxError
    var unmarshalTypeError *json.UnmarshalTypeError

    switch {
    case errors.As(err, &syntaxError):
        msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.Is(err, io.ErrUnexpectedEOF):
        msg := fmt.Sprintf("Request body contains badly-formed JSON")
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.As(err, &unmarshalTypeError):
        msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case strings.HasPrefix(err.Error(), "json: unknown field "):
        fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
        msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case errors.Is(err, io.EOF):
        msg := "Request body must not be empty"
        return &malformedRequest{status: http.StatusBadRequest, msg: msg}

    case err.Error() == "http: request body too large":
        msg := "Request body must not be larger than 1MB"
        return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

    default:
        return err
    }
  }

  err = dec.Decode(&struct{}{})
  if err != io.EOF {
    msg := "Request body must only contain a single JSON object"
    return &malformedRequest{status: http.StatusBadRequest, msg: msg}
  }

  return nil
}
