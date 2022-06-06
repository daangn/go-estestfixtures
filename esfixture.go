package esfixture

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"net/http"
	"os"
	"strings"
)

type Loader struct {
	dsn        string
	es         *elastic.Client
	limit      int
	dir        string
	searchFunc func(c *elastic.Client, indexes []string) *elastic.SearchService

	targetNames []string
}

func NewLoader(ctx context.Context, dsn string, options ...func(*Loader) error) (*Loader, error) {
	d := &Loader{}

	if dsn == "" {
		return nil, errors.New("need elasticsearch host dsn")
	}
	es, err := elastic.NewClient(
		elastic.SetURL(dsn),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetGzip(true),
	)
	if err != nil {
		return nil, err
	}
	d.dsn = dsn
	d.es = es

	for _, option := range options {
		if err = option(d); err != nil {
			return nil, err
		}
	}
	if len(d.targetNames) == 0 {
		return nil, fmt.Errorf("you didn't define targetNames. use WithTargetNames(...someOfIndice)")
	}
	if d.limit == 0 {
		d.limit = defaultLimit
	}
	if d.dir == "" {
		d.dir = defaultDir
	}

	if d.searchFunc == nil {
		d.searchFunc = func(c *elastic.Client, targetNames []string) *elastic.SearchService {
			return c.Search(targetNames...).From(0).Size(d.limit)
		}
	}
	if err := d.initFixturesSettings(ctx); err != nil {
		return d, err
	}
	return d, nil
}

// WithDirectory sets the directory where the fixtures files will be created.
func WithDirectory(dir string) func(*Loader) error {
	return func(d *Loader) error {
		d.dir = dir
		return nil
	}
}

func WithSearchFunc(searchFunc func(c *elastic.Client, targetNames []string) *elastic.SearchService) func(*Loader) error {
	return func(d *Loader) error {
		d.searchFunc = searchFunc
		return nil
	}
}

// WithLimit sets the limit to get document size, if SearchFunc is set, limit property ignored
func WithLimit(limit int) func(*Loader) error {
	return func(d *Loader) error {
		d.limit = limit
		return nil
	}
}

// WithTargetNames sets the directory where the fixtures files will be created. It doesn't care aliasing or not.
func WithTargetNames(targetNames ...string) func(*Loader) error {
	return func(d *Loader) error {
		d.targetNames = targetNames
		return nil
	}
}

func (d *Loader) Load(ctx context.Context) error {
	schemaFileNames, documentFileNames, err := getFixtureFiles(d.dir)
	if err != nil {
		return err
	}

	for _, fullPath := range schemaFileNames {
		info, subErr := getSchemaInfoFromFile(fullPath)
		if subErr != nil {
			return subErr
		}
		_, subErr = d.es.CreateIndex(info.ProvidedName).BodyString(string(info.JSONFromFile)).Do(ctx)
		if subErr != nil {
			return subErr
		}
	}

	for _, fullPath := range documentFileNames {
		bulk, subErr := readDocumentFile(fullPath)
		if subErr != nil {
			return subErr
		}
		_, subErr = http.Post(fmt.Sprintf("%s/_bulk", d.dsn), "application/x-ndjson", bytes.NewBuffer(bulk))
		if subErr != nil {
			return subErr
		}
	}

	return nil
}

// Dump get documents from you injected DSN. those are stored at you injected directory path.
func (d *Loader) Dump(ctx context.Context) error {
	// Step1, check to exist defined 'd.targetNames' indices at elasticsearch.
	indicesResponse, err := d.es.IndexGet(d.targetNames...).Do(ctx)
	if err != nil {
		return err
	}
	var providedNameAndTargetName = genProvidedNameAndTargetNameMapUsingElasticGetIndices(d.targetNames, indicesResponse)

	//
	doesntExistIndexNames := make([]string, len(d.targetNames))
	copy(doesntExistIndexNames, d.targetNames)

	for _, tgName := range providedNameAndTargetName {
		var deleteTargetIndex = -1
		for i, v := range doesntExistIndexNames {
			if v == tgName {
				deleteTargetIndex = i
				continue
			}
		}
		if deleteTargetIndex != -1 {
			doesntExistIndexNames = append(doesntExistIndexNames[:deleteTargetIndex], doesntExistIndexNames[deleteTargetIndex+1:]...)
		}
	}
	if len(doesntExistIndexNames) > 0 {
		return errors.New("failure loadSchema, dosen't exist [" + strings.Join(doesntExistIndexNames, ",") + "]")
	}

	// Step2, Creating Schema files
	for providedName, index := range indicesResponse {
		targetName := providedNameAndTargetName[providedName]
		b, subErr := getBytesToStoreFromIndex(targetName, providedName, index)
		if subErr != nil {
			return subErr
		}

		fullPath := genFullPathWithFileNameForSchema(d.dir, providedName)
		f, subErr := generateFile(fullPath)
		if subErr != nil {
			return subErr
		}

		buf := bufio.NewWriter(f)
		if _, err = buf.Write(b); err != nil {
			_ = f.Close()
			_ = os.Remove(fullPath)
			return err
		}
		_ = buf.Flush()
		_ = f.Close()
	}

	// Step3, Generating documents files with 'ndjson' format
	for providedName, _ := range indicesResponse {
		r, subErr := d.searchFunc(d.es, d.targetNames).Do(ctx)
		if subErr != nil {
			return subErr
		}

		fullPath := genFullPathWithFileNameForDocument(d.dir, providedName)
		f, subErr := generateFile(fullPath)
		if subErr != nil {
			return subErr
		}

		buf := bufio.NewWriter(f)
		for _, hit := range r.Hits.Hits {
			actionBytes, sourceBytes, subErr := getBytesFromHit(hit)
			if subErr != nil {
				return subErr
			}
			if _, subErr = buf.Write(actionBytes); subErr != nil {
				_ = f.Close()
				_ = os.Remove(fullPath)
				return subErr
			}
			_, _ = buf.Write([]byte("\n"))

			if _, subErr = buf.Write(sourceBytes); subErr != nil {
				_ = f.Close()
				_ = os.Remove(fullPath)
				return subErr
			}
			_, _ = buf.Write([]byte("\n"))
		}
		_ = buf.Flush()
		_ = f.Close()
	}
	return nil
}

// initFixturesSettings ready for loading and dumping process at machine
func (d *Loader) initFixturesSettings(ctx context.Context) error {
	// create dir if not exist
	_, err := os.Stat(d.dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(d.dir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// check health about es
	h, err := d.es.ClusterHealth().Do(ctx)
	if err != nil {
		return err
	}
	if h.Status != clusterHealthStatusValueGreen {
		return errors.New(fmt.Sprintf("cluster state is not %s, check your es cluster", clusterHealthStatusValueGreen))
	}
	return nil
}

// Deprecated: kr: os.Create 가 알아서 Overwrite 해주도록 만듦. 굳이 해당 메서드를 쓸 필요 없음
// Deprecated: en: os.Create can erase contents on target file. So, this method unusable stuff.
//func (d *Loader) clearFile() error {
//	// Step1. Read schema files at fixturedata directory
//	schemaFileFullPaths, _, err := getFixtureFiles(d.dir)
//	if err != nil {
//		return err
//	}
//
//	var deleteAbleFiles []string
//	// Step2. Validate to exist schema files following loader target names
//	for _, fullPath := range schemaFileFullPaths {
//		info, err := getSchemaInfoFromFile(fullPath)
//		if err != nil {
//			return fmt.Errorf("%w at '%s'", err, fullPath)
//		}
//		if ok := containStringArray(d.targetNames, info.TargetName); ok {
//			deleteAbleFiles = append(deleteAbleFiles, fullPath)
//		}
//	}
//	for _, fullPath := range deleteAbleFiles {
//		_ = os.Remove(fullPath)
//	}
//	return nil
//}

func (d *Loader) ClearElasticsearch(ctx context.Context) error {
	// Step1. get provided_name from elasticsearch
	indicesResponse, err := d.es.IndexGet(d.targetNames...).Do(ctx)
	if err != nil {
		return err
	}

	providedNameAndTargetNameMap := genProvidedNameAndTargetNameMapUsingElasticGetIndices(d.targetNames, indicesResponse)

	// Step2. validating able to delete each index
	var deleteAbleProvidedName []string
	for providedName, index := range indicesResponse {
		targetName := providedNameAndTargetNameMap[providedName]
		subErr := isAbleToDelete(index)
		if subErr != nil {
			return fmt.Errorf(
				"%v\ntarget_pame[%s] is not generated from esfixture you can check the '_meta' %s/%s",
				targetName,
				d.dsn,
				providedName,
				subErr,
			)
		}
		deleteAbleProvidedName = append(deleteAbleProvidedName, providedName)
	}

	_, err = d.es.DeleteIndex(deleteAbleProvidedName...).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
