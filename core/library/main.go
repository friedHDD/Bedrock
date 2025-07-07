package library

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// md5 is the key
type BookInfo struct {
	OriginPath string `yaml:"originPath"`
	Series     string `yaml:"series"`
	BookName   string `yaml:"bookName"`
}

type Data struct {
	//string as index, BookInfo as value
	Books map[string]BookInfo `yaml:"books"`
}

// return error and how many books have been added
func Add(data map[string]BookInfo) (error, int) {
	libraryYamlFile := "./data/index/library.yaml"
	newBooks := 0
	var libraryData Data

	//read
	yamlFile, err := os.ReadFile(libraryYamlFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read library data from %s: %w", libraryYamlFile, err), 0
	}

	//unmarshal
	if err := yaml.Unmarshal(yamlFile, &libraryData); err != nil {
		return fmt.Errorf("failed to parse library data"), 0
	}

	//init
	if libraryData.Books == nil {
		libraryData.Books = make(map[string]BookInfo)
	}

	//k:string(md5) v:BookInfo
	for k, v := range data {
		//!exists, so add this book
		if _, exists := libraryData.Books[k]; !exists {
			libraryData.Books[k] = v
			newBooks++
		}
	}

	//newYaml: after adding books
	var newYaml []byte
	newYaml, err = yaml.Marshal(&libraryData)
	if err != nil {
		return fmt.Errorf("failed to pack data before saving"), 0
	}

	//write
	if err := os.WriteFile(libraryYamlFile, newYaml, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %v", libraryYamlFile, err), 0
	}

	return nil, newBooks
}
