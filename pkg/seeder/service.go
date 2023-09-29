package seeder

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/mohamed-abdelrhman/pack-dispatch/internal/pack-size"
	"gorm.io/gorm"
)

const (
	SeedPackSizeTable = "SeedPackSizeTable"
)

type service struct {
}

type IService interface {
	Seeds() map[string]Definition
	Run()
}

var (
	seeders ISeeder
)

func NewService(dbClient *gorm.DB) IService {
	seeders = NewSeeder(dbClient)
	newService := new(service)
	return newService
}

func (s service) Seeds() map[string]Definition {
	return map[string]Definition{
		SeedPackSizeTable: {
			Run: func() error {
				pack250 := pack_size.PackSize{
					ID:   uuid.NewString(),
					Name: "Pack250",
					Size: 250,
				}
				pack500 := pack_size.PackSize{
					ID:   uuid.NewString(),
					Name: "Pack500",
					Size: 500,
				}
				pack1000 := pack_size.PackSize{
					ID:   uuid.NewString(),
					Name: "Pack1000",
					Size: 1000,
				}
				pack2000 := pack_size.PackSize{
					ID:   uuid.NewString(),
					Name: "Pack2000",
					Size: 2000,
				}

				pack5000 := pack_size.PackSize{
					ID:   uuid.NewString(),
					Name: "Pack5000",
					Size: 5000,
				}

				patch := []*pack_size.PackSize{&pack250, &pack500, &pack1000, &pack2000, &pack5000}
				return seeders.SeedPackSizeTable(patch)
			},
		},
	}
}

func (s service) Run() {
	err := s.Seeds()[SeedPackSizeTable].Run()
	if err != nil {
		go log.Error("SeedPackSizeTable", err.Error())
	}
}
