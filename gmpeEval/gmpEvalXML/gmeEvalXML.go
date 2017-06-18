package gmpEvalXML

import (
	"encoding/xml"
)

// ShakemapGrid struct will be an exported type
type ShakemapGrid struct {
	XMLName                  xml.Name                   `xml:"shakemap_grid"`
	ID                       string                     `xml:",omitempty"`
	Name                     string                     `xml:",omitempty"`
	EventID                  string                     `xml:"event_id,attr"`
	ShakemapID               string                     `xml:"shakemap_id,attr"`
	ShakemapVersion          string                     `xml:"shakemap_version,attr"`
	CodeVersion              string                     `xml:"code_version,attr"`
	ProcessTimestamp         string                     `xml:"process_timestamp,attr"`
	ShakemapOriginator       string                     `xml:"shakemap_originator,attr"`
	MapStatus                string                     `xml:"map_status,attr"`
	ShakemapEventType        string                     `xml:"shakemap_event_type,attr"`
	SMEvent                  SMEvent                    `xml:"event_info"`
	GridSpecification        []GridSpecification        `xml:"grid_specification"`
	EventSpecificUncertainty []EventSpecificUncertainty `xml:"event_specific_uncertainty"`
	GridField                []GridField                `xml:"grid_field"`
	GridData                 string                     `xml:"grid_data"`
}

// SMEvent struct will be an exported type
type SMEvent struct {
	EventID          string  `xml:"event_id"`
	Magnitude        string  `xml:"magnitude"`
	Depth            string  `xml:"depth"`
	Lat              float64 `xml:"lat"`
	Lon              float64 `xml:"lon"`
	EventTimestamp   string  `xml:"event_timestamp"`
	EventNetwork     string  `xml:"event_network"`
	EventDescription string  `xml:"event_description"`
}

// GridSpecification struct will be an exported type
type GridSpecification struct {
	LonMin            string `xml:"lon_min"`
	LatMin            string `xml:"lat_min"`
	LonMax            string `xml:"lon_max"`
	LatMax            string `xml:"lat_max"`
	NominalLonSpacing string `xml:"nominal_lon_spacing"`
	NominalLatSpacing string `xml:"nominal_lat_spacing"`
	Nlon              string `xml:"nlon"`
	Nlat              string `xml:"nlat"`
}

// EventSpecificUncertainty struct will be an exported type
type EventSpecificUncertainty struct {
	Name   string `xml:"name,attr"`
	Value  string `xml:"value,attr"`
	NumSta string `xml:"numsta,attr"`
}

// GridField struct will be an exported type
type GridField struct {
	Index string `xml:"index,attr"`
	Name  string `xml:"name,attr"`
	Units string `xml:"units,attr"`
}
