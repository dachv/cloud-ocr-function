package main

import (
	"bytes"
	"fmt"
	"github.com/dachv/cloud-ocr-function/cloudocr"
	"io/ioutil"
)

func main() {
	data, _ := ioutil.ReadFile("/Users/vitaliy.dach/Desktop/ADF-W8-BEN-E 2016-Signed.pdf")
	buffer := bytes.NewBuffer(data)
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
}
