package enums

const (
	PGW_SUCCESS          = "000"
	PGW_SUCCESS_BODY     = "TRANSACTION_PAID"
	PGW_FAILED_BODY      = "TRANSACTION_FAILED"
	PGW_FAILED           = "999"
	PGW_STATUS_SUCCESS   = 3
	PGW_STATUS_FAILED    = 7
	RESULT_CODE_SUCCESS  = "0"
	PAID_STATUS          = true
	PAID_MESSAGE_SUCCESS = "Process service request successfully."
	ERROR_MSG_UNMARSHL   = "Can't Unmarshal JSON Callback From TNM Malawi"
	ERROR_MSG_EMPTY_REQ  = "Body was not provided in the ApiGateway event"
)
