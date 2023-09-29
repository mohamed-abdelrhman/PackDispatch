package pack_size

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	pack_size "github.com/mohamed-abdelrhman/pack-dispatch/internal/pack-size"
	restErrors "github.com/mohamed-abdelrhman/pack-dispatch/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type packSizeServiceMock struct{}

var CalculatePacksFunc func(qty uint) ([]pack_size.PackCountResponseDto, restErrors.IRestErr)

func (p packSizeServiceMock) CalculatePacks(qty uint) ([]pack_size.PackCountResponseDto, restErrors.IRestErr) {
	return CalculatePacksFunc(qty)
}

func newFiberCtx(dto interface{}, method func(c *fiber.Ctx) error, locals map[string]interface{}) ([]byte, *http.Response) {
	app := fiber.New()
	app.Post("/test/", func(c *fiber.Ctx) error {
		for key, element := range locals {
			c.Locals(key, element)
		}
		return method(c)
	})

	marshaledDto, err := json.Marshal(dto)
	if err != nil {
		panic(err.Error())
	}

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(marshaledDto))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	return body, resp
}

func TestMain(m *testing.M) {
	packSizeService = &packSizeServiceMock{}
	code := m.Run()

	os.Exit(code)
}

func TestPacks(t *testing.T) {
	var dto = map[string]int{
		"quantity": 100,
	}

	t.Run("should throw bad request error", func(t *testing.T) {
		body, resp := newFiberCtx("", Packs, nil)
		var result restErrors.RestErr
		err := json.Unmarshal(body, &result)
		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
		assert.EqualValues(t, "invalid request body", result.Message)
	})

	t.Run("should throw error if service throws", func(t *testing.T) {
		CalculatePacksFunc = func(qty uint) ([]pack_size.PackCountResponseDto, restErrors.IRestErr) {
			return []pack_size.PackCountResponseDto{{}}, restErrors.NewInternalServerError("something went wrong")
		}

		body, resp := newFiberCtx(dto, Packs, nil)
		var result restErrors.RestErr
		err := json.Unmarshal(body, &result)
		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
		assert.EqualValues(t, "something went wrong", result.Message)
	})
	t.Run("should pass", func(t *testing.T) {
		CalculatePacksFunc = func(qty uint) ([]pack_size.PackCountResponseDto, restErrors.IRestErr) {
			return []pack_size.PackCountResponseDto{{}}, nil
		}
		_, resp := newFiberCtx(dto, Packs, nil)
		assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	})
}
