package esfixture

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func containStringArray(ary []string, item string) bool {
	for _, s := range ary {
		if s == item {
			return true
		}
	}
	return false
}

func getFixtureFiles(dirPath string) (returnSchemaFileFullPaths []string, returnDocumentFileFullPaths []string, returnErr error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		returnErr = err
		return
	}

	returnSchemaFileFullPaths = make([]string, 0)
	returnDocumentFileFullPaths = make([]string, 0)
	for _, f := range files {
		fileName := f.Name()
		if !f.IsDir() {
			splits := strings.Split(fileName, ".")
			var isIncludeSchema bool
			var isIncludeDocument bool
			for _, split := range splits {
				if split == postFixSchema {
					isIncludeSchema = true
					break
				}
				if split == postFixDocument {
					isIncludeDocument = true
					break
				}
			}
			fullPath := fmt.Sprintf("%s/%s", dirPath, fileName)
			if isIncludeSchema {
				returnSchemaFileFullPaths = append(returnSchemaFileFullPaths, fullPath)
			} else if isIncludeDocument {
				returnDocumentFileFullPaths = append(returnDocumentFileFullPaths, fullPath)
			}
		}
	}
	return
}

func readSchemaFile(fullPath string) ([]byte, error) {
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func readDocumentFile(fullPath string) ([]byte, error) {
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func genFullPathWithFileNameForSchema(dir, targetName string) string {
	return fmt.Sprintf("%s/%s.%s.%s", dir, targetName, postFixSchema, extJSON)
}

func genFullPathWithFileNameForDocument(dir, targetName string) string {
	return fmt.Sprintf("%s/%s.%s.%s", dir, targetName, postFixDocument, extNDJSON)
}

func generateFile(fullPath string) (*os.File, error) {
	f, err := os.Create(fullPath)
	if err != nil {
		if err == os.ErrExist {
			_ = os.Remove(fullPath)
			f, err = os.Create(fullPath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return f, nil
}
