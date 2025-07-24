package geodata

import (
	"context"
	"zulu_service/internal/models/geodata"
)

func (r *service) GetStates(ctx context.Context) ([]geodata.ObjectState, error) {
	states, err := r.geodataRepository.GetStates(ctx)
	if err != nil {
		return nil, err
	}

	return states, nil

}
