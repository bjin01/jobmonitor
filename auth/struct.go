package auth

import "encoding/xml"

type Suma_API_Error_MethodResponse struct {
	XMLName xml.Name `xml:"methodResponse"`
	Text    string   `xml:",chardata"`
	Fault   struct {
		Text  string `xml:",chardata"`
		Value struct {
			Text   string `xml:",chardata"`
			Struct struct {
				Text   string `xml:",chardata"`
				Member []struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name"`
					Value struct {
						Text   string `xml:",chardata"`
						Int    string `xml:"int"`
						String string `xml:"string"`
					} `xml:"value"`
				} `xml:"member"`
			} `xml:"struct"`
		} `xml:"value"`
	} `xml:"fault"`
}

type Fault struct {
	FaultCode   int    `xml:"fault>value>struct>member>value>int"`
	FaultString string `xml:"fault>value>struct>member>value>string"`
}

/* type Suma_API_Error_MethodResponse struct {
	Fault []Error_String `xml:"methodResponse>fault>value>struct>member"`
}

type Error_String struct {
	Name  string      `xml:"name"`
	Value interface{} `xml:"value"`
} */
