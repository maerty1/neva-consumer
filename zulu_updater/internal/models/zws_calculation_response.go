package models

import "encoding/xml"

type ZwsRunCalculationResponse struct {
	XMLName         xml.Name        `xml:"zwsResponse"`
	NetToolsTaskRun NetToolsTaskRun `xml:"NetToolsTaskRun"`
	RetVal          int             `xml:"RetVal"`
}

type NetToolsTaskRun struct {
	TaskHandle string `xml:"TaskHandle"`
}
