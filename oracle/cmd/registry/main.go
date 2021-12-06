package main

import (
	"os"
	"strings"

	"github.com/project-auxo/auxo/olympus/logging"
	"github.com/project-auxo/auxo/oracle/registry"
)

const relServicePath = "oracle/services"
const relRegistryPath = "oracle/registry/registry.xml"

var log = logging.Base()

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("current directory broken: %v", err)
	}
	if strings.HasSuffix(currDir, "auxo") {
		outputPath := currDir + "/" + relRegistryPath
		err = os.Chdir(currDir + "/" + relServicePath)
		if err != nil {
			log.Fatalf("failed to change to the service directory: %v", err)
		}
		currDir, _ = os.Getwd()
		registry.Update(currDir, outputPath)
	}
	log.Info("Successfully updated Oracle's registry.")
}
