package repository


import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
	"be/models"

	
)

var xataAPIKey = "xau_tDGvSkzW75qVQEdZhF1ETjXTZYVxrejT1"
var baseURL = "https://Mesach-Harmasendro-s-workspace-3c05vf.ap-southeast-2.xata.sh/db/tubes-stima"

func createRequest(method, url string, bodyData *bytes.Buffer) (*http.Request, error) {
    var req *http.Request
    var err error

    if method == "GET" || method == "DELETE" {
        req, err = http.NewRequest(method, url, nil)
    } else {
        req, err = http.NewRequest(method, url, bodyData)
    }

    if err != nil {
        return nil, err
    }

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", xataAPIKey))
    return req, nil
}

func makeRequest(req *http.Request, target interface{}) error {
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if target != nil {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return err
        }
        err = json.Unmarshal(body, target)
        if err != nil {
            return err
        }
    }
    return nil
}

func CreateData(newData *models.DataRequest) (*models.DataResponse, error){
	createData := models.DataResponse{}

	jsonData := models.Data{
		Parent: newData.Parent,
		Children: newData.Children,
	}

	postBody, _ := json.Marshal(jsonData)
	bodyData := bytes.NewBuffer(postBody)

	// fmt.Println(bodyData)

	fullURL := fmt.Sprintf("%s:main/tables/Article/data", baseURL)
	req, err := createRequest("POST", fullURL, bodyData)

	if err != nil {
		return nil, err
	}

	err = makeRequest(req, &createData)

	if err != nil {
		return nil, err
	}

	return &createData, nil
}

func GetData(parent string) (*models.Data, error) {
	details := models.Data{}

	jsonData := models.DataQuery{
		Columns: []string{"parent", parent},
	}

	postBody, _ := json.Marshal(jsonData)
	bodyData := bytes.NewBuffer(postBody)


	fullURL := fmt.Sprintf("%s:main/tables/Article/query", baseURL)

	req, err := createRequest("GET", fullURL, bodyData)

	if err != nil {
		return nil, err
	}

	err = makeRequest(req, &details)

	if err != nil {
        return nil, err
    }

	return &details, nil

}

func CheckData(parent string) bool {
	_, err := GetData(parent)

	if err != nil {
		return false
	}

	return true
}