package config

import (
	"log"

	"github.com/spf13/viper"
)

var E Env

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	YouthApiKey            string `mapstructure:"YOUTH_API_KEY"`
	WelfareApiKey          string `mapstructure:"WELFARE_API_KEY"`
	VolunteerApiKey        string `mapstructure:"VOLUNTEER_API_KEY"`
	MaternityApiKey        string `mapstructure:"MATERNITY_API_KEY"`
	KakaoMapApiKey         string `mapstructure:"KAKAO_MAP_API_KEY"`
	AIApiKey               string `mapstructure:"AI_API_KEY"`
}

func NewEnv() {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .E : ", err)
	}

	err = viper.Unmarshal(&E)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if E.AppEnv == "development" {
		log.Println("The App is running in development E")
	}
}
