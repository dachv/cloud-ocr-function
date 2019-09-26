package function

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/dachv/cloud-ocr-function/cloudocr"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"time"
)

var (
	storageClient *storage.Client
	globalCtx     context.Context
)

func init() {
	globalCtx = context.Background()
	var storageErr error
	storageClient, storageErr = storage.NewClient(globalCtx)
	if storageErr != nil {
		log.Fatalf("Failed to create storage client: %v", storageErr)
	}
}

//zip -r cloud-ocr-function.zip ./ -x \*.git\* -x \*.idea\* -x \*/.DS_Store\*
func HandleCloudStorageUpload(ctx context.Context, event GCSEvent) error {
	logEvent(event)
	objectData := readObjectData(event)
	buffer := bytes.NewBuffer(objectData)
	submitImageResp := cloudocr.SubmitImage(buffer, "test.pdf")
	fmt.Printf("SubmitImage response: %v\n", submitImageResp)
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
	fmt.Printf("ProcessFields response: %v\n", processFieldsResp)
	processFieldsResponse := cloudocr.GetProcessFieldsResponse(processFieldsResp.Task.Id)
	fmt.Printf("ProcessFieldsResponse response: %v", processFieldsResponse)
	return nil
}

func readObjectData(event GCSEvent) []byte {
	objReader, objErr := storageClient.Bucket(event.Bucket).Object(event.Name).NewReader(globalCtx)
	if objErr != nil {
		log.Fatalf("Failed to get reader: %v", objErr)
	}
	defer objReader.Close()
	data, readErr := ioutil.ReadAll(objReader)
	if readErr != nil {
		log.Fatalf("Failed to read object from storage: %v", readErr)
	}
	return data
}

func logEvent(event GCSEvent) {
	log.Printf("Bucket: %v\n", event.Bucket)
	log.Printf("Name: %v\n", event.Name)
	log.Printf("Metageneration: %v\n", event.Metageneration)
	log.Printf("Created: %v\n", event.TimeCreated)
	log.Printf("Updated: %v\n", event.Updated)
}

type GCSEvent struct {
	Id             string    `json:"id"`
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}
