package models

import (
  "github.com/golang/glog"

  "github.com/nats-io/nats.go"
)

func (db *DBConn) InitEvents(NatsUrl string) bool {
  if glog.V(2) {
    glog.Infof("LOG: Init Events")
  }
  var err error
  natInit := false
  db.ncNatsMsg, err = nats.Connect(NatsUrl)
  if err != nil {
    glog.Errorf("ERR: MODEL: NATS Connect(%s): %v", NatsUrl, err)
  } else {
    db.ecNatsMsg, err = nats.NewEncodedConn(db.ncNatsMsg, nats.JSON_ENCODER)
    if err != nil {
      glog.Errorf("ERR: MODEL: NATS NewEncodedConn: %v", err)
    } else {
      natInit = true
    }
  }
  return natInit
}

func contains(s []string, e string) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}

func (db *DBConn) SendNatsMsg(model string, event TypeActionDB, values *map[string][]string) {
  m, ok := db.mods[model]
  if !ok {
    glog.Errorf("ERR: Model(%s) not found\n", model)
    return
  }
  db.sendNatsMsg(&m, event, values)
}

func (db *DBConn) sendNatsMsg(model *ModelInfo, event TypeActionDB, values interface{}) {
  if model != nil  && db.ecNatsMsg != nil {
    if (model.EventsMask & event) != 0 {
      subject := model.CODE + "." + event.String()
      if glog.V(9) {
        glog.Infof("LOG: MODEL: sendNatsMsg(%s)", subject)
      }
      if err := db.ecNatsMsg.Publish(subject, values); err != nil {
        glog.Errorf("ERR: MODEL: sendNatsMsg(%s) err=%v", subject, err)
      }
    }
  }
}
