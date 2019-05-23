package models

// Envelope is the soap wrapper for all requests
type Envelope struct {
	XMLName   string      `xml:"soapenv:Envelope"`
	Body      interface{} `xml:"soapenv:Body"`
	Soapenv   string      `xml:"xmlns:soapenv,attr"`
	Namespace string      `xml:"xmlns:q0,attr"`
}

// Request has just the default auth fields on all requests
type Request struct {
	WebAuthenticationDetail WebAuthenticationDetail `xml:"q0:WebAuthenticationDetail"`
	ClientDetail            ClientDetail            `xml:"q0:ClientDetail"`
	TransactionDetail       *TransactionDetail      `xml:"q0:TransactionDetail,omitempty"`
	Version                 Version                 `xml:"q0:Version"`
}

type WebAuthenticationDetail struct {
	UserCredential UserCredential `xml:"q0:UserCredential"`
}

type UserCredential struct {
	Key      string `xml:"q0:Key"`
	Password string `xml:"q0:Password"`
}

type ClientDetail struct {
	AccountNumber string `xml:"q0:AccountNumber"`
	MeterNumber   string `xml:"q0:MeterNumber"`
}

type Version struct {
	ServiceID    string `xml:"q0:ServiceId"`
	Major        int    `xml:"q0:Major"`
	Intermediate int    `xml:"q0:Intermediate"`
	Minor        int    `xml:"q0:Minor"`
}
