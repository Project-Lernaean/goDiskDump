package goDiskDump

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/karrick/godirwalk"
)

// Contains a compressed buffer of results in JSON
type Results struct {
	zbuff bytes.Buffer
}

// Contains all the info related to an file entry
type Entry struct {
	IsDir    bool
	Name     string
	Ext      string
	FullPath string
}

// DumpDisk used godirwalk to rapidly scan a directory for files
func DumpDisk(basePaths []string) (r Results, err error) {
	zbuff := zlib.NewWriter(&r.zbuff)
	defer zbuff.Close()

	// Walk Disk
	for _, s := range basePaths {

		err = godirwalk.Walk(s, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) error {

				var e entry
				e.IsDir, _ = de.IsDirOrSymlinkToDir()
				e.FullPath = osPathname
				e.Name = de.Name()
				e.Ext = filepath.Ext(osPathname)
				j, _ := json.MarshalIndent(e, "", "  ")
				zbuff.Write([]byte(string(j)))

				return nil
			},
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				return godirwalk.SkipNode
			},
			Unsorted: false,
		})
	}

	if err != nil {
		return r, err
	}

	return r, nil
}

func (dd Results) Json() (string, error) {
	r, err := zlib.NewReader(&dd.zbuff)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var out strings.Builder
	io.Copy(&out, r)

	return out.String(), nil
}

// GetPaths returns the base paths for different OS's
func getPaths() (out []string, err error) {
	err = nil
	switch runtime.GOOS {

	case "linux":
		return []string{"/"}, err

	case "darwin":
		return []string{"/"}, err

	case "windows":
		return []string{"C:\\"}, err

	default:
		return []string{}, fmt.Errorf("Unsupported OS")
	}
}
