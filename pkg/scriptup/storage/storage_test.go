package storage

import (
	"os"
	"path"
	"runtime"
)

func init() {
	// set working directory to project root
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
