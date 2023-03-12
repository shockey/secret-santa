package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

func convertYamlBytesToJsonBytes(yamlBuf *[]byte) *[]byte {
	var tmp interface{}

	err := yaml.Unmarshal(*yamlBuf, &tmp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	// Context/author: https://stackoverflow.com/a/40737676
	tempMapS := dyno.ConvertMapI2MapS(tmp)

	jsonBytes, err := json.Marshal(tempMapS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	return &jsonBytes
}

func printValidationErrorsAndExit(errs []error) {
	fmt.Fprint(os.Stderr, "Your input configuration is invalid :(\nHere's what the validation engine said:\n")
	for i, err := range errs {
		fmt.Fprintf(os.Stderr, "%v: %v\n", i, err.Error())
	}
	os.Exit(1)
}
