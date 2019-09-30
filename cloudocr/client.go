package cloudocr

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var client *http.Client

var (
	apiUrl        = os.Getenv("CLOUD_OCR_API_URL")
	applicationId = os.Getenv("CLOUD_OCR_APPLICATION_ID")
	password      = os.Getenv("CLOUD_OCR_PASSWORD")
)

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 1,
		},
	}
}

func PerformOcr(documentData []byte) *ProcessFieldsResponse {
	submitImageResp := SubmitImage(documentData, "document.pdf")
	log.Printf("SubmitImage response: %v\n", submitImageResp)
	processFieldsReq := ProcessFieldsRequest{
		Xmlns:          "http://ocrsdk.com/schema/taskDescription-1.0.xsd",
		FieldTemplates: ReqFieldTemplates{},
		Page: []ReqPage{{
			ApplyTo:   "0",
			Text:      OcrTextFields,
			Barcode:   []ReqBarcode{},
			Checkmark: []ReqCheckmark{},
		}},
	}
	processFieldsResp := ProcessFields(submitImageResp.Task.Id, processFieldsReq)
	log.Printf("ProcessFields response: %v\n", processFieldsResp)
	processFieldsResponse := GetProcessFieldsResponse(processFieldsResp.Task.Id)
	log.Printf("ProcessFieldsResponse response: %v", processFieldsResponse)
	return processFieldsResponse
}

func SubmitImage(data []byte, name string) *TaskStatusResponse {
	buffer := bytes.NewBuffer(data)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filePart, _ := writer.CreateFormFile("upload", name)
	_, _ = io.Copy(filePart, buffer)
	_ = writer.Close()
	req := createHttpRequest("POST", apiUrl+"/submitImage", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	taskStatusResponse := &TaskStatusResponse{}
	executeAndGetResponse(req, taskStatusResponse)
	return taskStatusResponse
}

func ProcessFields(taskId string, processReq ProcessFieldsRequest) *TaskStatusResponse {
	xmlData, marshalErr := xml.MarshalIndent(processReq, " ", "  ")
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}
	body := bytes.NewBuffer(xmlData)
	req := createHttpRequest("POST", apiUrl+"/processFields?taskId="+taskId, body)
	taskStatusResponse := &TaskStatusResponse{}
	executeAndGetResponse(req, taskStatusResponse)
	return taskStatusResponse
}

func GetTaskStatus(taskId string) *TaskStatusResponse {
	req := createHttpRequest("GET", apiUrl+"/getTaskStatus?taskId="+taskId, nil)
	taskStatusResponse := &TaskStatusResponse{}
	executeAndGetResponse(req, taskStatusResponse)
	return taskStatusResponse
}

func GetProcessFieldsResponse(taskId string) *ProcessFieldsResponse {
	checkLim := 5
	processFieldsResponse := &ProcessFieldsResponse{}
	for attempt := 0; attempt < checkLim; attempt++ {
		taskStatusResponse := GetTaskStatus(taskId)
		switch taskStatusResponse.Task.Status {
		case "Completed":
			resultUrl := taskStatusResponse.Task.ResultUrl
			req := createHttpRequest("GET", resultUrl, nil)
			req.Header.Del("Authorization")
			executeAndGetResponse(req, processFieldsResponse)
			return processFieldsResponse
		case "Submitted", "Queued", "InProgress":
			time.Sleep(1 * time.Second)
		default:
			log.Fatal("Unexpected task status got for task: ", taskId)
		}
	}
	log.Fatalf("Status check limit %v exceeded for GetProcessFieldsResponse, taskId: %s", checkLim, taskId)
	return nil
}

func (response *ProcessFieldsResponse) GetTextData() map[string]string {
	metadata := make(map[string]string)
	for _, page := range response.Page {
		for _, text := range page.Text {
			metadata[text.Id] = text.Value
		}
	}
	log.Printf("Text data: %v", metadata)
	return metadata
}

func executeAndGetResponse(req *http.Request, model interface{}) {
	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}
	if res.StatusCode != 200 {
		log.Fatal("Non 200 http status returned in response: ", res.StatusCode)
	}
	respData, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	xmlErr := xml.Unmarshal(respData, model)
	if xmlErr != nil {
		log.Fatal(xmlErr)
	}
}

func createHttpRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(applicationId, password)
	return req
}
