package cloudocr

import "encoding/xml"

type TaskStatusResponse struct {
	XMLName xml.Name `xml:"response"`
	Task    Task     `xml:"task"`
}

type Task struct {
	Id                      string `xml:"id,attr"`
	Status                  string `xml:"status,attr"`
	RegistrationTime        string `xml:"registrationTime,attr"`
	StatusChangeTime        string `xml:"statusChangeTime,attr"`
	FilesCount              int    `xml:"filesCount,attr"`
	Credits                 int    `xml:"credits,attr"`
	EstimatedProcessingTime int    `xml:"estimatedProcessingTime,attr"`
	ResultUrl               string `xml:"resultUrl,attr"`
	Error                   string `xml:"error,attr"`
	Description             string `xml:"description,attr"`
}

type ProcessFieldsResponse struct {
	XMLName        xml.Name `xml:"document"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Page           []struct {
		Index string `xml:"index,attr"`
		Text  []struct {
			Id     string `xml:"id,attr"`
			Left   int    `xml:"left,attr"`
			Top    int    `xml:"top,attr"`
			Right  int    `xml:"right,attr"`
			Bottom int    `xml:"bottom,attr"`
			Value  string `xml:"value"`
			/*Line   []struct {
				Left   int `xml:"left,attr"`
				Top    int `xml:"top,attr"`
				Right  int `xml:"right,attr"`
				Bottom int `xml:"bottom,attr"`
				Char   []struct {
					Left       int    `xml:"left,attr"`
					Top        int    `xml:"top,attr"`
					Right      int    `xml:"right,attr"`
					Bottom     int    `xml:"bottom,attr"`
					Suspicious bool   `xml:"suspicious,attr"`
					Value      string `xml:",chardata"`
				} `xml:"char"`
			} `xml:"line"`*/
		} `xml:"text"`
		Checkmark []struct {
			Id     string `xml:"id,attr"`
			Left   int    `xml:"left,attr"`
			Top    int    `xml:"top,attr"`
			Right  int    `xml:"right,attr"`
			Bottom int    `xml:"bottom,attr"`
			Value  string `xml:"value"`
		} `xml:"checkmark"`
		Barcode []struct {
			Id     string `xml:"id,attr"`
			Left   int    `xml:"left,attr"`
			Top    int    `xml:"top,attr"`
			Right  int    `xml:"right,attr"`
			Bottom int    `xml:"bottom,attr"`
			Value  string `xml:"value"`
		} `xml:"barcode"`
	} `xml:"page"`
}

type ProcessFieldsRequest struct {
	XMLName        xml.Name          `xml:"document"`
	Xmlns          string            `xml:"xmlns,attr"`
	FieldTemplates ReqFieldTemplates `xml:"fieldTemplates"`
	Page           []ReqPage         `xml:"page"`
}

type ReqFieldTemplates struct {
	Text      []ReqTextTpl      `xml:"text"`
	Barcode   []ReqBarcodeTpl   `xml:"barcode"`
	Checkmark []ReqCheckmarkTpl `xml:"checkmark"`
}

type ReqTextTpl struct {
	Id          string `xml:"id,attr"`
	Bottom      int    `xml:"bottom,attr"`
	Left        int    `xml:"left,attr"`
	Right       int    `xml:"right,attr"`
	Top         int    `xml:"top,attr"`
	Language    string `xml:"language"`
	TextType    string `xml:"textType"`
	OneTextLine string `xml:"oneTextLine"`
}

type ReqBarcodeTpl struct {
	Bottom             string `xml:"bottom,attr"`
	Id                 string `xml:"id,attr"`
	Type               string `xml:"type"`
	ContainsBinaryData string `xml:"containsBinaryData"`
}

type ReqCheckmarkTpl struct {
	Id                string `xml:"id,attr"`
	Type              string `xml:"type"`
	CorrectionAllowed string `xml:"correctionAllowed"`
}

type ReqPage struct {
	ApplyTo   string         `xml:"applyTo,attr"`
	Text      []ReqText      `xml:"text"`
	Barcode   []ReqBarcode   `xml:"barcode"`
	Checkmark []ReqCheckmark `xml:"checkmark"`
}

type ReqText struct {
	Template string `xml:"template,attr"`
	Id       string `xml:"id,attr"`
	Left     int    `xml:"left,attr"`
	Top      int    `xml:"top,attr"`
	Right    int    `xml:"right,attr"`
	Bottom   int    `xml:"bottom,attr"`
	Language string `xml:"language"`
}

type ReqBarcode struct {
	Template string `xml:"template,attr"`
	Id       string `xml:"id,attr"`
	Left     int    `xml:"left,attr"`
	Top      int    `xml:"top,attr"`
	Right    int    `xml:"right,attr"`
	Bottom   int    `xml:"bottom,attr"`
	Type     string `xml:"type"`
}

type ReqCheckmark struct {
	Template string `xml:"template,attr"`
	Id       string `xml:"id,attr"`
	Left     int    `xml:"left,attr"`
	Top      int    `xml:"top,attr"`
	Right    int    `xml:"right,attr"`
	Bottom   int    `xml:"bottom,attr"`
}
