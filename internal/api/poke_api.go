package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NeriusZar/pokedexcli/internal/models"
)

type AreaResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const pokeApiBaseUrl = "https://pokeapi.co/api/v2"

func RetrieveAreas(url string) ([]models.Area, models.Pagination, error) {
	if url == "" {
		areasPath := "/location-area"
		url = pokeApiBaseUrl + areasPath
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []models.Area{}, models.Pagination{}, err
	}

	client := http.Client{}
	res, err := client.Do(req)
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

	pagination := models.Pagination{
		Next:     areaResponse.Next,
		Previous: areaResponse.Previous,
	}

	return mapAreasFromResponse(areaResponse), pagination, nil
}

func mapAreasFromResponse(res AreaResponse) []models.Area {
	areas := make([]models.Area, len(res.Results))

	for i, r := range res.Results {
		areas[i] = models.Area{
			Name: r.Name,
			Url:  r.URL,
		}
	}

	return areas
}
