package cmd

import (
	"fmt"
	"log"
	"regexp"

	"github.com/pkg/errors"

	redisClient "github.com/golang-base-template/util/cache/client"
	"github.com/golang-base-template/util/config"
	databaseClient "github.com/golang-base-template/util/database/client"
)

var (
	passwordRegex = regexp.MustCompile(`password=[^ ]*`)
)

// InitApp is func to init database connection and other external tools
func InitApp(appConfig config.Config, dbs databaseClient.DatabaseList, rds redisClient.RedisList) (err error) {
	//Init DB, fatal if critical db error but let app run if non-critical db error
	if len(dbs.CriticalDatabase) != 0 {
		log.Println("Initiating Critical DB:", dbs.CriticalDatabase)
		err = databaseClient.InitDB(dbs.CriticalDatabase, false)
		if err != nil {
			return handleErr("DB:", err)
		}
	}
	if len(dbs.NonCriticalDatabase) != 0 {
		log.Println("Initiating Non Critical DB:", dbs.NonCriticalDatabase)
		err = databaseClient.InitDB(dbs.NonCriticalDatabase, false)
		if err != nil {
			log.Println("[Init] Error when init", handleErr("nonCriticalDatabases", err))
		}
	}

	//Init Redis, fatal if critical redis error but let app run if non-critical redis error
	if len(rds.CriticalRedis) != 0 {
		log.Println("Initiating Critical Redis:", rds.CriticalRedis)
		err = redisClient.InitRedis(rds.CriticalRedis)
		if err != nil {
			return handleErr("Redis:", err)
		}
	}
	if len(rds.NonCriticalRedis) != 0 {
		log.Println("Initiating Non Critical Redis:", rds.CriticalRedis)
		err = redisClient.InitRedis(rds.CriticalRedis)
		if err != nil {
			log.Println("[Init] Error when init", handleErr("nonCriticalRedis", err))
		}
	}

	return nil
}

// handleErr func to mask db password in log
func handleErr(label string, e error) (err error) {
	if e != nil {
		err = fmt.Errorf(passwordRegex.ReplaceAllString(e.Error(), "password=*****"))
		err = errors.Wrapf(err, "%v ", label)
	}
	return
}
