package sqlclient

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/config"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var (
	dbConnection *gorm.DB
	clientOnce   sync.Once
	err          error
)

// OpenDBConnection when initiating new repository we should use this function
// we use OpenDBConnection which returns the original dbConnection variable from sqlclient pkg
func OpenDBConnection() *gorm.DB {
	clientOnce.Do(func() {
		dbConnection, err = gorm.Open(postgres.Open(config.Environment.DatabaseServerURL), &gorm.Config{
			Logger: glogger.Default.LogMode(glogger.Error),
		})
		if err != nil {
			go log.Warn("DATABASE_CONNECTION_ERROR", err)
			panic(err)
		}
	})
	return dbConnection
}
