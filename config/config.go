package config

import (
	"log"
	"os"
	"time"
)

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabaseUser        string
	DatabasePassword    string
	MigrateToVersion    string
	MigrationLocation   string
	FileStorageLocation string
	JwtSecret           string
	GoogleClientId      string
	GoogleClientSecret  string
	JwtTTL              time.Duration
}

func GetConfiguration() Configuration {
	return Configuration{
		DatabaseName:        getOrFail("DB_NAME"),
		DatabaseHost:        getOrFail("DB_HOST"),
		DatabaseUser:        getOrFail("DB_USER"),
		DatabasePassword:    getOrFail("DB_PASSWORD"),
		MigrateToVersion:    getOrDefault("MIGRATE", "latest"),
		MigrationLocation:   getOrDefault("MIGRATION_LOCATION", "migrations"),
		FileStorageLocation: getOrDefault("FILES_LOCATION", "storage/posts"),
		JwtSecret:           getOrFail("JWT_SECRET"),
		GoogleClientId:      getOrFail("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  getOrFail("GOOGLE_CLIENT_SECRET"),
		JwtTTL:              time.Hour * 24,
	}
}

func getOrFail(key string) string {
	env, set := os.LookupEnv(key)
	if !set || env == "" {
		log.Fatalf("%s env var is missing", key)
	}
	return env
}

func getOrDefault(key, defaultVal string) string {
	env, set := os.LookupEnv(key)
	if !set {
		return defaultVal
	}
	return env
}
