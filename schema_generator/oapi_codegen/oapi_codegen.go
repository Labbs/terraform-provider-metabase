package oapi_codegen

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type oapiCodegen struct {
	Package  string   `yaml:"package"`
	Generate []string `yaml:"generate"`
	Output   string   `yaml:"output"`
}

func NewOapiCodegenConfig(path string) {
	oapi := oapiCodegen{
		Package:  "schema",
		Generate: []string{"types", "client", "models"},
		Output:   path + "/client.gen.go",
	}

	// Marshal the struct to YAML
	data, err := yaml.Marshal(oapi)
	if err != nil {
		log.Fatalf("error while generating oapi-codegen configuration file: %v", err)
	}

	// Write the YAML to a file
	err = os.WriteFile(path+"/oapi-codegen-config.yaml", data, 0644)
	if err != nil {
		log.Fatalf("error writing oapi-codegen configuration file: %v", err)
	}
}

func ExecuteCodegen(path string) {
	// Execute the oapi-codegen command
	cmd := exec.Command("go", "run", "github.com/deepmap/oapi-codegen/cmd/oapi-codegen",
		"-config", path+"/oapi-codegen-config.yaml",
		path+"/updated-openapi.json")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error while running oapi-codegen: %v\n", err)
		fmt.Printf("Output:\n%s\n", output)
		return
	}
}
