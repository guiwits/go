This program compares the ground motion prediction of two shakemaps.
The shakemaps have a .xyz file associated with their data. This program 
reads in two files (one observed and one predicted) and then calculates
the receiver operating characteristic. 

https://en.wikipedia.org/wiki/Receiver_operating_characteristic

It uses the Haversine Formula to determing the length (in meters) between two
grid points, given in lat/lon.

https://en.wikipedia.org/wiki/Haversine_formula
