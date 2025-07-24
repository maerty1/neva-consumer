package geojson

// GeoJSONFeatureCollection представляет собой GeoJSON FeatureCollection
// swagger:model GeoJSONFeatureCollection
type GeoJSONFeatureCollection struct {
	Type     string           `json:"type" example:"FeatureCollection"` // Тип объекта, всегда "FeatureCollection"
	Features []GeoJSONFeature `json:"features"`                         // Массив объектов Feature
}

// GeoJSONFeature представляет собой отдельный GeoJSON Feature
// swagger:model GeoJSONFeature
type GeoJSONFeature struct {
	ID         int        `json:"id" example:"1320"`      // Уникальный идентификатор объекта
	Type       string     `json:"type" example:"Feature"` // Тип объекта, всегда "Feature"
	Geometry   Geometry   `json:"geometry"`               // Геометрия объекта
	Properties Properties `json:"properties"`             // Свойства объекта
}

// Geometry представляет собой геометрическую информацию GeoJSON Feature
// swagger:model Geometry
type Geometry struct {
	Type        string      `json:"type" example:"LineString"` // Тип геометрии, например, "LineString"
	Coordinates [][]float64 `json:"coordinates"`               // Координаты геометрии
}

// Properties представляет собой свойства GeoJSON Feature
// swagger:model Properties
type Properties struct {
	ElemID       int `json:"elem_id" example:"1320"`           // Идентификатор элемента
	ZwsMode      int `json:"zws_mode" example:"1"`             // Режим объекта
	ZwsType      int `json:"zws_type" example:"6"`             // Тип объекта
	ParentID     int `json:"parent_id" example:"3058"`         // Идентификатор родительского элемента
	ZwsLinecolor int `json:"zws_linecolor" example:"10966016"` // Цвет линии в формате целого числа
}
