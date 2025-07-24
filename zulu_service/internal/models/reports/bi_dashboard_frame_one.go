package reports

type BiDashboardFrame struct {
	BlockID      int     `json:"block_id"`
	BlockName    string  `json:"block_name" example:"Котельная №16"`
	Qsum         float64 `json:"qsum" example:"23.13"`
	SpecificQsum float64 `json:"specific_qsum"`
}
