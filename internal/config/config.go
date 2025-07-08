package config

import "github.com/Ahmed1monm/Axis-BE-assessment/pkg/utils"

type Config struct {
	MongoURI     string
	DatabaseName string
	Port         string
	Environment  string
}

func Load() *Config {
	return &Config{
		MongoURI:     utils.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: utils.GetEnv("DB_NAME", "axis_assessment"),
		Port:         utils.GetEnv("PORT", "8080"),
		Environment:  utils.GetEnv("ENV", "development"),
	}
}
