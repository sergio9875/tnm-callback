package models

type RedisMessage struct {
	RedisKey string `json:"redisKey"`
}

type SqsMessage struct {
	ReferenceId      string `json:"referenceId"`
	ServiceName      string `json:"serviceName"`
	PaymentReference string `json:"paymentReference"`
	QueueName        string `json:"queueName"`
	Ttl              string `json:"ttl"`
	MaxRetry         string `json:"maxRetry"`

	ConsumerKey    string `json:"consumerKey" redact:"complete"`
	ConsumerSecret string `json:"consumerSecret" redact:"complete"`
	AcquireRoute   string `json:"acquireRoute"`
	MnoApiUrl      string `json:"mnoApiUrl"`

	InputQueryReference               string `json:"inputQueryReference"`
	InputCountry                      string `json:"inputCountry"`
	InputThirdPartyConversationID     string `json:"inputThirdPartyConversationID"`
	InputVirtualPSPOrAquirerShortCode string `json:"inputVirtualPSPOrAquirerShortCode"`
	InputOrganisationShortCode        string `json:"inputOrganisationShortCode"`

	Counter    string `json:"counter"`
	TerminalId string `json:"terminalId"`
	TransId    string `json:"transId"`
}

type SumoPusherMessage struct {
	Category    string      `json:"category"`
	Fields      string      `json:"fields,omitempty"`
	SumoPayload SumoPayload `json:"sumoPayload,omitempty"`
}

type SumoPayload struct {
	Stack   string      `json:"stack"`
	Message string      `json:"message"`
	Params  interface{} `json:"params"`
}
