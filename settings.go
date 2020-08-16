package main

import (
	"os"
	"strconv"
)

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Region    string
}

type GpgKeyConfig struct {
	// (Bucket, Name) or LocalPath not empty
	Bucket    string
	Name      string
	LocalPath string
}

type AppConfig struct {
	LocalS3   S3Config
	TargetS3  S3Config
	Key       GpgKeyConfig
	DebugMode bool
}

func GetConfig() *AppConfig {
	return &AppConfig{
		LocalS3: S3Config{
			Endpoint:  "minio:9000",
			AccessKey: "minio",
			SecretKey: "miniostorage",
			UseSSL:    false,
			Region:    "us-east-1",
		},
		TargetS3: S3Config{
			Endpoint:  getEnv("TARGET_ENDPOINT", ""),
			AccessKey: getEnv("TARGET_ACCESS_KEY", ""),
			SecretKey: getEnv("TARGET_SECRET_KEY", ""),
			UseSSL:    getEnvAsBool("TARGET_SSL", false),
			Region:    getEnv("TARGET_REGION", "us-east-1"),
		},
		Key: GpgKeyConfig{
			Bucket:    getEnv("GPG_KEY_TARGET_BUCKET", ""),
			Name:      getEnv("GPG_KEY_TARGET_NAME", ""),
			LocalPath: getEnv("GPG_KEY_LOCAL_PATH", ""),
		},
		DebugMode: getEnvAsBool("DEBUG_MODE", false),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
