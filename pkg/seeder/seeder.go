package seeder

import (
	"errors"
	"github.com/jackc/pgconn"
	"github.com/mohamed-abdelrhman/pack-dispatch/internal/pack-size"
	"gorm.io/gorm"
)

type Definition struct {
	Run func() error
}

type seeder struct {
	dbClient *gorm.DB
}

type ISeeder interface {
	SeedPackSizeTable(sizes []*pack_size.PackSize) error
}

func NewSeeder(dbClient *gorm.DB) ISeeder {
	newSeeder := new(seeder)
	newSeeder.dbClient = dbClient
	return newSeeder
}

func (s seeder) SeedPackSizeTable(patch []*pack_size.PackSize) error {
	var count int64
	s.dbClient.Table("pack_sizes").Count(&count)
	if count > 0 {
		return nil
	}
	res := s.dbClient.CreateInBatches(patch, 3)
	if res.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(res.Error, &pgErr) {
			if pgErr.Code != "23505" {
				return res.Error
			}
		}
	}
	return nil
}
