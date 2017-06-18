package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"./gmpEval"
	"./gmpEvalXML"
)

func main() {

	/* Command line options */
	predFile := flag.String("p", "pred_grid.xml", "Predicted grid.xml file")
	obsFile := flag.String("o", "obs_grid.xml", "Observed grid.xml file")
	aTime := flag.Float64("a", 1.00, "Alert time")
	gmType := flag.String("g", "mmi", "mmi, pga, or pgv")
	tolerance := flag.Float64("t", 0.5, "Tolerance")
	threshold := flag.Float64("v", 4.0, "threshold value")

	flag.Parse()

	/* read in each file into array */
	pFile, err := os.Open(*predFile)
	if err != nil {
		fmt.Println(err)
	}
	defer pFile.Close()

	oFile, err := os.Open(*obsFile)
	if err != nil {
		fmt.Println(err)
	}
	defer oFile.Close()

	pLines := []string{} // slice of strings to hold the read in data from predicted.xml
	oLines := []string{} // slice of strings to hold the read in data from observed.xml

	// read our xml files as a byte array.
	pGridBytes, _ := ioutil.ReadAll(pFile)
	oGridBytes, _ := ioutil.ReadAll(oFile)

	// initialize our grid arrays
	var pGrid gmpEvalXML.ShakemapGrid
	var oGrid gmpEvalXML.ShakemapGrid

	// unmarshal our arrays which contain our xml data
	xml.Unmarshal(pGridBytes, &pGrid)
	xml.Unmarshal(oGridBytes, &oGrid)

	// populate the pLines and oLine slices with the newly read in grid data
	lines := strings.Split(pGrid.GridData, "\n")
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			pLines = append(pLines, lines[i])
		}
	}

	lines = strings.Split(oGrid.GridData, "\n")
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			oLines = append(oLines, lines[i])
		}
	}

	// set observer lat and lon variables to pass along
	oLat := oGrid.SMEvent.Lat
	oLon := oGrid.SMEvent.Lon

	// GmpEval call
	_, err = gmpEval.GmpEval(oLat, oLon, oLines, pLines, *tolerance, *aTime, *threshold, *gmType)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("gmpe evaluation complete.")
	}
}
