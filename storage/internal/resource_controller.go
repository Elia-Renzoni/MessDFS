package internal

import (
	"encoding/csv"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type ResourceController struct {
	mutex sync.Mutex
}

type PairChecker struct {
	columnName  string
	columnIndex int
}

var mainDir string

// change the worker dir to path specified
type changeWorkDir func(string) error

// goes back to renzofs main dir
type backToHomeDir func() error

// crud operations
const (
	insert string = "insert"
	update string = "update"
	delete string = "delete"
	read   string = "read"
)

func (r *ResourceController) CreateNewDir(dirname string) error {
	if err := os.Mkdir(filepath.Join(mainDir, "files", dirname), 0750); err != nil {
		return err
	}
	return nil
}

func (r *ResourceController) DeleteDir(dirname string) error {
	if err := os.Remove(filepath.Join(mainDir, "files", dirname)); err != nil {
		return err
	}

	return nil
}

func (r *ResourceController) DeleteFile(dirname, filename string) error {
	if err := os.Remove(filepath.Join("files", dir, filename)); err != nil {
		return err
	}

	return nil
}

func (r *ResourceController) WriteRemoteCSV(dir, filename, queryType string, query []string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if queryType != insert {
		return errors.New("Invalid Crud Operation")
	}

	file, err := os.OpenFile(filepath.Join("files", dir, filename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // defer-close pattern
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(query); err != nil {
		return err
	}

	return nil
}

func (r *ResourceController) ReadInRemoteCSV(dir, filename, queryType string, query url.Values) (map[string]string, error) {
	var (
		storeControlResult  []PairChecker     = make([]PairChecker, 0)
		idToSearch          string            = query.Get("id")
		readOperationResult map[string]string = make(map[string]string)
	)

	if queryType != read {
		return nil, errors.New("Invalid Crud Operation")
	}

	file, err := os.OpenFile(filepath.Join("files", dir, filename), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	fileContent, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for rows, partialFileContent := range fileContent {
		if rows == 0 {
			for index := range partialFileContent {
				if index > 0 {
					storeControlResult = append(storeControlResult, PairChecker{
						columnName:  partialFileContent[index],
						columnIndex: index,
					})
				}
			}
		} else {
			for index := range partialFileContent {
				if index == 0 {
					if partialFileContent[index] != idToSearch {
						break
					}
				}

				for _, value := range storeControlResult {
					if value.columnIndex == index {
						readOperationResult[value.columnName] = partialFileContent[index]
						break
					}
				}
			}
		}
	}

	return readOperationResult, nil
}

func (r *ResourceController) UpdateRemoteCSV(dir, filename, queryType string, query map[string][]string) error {
	var (
		storeControlResults []PairChecker = make([]PairChecker, 0)
		idList              []string      = make([]string, 0)
	)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if queryType != update {
		return errors.New("Invalid Crud Operation")
	}

	file, err := os.OpenFile(filepath.Join("files", dir, filename), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	fileContent, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// parse the file content and change file's informations
	for row, partialFileContent := range fileContent {
		if row == 0 {
			for index := range partialFileContent {
				for key := range query {
					if partialFileContent[index] == key {
						storeControlResults = append(storeControlResults, PairChecker{
							columnName:  key,
							columnIndex: index,
						})
					}
				}
			}

			// extraction
			for _, value := range query {
				idList = append(idList, value[0])
			}
		} else {
			for index := range partialFileContent {
				var exit bool
				if index == 0 {
					for indexIdList := range idList {
						if idList[indexIdList] == partialFileContent[index] {
							exit = true
						}
					}
					if !exit {
						break
					}
				} else {
					for _, value := range storeControlResults {
						if index == value.columnIndex {
							if queryValue, ok := query[value.columnName]; ok {
								if queryValue[1] == partialFileContent[index] {
									// now i change the value
									partialFileContent[index] = queryValue[2]
								}
							}
						}
					}
				}
			}
		}
	}

	// file truncation to show updates without redundancy
	if trErr := os.Truncate(filepath.Join("files", dir, filename), 0); trErr != nil {
		return trErr
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(fileContent); err != nil {
		return err
	}

	return nil
}

func (r *ResourceController) DeleteRemoteCSV(dir, filename, queryType string, query url.Values) error {
	var (
		emptyField          string        = "/"
		storeControlResults []PairChecker = make([]PairChecker, 0)
	)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if queryType != delete {
		return errors.New("Invalid Crud Operation")
	}

	file, err := os.OpenFile(filepath.Join("files", dir, filename), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	fileContent, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// delete informations algorithm
	for row, partialFileContent := range fileContent {
		if row == 0 {
			for index := range partialFileContent {
				column := query["column"]
				effectiveColumn := column[0]
				if partialFileContent[index] == strings.ToUpper(effectiveColumn[:0]) {
					storeControlResults = append(storeControlResults, PairChecker{
						columnName:  strings.ToUpper(effectiveColumn[:0]),
						columnIndex: index,
					})
				}
			}
		} else {
			for index := range partialFileContent {
				if index == 0 {
					var exit bool
					id := query["id"]
					if convertedId, _ := strconv.Atoi(id[0]); convertedId == index {
						exit = true
					}

					if !exit {
						break
					}
				} else {
					for _, value := range storeControlResults {
						if value.columnIndex == index {
							partialFileContent[index] = emptyField
						}
					}
				}

			}
		}
	}

	// file truncation to avoid redundancy
	if trErr := os.Truncate(filepath.Join("files", dir, filename), 0); trErr != nil {
		return trErr
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(fileContent); err != nil {
		return err
	}

	return nil
}
