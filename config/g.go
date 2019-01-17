package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	projectName := "plumber"
	getConfig(projectName)
}

/*
	Get config
 */
func getConfig(projectName string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", projectName))
	viper.AddConfigPath(fmt.Sprintf("/data/docker/config/%s", projectName))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s, please create config.yml ", err))
	}
}

/*
	Get postgresql connecting string param
 */
func GetPostgreSQLConnectingString() string {
	host := viper.GetString("postgresql.host")
	port := viper.GetInt("postgresql.port")
	usr := viper.GetString("postgresql.user")
	pwd := viper.GetString("postgresql.password")
	db := viper.GetString("postgresql.db")
	ssl := viper.GetString("postgresql.sslmode")

	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, usr, db, pwd, ssl)
}
