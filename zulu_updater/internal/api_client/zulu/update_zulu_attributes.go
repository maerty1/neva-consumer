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

func (a apiClient) UpdateAttribute(ctx context.Context, attributeName string, value float64, elemId int) error {
	xmlData := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
        <zulu-server service="zws" version="1.0.0">
          <Command>
            <UpdateElemAttributes>
              <Layer>%s</Layer>
              <Element>
                <Key>
                  <Name>Sys</Name>
                  <Value>%d</Value>
                </Key>
                <Field>
                  <Name>%s</Name>
                  <Value>%f</Value>
                </Field>
              </Element>
            </UpdateElemAttributes>
          </Command>
        </zulu-server>`, a.layer, elemId, attributeName, value)

	req, err := http.NewRequest("POST", a.baseUrl+"/zws", bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return errors.New(fmt.Sprintf("Невозможно создать запрос: %s", err.Error()))
	}

	a.setDefaultHeaders(req)

	resp, err := a.client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Невозможно отправить запрос: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("код ответа:" + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Невозможно прочесть ответ: %s", err.Error()))
	}

	var response models.ZWSUpdateAttributeResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {

		return errors.New(fmt.Sprintf("Ошибка парсинга XML ответа: %s XML: %s", err.Error(), string(body)))
	}

	fmt.Printf("Update ElemAttributes: %s\nRetVal: %d\n", response.UpdateElemAttributes, response.RetVal)
	return nil
}
