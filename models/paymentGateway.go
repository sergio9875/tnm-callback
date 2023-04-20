package models

type PaymentGatewayRequest struct {
	Paymentmethod  string         `json:"paymentmethod"`
	Transunq       interface{}    `json:"transunq"`
	Paymenttype    string         `json:"paymenttype"`
	Paymentdetails PaymentDetails `json:"paymentdetails"`
}

type PaymentDetails struct {
	Code              int           `json:"code"`
	Explanation       string        `json:"explanation"`
	Paymentreference  string        `json:"paymentreference"`
	Sequenceid        string        `json:"sequenceid"`
	Terminalmno       string        `json:"terminalmno"`
	Terminalsettings  []interface{} `json:"terminalsettings"`
	Amount            string        `json:"amount"`
	Currency          string        `json:"currency"`
	Msisdn            string        `json:"msisdn"`
	Approvalreference string        `json:"approvalreference"`
	Timestamp         string        `json:"timestamp"`
	Customerdetails   []interface{} `json:"customerdetails"`
	Mbtid             string        `json:"mbtid"`
}

type PaymentGatewayResponse struct {
	Code         string `json:"code"`
	Explanation  string `json:"explanation"`
	RedirectURL  string `json:"redirectURL"`
	Instructions string `json:"instructions"`
	Details      struct {
		ResultCode int    `json:"ResultCode"`
		StatusCode string `json:"StatusCode"`
	} `json:"details"`
}
