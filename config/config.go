package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Document struct {
	Version DocumentVersion
	Groups  []map[string]*Group
	Rules   []*Rule
}

type DocumentVersion string

const (
	DocumentVersion1_0 = "1.0"
)

type Group struct {
	Members []string `yaml:"members"`
}

// Rules

// TODO: sumtype instead? https://github.com/BurntSushi/go-sumtype
type Rule struct {
	NoMatchBetween *NoMatchNondirectionalCondition
	NoMatchTo      *NoMatchDirectionalCondition
}

type NoMatchNondirectionalCondition [2]*EntityMatcher

type NoMatchDirectionalCondition struct {
	From *EntityMatcher `yaml:"from"`
	To   *EntityMatcher `yaml:"to"`
}

type EntityMatcher struct {
	Groups *[]string
	People *[]string
}

func (e *EntityMatcher) DoesPersonMatch(personName string, groupName string) bool {
	if e.People != nil {
		for _, matchablePerson := range *e.People {
			if matchablePerson == personName {
				return true
			}
		}
	}

	if e.Groups != nil {
		for _, matchableGroup := range *e.Groups {
			if matchableGroup == groupName {
				return true
			}
		}
	}

	return false
}

func MustLoadConfigDocument(inputName string) *Document {
	filename := fmt.Sprintf("input/%v.yaml", inputName)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	buf, _ := ioutil.ReadAll(file)

	var document Document = Document{}

	err = yaml.Unmarshal(buf, &document)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("--- t:\n%+v\n\n", document)

	return &document
}
