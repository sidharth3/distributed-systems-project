package config

const (
	REPLICAS    = 3                // Number of replicas
	TIMEOUT     = 5                // Timeout when connecting to slaves in seconds
	HBINTERVAL  = 5                // Interval between heartbeats in seconds
	FLINTERVAL  = 10               // Interval between updating FileLocations in seconds
	GCINTERVAL  = 5                // Interval between sending namespace values to slaves for garbage collection in seconds
	REPINTERVAL = 15               // Interval between replication cycles
	MGCINTERVAL = 35               // Interval between master garbage collection cycles. Guarantee that namespace entries are at least MGCINTERVAL seconds old before being removed. MGCINTERVAL > HBINTERVAL + FLINTERVAL + 2*client http timeout + master http timeout
	LDINTERVAL  = 10               // Interval between load checking in seconds
	IP          = "127.0.0.1:8080" // IP address for the master
	// IPLIST      = []string{"127.0.0.1:8080","127.0.0.1:8081","127.0.0.1:8082"} // IP address for the master
)
