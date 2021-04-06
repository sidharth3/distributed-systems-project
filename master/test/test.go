package test

import (
	"ds-proj/master/structs"
)

// Pure empty init
func EmptyCase() *structs.Master {
	return structs.InitMaster()
}

// Initialize with single file in namespace
func SimpleCase() *structs.Master {
	m := structs.InitMaster()
	m.Namespace.SetHash("/foo/bar/test_file.txt", "d383caabf6289b8ad52e401dafb20fb301ec3b760d1708e2501e5a39f130a1fc")
	return m
}
