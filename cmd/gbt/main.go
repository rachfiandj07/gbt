package main

import (
	"fmt"
	"github.com/golang-base-template/cmd"
	redisClient "github.com/golang-base-template/util/cache/client"
	"github.com/golang-base-template/util/config"
	databaseClient "github.com/golang-base-template/util/database/client"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"log"

	gbthttp "github.com/golang-base-template/pkg/http"

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

	router := httprouter.New()
	gbthttp.Init()
	gbthttp.AssignRoutes(router)

	n := negroni.New()
	n.UseHandler(router)

	err = gbtserve.Serve(fmt.Sprintf(":%s", cfg.Port.GBT), n)
	if err != nil {
		log.Println("error when serve http app")
		return
	}
}
