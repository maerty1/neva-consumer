package geodata

// Данные для получения сырья
type GroupMeasurementsData struct {
	In  string `json:"in" example:"T_in"`
	Out string `json:"out" example:"T_out"`
}

// Данные из Зулу
type GroupMeasurementsCalculatedData struct {
	In  *float64 `json:"in"`
	Out *float64 `json:"out"`
}

type GroupMeasurement struct {
	Name           string                          `json:"name" example:"Температура"`
	Unit           string                          `json:"unit" example:"атм"`
	CalculatedData GroupMeasurementsCalculatedData `json:"calculated_data"`
	Data           GroupMeasurementsData           `json:"data"`
	ZuluCoeff      *float64                        `json:"zulu_coeff"`
	LersCoeff      *float64                        `json:"lers_coeff"`
	Rn             int                             `json:"rn"`
}

type GetPointDataByCategoryGroup struct {
	Measurements map[string]map[int]*GroupMeasurement `json:"measurements"`
}

type MeasurementKeyvalue struct {
	Name   string      `json:"name" example:"Температура"`
	Unit   string      `json:"unit" example:"атм"`
	Source string      `json:"source" example:"zulu/scada"`
	Value  interface{} `json:"value"`
	Rn     int         `json:"rn"`
}

type GetPointDataByCategoryKeyvalue struct {
	Measurements []MeasurementKeyvalue `json:"measurements"`
}
