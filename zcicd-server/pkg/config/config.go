package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	NATS     NATSConfig     `mapstructure:"nats"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Casbin   CasbinConfig   `mapstructure:"casbin"`
	Log      LogConfig      `mapstructure:"log"`
	Crypto   CryptoConfig   `mapstructure:"crypto"`
}

type ServerConfig struct {
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host            string          `mapstructure:"host"`
	Port            int             `mapstructure:"port"`
	User            string          `mapstructure:"user"`
	Password        string          `mapstructure:"password"`
	DBName          string          `mapstructure:"dbname"`
	SSLMode         string          `mapstructure:"sslmode"`
	MaxOpenConns    int             `mapstructure:"max_open_conns"`
	MaxIdleConns    int             `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int             `mapstructure:"conn_max_lifetime"`
	Replicas        []ReplicaConfig `mapstructure:"replicas"`
}

type ReplicaConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type NATSConfig struct {
	URL           string `mapstructure:"url"`
	StreamName    string `mapstructure:"stream_name"`
	MaxReconnects int    `mapstructure:"max_reconnects"`
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL    bool   `mapstructure:"use_ssl"`
	Bucket    string `mapstructure:"bucket"`
}

type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	ExpireHours        int    `mapstructure:"expire_hours"`
	RefreshExpireHours int    `mapstructure:"refresh_expire_hours"`
	Issuer             string `mapstructure:"issuer"`
}

type CasbinConfig struct {
	ModelPath string `mapstructure:"model_path"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type CryptoConfig struct {
	AESKey string `mapstructure:"aes_key"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
