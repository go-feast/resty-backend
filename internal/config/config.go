package config

import "os"

// EnvironmentProvider by default os.Getenv
type EnvironmentProvider interface {
	Getenv(key string) string
}

type envFunc func(string) string

func (f envFunc) Getenv(key string) string {
	return f(key)
}

var envProvider EnvironmentProvider = envFunc(os.Getenv)

func SetEnvironmentProvider(p func(string) string) {
	envProvider = envFunc(p)
}

func DBConn() string {
	return envProvider.Getenv("DB")
}

func Addr() string {
	return envProvider.Getenv("ADDR")
}
