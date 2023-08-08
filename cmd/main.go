package main

import (
	"github.com/Baldislayer/t-bmstu"
	"github.com/Baldislayer/t-bmstu/pkg/handler"
	"github.com/Baldislayer/t-bmstu/pkg/repository"
	"github.com/Baldislayer/t-bmstu/pkg/testsystems"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatalf("Error occured while reading config: %s", err.Error())
	}

	err := repository.CreateTables()

	if err != nil {
		log.Fatalf("Error occured while creating tables: %s", err.Error())
	}

	// normal code
	handlers := new(handler.Handler)

	srv := new(t_bmstu.Server)

	// запуск горутин проверки задач
	go testsystems.InitGorutines()

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
