package migration

import (
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

const (
	MigratePackSizeTable = "MigratePackSizeTable"
)

type service struct {
}

type IService interface {
	Migrations() map[string]Definition
	Run()
}

var migrator IMigration

func NewService(dbClient *gorm.DB) IService {
	migrator = NewMigration(dbClient)

	newService := new(service)
	return newService

}

func (service) Migrations() map[string]Definition {
	return map[string]Definition{
		MigratePackSizeTable: {
			Name: MigratePackSizeTable,
			Run: func() error {
				return migrator.CreatePackSizeTable()
			},
		},
	}
}

func (s service) Run() {
	for _, step := range s.Migrations() {
		if err := step.Run(); err != nil {
			go log.Error(step.Name, err.Error())
		}
	}
}
