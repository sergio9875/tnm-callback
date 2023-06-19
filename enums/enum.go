package enums

const (
	PGW_SUCCESS         = "000"
	PGW_SUCCESS_BODY    = "TRNSACTION_PAID"
	PGW_FAILED_BODY     = "TRNSACTION_FAILED"
	PGW_FAILED          = "999"
	PGW_STATUS_SUCCESS  = 3
	PGW_STATUS_FAILED   = 7
	RESULT_CODE_SUCCESS = "200"
	ERROR_MSG_UNMARSHL  = "Can't Unmarshal Json Message callback TNM Malawi"
	ERROR_MSG_EMPTY_REQ = "Body was not provided in the ApiGateway event"
)
