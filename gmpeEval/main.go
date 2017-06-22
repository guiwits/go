package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"./gmpEval"
	"./gmpEvalXML"
)

func main() {

	/* Command line options */
	obsFile := flag.String("o", "", "Observed grid.xml file")
	gmType := flag.String("g", "mmi", "mmi, pga, or pgv")
	tolerance := flag.Float64("t", 0.5, "Tolerance")
	threshold := flag.Float64("v", 4.0, "threshold value")
	alg := flag.String("a", "dm", "algorithm")
	evt := flag.String("e", "20140824_southnapa", "event name: eg 20140824_southnapa")
	dirPath := flag.String("d", "/app/shakemap/data", "Path to the directory where the shakemap event.xmls are stored")
	dt := flag.Float64("dt", 0.5, "Time increment in GMPE")

	/* parse all the flags */
	flag.Parse()

	/* check to make sure that an observed grid.xml file was entered as an command line argument */
	if *obsFile == "" {
		fmt.Println("You must specify a observed shakemap grid.xml file")
		os.Exit(1)
	}

	/* Write the observed grid.xml into a proper shakemap xml struct */
	/* read in the observed grid.xml file */
	oFile, err := os.Open(*obsFile)
	if err != nil {
		fmt.Println(err)
	}
	defer oFile.Close()
	oLines := []string{} // slice of strings to hold the read in data from observed.xml
	// read our xml files as a byte array and create a variable of type gmpEvalXML.ShakemapGrid.
	oGridBytes, _ := ioutil.ReadAll(oFile)
	var oGrid gmpEvalXML.ShakemapGrid
	/* write the xml from oGridBytes to oGrid */
	xml.Unmarshal(oGridBytes, &oGrid)
	lines := strings.Split(oGrid.GridData, "\n")
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			oLines = append(oLines, lines[i])
		}
	}

	// set observed lat and lon variables to pass along to gmpEval
	oLat := oGrid.SMEvent.Lat
	oLon := oGrid.SMEvent.Lon
	/* Done setting up the observed.xml data structures */

	/* string pointing to where we can find the algorithm grid.xml data */
	xmlData := *dirPath + "/" + *alg + "_" + *evt + "_*"
	fmt.Println("XML Data dir:", xmlData)

	/* find all the directories in the 'xmlData' location */
	files, _ := filepath.Glob(xmlData)
	sort.Strings(files) // sorting the files so we know how many event.xml's we have

	/* create slice to hold all the directory names where the grid.xml file are located */
	gridDirs := []string{}

	/* create a map to map dirname to slice of bytes (to store the shakemap xml data) */
	gridMap := make(map[string][]gmpEvalXML.ShakemapGrid)

	/* map to hold alert times */
	alertTimes := make(map[string][]float64)

	/* loop over the files and populate the data structures with the grid.xml files */
	for i, f := range files {
		tmp := *alg + "_" + *evt + "_" + string(f[len(f)-3:]) // removing the _run_x from the name
		gridDirs = append(gridDirs, tmp)
		res, err := strconv.ParseFloat((f[len(f)-3:]), 64)
		if err != nil {
			panic(err)
		}

		xmlFileName := f + "/output/grid.xml"
		// read in grid.xml files from each directory
		fi, err := os.Open(xmlFileName)
		if err != nil {
			panic(err)
		}
		// close fi on exit and check for its returned error
		defer func() {
			if err := fi.Close(); err != nil {
				panic(err)
			}
		}()
		// make a read buffer for the grid data (byte array)
		r, _ := ioutil.ReadAll(fi)
		alertTime := res * (*dt)
		var tmpSMXML gmpEvalXML.ShakemapGrid // temporary shakemap struct to hold the grid.xml data.
		xml.Unmarshal(r, &tmpSMXML)
		gridMap[gridDirs[i]] = append(gridMap[gridDirs[i]], tmpSMXML)
		alertTimes[gridDirs[i]] = append(alertTimes[gridDirs[i]], alertTime)
	}

	/* Now that we have all the grid.xml files stored in a slice of string arrays *
	 * we can send the gridMap and gridDirs over to the GmpEval function.         *
	 * GridDirs is needed fo the keys into gridMap                                */
	_, err = gmpEval.GmpEval(oLat, oLon, oLines, gridDirs, gridMap, alertTimes, *tolerance, *threshold, *gmType)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("gmpe evaluation complete.")
	}
}
