package zulu

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"zulu_updater/internal/models"
)

func (a apiClient) GetZuluTaskStatus(ctx context.Context, taskHandle string) (string, error) {
	xmlData := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<zulu-server service='zws' version='1.0.0'>
		<Command>
			<NetToolsTaskGetStatus>
				<TaskHandle>%s</TaskHandle>
			</NetToolsTaskGetStatus>
		</Command>
	</zulu-server>`, taskHandle)

	req, err := http.NewRequest("POST", a.baseUrl+"/zws", bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Невозможно создать запрос: %s", err.Error()))
	}

	a.setDefaultHeaders(req)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Невозможно отправить запрос: %s", err.Error()))
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("код ответа:" + resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Невозможно прочесть ответ: %s", err.Error()))
	}

	var response models.ZWSStatusResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка парсинга XML ответа: %s XML: %s", err.Error(), string(body)))
	}

	return response.NetToolsTaskGetStatus.Status, nil
}
