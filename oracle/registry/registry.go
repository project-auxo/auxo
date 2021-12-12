package registry

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/project-auxo/auxo/olympus/logging"
)

const (
	serviceManifest = "ServiceManifest.xml"
	registryFile    = "registry.xml"
	relServicePath  = "oracle/services"
	relRegistryPath = "oracle/registry/registry.xml"
)

var log = logging.Base()
var servicesMap = make(map[string]Manifest)

type Manifest struct {
	XMLName     xml.Name `xml:"manifest"`
	ServiceName string   `xml:"service,attr"`
	Deployable  bool     `xml:"deployable,attr"`
}

type Registry struct {
	XMLName  xml.Name  `xml:"services"`
	Services []Service `xml:"service"`
}

type Service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`
}

func init() {
	currDir, err := os.Getwd()
	prevDir := currDir
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
		update(currDir, outputPath)
	}
	log.Info("Successfully updated Oracle's registry.")
	os.Chdir(prevDir)
}

// Updates the registry by looking through Oracle/services for any additional
// services.
func update(servicePath string, outputPath string) {
	registry := &Registry{}

	filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == serviceManifest {
			manifest := parseServiceManifest(path)
			registry.Services = append(registry.Services, Service{Name: manifest.ServiceName})
			servicesMap[manifest.ServiceName] = manifest
			return nil
		}
		return nil
	})

	if len(registry.Services) > 0 {
		sort.Slice(registry.Services, func(i, j int) bool {
			return registry.Services[i].Name < registry.Services[j].Name
		})

		f, _ := os.Create(outputPath)
		xmlWriter := io.Writer(f)
		enc := xml.NewEncoder(xmlWriter)
		enc.Indent(" ", "  ")
		if err := enc.Encode(registry); err != nil {
			log.Fatal(err)
		}
	}
}

func parseServiceManifest(path string) Manifest {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to parse the ServiceManifest.xml: %v", err)
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	var manifest Manifest
	xml.Unmarshal(byteValue, &manifest)
	return manifest
}

func ServiceExists(query string) bool {
	_, found := servicesMap[query]
	return found
}
