package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config 는 애플리케이션 전역 설정이다. 환경변수에서 로딩한다.
type Config struct {
	AppEnv   string
	HTTPPort string
	DB       DBConfig
}

// DBConfig 는 데이터베이스 연결 설정이다.
type DBConfig struct {
	Driver     string // sqlite | postgres
	DSN        string // postgres 용 연결 문자열
	SQLitePath string // sqlite 용 파일 경로
}

// Load 는 .env 파일(있으면)과 환경변수를 읽어 Config 를 구성한다.
func Load() *Config {
	// .env 가 없어도 에러로 취급하지 않는다(운영 환경에서는 실제 환경변수를 사용).
	_ = godotenv.Load()

	return &Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DB: DBConfig{
			Driver:     getEnv("DB_DRIVER", "sqlite"),
			DSN:        getEnv("DB_DSN", ""),
			SQLitePath: getEnv("SQLITE_PATH", "cms.db"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
