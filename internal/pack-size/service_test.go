package pack_size

import (
	restErrors "github.com/mohamed-abdelrhman/pack-dispatch/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"sort"
	"testing"
)

var (
	packSizeService IService
	OrderManyFunc   func(field string, direction string) ([]PackSize, restErrors.IRestErr)
)

type repoMock struct{}

func (r repoMock) OrderMany(field string, direction string) ([]PackSize, restErrors.IRestErr) {
	return OrderManyFunc(field, direction)
}

func TestMain(m *testing.M) {
	repo = &repoMock{}
	packSizeService = NewService()
	code := m.Run()
	os.Exit(code)
}

func TestService_CalculatePacks(t *testing.T) {
	t.Run("calculate packs should return empty list if qty = 0", func(t *testing.T) {
		res, err := packSizeService.CalculatePacks(0)
		assert.Nil(t, err)
		assert.EqualValues(t, 0, len(res))
	})

	t.Run("calculate packs should return err if repo returns an err", func(t *testing.T) {
		OrderManyFunc = func(field string, direction string) ([]PackSize, restErrors.IRestErr) {
			return nil, restErrors.NewInternalServerError("something went wrong")
		}
		_, err := packSizeService.CalculatePacks(1)
		assert.EqualValues(t, "something went wrong", err.Error())
	})

	t.Run("calculate packs should return empty list if there no packSizes in the db", func(t *testing.T) {
		OrderManyFunc = func(field string, direction string) ([]PackSize, restErrors.IRestErr) {
			return []PackSize{}, nil
		}
		res, err := packSizeService.CalculatePacks(1)
		assert.Nil(t, err)
		assert.EqualValues(t, 0, len(res))
	})

	t.Run("calculate packs should pass", func(t *testing.T) {
		OrderManyFunc = func(field string, direction string) ([]PackSize, restErrors.IRestErr) {
			return []PackSize{
				{Name: "5000", Size: 5000},
				{Name: "2000", Size: 2000},
				{Name: "1000", Size: 1000},
				{Name: "500", Size: 500},
				{Name: "250", Size: 250},
			}, nil
		}
		tests := []struct {
			name string
			qty  uint
			want []PackCountResponseDto
		}{
			{"Test Case 1", 1, []PackCountResponseDto{{Size: 250, Count: 1}}},
			{"Test Case 2", 250, []PackCountResponseDto{{Size: 250, Count: 1}}},
			{"Test Case 3", 251, []PackCountResponseDto{{Size: 500, Count: 1}}},
			{"Test Case 4", 500, []PackCountResponseDto{{Size: 500, Count: 1}}},
			{"Test Case 5", 501, []PackCountResponseDto{{Size: 500, Count: 1}, {Size: 250, Count: 1}}},
			{"Test Case 6", 12001, []PackCountResponseDto{{Size: 250, Count: 1}, {Size: 5000, Count: 2}, {Size: 2000, Count: 1}}},
		}

		sortPacks := func(packs []PackCountResponseDto) []PackCountResponseDto {
			sort.Slice(packs, func(i, j int) bool {
				return packs[i].Size < packs[j].Size
			})
			return packs

		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := packSizeService.CalculatePacks(tt.qty)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(sortPacks(got), sortPacks(tt.want)) {
					t.Errorf("CalculatePacks() = %v, want %v", got, tt.want)
				}
			})
		}

	})

}
