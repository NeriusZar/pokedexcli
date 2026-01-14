package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NeriusZar/pokedexcli/internal/models"
	"github.com/NeriusZar/pokedexcli/internal/pokecache"
)

const pokeApiBaseUrl = "https://pokeapi.co/api/v2"
const cacheInterval = time.Second * 5

type AreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokeApi struct {
	cache  pokecache.Cache
	client http.Client
}

func NewPokeApi() PokeApi {
	return PokeApi{
		cache:  pokecache.NewCache(cacheInterval),
		client: http.Client{},
	}
}

func (api *PokeApi) RetrieveAreas(pageUrl *string) ([]models.Area, models.Pagination, error) {
	areasPath := "/location-area"
	url := pokeApiBaseUrl + areasPath
	if pageUrl != nil {
		url = *pageUrl
	}

	if entry, ok := api.cache.Get(url); ok {
		var areaResponse AreaResponse
		if err := json.Unmarshal(entry, &areaResponse); err == nil {
			areas, pagination := mapResponse(areaResponse)
			return areas, pagination, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []models.Area{}, models.Pagination{}, err
	}

	res, err := api.client.Do(req)
	if err != nil {
		return []models.Area{}, models.Pagination{}, err
	}
	if res.StatusCode != http.StatusOK {
		return []models.Area{}, models.Pagination{}, fmt.Errorf("Failed to fetch areas. Status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	var areaResponse AreaResponse
	if err := json.NewDecoder(res.Body).Decode(&areaResponse); err != nil {
		return []models.Area{}, models.Pagination{}, err
	}

	entry, err := json.Marshal(areaResponse)
	if err != nil {
		return []models.Area{}, models.Pagination{}, err
	}
	api.cache.Add(url, entry)

	areas, pagination := mapResponse(areaResponse)
	return areas, pagination, nil
}

func mapResponse(res AreaResponse) ([]models.Area, models.Pagination) {
	areas := make([]models.Area, len(res.Results))

	for i, r := range res.Results {
		areas[i] = models.Area{
			Name: r.Name,
			Url:  r.URL,
		}
	}

	pagination := models.Pagination{
		Next:     res.Next,
		Previous: res.Previous,
	}

	return areas, pagination
}
