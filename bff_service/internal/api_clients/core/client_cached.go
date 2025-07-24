package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheApiClient struct {
	underlying ApiClient
	cache      *cache.Cache
}

func NewCacheApiClient(client ApiClient, ttl time.Duration) ApiClient {
	return &CacheApiClient{
		underlying: client,
		cache:      cache.New(ttl, 10*time.Minute), // TTL и интервал очистки
	}
}

func (c *CacheApiClient) GetElementIDs() ([]int, error) {
	key, err := generateKey("GetElementIDs", nil)
	if err != nil {
		return nil, err
	}

	// Попытка получить данные из кеша
	if cached, found := c.cache.Get(key); found {
		if elementIDs, ok := cached.([]int); ok {
			fmt.Println("Кеш использован для GetElementIDs")
			return elementIDs, nil
		}
		// Если тип данных не совпадает, игнорируем кеш и продолжаем
	}

	fmt.Println("Кеш пропущен для GetElementIDs")

	// Вызов исходного клиента для получения данных
	elementIDs, err := c.underlying.GetElementIDs()
	if err != nil {
		return nil, err
	}

	// Сохранение результата в кеш
	c.cache.Set(key, elementIDs, cache.DefaultExpiration)
	fmt.Println("Данные сохранены в кеш для GetElementIDs")

	return elementIDs, nil
}

func (c *CacheApiClient) GetPointsData(reqData []GetPointsDataRequest, timestamp string) ([]GetPointsDataResponse, error) {
	key, err := generateKey("GetPointsData", reqData, timestamp)
	if err != nil {
		return nil, err
	}

	// Попытка получить данные из кеша
	if cached, found := c.cache.Get(key); found {
		if pointsData, ok := cached.([]GetPointsDataResponse); ok {
			fmt.Println("Кеш использован для GetPointsData")
			return pointsData, nil
		}
		// Если тип данных не совпадает, игнорируем кеш и продолжаем
	}

	fmt.Println("Кеш пропущен для GetPointsData")

	// Вызов исходного клиента для получения данных
	pointsData, err := c.underlying.GetPointsData(reqData, timestamp)
	if err != nil {
		return nil, err
	}

	// Сохранение результата в кеш
	c.cache.Set(key, pointsData, cache.DefaultExpiration)
	fmt.Println("Данные сохранены в кеш для GetPointsData")
	// response, err := c.underlying.GetPointsData(reqData, timestamp)
	return pointsData, err
}
func (c *CacheApiClient) GetPointsDataHistory(reqData []GetPointsDataHistoryRequest, nDays int, timestamp string) (GetPointsDataHistoryResponse, error) {
	key, err := generateKey("GetPointsDataHistory", reqData, nDays)
	if err != nil {
		return GetPointsDataHistoryResponse{}, err
	}

	// Попытка получить данные из кеша
	if cached, found := c.cache.Get(key); found {
		if response, ok := cached.(GetPointsDataHistoryResponse); ok {
			fmt.Println("Кеш использован для GetPointsDataHistory")
			return response, nil
		}
		// Если тип данных не совпадает, игнорируем кеш и продолжаем
	}
	fmt.Println("Кеш пропущен для GetPointsDataHistory")

	// Вызов исходного клиента для получения данных
	response, err := c.underlying.GetPointsDataHistory(reqData, nDays, timestamp)
	if err != nil {
		return GetPointsDataHistoryResponse{}, err
	}

	// Сохранение результата в кеш
	c.cache.Set(key, response, cache.DefaultExpiration)
	fmt.Println("Данные сохранены в кеш для GetPointsDataHistory")

	return response, nil
}

func generateKey(method string, params ...interface{}) (string, error) {
	var paramStrings []string

	for _, param := range params {
		paramBytes, err := json.Marshal(param)
		if err != nil {
			return "", fmt.Errorf("ошибка маршалинга параметра (%v): %w", param, err)
		}
		paramStrings = append(paramStrings, string(paramBytes))
	}

	joinedParams := strings.Join(paramStrings, ":")

	return fmt.Sprintf("%s:%s", method, joinedParams), nil
}
