package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	configFile = "config.toml"
	Log        = "log"
	Params     = "params"
	Template   = "template"
	paramFile  = "params.yaml"
)

var HomeDirectory string

func init() {
	home, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	HomeDirectory = home
}

type Dir struct {
	Log      string `json:"log"`
	Params   string `json:"params"`
	Template string `json:"template"`
}

type Database struct {
	IP          string `json:"ip"`
	ServiceName string `json:"service_name"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	Port        string `json:"port"`
}

//Check existence of the configuration file
func IsConfigFileExist() bool {
	if fi, err := os.Stat(path.Join(HomeDirectory, configFile)); err != nil || fi.IsDir() {
		return false
	}
	return true
}

//Read the configuration
func ReadConfigFile() (*Dir, *Database, error) {
	viper.SetConfigFile(path.Join(HomeDirectory, configFile))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, err
	}

	return &Dir{
			Log:      viper.GetString("dir.log"),
			Params:   viper.GetString("dir.params"),
			Template: viper.GetString("dir.template"),
		}, &Database{
			IP:          viper.GetString("database.ip"),
			Port:        viper.GetString("database.port"),
			ServiceName: viper.GetString("database.service_name"),
			UserName:    viper.GetString("database.user_name"),
			Password:    viper.GetString("database.password"),
		},
		nil
}

//Check valid data of the configuration after read successfully
func (cfg *Dir) IsValidConfig() error {
	if err := IsValidLog(cfg.Log); err != nil {
		return err
	}
	if err := IsValidLog(cfg.Params); err != nil {
		return err
	}
	if err := IsValidLog(cfg.Template); err != nil {
		return err
	}
	return nil
}

func IsValidLog(fp string) error {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return nil
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return nil
	}

	return fmt.Errorf("this %s is not a path", fp)
}

//Make download directory if it not exist
func (cfg *Dir) CreateFolderIfNotExist(folder string) {
	if _, err := os.Stat(folder); err != nil {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func (cfg *Dir) CheckFolder() {
	cfg.CreateFolderIfNotExist(cfg.Log)
	cfg.CreateFolderIfNotExist(cfg.Params)
	cfg.CreateFolderIfNotExist(cfg.Template)
}

//Write data of config to the configuration file
func WriteConfigFile(cfg *Dir, db *Database) error {
	_, err := os.Create(path.Join(HomeDirectory, configFile))
	if err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(HomeDirectory)
	viper.Set("dir.log", cfg.Log)
	viper.Set("dir.params", cfg.Params)
	viper.Set("dir.template", cfg.Template)
	viper.Set("database.ip", db.IP)
	viper.Set("database.port", db.Port)
	viper.Set("database.service_name", db.ServiceName)
	viper.Set("database.user_name", db.UserName)
	viper.Set("database.password", db.Password)
	return viper.WriteConfig()
}

//Write default config into the configuration file
func WriteDefaultConfig() error {
	return WriteConfigFile(&Dir{
		Log:      Log,
		Params:   Params,
		Template: Template,
	}, &Database{
		IP:          "",
		Port:        "1521",
		ServiceName: "",
		UserName:    "",
		Password:    "",
	})
}
