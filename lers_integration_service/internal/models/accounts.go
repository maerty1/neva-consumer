package models

type AccountToSync struct {
	ID         int    `json:"id"`
	Token      string `json:"token"`
	ServerHost string `json:"server_host"`
}
