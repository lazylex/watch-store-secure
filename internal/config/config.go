package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const (
	EnvironmentLocal      = "local"
	EnvironmentDebug      = "debug"
	EnvironmentProduction = "production"
)

type Config struct {
	Instance          string `yaml:"instance" env:"INSTANCE" env-required:"true"`
	Env               string `yaml:"env" env:"ENV" env-required:"true"`
	UseKafka          bool   `yaml:"use_kafka" env:"USE_KAFKA"`
	Redis             `yaml:"redis"`
	HttpServer        `yaml:"http_server"`
	PersistentStorage `yaml:"persistent_storage"`
	Secure            `yaml:"secure"`
	Kafka             `yaml:"kafka"`
	Prometheus        `yaml:"prometheus"`
}

type Secure struct {
	RootLogin    string `yaml:"root_login" env:"ROOT_LOGIN" env-required:"true"`
	RootPassword string `yaml:"root_password" env:"ROOT_PWD" env-required:"true"`
	Salt         string `yaml:"salt" env:"SALT" env-required:"true"`
}

type HttpServer struct {
	Address         string        `yaml:"address" env:"ADDRESS" env-required:"true"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-required:"true"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-required:"true"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-required:"true"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-required:"true"`
}

type PersistentStorage struct {
	DatabaseLogin              string `yaml:"database_login" env:"DATABASE_LOGIN" env-required:"true"`
	DatabasePassword           string `yaml:"database_password" env:"DATABASE_PASSWORD" env-required:"true"`
	DatabaseAddress            string `yaml:"database_address" env:"DATABASE_ADDRESS" env-required:"true"`
	DatabasePort               int    `yaml:"database_port" env:"DATABASE_PORT" env-required:"true"`
	DatabaseName               string `yaml:"database_name" env:"DATABASE_NAME" env-required:"true"`
	DatabaseMaxOpenConnections int    `yaml:"database_max_open_connections" env:"DATABASE_MAX_OPEN_CONNECTIONS" env-required:"true"`

	QueryTimeout time.Duration `yaml:"query_timeout" env:"QUERY_TIMEOUT" env-required:"true"`
}

type Kafka struct {
	Brokers                []string `yaml:"kafka_brokers" env:"KAFKA_BROKERS"`
	NeedToUpdateTokenTopic string   `yaml:"kafka_topic_need_update_token" env:"KAFKA_TOPIC_NEED_UPDATE_TOKEN"`
}

type Prometheus struct {
	PrometheusPort       string `yaml:"prometheus_port" env:"PROMETHEUS_PORT"`
	PrometheusMetricsURL string `yaml:"prometheus_metrics_url" env:"PROMETHEUS_METRICS_URL"`
}

type Redis struct {
	RedisAddress  string `yaml:"redis_address" env:"REDIS_ADDRESS" env-required:"true"`
	RedisUser     string `yaml:"redis_user" env:"REDIS_USER"`
	RedisPassword string `yaml:"redis_password" env:"REDIS_PWD"`
	RedisDB       int    `yaml:"redis_db" env:"REDIS_DB"`
}

// MustLoad возвращает конфигурацию, считанную из файла, путь к которому передан из командной строки по флагу config или
// содержится в переменной окружения SECURE_CONFIG_PATH. Переопределение конфигурационных значений может находится в
// соответсвующих переменных окружения (описанных в структурах данных в этом файле)
func MustLoad() *Config {
	flag.Parse()

	var configPath = flag.String("config", "", "путь к файлу конфигурации")
	var cfg Config

	if *configPath == "" {
		*configPath = os.Getenv("SECURE_CONFIG_PATH")
	}
	if *configPath == "" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", *configPath)
	}

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
