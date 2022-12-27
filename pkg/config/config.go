package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	TelegramToken  string
	PostgresLogin  string
	PostgresPass   string
	PostgresDBName string
	AdminChatId    string

	Messages Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start       string `mapstructure:"start"`
	StartGifts  string `mapstructure:"start_gifts"`
	AddWish     string `mapstructure:"add_wish"`
	EmptyWishes string `mapstructure:"empty_wishes"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	os.Setenv("TOKEN", "5863235609:AAGcXeB3yTXJozfbauwg_13hYLmwVp7Y9As")
	os.Setenv("POSTGRES_LOGIN", "secretsantabot")
	os.Setenv("POSTGRES_PASS", "q125v450z345")
	os.Setenv("POSTGRES_DB_NAME", "newYearBot")
	os.Setenv("ADMIN_CHAT_ID", "398382229")
	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("postgres_login"); err != nil {
		return err
	}

	if err := viper.BindEnv("postgres_pass"); err != nil {
		return err
	}

	if err := viper.BindEnv("postgres_db_name"); err != nil {
		return err
	}

	if err := viper.BindEnv("admin_chat_id "); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("token")
	cfg.PostgresLogin = viper.GetString("postgres_login")
	cfg.PostgresPass = viper.GetString("postgres_pass")
	cfg.PostgresDBName = viper.GetString("postgres_db_name")
	cfg.AdminChatId = viper.GetString("admin_chat_id")

	return nil
}
