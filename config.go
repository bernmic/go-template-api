package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

const (
	defaultHttpPort = 8080
)

type Configuration struct {
	Http HttpConfiguration `json:"http,omitempty", yaml:"http,omitempty"`
}

type HttpConfiguration struct {
	Port int `json:"port,omitempty", yaml:"port,omitempty"`
}

func NewConfig() *Configuration {
	c := Configuration{}
	flag.IntVar(&c.Http.Port, "port", defaultHttpPort, "http port to listen")
	flag.Parse()
	arguments(&c)
	return &c
}

func arguments(c *Configuration) {
	c.Http.Port = intConfig("port", c.Http.Port, "PORT", 8080)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func intConfig(parmName string, parmValue int, envName string, defaultValue int) int {
	if isFlagPassed(parmName) {
		return parmValue
	} else if val, ok := os.LookupEnv(envName); ok {
		p, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("format for int is wrong: %s", val)
		}
		return p
	}
	return defaultValue
}

func stringConfig(parmName string, parmValue string, envName string, defaultValue string) string {
	if isFlagPassed(parmName) {
		return parmValue
	} else if val, ok := os.LookupEnv(envName); ok {
		return val
	}
	return defaultValue
}

func boolConfig(parmName string, parmValue bool, envName string, defaultValue bool) bool {
	if isFlagPassed(parmName) {
		return parmValue
	} else if val, ok := os.LookupEnv(envName); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("format for bool is wrong: %s", val)
		}
		return b
	}
	return defaultValue
}
