package main

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() (Env, error) {
	var envMap map[string]string
	envMap, err := godotenv.Read()
	if err != nil {
		return nil, fmt.Errorf("godotenv failed to load environment variables: %s", err)
	}

	port, err := strconv.ParseUint(envMap["PORT"], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address port: %s", err)
	}

	e := env{
		host:   envMap["HOST"],
		port:   uint16(port),
		db_uri: envMap["DB_URI"],
	}

	return e, nil
}

type Env interface {
	Host() string
	Port() uint16
	DatabaseUri() string
}

type env struct {
	host   string
	port   uint16
	db_uri string
}

func (e env) Host() string {
	return e.host
}

func (e env) Port() uint16 {
	return e.port
}

func (e env) DatabaseUri() string {
	return e.db_uri
}
