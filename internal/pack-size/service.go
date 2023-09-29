package pack_size

import (
	restErrors "github.com/mohamed-abdelrhman/pack-dispatch/pkg/errors"
	"gorm.io/gorm"
	"sort"
)

type service struct{}

type IService interface {
	WithTransaction(txHandle *gorm.DB) IService
	WithoutTransaction() IService
	CalculatePacks(qty uint) ([]PackCountResponseDto, restErrors.IRestErr)
}

var (
	repo = NewRepository()
)

func NewService() IService {
	newService := &service{}
	return newService
}

func (s service) WithTransaction(txHandle *gorm.DB) IService {
	repo = repo.WithTransaction(txHandle)
	return s
}
func (s service) WithoutTransaction() IService {
	repo = repo.WithoutTransaction()
	return s
}

func (s service) CalculatePacks(qty uint) ([]PackCountResponseDto, restErrors.IRestErr) {
	packSizes, err := repo.OrderMany("size", "desc")
	if err != nil {
		return []PackCountResponseDto{}, err
	}

	packsMap := map[uint]PackCountResponseDto{}

	for _, pack := range packSizes {
		if qty == 0 {
			break // Stop iterating if the order quantity has been fulfilled
		}

		packsCount := qty / pack.Size // Calculate the number of packs needed as an integer
		if packsCount > 0 {
			v, ok := packsMap[pack.Size]
			if !ok {
				packsMap[pack.Size] = PackCountResponseDto{
					Size:  pack.Size,
					Count: packsCount,
				}
			} else {
				v.Count++
				packsMap[pack.Size] = v
			}
			qty -= packsCount * pack.Size // Reduce the remaining quantity by the packs sent
		}
	}

	if qty > 0 {
		// Find the smallest pack size available that is larger than the remaining quantity
		nextPackSize := uint(0)
		for _, pack := range packSizes {
			if pack.Size > qty {
				nextPackSize = pack.Size
			}
		}

		if nextPackSize > 0 {
			v, ok := packsMap[nextPackSize]
			if !ok {
				packsMap[nextPackSize] = PackCountResponseDto{
					Size:  nextPackSize,
					Count: 1,
				}
			} else {
				v.Count++
				packsMap[nextPackSize] = v
			}
		}
	}

	//Consolidate and combine Packs//

	// Sort the pack sizes in ascending order
	sort.Slice(packSizes, func(i, j int) bool { return packSizes[i].Size < packSizes[j].Size })

	for i := 0; i < len(packSizes)-1; i++ {
		smallerPackSize, largerPackSize := packSizes[i].Size, packSizes[i+1].Size

		// If we can combine smaller packs to form a larger pack, do so
		if packsMap[smallerPackSize].Count*smallerPackSize >= largerPackSize {
			additionalLargerPacks := (packsMap[smallerPackSize].Count * smallerPackSize) / largerPackSize
			packsMap[largerPackSize] = PackCountResponseDto{
				Size:  largerPackSize,
				Count: packsMap[largerPackSize].Count + additionalLargerPacks,
			}
			packsMap[smallerPackSize] = PackCountResponseDto{
				Size:  smallerPackSize,
				Count: (packsMap[smallerPackSize].Count * smallerPackSize) % largerPackSize / smallerPackSize,
			}
		}
	}

	// Convert the map back to a slice
	consolidatedPacks := make([]PackCountResponseDto, 0, len(packsMap))
	for _, pack := range packsMap {
		if pack.Count > 0 {
			consolidatedPacks = append(consolidatedPacks, pack)
		}
	}

	return consolidatedPacks, nil
}

