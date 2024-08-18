package config

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type Config struct {
	Host                   string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Host:                   getEnv("HOST", "http://localhost"),
		Port:                   getEnv("PORT", "42069"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "ecom"),
		JWTExpirationInSeconds: getEnvInt("JWT_EXP", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not_secret_anymore_i_guess?"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if val, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
