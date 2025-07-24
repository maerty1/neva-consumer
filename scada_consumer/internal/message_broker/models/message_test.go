package models_test

import (
	"encoding/json"
	"testing"

	"scada_consumer/internal/message_broker/models"

	"github.com/stretchr/testify/require"
)

func TestMessage_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected string
	}{
		{
			name:     "String value",
			jsonData: `{"value": "stringValue", "dataType": "string", "lastChanged": "2024-04-04T07:04:15.5548405Z"}`,
			expected: "stringValue",
		},
		{
			name:     "Float value",
			jsonData: `{"value": 123.456, "dataType": "float", "lastChanged": "2024-04-04T07:04:15.5548405Z"}`,
			expected: "123.456",
		},
		{
			name:     "Integer value",
			jsonData: `{"value": 789, "dataType": "int", "lastChanged": "2024-04-04T07:04:15.5548405Z"}`,
			expected: "789",
		},
		{
			name:     "Boolean value",
			jsonData: `{"value": true, "dataType": "bool", "lastChanged": "2024-04-04T07:04:15.5548405Z"}`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg models.Message
			err := json.Unmarshal([]byte(tt.jsonData), &msg)
			require.NoError(t, err)
			require.Equal(t, tt.expected, msg.Value)
		})
	}
}

func TestParseHash(t *testing.T) {
	hash := "V1_CTP6_DAMN"

	message := models.Message{
		Hash: hash,
	}

	err := message.ParseHash()

	if err != nil {
		t.Errorf("ParseHash() вернул ошибку: %v", err)
	}

	if message.Version != "V1" {
		t.Errorf("Ожидалось Version = 'V1', получено '%s'", message.Version)
	}

	if message.DataSourceName != "CTP6" {
		t.Errorf("Ожидалось DataSourceName = 'CTP6', получено '%s'", message.DataSourceName)
	}

	if message.Variable != "DAMN" {
		t.Errorf("Ожидалось Variable = 'DAMN', получено '%s'", message.Variable)
	}
}
