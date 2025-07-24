package message_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"scada_consumer/internal/message_broker/models"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func TestWriteBatch(t *testing.T) {
	pgConn := appInstance.S().PostgresDB().DB()

	messages := []models.Message{
		{
			Value:              "100",
			DataType:           "DINT",
			LastChanged:        "2024-04-06T07:05:15.5548405Z",
			StatusCodes:        12,
			NodeId:             "f80f7085-1dbf-4f49-9e63-a279a4a3227e",
			NodeName:           "RabbitMQ super",
			OwnerId:            1,
			Hash:               "V1_CTP6_var1",
			DataPointClassEnum: "Input",
			DataSourceName:     "datasource1",
			Variable:           "var1",
		},
		{
			Value:              "200",
			DataType:           "DINT",
			LastChanged:        "2024-04-06T07:05:15.5548405Z",
			StatusCodes:        12,
			NodeId:             "f80f7085-1dbf-4f49-9e63-a279a4a3227f",
			NodeName:           "RabbitMQ ultra",
			OwnerId:            1,
			Hash:               "V1_CTP6_var2",
			DataPointClassEnum: "Output",
			DataSourceName:     "datasource2",
			Variable:           "var2",
		},
	}

	// Сериализация сообщения в сырой JSON-пакет
	for i := range messages {
		rawPacket, err := json.Marshal(map[string]interface{}{
			"value":              messages[i].Value,
			"dataType":           messages[i].DataType,
			"lastChanged":        messages[i].LastChanged,
			"statusCodes":        messages[i].StatusCodes,
			"nodeId":             messages[i].NodeId,
			"nodeName":           messages[i].NodeName,
			"ownerId":            messages[i].OwnerId,
			"hash":               messages[i].Hash,
			"dataPointClassEnum": messages[i].DataPointClassEnum,
		})
		require.NoError(t, err, "ошибка сериализации сообщения в сырой пакет JSON")
		messages[i].RabbitMQMessage = amqp091.Delivery{Body: rawPacket}
	}

	err := repositoryTest.WriteBatch(context.Background(), messages)
	require.NoError(t, err, "функция WriteBatch вернула ошибку")

	var count int
	err = pgConn.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM scada_rawdata
		WHERE varname = ANY($1)
	`, []string{"var1", "var2"}).Scan(&count)
	require.NoError(t, err, "ошибка выполнения запроса к scada_rawdata")
	require.Equal(t, 2, count, "ожидалось 2 записи в scada_rawdata, получено %d", count)

	err = pgConn.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM scada_measure_points
		WHERE title = ANY($1) AND account_id = $2
	`, []string{"datasource1", "datasource2"}, 1).Scan(&count)
	require.NoError(t, err, "ошибка выполнения запроса к scada_measure_points")
	require.Equal(t, 2, count, "ожидалось 2 записи в scada_measure_points, получено %d", count)

	for _, msg := range messages {
		var scadaMeasurePointID int
		var timestamp time.Time
		var varName, value, rawPacket string

		err := pgConn.QueryRow(context.Background(), `
			SELECT scada_measure_point_id, timestamp, varname, value, raw_packet
			FROM scada_rawdata
			WHERE varname = $1
		`, msg.Variable).Scan(&scadaMeasurePointID, &timestamp, &varName, &value, &rawPacket)
		require.NoError(t, err, "ошибка получения данных из scada_rawdata для переменной %s", msg.Variable)

		var expectedMeasurePointID int
		err = pgConn.QueryRow(context.Background(), `
			SELECT id FROM scada_measure_points
			WHERE title = $1 AND account_id = $2
		`, msg.DataSourceName, msg.OwnerId).Scan(&expectedMeasurePointID)
		require.NoError(t, err, "ошибка получения scada_measure_point_id для title %s", msg.DataSourceName)

		require.Equal(t, expectedMeasurePointID, scadaMeasurePointID, "scada_measure_point_id не совпадает для переменной %s", msg.Variable)
		require.Equal(t, msg.Variable, varName, "varname не совпадает для переменной %s", msg.Variable)
		require.Equal(t, msg.Value, value, "value не совпадает для переменной %s", msg.Variable)

		// Десериализуем оба JSON пакета и сравним их как структуры
		var actualRawPacket map[string]interface{}
		err = json.Unmarshal([]byte(rawPacket), &actualRawPacket)
		require.NoError(t, err, "ошибка десериализации raw_packet для переменной %s", msg.Variable)

		var expectedRawPacket map[string]interface{}
		err = json.Unmarshal(msg.RabbitMQMessage.Body, &expectedRawPacket)
		require.NoError(t, err, "ошибка десериализации ожидаемого сырого пакета для переменной %s", msg.Variable)

		require.Equal(t, expectedRawPacket, actualRawPacket, "raw_packet не совпадает для переменной %s", msg.Variable)
	}
}
