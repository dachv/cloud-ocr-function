package cloudstorage

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"io"
	"io/ioutil"
	"log"
)

var (
	storageClient *storage.Client
	globalCtx     context.Context
)

func init() {
	globalCtx = context.Background()
	var err error
	storageClient, err = storage.NewClient(globalCtx)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
}

func ReadObjectData(bucket string, name string) []byte {
	objReader, err := storageClient.Bucket(bucket).Object(name).NewReader(globalCtx)
	if err != nil {
		log.Fatalf("Failed to get reader: %v", err)
	}
	defer objReader.Close()
	data, err := ioutil.ReadAll(objReader)
	if err != nil {
		log.Fatalf("Failed to read object from storage: %v", err)
	}
	return data
}

func CreateNewObject(bucket string, name string, data []byte, metadata map[string]string) {
	objectHandle := storageClient.Bucket(bucket).Object(name)
	objWriter := objectHandle.NewWriter(globalCtx)
	if metadata != nil {
		objWriter.Metadata = metadata
	}
	defer objWriter.Close()
	_, err := io.Copy(objWriter, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to write object to storage: %v", err)
	}
}
