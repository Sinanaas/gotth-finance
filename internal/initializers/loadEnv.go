package initializers

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	SessionSecretKey       string        `mapstructure:"SESSION_SECRET_KEY"`

	// Hosting / security flags (default false when absent).
	AllowRegistration bool `mapstructure:"ALLOW_REGISTRATION"`
	CookieSecure      bool `mapstructure:"COOKIE_SECURE"`
}

// Validate fails fast on missing required configuration so the server never
// boots half-configured (e.g. no session secret or signing keys).
func (c Config) Validate() error {
	required := map[string]string{
		"POSTGRES_HOST":             c.DBHost,
		"POSTGRES_USER":             c.DBUserName,
		"POSTGRES_PASSWORD":         c.DBUserPassword,
		"POSTGRES_DB":               c.DBName,
		"POSTGRES_PORT":             c.DBPort,
		"PORT":                      c.ServerPort,
		"SESSION_SECRET_KEY":        c.SessionSecretKey,
		"ACCESS_TOKEN_PRIVATE_KEY":  c.AccessTokenPrivateKey,
		"ACCESS_TOKEN_PUBLIC_KEY":   c.AccessTokenPublicKey,
		"REFRESH_TOKEN_PRIVATE_KEY": c.RefreshTokenPrivateKey,
		"REFRESH_TOKEN_PUBLIC_KEY":  c.RefreshTokenPublicKey,
	}
	var missing []string
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}
	return nil
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			config = parseConfigFromEnv(config)
		}
		err = nil
	}

	err = viper.Unmarshal(&config)
	return
}

func parseConfigFromEnv(config Config) Config {
	r := reflect.TypeOf(config)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	for i := 0; i < r.NumField(); i++ {
		env := r.Field(i).Tag.Get("mapstructure")
		if err := viper.BindEnv(env); err != nil {
			log.Fatal("? Failed to bind env variable:", err)
		}
	}
	return config
}
