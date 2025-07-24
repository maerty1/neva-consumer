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

func (a apiClient) RunZuluCalculation(ctx context.Context) (string, error) {
	elemId := 3580
	xmlData := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
            <zulu-server service='zws' version='1.0.0'>
                <Command>
                    <NetToolsTaskRun>
                        <Layer>%s</Layer>
                        <Model>thermo</Model>
                        <Task>2</Task>
                        <Settings>
                            <Sources>%d</Sources>
                            <DHW>no</DHW>
                            <Leaks>no</Leaks>
                            <HeatLoss>1</HeatLoss>
                            <Nozzle>no</Nozzle>
                        </Settings>
                    </NetToolsTaskRun>
                </Command>
            </zulu-server>`, a.layer, elemId)

	req, err := http.NewRequest("POST", a.baseUrl+"/zws", bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Невозможно создать запрос: %s", err.Error()))
	}

	a.setDefaultHeaders(req)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Невозможно отправить запрос: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("код ответа:" + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Невозможно прочесть ответ: %s", err.Error()))
	}

	var response models.ZwsRunCalculationResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ошибка парсинга XML ответа: %s XML: %s", err.Error(), string(body)))
	}

	return response.NetToolsTaskRun.TaskHandle, nil
}
