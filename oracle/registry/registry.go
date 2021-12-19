package registry

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/project-auxo/auxo/olympus/logging"
)

const (
	serviceManifest = "ServiceManifest.xml"
	relServicePath  = "oracle/services"
	relRegistryPath = "registry.xml"
	debug           = false
)

var log = logging.Base()
var servicesMap = make(map[string]Manifest)
var once sync.Once

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
	once.Do(func() {
		update(relServicePath, relRegistryPath)
	})
}

// Updates the servicesMap by looking through Oracle/services for any additional
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

	if debug {
		outToXml(*registry, outputPath)
	}
}

func outToXml(registry Registry, outputPath string) {
	if len(registry.Services) > 0 {
		sort.Slice(registry.Services, func(i, j int) bool {
			return registry.Services[i].Name < registry.Services[j].Name
		})

		f, err := os.Open(outputPath)
		if err != nil {
			log.Fatalf("failed to create output path for the registry.xml: %v", err)
		}
		xmlWriter := io.Writer(f)
		enc := xml.NewEncoder(xmlWriter)
		enc.Indent(" ", "  ")
		if err := enc.Encode(registry); err != nil {
			log.Fatalln(err)
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
