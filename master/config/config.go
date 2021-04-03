package config

const (
<<<<<<< HEAD
	REPLICAS    = 3                // Number of replicas
	TIMEOUT     = 5                // Timeout when connecting to slaves in seconds
	HBINTERVAL  = 5                // Interval between heartbeats in seconds
	FLINTERVAL  = 10               // Interval between updating FileLocations in seconds
	GCINTERVAL  = 5                // Interval between sending namespace values to slaves for garbage collection in seconds
	REPINTERVAL = 15               // Interval between replication cycles
	DQINTERVAL  = 30               // Timeout to delete the UID from Queue
	IP          = "127.0.0.1:8080" // IP address for the master
=======
  REPLICAS    = 3  // Number of replicas
	TIMEOUT     = 5  // Timeout when connecting to slaves in seconds
	HBINTERVAL  = 5  // Interval between heartbeats in seconds
	FLINTERVAL  = 10 // Interval between updating FileLocations in seconds
	GCINTERVAL  = 5  // Interval between sending namespace values to slaves for garbage collection in seconds
	REPINTERVAL = 15 // Interval between replication cycles
  DQINTERVAL = 30 // Timeout to delete the UID from Queue
>>>>>>> 23718734330cfc879c3e2cce57dc9c6dc6f21b81
)
