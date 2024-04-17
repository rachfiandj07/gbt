package main

import (
	"fmt"
	"log"

	"github.com/urfave/negroni"

	"github.com/golang-base-template/cmd"
	gbtconsumer "github.com/golang-base-template/pkg/consumer"
	redisClient "github.com/golang-base-template/util/cache/client"
	"github.com/golang-base-template/util/config"
	databaseClient "github.com/golang-base-template/util/database/client"
	gbtserve "github.com/golang-base-template/util/serve"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		msg := fmt.Sprintf("error when init http config: %+v", err)
		log.Fatalln(msg)
	}

	cfg := config.Get()

	err = cmd.InitApp(cfg,
		databaseClient.DatabaseList{
			CriticalDatabase:    []string{},
			NonCriticalDatabase: []string{},
		},
		redisClient.RedisList{
			CriticalRedis:    []string{},
			NonCriticalRedis: []string{},
		})
	if err != nil {
		msg := fmt.Sprintf("error when init http app: %+v", err)
		log.Fatalln(msg)
	}

	gbtconsumer.Init(&cfg)

	n := negroni.New()

	err = gbtserve.Serve(fmt.Sprintf(":%s", cfg.Port.Bg), n)
	if err != nil {
		log.Println("error when serve http app")
		return
	}
}
