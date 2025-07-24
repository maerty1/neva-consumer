package zulu

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

func (c *CacheApiClient) GetFilteredPoints(elementIds []int, zwsTypeIDs []int, timestamp string) ([]Point, error) {
	key, err := generateKey("GetFilteredPoints", elementIds, zwsTypeIDs, timestamp)
	if err != nil {
		return nil, err
	}

	if cached, found := c.cache.Get(key); found {
		if points, ok := cached.([]Point); ok {
			fmt.Println("Кеш использован для GetFilteredPoints")
			return points, nil
		}
	}

	fmt.Println("Кеш пропущен для GetFilteredPoints")

	points, err := c.underlying.GetFilteredPoints(elementIds, zwsTypeIDs, timestamp)
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, points, cache.DefaultExpiration)
	fmt.Println("Данные сохранены в кеш для GetFilteredPoints")

	return points, nil
}

// Реализация остальных методов ApiClient без кеширования
func (c *CacheApiClient) GetPoints(zwsTypeIds []int) ([]Point, error) {
	return c.underlying.GetPoints(zwsTypeIds)
}

func (c *CacheApiClient) GetFullPoint(elemID int, nDays int) (FullElementData, error) {
	return c.underlying.GetFullPoint(elemID, nDays)
}

func (c *CacheApiClient) GetPointCategoryDataGroup(elemID int, categoryID int, timestamp string) (GetPointDataByCategoryGroup, error) {
	return c.underlying.GetPointCategoryDataGroup(elemID, categoryID, timestamp)
}

func (c *CacheApiClient) GetPointCategoryDataKeyvalue(elemID int, categoryID int) (GetPointDataByCategoryKeyvalue, error) {
	return c.underlying.GetPointCategoryDataKeyvalue(elemID, categoryID)
}

func (c *CacheApiClient) GetPointsDataCategoryByZwsType(zwsTypeIds []int, categoryID int, timestamp string) (GetPointsDataCategoryResponse, error) {
	return c.underlying.GetPointsDataCategoryByZwsType(zwsTypeIds, categoryID, timestamp)

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
