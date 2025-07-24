package models

import (
	"encoding/xml"
)

type ZWSUpdateAttributeResponse struct {
	XMLName              xml.Name `xml:"zwsResponse"`
	UpdateElemAttributes string   `xml:"UpdateElemAttributes,omitempty"`
	RetVal               int      `xml:"RetVal"`
}

type RecordsJsonStruct struct {
	Parameter string `json:"parameter"`
	Val       string `json:"val"`
	ElemID    string `json:"elem_id"`
	LersTS    int64  `json:"lers_ts"`
	ParentID  int    `json:"parent_id"`
}
