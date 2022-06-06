package esfixture

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
)

type SchemaInfoFromFile struct {
	ProvidedName    string
	TargetName      string
	IsFromESFixture bool
	JSONFromFile    []byte
}

func getSchemaInfoFromFile(fullPath string) (returnInfo *SchemaInfoFromFile, returnError error) {
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		returnError = err
		return
	}
	schemaMap := make(map[string]interface{})
	err = json.Unmarshal(b, &schemaMap)
	if err != nil {
		returnError = err
		return
	}

	mappings, ok := schemaMap["mappings"]
	if !ok || mappings == nil {
		returnError = errors.New("(mappings) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'")
		return
	}

	mappingsMap, ok := mappings.(map[string]interface{})
	if !ok {
		returnError = errors.New("(mappings) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'")
		return
	}

	meta, ok := mappingsMap["_meta"]
	if !ok || meta == nil {
		returnError = errors.New("(mappings._meta) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'")
		return
	}

	mataMap, ok := meta.(map[string]interface{})
	if !ok {
		returnError = errors.New("(mappings._meta) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'")
		return
	}

	from, ok := mataMap[metaKeyFrom]
	if !ok || from == nil {
		returnError = errors.New(fmt.Sprintf("(mappings._meta.%s) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'", metaKeyFrom))
		return
	}
	fromStr, ok := from.(string)
	if !ok {
		returnError = errors.New(fmt.Sprintf("(mappings._meta.%s) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'", metaKeyFrom))
		return
	}
	if fromStr != metaValueFrom {
		returnError = errors.New(fmt.Sprintf("(mappings._meta.%s) doesn't match value %s. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'", metaKeyFrom, metaValueFrom))
		return
	}

	provideName, ok := mataMap[metaKeyProvidedName]
	if !ok || provideName == nil {
		returnError = errors.New(fmt.Sprintf("(mappings._meta.%s) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'", metaKeyProvidedName))
		return
	}
	provideNameStr, ok := provideName.(string)
	if !ok {
		returnError = errors.New(fmt.Sprintf("(mappings._meta.%s) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'", metaKeyProvidedName))
		return
	}

	targetName, ok := mataMap[metaKeyGetTargetName]
	if !ok || targetName == nil {
		returnError = errors.New(fmt.Sprintf("(mappings._meta.%s) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'", metaKeyGetTargetName))
		return
	}
	targetNameStr, ok := targetName.(string)
	if !ok {
		returnError = errors.New(fmt.Sprintf("(mappings._meta.%s) doesn't have origin index name. Maybe It's not generated 'esfixture' golang package. Only accepted generated by 'exfixture'", metaKeyGetTargetName))
		return
	}
	returnInfo = &SchemaInfoFromFile{
		ProvidedName:    provideNameStr,
		TargetName:      targetNameStr,
		IsFromESFixture: true,
		JSONFromFile:    b,
	}
	return
}
func getBytesFromHit(hit *elastic.SearchHit) ([]byte, []byte, error) {
	b, err := json.Marshal(hit)
	if err != nil {
		return nil, nil, err
	}
	hitMap := make(map[string]interface{})
	err = json.Unmarshal(b, &hitMap)
	if err != nil {
		return nil, nil, err
	}

	action := make(map[string]interface{})
	action["index"] = map[string]interface{}{
		"_index": hitMap["_index"],
		"_id":    hitMap["_id"],
	}
	actionBytes, subErr := json.Marshal(&action)
	if subErr != nil {
		return nil, nil, subErr
	}
	sourceBytes, subErr := json.Marshal(hitMap["_source"])
	if subErr != nil {
		return nil, nil, subErr
	}
	return actionBytes, sourceBytes, nil
}

func getBytesToStoreFromIndex(targetName, providedName string, index *elastic.IndicesGetResponse) ([]byte, error) {
	settingsIndex, ok := index.Settings["index"]
	if !ok || settingsIndex == nil {
		return nil, errors.New(targetName + " has not settings.index property at schema json")
	}

	settingsIndexsMap, ok := settingsIndex.(map[string]interface{})
	if !ok {
		return nil, errors.New(targetName + " has not settings.index property at schema json")
	}

	delete(settingsIndexsMap, "creation_date")
	delete(settingsIndexsMap, "provided_name")
	delete(settingsIndexsMap, "uuid")
	delete(settingsIndexsMap, "version")

	indexMappingsMeta, ok := index.Mappings["_meta"]
	if !ok || indexMappingsMeta == nil {
		index.Mappings["_meta"] = map[string]interface{}{
			metaKeyFrom:          metaValueFrom,
			metaKeyProvidedName:  providedName,
			metaKeyGetTargetName: targetName,
		}
	} else {
		if indexMappingsMetaMap, subOk := indexMappingsMeta.(map[string]interface{}); subOk {
			indexMappingsMetaMap[metaKeyFrom] = metaValueFrom
			indexMappingsMetaMap[metaKeyProvidedName] = providedName
			indexMappingsMetaMap[metaKeyGetTargetName] = targetName
			index.Mappings["_meta"] = indexMappingsMetaMap
		}
	}

	b, err := json.Marshal(index)
	if err != nil {
		return nil, err
	}
	indexMapTemp := make(map[string]interface{})
	err = json.Unmarshal(b, &indexMapTemp)
	if err != nil {
		return nil, err
	}
	delete(indexMapTemp, "routing")
	delete(indexMapTemp, "warmers")
	//delete(indexMapTemp, "aliases")

	b, subErr := json.MarshalIndent(&indexMapTemp, "", "\t")
	if subErr != nil {
		return nil, subErr
	}
	return b, nil
}

func isAbleToDelete(index *elastic.IndicesGetResponse) error {
	settingsIndex, ok := index.Settings["index"]
	if !ok || settingsIndex == nil {
		return fmt.Errorf("doesn't have settings.index")
	}
	settingsIndexMap, ok := settingsIndex.(map[string]interface{})
	if !ok {
		return fmt.Errorf("doesn't have settings.index")
	}

	actualProvidedName, ok := settingsIndexMap["provided_name"]
	if !ok || actualProvidedName == nil {
		return fmt.Errorf("doesn't have settings.index.provided_name")
	}
	actualProvidedNameStr, ok := actualProvidedName.(string)
	if !ok {
		return fmt.Errorf("doesn't have settings.index.provided_name")
	}

	mappingsMeta, ok := index.Mappings["_meta"]
	if !ok || mappingsMeta == nil {
		return fmt.Errorf("doesn't have mappings._meta")
	}
	mappingsMetaMap, ok := mappingsMeta.(map[string]interface{})
	if !ok {
		return fmt.Errorf("doesn't have mappings._meta")
	}

	from, ok := mappingsMetaMap[metaKeyFrom]
	if !ok || from == nil {
		return fmt.Errorf("doesn't have mappings._meta.%s", metaKeyFrom)
	}
	fromStr, ok := from.(string)
	if !ok || from == nil {
		return fmt.Errorf("doesn't have mappings._meta.%s", metaKeyFrom)
	}
	if fromStr != metaValueFrom {
		return fmt.Errorf("it's weired! It's different [expect] %s and [actual] %s have mappings._meta.%s", metaValueFrom, fromStr, metaKeyFrom)
	}

	providedName, ok := mappingsMetaMap[metaKeyProvidedName]
	if !ok || providedName == nil {
		return fmt.Errorf("doesn't have mappings._meta.%s", metaKeyProvidedName)
	}
	providedNameStr, ok := providedName.(string)
	if !ok || from == nil {
		return fmt.Errorf("doesn't have mappings._meta.%s", metaKeyProvidedName)
	}
	if providedNameStr != actualProvidedNameStr {
		return fmt.Errorf("it's weired! It's different [expect] %s and [actual] %s have mappings._meta.%s", providedNameStr, actualProvidedNameStr, metaKeyProvidedName)
	}

	targetName, ok := mappingsMetaMap[metaKeyGetTargetName]
	if !ok || targetName == nil {
		return fmt.Errorf("doesn't have mappings._meta.%s", metaKeyGetTargetName)
	}
	_, ok = targetName.(string)
	if !ok || from == nil {
		return fmt.Errorf("doesn't have mappings._meta.%s", metaKeyGetTargetName)
	}

	return nil
}

func genProvidedNameAndTargetNameMapUsingElasticGetIndices(targetNames []string, indicesResponse map[string]*elastic.IndicesGetResponse) map[string]string {
	var providedNameAndTargetNameMap = map[string]string{}
	for providedName, index := range indicesResponse {
		for alias, _ := range index.Aliases {
			if ok := containStringArray(targetNames, alias); ok {
				providedNameAndTargetNameMap[providedName] = alias
				continue
			}
		}
		if ok := containStringArray(targetNames, providedName); ok {
			providedNameAndTargetNameMap[providedName] = providedName
		}
	}
	return providedNameAndTargetNameMap
}