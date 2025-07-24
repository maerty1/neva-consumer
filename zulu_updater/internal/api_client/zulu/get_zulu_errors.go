package zulu

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"zulu_updater/internal/models"
)

func (a apiClient) GetZuluErrors(ctx context.Context, taskHandle string) (int, []models.Error, error) {
	xmlData := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
            <zulu-server service='zws' version='1.0.0'>
                <Command>
                    <NetToolsTaskGetErrors>
                        <TaskHandle>%s</TaskHandle>
                    </NetToolsTaskGetErrors>
                </Command>
            </zulu-server>`, taskHandle)

	req, err := http.NewRequest("POST", a.baseUrl+"/zws", bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("Невозможно создать запрос: %s", err.Error()))
	}

	a.setDefaultHeaders(req)

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("Невозможно отправить запрос: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, nil, errors.New("код ответа:" + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("Невозможно прочесть ответ: %s", err.Error()))
	}

	var response models.ZWSErrorResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return 0, nil, errors.New(fmt.Sprintf("Ошибка парсинга XML ответа: %s XML: %s", err.Error(), string(body)))
	}

	return response.NetToolsTaskGetErrors.Errors.Count, response.NetToolsTaskGetErrors.Errors.Errs, nil
}
