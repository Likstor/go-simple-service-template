package config

import (
	"context"
	"fmt"
	"os"
	"service/internal/pkg/logs"
	"strconv"
	"sync"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Server struct {
		Host              string
		Port              string `validate:"number,gte=0,lte=65535"`
		Domain            string
		ReadTimeout       time.Duration
		ReadHeaderTimeout time.Duration
		WriteTimeout      time.Duration
		IdleTimeout       time.Duration
	}
	Storage struct {
		Host     string
		Port     string `validate:"number,gte=0,lte=65535"`
		Name     string
		User     string
		Password string
	}
	TLS struct {
		Certificate string `validate:"filepath"` // Path to tls certificate
		Key         string `validate:"filepath"` // Path to tls key
	}
}

const opCfg = "config.GetConfig"

func defError(ctx context.Context, key string, def any) {
	logs.Warn(
		ctx,
		fmt.Sprintf("An error occurred or the %v environment variable was not found. The default value is %v", key, def),
		opCfg,
	)
}

func getEnvAsString(ctx context.Context, key, def string) string {
	envVal, ok := os.LookupEnv(key)
	if !ok {
		defError(ctx, key, def)
		return def
	}

	return envVal
}

func getEnvAsUint64(ctx context.Context, key string, def uint64) uint64 {
	envVal, ok := os.LookupEnv(key)
	if !ok {
		defError(ctx, key, def)
		return def
	}

	val, err := strconv.ParseUint(envVal, 10, 64)
	if err != nil {
		defError(ctx, key, def)
		return def
	}

	return val
}

func getEnvAsTime(ctx context.Context, key string, def time.Duration) time.Duration {
	envVal, ok := os.LookupEnv(key)
	if !ok {
		defError(ctx, key, def)
		return def
	}

	pos := 0
	for i, c := range envVal {
		if !unicode.IsDigit(c) {
			pos = i
			break
		}

		if i == len(envVal)-1 {
			defError(ctx, key, def)
			return def
		}
	}

	num, err := strconv.Atoi(envVal[:pos])
	if err != nil {
		defError(ctx, key, def)
		return def
	}

	val := time.Duration(num)
	switch envVal[pos:] {
	case "ms":
		val *= time.Millisecond
	case "s":
		val *= time.Second
	case "m":
		val *= time.Minute
	case "h":
		val *= time.Hour
	default:
		defError(ctx, key, def)
		return def
	}

	return val
}

var instance *Config
var once sync.Once

func GetConfig(ctx context.Context) (*Config, error) {
	var configErr error
	once.Do(func() {
		instance = &Config{}

		serverCfg := &instance.Server
		serverCfg.Host = getEnvAsString(ctx, "BACKEND_HOST", "127.0.0.1")
		serverCfg.Port = getEnvAsString(ctx, "BACKEND_PORT", "8080")
		serverCfg.Domain = getEnvAsString(ctx, "BACKEND_DOMAIN", "localhost.com") // need to be added to hosts
		serverCfg.ReadTimeout = getEnvAsTime(ctx, "BACKEND_READ_TIMEOUT", 250*time.Millisecond)
		serverCfg.ReadHeaderTimeout = getEnvAsTime(ctx, "BACKEND_READ_HEADER_TIMEOUT", 250*time.Millisecond)
		serverCfg.WriteTimeout = getEnvAsTime(ctx, "BACKEND_WRITE_TIMEOUT", 1*time.Second)
		serverCfg.IdleTimeout = getEnvAsTime(ctx, "BACKEND_IDLE_TIMEOUT", 5*time.Minute)

		storageCfg := &instance.Storage
		storageCfg.Host = getEnvAsString(ctx, "STORAGE_HOST", "127.0.0.1")
		storageCfg.Port = getEnvAsString(ctx, "STORAGE_PORT", "5432")
		storageCfg.Name = getEnvAsString(ctx, "STORAGE_NAME", "storage")
		storageCfg.User = getEnvAsString(ctx, "STORAGE_USER", "user")
		storageCfg.Password = getEnvAsString(ctx, "STORAGE_PASSWORD", "pwd")

		tlsCfg := &instance.TLS
		tlsCfg.Certificate = getEnvAsString(ctx, "TLS_CERTIFICATE", "tls.crt")
		tlsCfg.Key = getEnvAsString(ctx, "TLS_KEY", "tls.key")

		// Validate config
		validate := validator.New()
		err := validate.Struct(instance)
		if err != nil {
			configErr = err
		}
	})
	return instance, configErr
}
