package pack_size

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	restErrors "github.com/mohamed-abdelrhman/pack-dispatch/pkg/errors"
	"github.com/mohamed-abdelrhman/pack-dispatch/pkg/sqlclient"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type IRepository interface {
	WithTransaction(txHandle *gorm.DB) IRepository
	WithoutTransaction() IRepository
	OrderMany(field string, direction string) ([]PackSize, restErrors.IRestErr)
}

func NewRepository() IRepository {
	newRepo := repository{}
	newRepo.db = sqlclient.OpenDBConnection()
	return newRepo
}

func (r repository) WithTransaction(txHandle *gorm.DB) IRepository {
	r.db = txHandle
	return r
}
func (r repository) WithoutTransaction() IRepository {
	r.db = sqlclient.OpenDBConnection()
	return r
}

func (r repository) OrderMany(field string, direction string) ([]PackSize, restErrors.IRestErr) {
	var records []PackSize
	result := r.db.Order(fmt.Sprintf("%s %s", field, direction)).Find(&records)
	if result.Error != nil {
		go log.Error("OrderMany", result.Error)
		return nil, restErrors.NewInternalServerError("something went wrong")
	}
	return records, nil
}
