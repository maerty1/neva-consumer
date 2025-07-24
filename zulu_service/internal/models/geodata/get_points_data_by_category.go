package geodata

type GetPointsDataByCategoryGroup struct {
	IsCopied     bool                      `json:"iscopied"` // Новое поле
	Measurements map[int]*GroupMeasurement `json:"measurements"`
}
