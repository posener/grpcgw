package swagger

import (
	"path/filepath"
	"os"
	"log"
	"fmt"
	"strings"
	"io"
)

func GenerateSwaggerGoFile(swaggerOutFile string) {
	swaggerOutDir := filepath.Dir(swaggerOutFile)

	err := os.MkdirAll(swaggerOutDir, os.ModePerm)
	if err != nil {
		log.Panicf("Failed creating swagger go directory: %s", err)
	}
	out, err := os.Create(swaggerOutFile)
	if err != nil {
		log.Panicf("Failed creating swagger go file: %s", err)
	}

	packageName := filepath.Base(swaggerOutDir)
	fmt.Fprintf(out, swaggerGoFileHeader, packageName)

	// Search files created in the out directory
	// the swagger json files were created in dir/.../<package>/<service>.swagger.json
	// add each file content to the swagger go file
	filepath.Walk(swaggerOutDir, generateSwaggerGoEachFile(out))
	out.Write([]byte("}\n"))
}

func generateSwaggerGoEachFile (out *os.File) filepath.WalkFunc {
	return func (path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			log.Printf("Processing %s", path)
			name  := filepath.Base(filepath.Dir(path))
			out.Write([]byte("\"" + name + "\": `"))
			f, _ := os.Open(path)
			io.Copy(out, f)
			out.Write([]byte("`,\n"))
		}
		return nil
	}
}

const swaggerGoFileHeader = `package %s

// Autogenerated file

var JSONs = map[string]string {
`
