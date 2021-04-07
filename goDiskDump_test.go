package goDiskDump_test

import (
	"fmt"
	"testing"

	"github.com/Project-Lernaean/goDiskDump"
)

func TestMain(m *testing.M) {
	disks, err := goDiskDump.GetPaths()
	if err != nil {
		panic(err)
	}
	fmt.Println(disks)

	buff, err := goDiskDump.DumpDisk(disks)
	if err != nil {
		panic(err)
	}

	fmt.Println(buff.Json())
}
