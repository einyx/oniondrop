package config

import (
  "os"
  "github.com/einyx/oniondrop/pkg"
)

func GetConfig() *root.Config {
  return &root.Config {
    Mongo: &root.MongoConfig {
      Ip: envOrDefaultString("github.com/einyx/oniondrop:mongo:ip", "127.0.0.1:27017"),
      DbName: envOrDefaultString("github.com/einyx/oniondrop:mongo:dbName", "myDb")},
    Server: &root.ServerConfig { Port: envOrDefaultString("github.com/einyx/oniondrop:server:port", ":1377")},
    Auth: &root.AuthConfig { Secret: envOrDefaultString("github.com/einyx/oniondrop:auth:secret", "mysecret")}}
}

func envOrDefaultString(envVar string, defaultValue string) string {
  value := os.Getenv(envVar)
  if value == "" {
    return defaultValue;
  }
  
  return value
}