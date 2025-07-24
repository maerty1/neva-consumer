package zulu

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"strconv"
	"time"
	"zulu_updater/internal/models"
)

type ObjectRecordsFromJsonResponse struct {
	Finished bool
	ErrMsg   string
}

func (r repository) InsertObjectRecordsJson(
	ctx context.Context,
	data models.Records,
	updateType string,
	elemID int,
	lersTS time.Time) (*ObjectRecordsFromJsonResponse, error) {

	var output []models.RecordsJsonStruct
	for _, record := range data.Record[0].Field {
		output = append(output, models.RecordsJsonStruct{
			Parameter: record.Name,
			Val:       record.Value,
			ElemID:    strconv.Itoa(elemID),
			LersTS:    lersTS.Unix(),
			ParentID:  elemID,
		})
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		fmt.Println("Ошибка маршалинга JSON:", err)
		return nil, err
	}

	query := `SELECT * FROM zulu.object_records_fromjson_test($1::json, $2::text)`

	row := r.db.DB().QueryRow(ctx, query, string(jsonData), updateType)
	if err != nil {
		return nil, err
	}
	var res ObjectRecordsFromJsonResponse
	err = row.Scan(&res.Finished, &res.ErrMsg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
