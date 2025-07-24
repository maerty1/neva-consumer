package models

import "encoding/xml"

type ZWSErrorResponse struct {
	XMLName               xml.Name           `xml:"zwsResponse"`
	NetToolsTaskGetErrors NetToolsTaskErrors `xml:"NetToolsTaskGetErrors"`
	RetVal                int                `xml:"RetVal"`
}

type NetToolsTaskErrors struct {
	TaskHandle string        `xml:"TaskHandle"`
	Errors     ErrorsSection `xml:"Errors"`
}

type ErrorsSection struct {
	Count int     `xml:"Count"`
	Errs  []Error `xml:"Err"`
}

type Error struct {
	Code   int     `xml:"Code"`
	ElemID int     `xml:"ElemID"`
	Type   int     `xml:"Type"`
	Param1 float64 `xml:"Param1"`
	Param2 float64 `xml:"Param2"`
	Name   string  `xml:"Name"`
	Text   string  `xml:"Text"`
}
