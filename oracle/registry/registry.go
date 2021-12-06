package registry

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/project-auxo/auxo/olympus/logging"
)

const serviceManifest = "ServiceManifest.xml"

var log = logging.Base()

type Manifest struct {
	XMLName     xml.Name `xml:"manifest"`
	ServiceName string   `xml:"service,attr"`
}

type Registry struct {
	XMLName  xml.Name  `xml:"services"`
	Services []Service `xml:"service"`
}

type Service struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`
}

func Update(servicePath string, outputPath string) {
	registry := &Registry{}

	filepath.Walk(servicePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == serviceManifest {
			registry.Services = append(
				registry.Services, Service{Name: parseServiceManifest(path)})
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

func parseServiceManifest(path string) (serviceName string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to parse the ServiceManifest.xml: %v", err)
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	var manifest Manifest
	xml.Unmarshal(byteValue, &manifest)
	return manifest.ServiceName
}
