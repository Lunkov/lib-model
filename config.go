package models

import (
  "fmt"
)


type PostgreSQLInfo struct {
  Host       string    `yaml:"host"`
  Port       int       `yaml:"port"`
  Name       string    `yaml:"name"`
  User       string    `yaml:"user"`
  Password   string    `yaml:"password"`
  ConnectStr string    `yaml:"connect_string"`
}

func ConnectStr(cfg PostgreSQLInfo) string {
  return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
                          cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}
