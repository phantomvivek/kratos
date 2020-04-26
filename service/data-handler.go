package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/phantomvivek/kratos/models"
)

//DataHandler handles all data
type DataHandler struct{}

//GetCSVData reads the data file and populates a string array with the values
func (d *DataHandler) GetCSVData(path string, rows int) ([][]string, error) {

	fileRef, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fileRef.Close()

	counter := 0

	//Assume that socket connections are more in number
	data := make([][]string, 0)

	csrv := csv.NewReader(fileRef)
OuterLoop:
	for {
		record, err := csrv.Read()
		if err != nil {

			//If file is finished reading
			if err == io.EOF {
				return data, nil
			}

			return data, err
		}

		for _, val := range record {
			if val == "" {
				//Value was nil, hence this row cannot be used
				continue OuterLoop
			}
		}

		//Append record to data
		data = append(data, record)

		counter++

		//We don't need more rows than sockets to be opened
		if counter > rows {
			break
		}
	}

	return data, nil
}

//ConstructDataConfig finds the index and CSV column frpom which data needs to be replaced in the test message string
func (d *DataHandler) ConstructDataConfig(message json.RawMessage) []*models.TestDataConfig {

	reg := regexp.MustCompile(`\${(.*?)}`)
	byteArr := reg.FindAll(message, -1)

	if len(byteArr) == 0 {
		//Return
		return nil
	}

	configs := make([]*models.TestDataConfig, 0)

	for _, val := range byteArr {
		//Get index from the byte
		columnStr := strings.TrimLeft(strings.TrimRight(string(val), "}"), "${")
		if columnIdx, err := strconv.Atoi(columnStr); err != nil {
			continue
		} else {
			dataConfig := models.TestDataConfig{
				ColumnIdx: columnIdx,
				TextBytes: val,
			}

			configs = append(configs, &dataConfig)
		}
	}

	return configs
}

//PrepareTestData prepares the test data for each test
func (d *DataHandler) PrepareTestData(file string, connCount int, tests []*models.Test) int {

	//Get file data
	data, err := d.GetCSVData(file, connCount)
	if err != nil {
		panic(err)
	}

	var maxLen int

	for _, test := range tests {
		if test.ReplaceStr {

			jsonMessages := make([]json.RawMessage, 0)

			//Prepare the config first
			configs := d.ConstructDataConfig(test.SendJSON)

			//Iterate over the configs and data to get the final array of messages
			for _, strData := range data {

				msg := test.SendJSON

				for _, config := range configs {
					if len(strData) > config.ColumnIdx {
						msg = bytes.Replace(msg, config.TextBytes, []byte(strData[config.ColumnIdx]), -1)
					}
				}

				jsonMessages = append(jsonMessages, msg)
			}

			test.Data = &models.TestData{
				Counter:   0,
				DataArray: jsonMessages,
			}

			maxLen = len(jsonMessages)
		}
	}

	return maxLen
}
