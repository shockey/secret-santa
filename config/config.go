package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/shockey/secret-santa/rules"
	"gopkg.in/yaml.v2"
)

type Document struct {
	Version DocumentVersion
	Groups  []map[string]*Group
	Rules   []*rules.Rule
}

type DocumentVersion string

const (
	DocumentVersion1_0 = "1.0"
)

type Group struct {
	Members []string `yaml:"members"`
}

func MustLoadConfigDocument(inputName string) *Document {
	filename := fmt.Sprintf("input/%v.yaml", inputName)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	buf, _ := ioutil.ReadAll(file)

	mustValidateConfigStructure(&buf)

	var document Document = Document{}

	err = yaml.Unmarshal(buf, &document)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	mustValidateConfigDocument(&document)

	return &document
}
