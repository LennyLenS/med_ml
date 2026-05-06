package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Config struct {
	App      App
	Adapters Adapters
	S3       S3
	Dbus     Dbus
	JWT      JWT
}

type App struct {
	Url string `env:"APP_URL" env-default:"localhost:8080"`
	// CORS_ALLOWED_ORIGINS — список origin через запятую (например http://localhost:3000,http://127.0.0.1:49118).
	// Пустое значение: разрешить любой origin (*), удобно при сервисах на одном хосте с разными портами.
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" env-default:""`
}

type Adapters struct {
	UziUrl      string `env:"ADAPTERS_UZIURL" env-required:"true"`
	AuthUrl     string `env:"ADAPTERS_AUTHURL" env-required:"true"`
	MedUrl      string `env:"ADAPTERS_MEDURL" env-required:"true"`
	BillingUrl  string `env:"ADAPTERS_BILLINGURL" env-required:"true"`
	CytologyUrl string `env:"ADAPTERS_CYTOLOGYURL" env-required:"true"`
	TilerUrl    string `env:"ADAPTERS_TILERURL" env-default:"http://localhost:50080"`
}

type S3 struct {
	Endpoint     string `env:"S3_ENDPOINT" env-required:"true"`
	Access_Token string `env:"S3_TOKEN_ACCESS" env-required:"true"`
	Secret_Token string `env:"S3_TOKEN_SECRET" env-required:"true"`
}

type Dbus struct {
	Addrs []string `env:"DBUS_ADDRS" env-required:"true"`
}

type JWT struct {
	RsaPublicKey string `env:"JWT_KEY_PUBLIC" env-required:"true"`
}

func (c *Config) ParseRsaPublicKey() (*rsa.PublicKey, error) {
	publicBlock, _ := pem.Decode([]byte(c.JWT.RsaPublicKey))
	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	return publicKey.(*rsa.PublicKey), nil
}
