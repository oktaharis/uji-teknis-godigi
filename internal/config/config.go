package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv     string
	Port       string
	DBDSN      string
	JWTSecret  string
	JWTExpires int64
}

func Load() *Config {
	return &Config{
		AppEnv:     get("APP_ENV", "development"),
		Port:       get("PORT", "8080"),
		DBDSN:      get("DB_DSN", "root:@tcp(127.0.0.1:3306)/godigi?parseTime=true&loc=Local"),
		JWTSecret:  get("JWT_SECRET", "supersecret_change_me"),
		JWTExpires: toInt64(get("JWT_EXPIRES_IN", "3600")),
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func toInt64(s string) int64 {
	var n int64 = 3600
	_, _ = fmt.Sscan(s, &n)
	return n
}
