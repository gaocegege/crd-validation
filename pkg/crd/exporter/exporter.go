package exporter

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// Exporter exports the yaml config to file.
type Exporter struct {
	outputDir string
	filename  string
	writer    *bufio.Writer
}

func NewFileExporter(outputDir, filename string) *Exporter {
	file, err := filepath.Abs(filepath.Join(e.outputDir, e.filename))
	if err != nil {
		log.Fatalf("Failed to open the output file %s: %v", file, err)
	}

	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Failed to open the output file %s: %v", file, err)
	}

	return &Exporter{
		outputDir: outputDir,
		filename:  filename,
		writer: bufio.NewWriter(f),
	}
}

func NewStdoutExporter() *Exporter {
	return &Exporter{
		writer: bufio.NewWriter(os.Stdout),
	}
}

// Export exports the yaml config to file.
func (e *Exporter) Export(final *apiextensions.CustomResourceDefinition) {
	
	e.writer.WriteString("# The code is generated by crd-validation\n\n")
	e.marshallCrd(final, "yaml")
	e.writer.Flush()
}

func (e Exporter) marshallCrd(crd *apiextensions.CustomResourceDefinition, outputFormat string) {
	jsonBytes, err := json.MarshalIndent(crd, "", "    ")
	if err != nil {
		log.Fatal("error:", err)
	}

	// Doing the following because the status section should not exist in the CRD yaml, but because the type definition of
	// CustomResourceDefinition, the field Status is a struct, which omitempty does not apply properly to, the status
	// section will still be generated when we marshal the CRD object to the yaml. What we are doing here is we take
	// an extra step, unmarshaling the jsonBytes to a map[string]interface{}, and delete the key "status" from the map,
	// and then marshal the redacted map to Json byte array, and then convert that to YAML.
	var redactedMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &redactedMap)
	if err != nil {
		log.Fatal("error:", err)
	}
	delete(redactedMap, "status")
	redactedJsonBytes, err := json.MarshalIndent(redactedMap, "", "    ")


	if outputFormat == "json" {
		e.writer.Write(redactedJsonBytes)
	} else {
		yamlBytes, err := yaml.JSONToYAML(redactedJsonBytes)
		if err != nil {
			log.Fatal("error:", err)
		}
		e.writer.WriteString("---\n")
		e.writer.Write(yamlBytes)
	}
}
