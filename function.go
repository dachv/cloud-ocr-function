package function

import (
	"context"
	"github.com/dachv/cloud-ocr-function/cloudocr"
	"github.com/dachv/cloud-ocr-function/cloudsql"
	"github.com/dachv/cloud-ocr-function/cloudstorage"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path"
	"strconv"
)

const (
	formTypeFieldId string = "FormType"
	validFormType   string = "0W-8BEN-E"
	processedDir    string = "processed"
	unprocessedDir  string = "unprocessed"
)

var destinationBucket string

func init() {
	destinationBucket = os.Getenv("CLOUD_OCR_DESTINATION_BUCKET")
}

//zip -r cloud-ocr-function.zip ./ -x \*.git\* -x \*.idea\* -x \*/.DS_Store\*
//gsutil rm "gs://cloud-ocr-source/**"
//gsutil -h "x-goog-meta-company-name:TestCompany1" -h "x-goog-meta-app-id:1" cp ~/Desktop/ADF-W8-BEN-E\ 2016-Signed.pdf gs://cloud-ocr-source/
func HandleCloudStorageUpload(ctx context.Context, event cloudstorage.StorageEvent) error {
	log.Printf("StorageEvent: %v\n", event)
	companyName := event.Metadata["company-name"]
	appId := event.Metadata["app-id"]
	objectData := cloudstorage.ReadObjectData(event.Bucket, event.Name)
	ocrResult := cloudocr.PerformOcr(objectData)
	ocrData := ocrResult.GetTextData()
	validDoc := isDocumentValid(ocrData)
	if validDoc {
		_, version := cloudsql.ExistsWith(companyName, validFormType)
		version = version + 1
		objectName := path.Join(companyName, processedDir, appId, strconv.Itoa(version), event.Name)
		cloudstorage.CreateNewObject(destinationBucket, objectName, objectData, ocrData)
		insertId := cloudsql.InsertDocument(objectName, companyName, version, ocrData)
		log.Printf("Document metadata insert id: %v", insertId)
	} else {
		objectName := path.Join(companyName, unprocessedDir, appId, event.Name)
		cloudstorage.CreateNewObject(destinationBucket, objectName, objectData, ocrData)
	}
	return nil
}

func isDocumentValid(ocrData map[string]string) bool {
	for id, value := range ocrData {
		if id == formTypeFieldId && value == validFormType {
			log.Printf("Document is valid: %v", value)
			return true
		}
	}
	log.Printf("Document is invalid")
	return false
}
