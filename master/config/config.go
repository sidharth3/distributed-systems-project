package config

const (
	TIMEOUT    = 5  // Timeout when connecting to slaves in seconds
	HBINTERVAL = 5  // Interval between heartbeats in seconds
	FLINTERVAL = 10 // Interval between updating FileLocations in seconds
	GCINTERVAL = 5 // Interval between sending namespace values to slaves for garbage collection in seconds
)
