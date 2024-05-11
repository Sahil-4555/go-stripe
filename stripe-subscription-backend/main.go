package main

import (
	"context"
	"stripe-subscription/configs"
	"stripe-subscription/routes"
	"stripe-subscription/shared/log"
	"stripe-subscription/shared/utils"
)

func main() {
	log.Init()
	configs.InitDB()
	log.GetLog().Info("", "DB connected")

	configs.MigrateModels()
	log.GetLog().Info("", "Successfully migrate the models")

	go routes.Run()

	utils.GracefulStop(log.GetLog(), func(ctx context.Context) error {
		var err error
		if err = routes.Close(ctx); err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			return err
		}
		if err = configs.Close(); err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			return err
		}
		return nil
	})

}
