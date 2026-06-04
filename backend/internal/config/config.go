package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config 는 애플리케이션 전역 설정이다. 환경변수에서 로딩한다.
type Config struct {
	AppEnv   string
	HTTPPort string
	DB       DBConfig
	Admin    AdminConfig
}

// DBConfig 는 데이터베이스 연결 설정이다.
type DBConfig struct {
	Driver     string // sqlite | postgres
	DSN        string // postgres 용 연결 문자열
	SQLitePath string // sqlite 용 파일 경로
}

// AdminConfig 는 관리자(admin) API 설정이다.
type AdminConfig struct {
	// CORSOrigins 는 허용할 Origin 목록이다. ["*"] 이면 전체 허용.
	CORSOrigins []string
	// JWTSecret 은 액세스/리프레시 토큰 서명 키이다(운영에서는 반드시 강력한 값으로 설정).
	JWTSecret string
	// AccessTTL/RefreshTTL 은 토큰 만료 기간이다.
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	// SeedPassword 는 admin_users 가 비어있을 때 생성하는 개발용 기본 계정 비밀번호이다.
	SeedPassword string
}

// devJWTSecret 은 JWT_SECRET 미설정 시 개발용 폴백이다(운영 사용 금지).
const devJWTSecret = "dev-insecure-secret-change-me"

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
		Admin: AdminConfig{
			CORSOrigins:  splitCSV(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5173,http://localhost:5174,http://localhost:9000,http://localhost:3000")),
			JWTSecret:    getEnv("JWT_SECRET", devJWTSecret),
			AccessTTL:    time.Duration(getEnvInt("JWT_ACCESS_TTL_MIN", 60)) * time.Minute,
			RefreshTTL:   time.Duration(getEnvInt("JWT_REFRESH_TTL_HOURS", 168)) * time.Hour,
			SeedPassword: getEnv("ADMIN_SEED_PASSWORD", "secret"),
		},
	}
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// splitCSV 는 콤마 구분 문자열을 트림된 비어있지 않은 요소 슬라이스로 변환한다.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
