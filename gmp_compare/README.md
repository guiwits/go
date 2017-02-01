This program compares the ground motion prediction of two ShakeMaps.
"ShakeMap is a product of the U.S. Geological Survey Earthquake Hazards Program 
in conjunction with regional seismic network operators. ShakeMap sites provide 
near-real-time maps of ground motion and shaking intensity following significant earthquakes."

[ShakeMap home](http://earthquake.usgs.gov/earthquakes/shakemap/)

A ShakeMap generates an .xyz file associated with its data. This program 
reads in two of these files files (one observed and one predicted) and then calculates
the [receiver operating characteristics](https://en.wikipedia.org/wiki/Receiver_operating_characteristic)
that help us determine how close our predicted results were to the observed results.

This program uses the [Haversine Formula](https://en.wikipedia.org/wiki/Haversine_formula) to determine 
the length (in meters) between two grid points in the .xyz files.

This program will take timeliness into account as well as the calculation that ignores the timeliness.

