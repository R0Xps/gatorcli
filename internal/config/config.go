package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	config := Config{}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return config, nil
	}
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (c Config) SetUser(new_user_name string) error {
	c.Current_user_name = new_user_name
	return write(c)
}

func getConfigFilePath() (string, error) {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFilePath := homeDirPath + "/" + configFileName
	return configFilePath, nil
}

func write(config Config) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	os.WriteFile(configFilePath, jsonBytes, os.ModePerm)
	return nil
}
