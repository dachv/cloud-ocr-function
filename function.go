package function

import (
	"context"
	"fmt"
	"github.com/dachv/cloud-ocr-function/cloudocr"
	"github.com/dachv/cloud-ocr-function/cloudsql"
	"github.com/dachv/cloud-ocr-function/cloudstorage"
	_ "github.com/lib/pq"
	"log"
)

const formTypeFieldId string = "FormType"
const validFormType string = "0W-8BEN-E"
const destinationBucket string = "cloud-ocr-destination"

func init() {

}

//zip -r cloud-ocr-function.zip ./ -x \*.git\* -x \*.idea\* -x \*/.DS_Store\*
func HandleCloudStorageUpload(ctx context.Context, event cloudstorage.StorageEvent) error {
	log.Printf("StorageEvent: %v\n", event)
	objectData := cloudstorage.ReadObjectData(event.Bucket, event.Name)
	ocrResult := performOcr(objectData)
	metadata := extractDocumentMetadata(ocrResult)
	if isDocumentValid(metadata) {
		destObjectName := fmt.Sprintf("TestCompany/processed/1/application/1.0/%v", event.Name)
		cloudstorage.CreateNewObject(destinationBucket, destObjectName, objectData, metadata)
		log.Printf("Storage destination object name: %v", destObjectName)
		insertId := cloudsql.InsertDocumentMetadata(metadata)
		log.Printf("Document metadata insert id: %v", insertId)

	}
	return nil
}

func performOcr(documentData []byte) *cloudocr.ProcessFieldsResponse {
	submitImageResp := cloudocr.SubmitImage(documentData, "document.pdf")
	log.Printf("SubmitImage response: %v\n", submitImageResp)
	processFieldsReq := cloudocr.ProcessFieldsRequest{
		Xmlns:          "http://ocrsdk.com/schema/taskDescription-1.0.xsd",
		FieldTemplates: cloudocr.ReqFieldTemplates{},
		Page: []cloudocr.ReqPage{{
			ApplyTo:   "0",
			Text:      cloudocr.OcrTextFields,
			Barcode:   []cloudocr.ReqBarcode{},
			Checkmark: []cloudocr.ReqCheckmark{},
		}},
	}
	processFieldsResp := cloudocr.ProcessFields(submitImageResp.Task.Id, processFieldsReq)
	log.Printf("ProcessFields response: %v\n", processFieldsResp)
	processFieldsResponse := cloudocr.GetProcessFieldsResponse(processFieldsResp.Task.Id)
	log.Printf("ProcessFieldsResponse response: %v", processFieldsResponse)
	return processFieldsResponse
}

func extractDocumentMetadata(response *cloudocr.ProcessFieldsResponse) map[string]string {
	metadata := make(map[string]string)
	for _, page := range response.Page {
		for _, text := range page.Text {
			metadata[text.Id] = text.Value
		}
	}
	log.Printf("Document metadata: %v", metadata)
	return metadata
}

func isDocumentValid(metadata map[string]string) bool {
	for id, value := range metadata {
		if id == formTypeFieldId && value == validFormType {
			log.Printf("Document is valid: %v", value)
			return true
		}
	}
	log.Printf("Document is invalid")
	return false
}
