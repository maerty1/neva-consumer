package geojson

// GeoJSONFeatureCollectionV2 представляет собой GeoJSON FeatureCollection
// swagger:model GeoJSONFeatureCollection
type GeoJSONFeatureCollectionV2 struct {
	Type     string             `json:"type" example:"FeatureCollection"` // Тип объекта, всегда "FeatureCollection"
	Features []GeoJSONFeatureV2 `json:"features"`                         // Массив объектов Feature
}

// GeoJSONFeatureV2 представляет собой отдельный GeoJSON Feature
// swagger:model GeoJSONFeature
type GeoJSONFeatureV2 struct {
	ID         int          `json:"id" example:"1320"`      // Уникальный идентификатор объекта
	Type       string       `json:"type" example:"Feature"` // Тип объекта, всегда "Feature"
	Geometry   Geometry     `json:"geometry"`               // Геометрия объекта
	Properties PropertiesV2 `json:"properties"`             // Свойства объекта
}

// PropertiesV2 представляет собой свойства GeoJSON Feature
// swagger:model Properties
type PropertiesV2 struct {
	Name         string `json:"Name" example:"П-427"`
	Adres        string `json:"Adres" example:"Пушкина 15"`
	ElemID       int    `json:"elem_id" example:"1320"`           // Идентификатор элемента
	ZwsMode      int    `json:"zws_mode" example:"1"`             // Режим объекта
	ZwsType      int    `json:"zws_type" example:"6"`             // Тип объекта
	ParentID     int    `json:"parent_id" example:"3058"`         // Идентификатор родительского элемента
	ZwsLinecolor int    `json:"zws_linecolor" example:"10966016"` // Цвет линии в формате целого числа
}
