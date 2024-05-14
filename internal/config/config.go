/*
Package config.
Пакет содержит объявления всех структур, содержащих конфигурацию компонентов приложения. Заполнение данными этих
структур производится вызовом метода MustLoad, который производит чтение файла конфигурации. Путь к файлу указывается
в командной строке по ключу 'config' или считывается из переменной окружения 'SECURE_CONFIG_PATH'

# Структуры конфигурации

Config - структура, содержащая все остальные конфигурации и поля Instance (название экземпляра приложения), Env
(уровень запуска приложения EnvironmentLocal, EnvironmentDebug или EnvironmentProduction), UseKafka (логическое
значение, указывающее, нужно ли использовать Кафку)

Redis - конфигурация redis-сервера

HttpServer - конфигурация http-сервера

PersistentStorage - настройки реляционной СУБД, используемой в качестве постоянного хранилища

Kafka - конфигурация для работы с Apache Кafka

Prometheus - конфигурация

TTL - настройки времени жизни сессий и прочих хранящихся в памяти данных

Secure - настройки времени жизни и длины токена, стоимости создания хэша пароля
*/
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
	Kafka             `yaml:"kafka"`
	Prometheus        `yaml:"prometheus"`
	TTL               `yaml:"ttl"`
	Secure            `yaml:"secure"`
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
	DatabaseSchema             string `yaml:"database_schema" env:"DATABASE_SCHEMA"`
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

type TTL struct {
	SessionTTL               time.Duration `yaml:"session_ttl" env:"TTL_SESSION_TTL" env-required:"true"`
	UserIdAndPasswordHashTTL time.Duration `yaml:"user_id_and_password_hash_ttl" env:"TTL_USER_ID_AND_PASSWORD_HASH_TTL" env-required:"true"`
	AccountStateTTL          time.Duration `yaml:"account_state_ttl" env:"TTL_ACCOUNT_STATE_TTL" env-required:"true"`
	PermissionsNumbersTTL    time.Duration `yaml:"permissions_numbers_ttl" env:"TTL_PERMISSIONS_TTL" env-required:"true"`
	InstanceDataTTL          time.Duration `yaml:"instance_data_ttl" env:"INSTANCE_DATA_TTL" env-required:"true"`
}

type Secure struct {
	LoginTokenLength     int           `yaml:"login_token_length" env:"LOGIN_TOKEN_LENGTH" env-required:"true"`
	PasswordCreationCost int           `yaml:"password_creation_cost" env:"PASSWORD_CREATION_COST" env-required:"true"`
	TokenTTL             time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-required:"true"`
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
