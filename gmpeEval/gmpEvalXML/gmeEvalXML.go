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
	SMEvent                  SMEvent                    `xml:"event"`
	GridSpecification        []GridSpecification        `xml:"grid_specification"`
	EventSpecificUncertainty []EventSpecificUncertainty `xml:"event_specific_uncertainty"`
	GridField                []GridField                `xml:"grid_field"`
	GridData                 string                     `xml:"grid_data"`
}

// SMEvent struct will be an exported type
type SMEvent struct {
	EventID          string  `xml:"event_id,attr"`
	Magnitude        string  `xml:"magnitude,attr"`
	Depth            string  `xml:"depth,attr"`
	Lat              float64 `xml:"lat,attr"`
	Lon              float64 `xml:"lon,attr"`
	EventTimestamp   string  `xml:"event_timestamp,attr"`
	EventNetwork     string  `xml:"event_network,attr"`
	EventDescription string  `xml:"event_description,attr"`
}

// GridSpecification struct will be an exported type
type GridSpecification struct {
	LonMin            string `xml:"lon_min,attr"`
	LatMin            string `xml:"lat_min,attr"`
	LonMax            string `xml:"lon_max,attr"`
	LatMax            string `xml:"lat_max,attr"`
	NominalLonSpacing string `xml:"nominal_lon_spacing,attr"`
	NominalLatSpacing string `xml:"nominal_lat_spacing,attr"`
	Nlon              string `xml:"nlon,attr"`
	Nlat              string `xml:"nlat,attr"`
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
