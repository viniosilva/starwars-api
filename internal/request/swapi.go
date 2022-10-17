package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/viniosilva/starwars-api/internal/model"
)

//go:generate mockgen -destination=../../mock/swapi_request_mock.go -package=mock . SwapiRequest
type SwapiRequest interface {
	GetPlanets(ctx context.Context, page int) (*model.SwapiPlanetsResponse, error)
	GetFilms(ctx context.Context, page int) (*model.SwapiFilmsResponse, error)
}

const SWAPI_URL = "https://swapi.dev/api"

type ISwapiRequest struct{}

func (impl *ISwapiRequest) GetPlanets(ctx context.Context, page int) (*model.SwapiPlanetsResponse, error) {
	body, err := impl.GetPath("/planets", page)
	if err != nil {
		return nil, err
	}

	var swapiRes model.SwapiPlanetsResponse
	err = json.Unmarshal(body, &swapiRes)
	if err != nil {
		return nil, err
	}

	return &swapiRes, nil
}

func (impl *ISwapiRequest) GetFilms(ctx context.Context, page int) (*model.SwapiFilmsResponse, error) {
	body, err := impl.GetPath("/films", page)
	if err != nil {
		return nil, err
	}

	var swapiRes model.SwapiFilmsResponse
	err = json.Unmarshal(body, &swapiRes)
	if err != nil {
		return nil, err
	}

	return &swapiRes, nil
}

func (impl *ISwapiRequest) GetPath(path string, page int) ([]byte, error) {
	paramPage := ""
	if page > 1 {
		paramPage = fmt.Sprintf("?page=%d", page)
	}

	url := fmt.Sprintf("%s/%s%s", SWAPI_URL, path, paramPage)
	httpRes, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer httpRes.Body.Close()
	body, err := ioutil.ReadAll(httpRes.Body)

	return body, err
}
