package scriptup

import (
	"os"
	"path"
	"runtime"
)

// createdTestFile describes the name of the file created in our dummy migrations.
const createdTestFile = ".test/foo.txt"

// testCfg is the Config to be used in testing.
var testCfg = &Config{
	Directory: ".test/migrations",
	Executor:  "/bin/bash",
	FileDB:    ".test/filedb.txt",
}

func init() {
	// set working directory to project root
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
