package configreader

import (
	"fmt"
	"io"
	"os"

	"github.com/shockey/secret-santa/rules"
	"gopkg.in/yaml.v2"
)

type Document struct {
	Version                 DocumentVersion
	Groups                  []map[string]*Group
	Rules                   []*rules.Rule
	InactiveGroupMembersMap map[[2]string]bool `yaml:"-"`
	InactiveGroupMembers    [][]string         `yaml:"inactiveGroupMembers"`
}

type DocumentVersion string

const (
	DocumentVersion1_0 = "1.0"
	DocumentVersion1_1 = "1.1"
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

	buf, _ := io.ReadAll(file)

	var document Document = Document{}

	err = yaml.Unmarshal(buf, &document)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	mustValidateConfigStructure(document.Version, &buf)

	mustValidateConfigDocument(&document)

	// Saves consumers the trouble of having to create the map themselves
	document.InactiveGroupMembersMap = make(map[[2]string]bool)
	for _, inactiveGroupMember := range document.InactiveGroupMembers {
		groupName, memberName := inactiveGroupMember[0], inactiveGroupMember[1]
		document.InactiveGroupMembersMap[[2]string{groupName, memberName}] = true
	}

	return &document
}
