package gen

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

/*
   Creation Time: 2020 - Apr - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Descriptor
type Descriptor struct {
	Name   string  `yaml:"name" json:"name"`
	Models []Model `yaml:"models" json:"models"`
	APIs   []API   `yaml:"apis" json:"apis"`
}

// Primitive
type Primitive string

const (
	Int64  Primitive = "int64"
	Int32  Primitive = "int32"
	UInt64 Primitive = "uint64"
	UInt32 Primitive = "uint32"
	String Primitive = "string"
	Bytes  Primitive = "bytes"
	Byte   Primitive = "byte"
)

func (p Primitive) String() string {
	return string(p)
}

// Model
type Model struct {
	Name       string      `yaml:"name" json:"name"`
	Manager    string      `yaml:"manager" json:"manager"`
	Comments   []string    `yaml:"comments" json:"comment"`
	Properties []Property  `yaml:"properties" json:"properties"`
	PrimaryKey FilterKey   `yaml:"unique_key" json:"unique_key"`
	FilterKeys []FilterKey `yaml:"filter_keys" json:"filter_keys"`
}

func (m *Model) GetProperty(name string) (Property, error) {
	for _, p := range m.Properties {
		if p.Name == name {
			return p, nil
		}
	}
	return Property{}, fmt.Errorf("could not find %s in model %s", name, m.Name)
}

// FilterKey the model will be queried by this key. The key is built of Properties. UniqueCombination is optional and
// defines that which part of this key create a unique
type FilterKey struct {
	Name           string   `yaml:"name" json:"name"`
	PartitionKeys  []string `yaml:"partition_keys" json:"partition_keys"`
	ClusteringKeys []string `yaml:"clustering_keys" json:"clustering_keys"`
}

// Property
type Property struct {
	Name    string           `yaml:"name" json:"name"`
	Type    string           `yaml:"type" json:"type"`
	Comment string           `yaml:"comment" json:"comment"`
	Tags    []string         `yaml:"tags" json:"tags"`
	Options []PropertyOption `yaml:"options,flow" json:"options"` // e.g. [unique, slice, optional]
}

func (p Property) CheckOption(opt PropertyOption) bool {
	for idx := range p.Options {
		if p.Options[idx] == opt {
			return true
		}
	}
	return false
}

func (p Property) ToScyllaType() string {
	var t string
	switch Primitive(p.Type) {
	case Int32, UInt32:
		t = "int"
	case Int64, UInt64:
		t = "bigint"
	case Byte:
		t = "smallint"
	case String:
		t = "text"
	case Bytes:
		t = "blob"
	default:
		panic("invalid primitive type")
	}
	if p.CheckOption(Slice) {
		return fmt.Sprintf("frozen<list<%s>>", t)
	} else {
		return t
	}
}

// PropertyOption
type PropertyOption string

func (po PropertyOption) String() string {
	return string(po)
}

const (
	Slice    PropertyOption = "slice"
	Unique   PropertyOption = "unique"
	Optional PropertyOption = "optional"
)

// API
type API struct {
	Input   Message   `yaml:"input" json:"input"`
	Outputs []Message `yaml:"outputs" json:"outputs"`
}

// Message
type Message struct {
	Constructor int64 `yaml:"constructor" json:"constructor"`
}

var (
	descriptors []Descriptor
)

func LoadDescriptors(dirPath string) (err error) {
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		desc, err := loadDescriptor(path)
		if err == nil {
			descriptors = append(descriptors, desc)
		}
		return nil
	})
	return
}

func loadDescriptor(filePath string) (desc Descriptor, err error) {
	return readFromYaml(filePath)
}

func readFromYaml(filePath string) (desc Descriptor, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	yd := yaml.NewDecoder(f)
	err = yd.Decode(&desc)
	return
}
