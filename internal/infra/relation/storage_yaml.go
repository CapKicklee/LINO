package relation

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"makeit.imfr.cgi.com/lino/pkg/relation"
)

// Version of the YAML strcuture.
const Version string = "v1"

// YAMLStructure of the file.
type YAMLStructure struct {
	Version   string         `yaml:"version"`
	Relations []YAMLRelation `yaml:"relations,omitempty"`
}

// YAMLRelation defines how to store a relation in YAML format.
type YAMLRelation struct {
	Name   string    `yaml:"name"`
	Parent YAMLTable `yaml:"parent"`
	Child  YAMLTable `yaml:"child"`
}

// YAMLTable defines how to store a relation in YAML format.
type YAMLTable struct {
	Name string   `yaml:"name"`
	Keys []string `yaml:"keys"`
}

// YAMLStorage provides storage in a local YAML file
type YAMLStorage struct{}

// NewYAMLStorage create a new YAML storage
func NewYAMLStorage() *YAMLStorage {
	return &YAMLStorage{}
}

// List all relations stored in the YAML file
func (s YAMLStorage) List() ([]relation.Relation, *relation.Error) {
	list, err := readFile()
	if err != nil {
		return nil, err
	}
	result := []relation.Relation{}

	for _, ym := range list.Relations {
		m := relation.Relation{
			Name: ym.Name,
			Parent: relation.Table{
				Name: ym.Parent.Name,
				Keys: ym.Parent.Keys,
			},
			Child: relation.Table{
				Name: ym.Child.Name,
				Keys: ym.Child.Keys,
			},
		}
		result = append(result, m)
	}

	return result, nil
}

// Store relations in the YAML file
func (s YAMLStorage) Store(relations []relation.Relation) *relation.Error {
	list := YAMLStructure{
		Version: Version,
	}

	for _, r := range relations {
		yml := YAMLRelation{
			Name: r.Name,
			Parent: YAMLTable{
				Name: r.Parent.Name,
				Keys: r.Parent.Keys,
			},
			Child: YAMLTable{
				Name: r.Child.Name,
				Keys: r.Child.Keys,
			},
		}
		list.Relations = append(list.Relations, yml)
	}

	err := writeFile(&list)
	if err != nil {
		return err
	}

	return nil
}

func readFile() (*YAMLStructure, *relation.Error) {
	list := &YAMLStructure{
		Version: Version,
	}

	dat, err := ioutil.ReadFile("relations.yaml")
	if err != nil {
		return nil, &relation.Error{Description: err.Error()}
	}

	err = yaml.Unmarshal(dat, list)
	if err != nil {
		return nil, &relation.Error{Description: err.Error()}
	}

	if list.Version != Version {
		return nil, &relation.Error{Description: "invalid version in ./relations.yaml (" + list.Version + ")"}
	}

	return list, nil
}

func writeFile(list *YAMLStructure) *relation.Error {
	out, err := yaml.Marshal(list)
	if err != nil {
		return &relation.Error{Description: err.Error()}
	}

	err = ioutil.WriteFile("relations.yaml", out, 0640)
	if err != nil {
		return &relation.Error{Description: err.Error()}
	}

	return nil
}