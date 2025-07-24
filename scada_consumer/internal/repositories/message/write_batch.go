package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"scada_consumer/internal/message_broker/models"
	"strconv"
	"time"
)

func (r *repository) WriteBatch(ctx context.Context, messages []models.Message) error {
	if len(messages) == 0 {
		return nil
	}

	titlesSet := make(map[string]struct{})
	for _, msg := range messages {
		titlesSet[msg.DataSourceName] = struct{}{}
	}
	titles := make([]string, 0, len(titlesSet))
	for title := range titlesSet {
		titles = append(titles, title)
	}

	titleToID, err := r.getMeasurePointIDs(ctx, titles)
	if err != nil {
		return fmt.Errorf("ошибка получения measure point IDs: %w", err)
	}

	for title := range titlesSet {
		if _, exists := titleToID[title]; !exists {
			id, err := r.ensureMeasurePoint(ctx, title)
			if err != nil {
				log.Printf("не удалось обеспечить scada_measure_point для заголовка '%s': %v", title, err)
				continue
			}
			titleToID[title] = id
		}
	}

	tx, err := r.db.DB().Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("не удалось отменить транзакцию: %v", rbErr)
			}
		}
	}()

	insertQuery := `
		INSERT INTO scada_rawdata (scada_measure_point_id, timestamp, varname, value, raw_packet, measurement_type_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, msg := range messages {
		log.Println(msg.MeasurementTypeID)
		scadaMeasurePointID, ok := titleToID[msg.DataSourceName]
		if !ok {
			log.Printf("scada_measure_point с названием '%s' не был найден после ensuring", msg.DataSourceName)
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, msg.LastChanged)
		if err != nil {
			log.Printf("неверный формат timestamp для сообщения с хешем '%s': %v", msg.Hash, err)
			continue
		}

		if !json.Valid(msg.RabbitMQMessage.Body) {
			log.Printf("недопустимый формат JSON для raw_packet в сообщении с хешем '%s'", msg.Hash)
			continue
		}

		rawPacket := string(msg.RabbitMQMessage.Body)
		var measurementPointID *int

		if len(msg.MeasurementTypeID) == 0 {
			measurementPointID = nil
		} else {
			tmp, _ := strconv.Atoi(msg.MeasurementTypeID)
			measurementPointID = &tmp

		}
		_, err = tx.Exec(ctx, insertQuery, scadaMeasurePointID, timestamp, msg.Variable, msg.Value, rawPacket, measurementPointID)
		if err != nil {
			log.Printf("не удалось вставить в scada_rawdata сообщение с хешем '%s': %v", msg.Hash, err)
			continue
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("не удалось подтвердить транзакцию: %w", err)
	}

	return nil
}

func (r *repository) getMeasurePointIDs(ctx context.Context, titles []string) (map[string]int, error) {
	if len(titles) == 0 {
		return nil, nil
	}

	query := `
		SELECT title, id FROM scada_measure_points
		WHERE title = ANY($1) AND account_id = $2
	`

	rows, err := r.db.DB().Query(ctx, query, titles, 1)
	if err != nil {
		return nil, fmt.Errorf("не удалось запросить scada_measure_points: %w", err)
	}
	defer rows.Close()

	titleToID := make(map[string]int)
	for rows.Next() {
		var title string
		var id int
		if err := rows.Scan(&title, &id); err != nil {
			return nil, fmt.Errorf("не удалось просканировать строку scada_measure_point: %w", err)
		}
		titleToID[title] = id
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации строк scada_measure_points: %w", err)
	}

	return titleToID, nil
}

func (r *repository) ensureMeasurePoint(ctx context.Context, title string) (int, error) {
	var id int
	err := r.db.DB().QueryRow(ctx, `
		INSERT INTO scada_measure_points (account_id, title)
		VALUES ($1, $2)
		ON CONFLICT (account_id, title) DO UPDATE SET title = EXCLUDED.title
		RETURNING id
	`, 1, title).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось вставить или получить scada_measure_point: %w", err)
	}
	return id, nil
}
