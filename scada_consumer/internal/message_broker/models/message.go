package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Version            string `json:"-"`
	DataSourceName     string `json:"-"`
	Variable           string `json:"-"`
	MeasurementTypeID  string `json:"measurement_type_id"`
	Value              string `json:"value"`
	DataType           string `json:"dataType"`
	LastChanged        string `json:"lastChanged"`
	StatusCodes        int    `json:"statusCodes"`
	NodeId             string `json:"nodeId"`
	NodeName           string `json:"nodeName"`
	OwnerId            int    `json:"ownerId"`
	Hash               string `json:"hash"`
	DataPointClassEnum string `json:"dataPointClassEnum"`
	RabbitMQMessage    amqp091.Delivery
}

func (m *Message) UnmarshalJSON(data []byte) error {
	// Временная структура, чтобы избежать рекурсии
	type Alias Message
	aux := &struct {
		Value interface{} `json:"value"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.Value.(type) {
	case string:
		m.Value = v
	case float64:
		m.Value = fmt.Sprintf("%v", v)
	case int:
		m.Value = fmt.Sprintf("%v", v)
	case bool:
		m.Value = fmt.Sprintf("%v", v)
	default:
		return fmt.Errorf("неожиданный тип для поля значения")
	}

	return nil
}

func (m *Message) ParseHash() error {
	if strings.HasPrefix(m.Hash, "V2_") {
		// Формат V2: V2_CTP14_4 — в конце measurement_type_id
		parts := strings.SplitN(m.Hash, "_", 3)
		if len(parts) != 3 {
			return fmt.Errorf("неверный формат хеша для V2: %s", m.Hash)
		}
		m.Version = parts[0]
		m.DataSourceName = parts[1]
		m.MeasurementTypeID = parts[2]
		log.Println("V2")
	} else if strings.HasPrefix(m.Hash, "V1_") {
		log.Println("V1")
		// Формат V1: V1_CTP14_10MIN_1_Nasos_perepuska_Narabotka
		parts := strings.SplitN(m.Hash, "_", 3)
		if len(parts) != 3 {
			return fmt.Errorf("неверный формат хеша для V1: %s", m.Hash)
		}
		m.Version = parts[0]
		m.DataSourceName = parts[1]
		m.Variable = parts[2]
	} else {
		return fmt.Errorf("неизвестная версия хеша: %s", m.Hash)
	}

	return nil
}

// func (m *Message) ParseHash() error {
// 	parts := strings.SplitN(m.Hash, "_", 3)
// 	if len(parts) != 3 {
// 		return fmt.Errorf("неверный формат хеша: %s", m.Hash)
// 	}
// 	m.Version = parts[0]
// 	m.DataSourceName = parts[1]
// 	m.Variable = parts[2]
// 	return nil
// }

// ALTER TABLE public.scada_rawdata ADD COLUMN measurement_type_id integer NULL DEFAULT NULL;
