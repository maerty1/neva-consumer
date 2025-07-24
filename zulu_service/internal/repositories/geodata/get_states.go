package geodata

import (
	"context"
	"log"
	"zulu_service/internal/models/geodata"
)

func (r *repository) GetStates(ctx context.Context) ([]geodata.ObjectState, error) {
	query := `
        SELECT zws_type, zws_mode, title, image 
        FROM zulu.dict_object_states;
    `

	rows, err := r.db.DB().Query(ctx, query)
	if err != nil {
		log.Printf("Ошибка при получении состояний объектов: %v", err)
		return nil, err
	}
	defer rows.Close()

	var states []geodata.ObjectState
	for rows.Next() {
		var state geodata.ObjectState
		if err := rows.Scan(&state.ZwsType, &state.ZwsMode, &state.Title, &state.Image); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}
		states = append(states, state)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, err
	}

	return states, nil
}
