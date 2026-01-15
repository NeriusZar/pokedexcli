package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NeriusZar/pokedexcli/internal/models"
	"github.com/NeriusZar/pokedexcli/internal/pokecache"
)

const pokeApiBaseUrl = "https://pokeapi.co/api/v2"
const locationAreasPath = "/location-area"
const pokemonDetailsPath = "/pokemon"
const cacheInterval = time.Second * 5

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
	url := pokeApiBaseUrl + locationAreasPath
	if pageUrl != nil {
		url = *pageUrl
	}

	if entry, ok := api.cache.Get(url); ok {
		var areaResponse AreaResponse
		if err := json.Unmarshal(entry, &areaResponse); err == nil {
			areas, pagination := mapAreasResponse(areaResponse)
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

	areas, pagination := mapAreasResponse(areaResponse)
	return areas, pagination, nil
}

func mapAreasResponse(res AreaResponse) ([]models.Area, models.Pagination) {
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

func (api *PokeApi) RetrievePokemonsInArea(area string) ([]models.PokemonShortInfo, error) {
	url := pokeApiBaseUrl + locationAreasPath + "/" + area

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []models.PokemonShortInfo{}, err
	}

	if entry, ok := api.cache.Get(url); ok {
		var areaDetailsResponse AreaDetailsResponse
		if err := json.Unmarshal(entry, &areaDetailsResponse); err == nil {
			return mapPokemonsResponse(areaDetailsResponse), nil
		}
	}

	res, err := api.client.Do(req)
	if err != nil {
		return []models.PokemonShortInfo{}, err
	}
	if res.StatusCode != http.StatusOK {
		return []models.PokemonShortInfo{}, errors.New("Failed to fetch pokemons in area")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []models.PokemonShortInfo{}, errors.New("Failed to fetch pokemons in area")
	}

	var areaDetailsResponse AreaDetailsResponse
	if err := json.Unmarshal(data, &areaDetailsResponse); err != nil {
		return []models.PokemonShortInfo{}, err
	}

	api.cache.Add(url, data)

	return mapPokemonsResponse(areaDetailsResponse), nil
}

func mapPokemonsResponse(res AreaDetailsResponse) []models.PokemonShortInfo {
	pokemons := make([]models.PokemonShortInfo, len(res.PokemonEncounters))

	for i, p := range res.PokemonEncounters {
		pokemons[i] = models.PokemonShortInfo{
			Name: p.Pokemon.Name,
			Url:  p.Pokemon.URL,
		}
	}

	return pokemons
}

func (api *PokeApi) GetPokemonDetails(name string) (models.Pokemon, error) {
	url := pokeApiBaseUrl + pokemonDetailsPath + "/" + name

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Pokemon{}, err
	}

	if entry, ok := api.cache.Get(url); ok {
		var pokemonDetailsResponse PokemonDetailsResponse
		if err := json.Unmarshal(entry, &pokemonDetailsResponse); err == nil {
			return mapPokemonDetailsResponse(pokemonDetailsResponse), nil
		}
	}

	res, err := api.client.Do(req)
	if err != nil {
		return models.Pokemon{}, err
	}
	if res.StatusCode != http.StatusOK {
		return models.Pokemon{}, errors.New("Failed to fetch pokemon details")
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Pokemon{}, errors.New("Failed to fetch pokemon details")
	}

	var pokemonDetailsResponse PokemonDetailsResponse
	if err := json.Unmarshal(data, &pokemonDetailsResponse); err != nil {
		return models.Pokemon{}, err
	}

	api.cache.Add(url, data)

	return mapPokemonDetailsResponse(pokemonDetailsResponse), nil
}

func mapPokemonDetailsResponse(res PokemonDetailsResponse) models.Pokemon {
	pokemon := models.Pokemon{
		ID:             res.ID,
		Name:           res.Name,
		BaseExperience: res.BaseExperience,
		Weight:         res.Weight,
		Height:         res.Height,
	}

	stats := make([]models.PokemonStat, len(res.Stats))
	for i, s := range res.Stats {
		stats[i] = models.PokemonStat{
			Name:     s.Stat.Name,
			BaseStat: s.BaseStat,
		}
	}

	types := make([]string, len(res.Types))
	for i, t := range res.Types {
		types[i] = t.Type.Name
	}

	pokemon.Stats = stats
	pokemon.Types = types

	return pokemon
}
