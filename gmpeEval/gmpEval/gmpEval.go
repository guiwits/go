package gmpEval

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"../gmpEvalXML"
)

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// http://en.wikipedia.org/wiki/Haversine_formula
func distance(lat1, lon1, lat2, lon2 float64) float64 {
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

// GmpEval -> main function for comparing 2 shakemap files
func GmpEval(oLat float64, oLon float64, observedData []string, gridDirs []string,
	gridMap map[string][]gmpEvalXML.ShakemapGrid, tolerance float64, aTime float64,
	threshold float64, gmType string) (int, error) {
	// default value for sWaveVelocity
	var sWaveVelocity float64 = 3.0

	fmt.Println(oLat)
	fmt.Println(oLon)
	fmt.Println(sWaveVelocity)
	fmt.Println(tolerance)
	fmt.Println(aTime)
	fmt.Println(threshold)
	fmt.Println(gmType)

	for i := 0; i < len(gridDirs); i++ {
		smevt := gridMap[gridDirs[i]]
		fmt.Println(smevt[0].ShakemapEventType)
		fmt.Println(smevt[0].SMEvent.EventID)
		fmt.Println(smevt[0].SMEvent.Magnitude)
		fmt.Println(smevt[0].SMEvent.Depth)
		fmt.Println(smevt[0].SMEvent.Lat)
		fmt.Println(smevt[0].SMEvent.Lon)
		fmt.Println(smevt[0].SMEvent.EventTimestamp)
		fmt.Println(smevt[0].SMEvent.EventNetwork)
		fmt.Println(smevt[0].SMEvent.EventDescription)
		fmt.Println("--------------")
		//fmt.Println(smevt[0].GridData)
		fmt.Println()
	}

	// array creations to hold elements from obs and pred grids
	//pRows := len(pData)
	oRows := len(observedData)
	//pVals := make([][]float64, pRows)
	oVals := make([][]float64, oRows)
	//numTP := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//numFP := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//numFN := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//numTN := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	truePR := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//numTPT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//numFPT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//numFNT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//numTNT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	falsePR := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	truePRT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	falsePRT := [6]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	//distGrid := make([]float64, pRows)

	//// We want to generalize from mmi to gm (ground motion)
	//// Either make the thresholds a command line input, or have
	//// three different default arrays for mmi, pga, pgv
	//mmiThresholds := [6]float64{2.0, 3.0, 4.0, 5.0, 6.0, 7.0}

	// split predicted file lines into array
	/*
		for i := 0; i < len(pData); i++ {
			lat, _ := strconv.ParseFloat(strings.Split(pData[i], " ")[1], 64)
			lon, _ := strconv.ParseFloat(strings.Split(pData[i], " ")[0], 64)
			//// Generalize to gm and [4] for mmi, [2] for pga, [3] for pgv
			//// This will be true for Shakemap and the eqinfo2gm grid.xml product
			mmi, _ := strconv.ParseFloat(strings.Split(pData[i], " ")[4], 64)
			pVals[i] = make([]float64, 3) // 3 elements in each pVal
			pVals[i][0] = lat
			pVals[i][1] = lon
			pVals[i][2] = mmi
			//fmt.Println("pData: ", i, ": ", lat, lon, mmi)
		}
	*/

	// split observed lines into array
	for i := 0; i < len(observedData); i++ {
		lat, _ := strconv.ParseFloat(strings.Split(observedData[i], " ")[1], 64)
		lon, _ := strconv.ParseFloat(strings.Split(observedData[i], " ")[0], 64)

		//// Again generalize to gm and [4] for mmi, [2] for pga, [3] for pgv
		mmi, _ := strconv.ParseFloat(strings.Split(observedData[i], " ")[4], 64)

		oVals[i] = make([]float64, 3) // 3 elements in each oVal
		oVals[i][0] = lat
		oVals[i][1] = lon
		oVals[i][2] = mmi
	}

	// populate distance grid
	/*
		for i := 0; i < len(pVals); i++ {
			distGrid[i] = distance(oLat, oLon, pVals[i][0], pVals[i][1])
		}
	*/

	// loop through each array and compare/populate the ROC arrays
	// For more on ROC, read https://en.wikipedia.org/wiki/Receiver_operating_characteristic
	/*
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
	*/

	// loop through each array and compare/populate the ROC arrays
	// For more on ROC, read https://en.wikipedia.org/wiki/Receiver_operating_characteristic
	// Timeliness calculated

	//// -This is the one we wish to modify because it takes into account timeliness
	//// -The first loop can be generalized to gmthresholds or simplified to just one threshold
	//// -I think the test we wish to run in the third loop is if the pVals for any version at that
	//// grid point are above the threshold rather than oVals.
	//// - Then test if the observed value was above the threshold+tolerance, below threshold-tolerance, etc
	//// -If it passes these tests, then consider whether or not it is timely.  What we really want is the
	//// difference between alert time and ground truth origin time.  This is something that is saved in
	//// the analysis script.  We want to see if this time is <= groundMT.
	//// etc...to be continued
	/*
		for i := 0; i < len(mmiThresholds); i++ {
			for j := 0; j < len(oVals); j++ {
				distKM := distGrid[j] / 1000.0
				groundMT := distKM / sWaveVelocity
				if oVals[j][2] > mmiThresholds[i] {
					if (pVals[j][2] + tolerance) > mmiThresholds[i] {
						if aTime <= groundMT {
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
	*/

	fmt.Println("Timeliness not considered")
	fmt.Println(truePR)
	fmt.Println(falsePR)
	fmt.Println("\nTimeliness considered")
	fmt.Println(truePRT)
	fmt.Println(falsePRT)

	// remove references to slices; memory allocated earlier with make can now be garbage collected
	//pVals = nil
	oVals = nil

	return 0, nil
}
