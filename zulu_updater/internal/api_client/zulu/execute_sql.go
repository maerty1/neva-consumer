package zulu

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
	"zulu_updater/internal/models"
)

func (a apiClient) ExecuteSqlGetParametersByValNames(ctx context.Context, valNames []string, elemId int) (*models.Records, error) {
	valNamesStr := strings.Join(valNames, ",")
	query := fmt.Sprintf(`SELECT %s where Sys=%d`, valNamesStr, elemId)

	xmlData := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<zulu-server service='zws' version='1.0.0'>
    <Command lang='ru'>
        <LayerExecSql>
            <Layer>%s</Layer>
            <Query>%s</Query>
            <CRS>EPSG:4326</CRS>
        </LayerExecSql>
    </Command>
</zulu-server>`, a.layer, query)

	req, err := http.NewRequest("POST", a.baseUrl+"/zws", bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Невозможно создать запрос: %s", err.Error()))
	}

	a.setDefaultHeaders(req)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Невозможно отправить запрос: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("код ответа:" + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Невозможно прочесть ответ: %s", err.Error()))
	}

	var response models.ZwsSqlResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Ошибка парсинга XML ответа: %s XML: %s", err.Error(), string(body)))
	}
	return &response.LayerExecSql.Records, nil
}
