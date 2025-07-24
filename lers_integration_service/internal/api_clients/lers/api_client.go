package lers

import "time"

type ApiClient interface {
	GetMeasurePoints(token string, serverHost string) ([]MeasurePoint, error)
	GetConsumptionData(accountID int, token string, serverHost string, measurePointID int, startDate, endDate string) (*ConsumptionResponse, error)
	PollMeasurePoints(token string, serverHost string, measurePointIDs []int, startDate, endDate string, timeout time.Duration) (*PollMeasurePointsResponse, error)
	ArePollingsCurrentlyRunning(token string, serverHost string) (bool, error)
	GetPollSessions(token string, serverHost string, startDate string, endDate string, timeout time.Duration) (map[int]string, error)
}

var _ ApiClient = (*apiClient)(nil)

type apiClient struct {
}

func NewApiClient() ApiClient {
	return &apiClient{}
}
