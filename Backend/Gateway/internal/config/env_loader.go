package config

import (
    "github.com/joho/godotenv"
    "log"
)

func LoadDotEnv(path string) {
    err := godotenv.Load(path)
    if err != nil {
        log.Printf("No .env file found at %s (ignoring)", path)
    }
}
