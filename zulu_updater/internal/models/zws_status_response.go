package models

import "encoding/xml"

type ZWSStatusResponse struct {
	XMLName               xml.Name           `xml:"zwsResponse"`
	NetToolsTaskGetStatus NetToolsTaskStatus `xml:"NetToolsTaskGetStatus"`
	RetVal                int                `xml:"RetVal"`
}

type NetToolsTaskStatus struct {
	TaskHandle     string `xml:"TaskHandle"`
	Status         string `xml:"Status"`
	NetToolsRetVal int    `xml:"NetToolsRetVal"`
}
