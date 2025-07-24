package geodata

import "time"

// ObjectState представляет собой состояние объекта из словаря
// swagger:model ObjectState
type ObjectState struct {
	ZwsType int    `json:"zws_type" example:"2"`
	ZwsMode int    `json:"zws_mode" example:"2"`
	Title   string `json:"title" example:"Разветвление"`
	Image   string `json:"image" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEgAAABICAYA"` // Изображение состояния в формате data URL (base64)
}

// ObjectStatesResponse представляет собой массив состояний объектов
// swagger:model ObjectStatesResponse
type ObjectStatesResponse struct {
	States []ObjectState `json:"states"` // Массив состояний объектов
}

// ElementData представляет собой массив значений объектов
// swagger:model ElementData
type ElementData struct {
	Parameter  string    `json:"parameter" example:"sys"`
	Val        string    `json:"val" example:"1"`
	RecordType string    `json:"record_type" example:"static"`
	InsertedTS time.Time `json:"inserted_ts" example:"2024-10-08T10:56:45.531005Z"`
}

type MeasurementGroup struct {
	Coeff *float64 `json:"coeff"`
	In    string   `json:"i" example:"T_in"`
	Out   string   `json:"o" example:"T_out"`
}

type Point struct {
	ElemID            int                      `json:"elem_id"`
	Title             *string                  `json:"title" example:"Котельная 22"`
	Address           *string                  `json:"address" example:"Улица Пушкина 12"`
	MeasurementGroups map[int]MeasurementGroup `json:"measurement_groups"`
	Coordinates       []float64                `json:"coordinates" description:"[lat, lon]" example:"55.751244,37.618423"`
	HasAccident       bool                     `json:"has_accident"`
	Type              int                      `json:"type"`
	IsCopied          bool                     `json:"iscopied"`
}

type FullElementData struct {
	Address string                          `json:"address"`
	Title   string                          `json:"title"`
	Packets map[string]map[int]*Measurement `json:"packets"`
}

type Measurement struct {
	Name           string                    `json:"name"`
	Unit           string                    `json:"unit"`
	ZuluCoeff      *float64                  `json:"zulu_coeff"`
	LersCoeff      *float64                  `json:"lers_coeff"`
	Data           MeasurementData           `json:"data"`
	CalculatedData MeasurementCalculatedData `json:"calculated_data"`
}

type MeasurementData struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

type MeasurementCalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type Category struct {
	Name      string `json:"name" example:"Метрики"`
	Type      string `json:"type" example:"group"`
	IsOpen    bool   `json:"is_open"`
	MaxValues *int   `json:"max_values"`
	ID        int    `json:"id"`
}

type PointWithCategories struct {
	Title      *string    `json:"title" example:"Котельная 22"`
	Address    *string    `json:"address" example:"Улица Пушкина 12"`
	Type       int        `json:"type"`
	Categories []Category `json:"categories"`
}
