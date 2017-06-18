package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180.0
	lo1 = lon1 * math.Pi / 180.0
	la2 = lat2 * math.Pi / 180.0
	lo2 = lon2 * math.Pi / 180.0
	r = 6378100.0 // Earth radius in meters

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func main() {
	predFile := flag.String("predicted", "pred_grid.xyz", "predicted grid file")
	obsFile := flag.String("observed", "obs_grid.xyz", "observed grid file")
	aTime := flag.Float64("atime", 3.59, "alert time")
	flag.Parse()

	// variables for computations
	var sWaveVelocity float64 = 3.0
	var tolerance float64 = 0.5

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

	pLines := []string{} // slice of strings to hold the read in data from predicted.xyz
	oLines := []string{} // slice of strings to hold the read in data from observed.xyz

	// read in pred grid
	pScanner := bufio.NewScanner(pFile)
	for pScanner.Scan() {
		pLines = append(pLines, pScanner.Text())
	}

	// read in obs grid
	oScanner := bufio.NewScanner(oFile)
	for oScanner.Scan() {
		oLines = append(oLines, oScanner.Text())
	}

	// observed values
	oLat, _ := strconv.ParseFloat(strings.Split(oLines[0], " ")[2], 64)
	oLon, _ := strconv.ParseFloat(strings.Split(oLines[0], " ")[3], 64)

	// array creations to hold elements from obs and pred grids
	pRows := len(pLines)
	oRows := len(oLines)
	pVals := make([][]float64, pRows)
	oVals := make([][]float64, oRows)
	numTP := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numFP := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numFN := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numTN := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	truePR := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numTPT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numFPT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numFNT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	numTNT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	falsePR := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	truePRT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	falsePRT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	distGrid := make([]float64, pRows)
	mmiThresholds := [6]float64{2.0, 3.0, 4.0, 5.0, 6.0, 7.0}

	// split predicted file lines into array
	for i := 0; i < len(pLines); i++ {
		lat, _ := strconv.ParseFloat(strings.Split(pLines[i], " ")[1], 64)
		lon, _ := strconv.ParseFloat(strings.Split(pLines[i], " ")[0], 64)
		mmi, _ := strconv.ParseFloat(strings.Split(pLines[i], " ")[4], 64)
		pVals[i] = make([]float64, 3) // 3 elements in each pVal
		pVals[i][0] = lat
		pVals[i][1] = lon
		pVals[i][2] = mmi
		//fmt.Println(i, ": ", lat, lon, mmi)
	}

	// split observed lines into array
	for i := 0; i < len(oLines); i++ {
		lat, _ := strconv.ParseFloat(strings.Split(oLines[i], " ")[1], 64)
		lon, _ := strconv.ParseFloat(strings.Split(oLines[i], " ")[0], 64)
		mmi, _ := strconv.ParseFloat(strings.Split(oLines[i], " ")[4], 64)
		oVals[i] = make([]float64, 3) // 3 elements in each oVal
		oVals[i][0] = lat
		oVals[i][1] = lon
		oVals[i][2] = mmi
	}

	// populate distance grid
	for i := 0; i < len(pVals); i++ {
		distGrid[i] = Distance(oLat, oLon, pVals[i][0], pVals[i][1])
	}

	// loop through each array and compare/populate the ROC arrays
	// For more on ROC, read https://en.wikipedia.org/wiki/Receiver_operating_characteristic
	for i := 0; i < len(mmiThresholds); i++ {
		for j := 0; j < len(oVals); j++ {
			if oVals[j][2] > mmiThresholds[i] {
				if (pVals[j][2] + tolerance) > mmiThresholds[i] {
					numTP[i] = numTP[i] + 1.0
				} else {
					numFN[i] = numFN[i] + 1.0
				}
			} else {
				if (pVals[j][2] - tolerance) > mmiThresholds[i] {
					numFP[i] = numFP[i] + 1.0
				} else {
					numTN[i] = numTN[i] + 1.0
				}
			}
			// populate the true and false positive rates
			if (numTP[i] + numFN[i]) != 0.0 {
				truePR[i] = numTP[i] / (numTP[i] + numFN[i])
				falsePR[i] = numFP[i] / (numTP[i] + numFN[i])
			}
		}
	}

	// loop through each array and compare/populate the ROC arrays
	// For more on ROC, read https://en.wikipedia.org/wiki/Receiver_operating_characteristic
	// Timeliness calculated
	for i := 0; i < len(mmiThresholds); i++ {
		for j := 0; j < len(oVals); j++ {
			distKM := distGrid[j] / 1000.0
			groundMT := distKM / sWaveVelocity
			if oVals[j][2] > mmiThresholds[i] {
				if (pVals[j][2] + tolerance) > mmiThresholds[i] {
					if *aTime <= groundMT {
						numTPT[i] = numTPT[i] + 1.0
					} else {
						numFNT[i] = numFNT[i] + 1.0
					}
				} else {
					numFNT[i] = numFNT[i] + 1.0
				}
			} else {
				if (pVals[j][2] - tolerance) > mmiThresholds[i] {
					numFPT[i] = numFPT[i] + 1.0
				} else {
					numTNT[i] = numTNT[i] + 1.0
				}
			}
			// populate the true and false positive rates
			if (numTPT[i] + numFNT[i]) != 0.0 {
				truePRT[i] = numTPT[i] / (numTPT[i] + numFNT[i])
				falsePRT[i] = numFPT[i] / (numTPT[i] + numFNT[i])
			}
		}
	}

	fmt.Println("Timeliness not considered")
	fmt.Println(truePR)
	fmt.Println(falsePR)
	fmt.Println("\nTimeliness considered")
	fmt.Println(truePRT)
	fmt.Println(falsePRT)

	// remove references to slices; memory allocated earlier with make can now be garbage collected
	pVals = nil
	oVals = nil
}
