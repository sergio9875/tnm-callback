package models

type PaymentGatewayRequest struct {
	Paymentmethod  string         `json:"paymentmethod"`
	Paymenttype    string         `json:"paymenttype"`
	Paymentdetails PaymentDetails `json:"paymentdetails"`
}

type PaymentDetails struct {
	Code              int    `json:"code"`
	Explanation       string `json:"explanation"`
	Success           bool   `json:"success"`
	Sequenceid        string `json:"sequenceid"`
	Terminalmno       string `json:"terminalmno"`
	Approvalreference string `json:"approvalreference"`
	Paymentreference  string `json:"paymentreference,omitempty"`
	Timestamp         string `json:"timestamp"`
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
