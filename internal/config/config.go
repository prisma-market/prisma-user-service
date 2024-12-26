package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort     string   `mapstructure:"SERVER_PORT"`
	MongoURI       string   `mapstructure:"MONGO_URI"`
	JWTSecret      string   `mapstructure:"JWT_SECRET"`       // Auth Service와 동일한 시크릿 사용
	AuthServiceURL string   `mapstructure:"AUTH_SERVICE_URL"` // Auth Service 연동용
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`  // CORS 허용 도메인
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// 기본값 설정
	viper.SetDefault("SERVER_PORT", "8002")
	viper.SetDefault("AUTH_SERVICE_URL", "http://auth-service:8001")

	if err := viper.ReadInConfig(); err != nil {
		// .env 파일이 없어도 환경변수로 실행 가능하게
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
