package migration

import (
	pack_size "github.com/mohamed-abdelrhman/pack-dispatch/internal/pack-size"
	"gorm.io/gorm"
	"log"
)

type Definition struct {
	Name string
	Run  func() error
}

type migration struct {
	dbClient *gorm.DB
}

type IMigration interface {
	CreatePackSizeTable() error
}

func NewMigration(dbClient *gorm.DB) IMigration {
	newMigration := new(migration)
	newMigration.dbClient = dbClient
	return newMigration
}

func (m migration) CreatePackSizeTable() error {
	err := m.dbClient.Migrator().AutoMigrate(pack_size.PackSize{})
	if err != nil {
		go log.Println("CreatePackSizeTable", err.Error())
		return err
	}
	log.Println("Success CreatePackSizeTable")
	return nil
}
