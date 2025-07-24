package models

import (
	"encoding/xml"
)

type ZwsSqlResponse struct {
	XMLName      xml.Name     `xml:"zwsResponse"`
	LayerExecSql LayerExecSql `xml:"LayerExecSql"`
	RetVal       int          `xml:"RetVal"`
}

type LayerExecSql struct {
	Records Records `xml:"Records"`
}

type Records struct {
	Record []Record `xml:"Record"`
}

type Record struct {
	Field []Field `xml:"Field"`
}

type Field struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
}
