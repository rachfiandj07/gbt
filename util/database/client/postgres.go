package client

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres
	"github.com/pkg/errors"

	"github.com/golang-base-template/util/config"
)

var (
	dbConn map[string]dbReplication
)

type (
	dbReplication struct {
		Master *sqlx.DB
		Slave  *sqlx.DB

		FailOverMaster *sqlx.DB
		FailOverSlave  *sqlx.DB

		RerouteMaster *sqlx.DB
		RerouteSlave  *sqlx.DB
	}

	DatabaseList struct {
		CriticalDatabase    []string
		NonCriticalDatabase []string
	}
)

func InitDB(Databases []string, readOnly bool) (err error) {
	if dbConn == nil || len(dbConn) == 0 {
		dbConn = make(map[string]dbReplication)
	}
	cfg := config.Get()

	// loop each config into db connection
	for _, name := range Databases {
		conns, exist := cfg.Database[name]
		if !exist {
			continue
		}

		db := dbReplication{}

		// Connect to master
		masterMaxConn := conns.MasterMaxConn
		if masterMaxConn == 0 {
			masterMaxConn = 30
		}
		masterMaxIdle := conns.MasterMaxIdle
		if masterMaxIdle == 0 {
			masterMaxIdle = 10
		}

		masterHost := conns.MasterHost
		if readOnly {
			masterHost = conns.SlaveHost
			log.Println("[Database] Readonly has been activated, Connection master will use slave instead")
		}

		connMaster := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", conns.User, conns.Password, conns.DbName, masterHost, conns.PortMaster)
		db.Master, err = openConnection(connMaster, masterMaxConn, masterMaxIdle)
		if err != nil {
			return fmt.Errorf("error open connection db master: %s", err.Error())
		}

		// Connect to slave
		slaveMaxConn := conns.SlaveMaxConn
		if slaveMaxConn == 0 {
			slaveMaxConn = 30
		}
		slaveMaxIdle := conns.SlaveMaxIdle
		if slaveMaxIdle == 0 {
			slaveMaxIdle = 10
		}

		v := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", conns.User, conns.Password, conns.DbName, conns.SlaveHost, conns.PortSlave)
		db.Slave, err = openConnection(v, slaveMaxConn, slaveMaxIdle)
		if err != nil {
			return fmt.Errorf("error open connection db slave: %s", err.Error())
		}

		dbConn[name] = db
	}

	return nil
}

func openConnection(connString string, maxConn, maxIdle int) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return db, errors.Wrapf(err, "connect to %v", connString)
	}
	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(9 * time.Second)

	return db, err
}
