package models

import (
	"encoding/json"
	log "malawi-callback/logger"
	"time"
)

type IncomingRequest struct {
	ConversationId string `json:"conversation_id,omitempty"`
	ResultCode     string `json:"result_code,omitempty"`
	ResultDesc     string `json:"result_desc,omitempty"`
	TransactionId  string `json:"transaction_id,omitempty"`
	ExternalRef    string `json:"external_ref,omitempty"`
	ResponseTime   string `json:"response_time,omitempty"`
}
type USSDRedirectParams struct {
	//PaymentMethod string `json:"paymentmethod"`
	TransId string `json:"transId"`
	//TransUnq string `json:"transunq"`
	//PaymentType   string `json:"paymenttype"`
	//TransTransstId string `json:"transtransstid"`
	//TransIp            string `json:"transip"`
	FinalPaymentAmount string `json:"finalpaymentamount"`
	//Currency           string `json:"currency"`
	Explanation string `json:"explanation"`

	MerchantId      string `json:"merchantId"`
	TerminalId      string `json:"terminalId"`
	SubMerchantName string `json:"SubMerchantName"`
	TransactionId   string `json:"TransactionId"`
	TraceId         string `json:"TraceId"`
	CustomerId      string `json:"CustomerId"`
	Reference       string `json:"Reference"`
	Amount          string `json:"Amount"`
	BankCode        string `json:"BankCode"`
	BankName        string `json:"BankName"`
	CustomerMobile  string `json:"CustomerMobile"`
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	TimeStamp       int    `json:"TimeStamp"`
	Signature       string `json:"Signature"`
}

type TransactionStatusParsedResponse struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	Reference       string `json:"Reference"`
	Amount          int    `json:"Amount"`
	TerminalId      string `json:"TerminalId"`
	MerchantId      string `json:"MerchantId"`
	BankCode        string `json:"BankCode"`
	BankName        string `json:"BankName"`
	CustomerMobile  string `json:"CustomerMobile"`
	SubMerchantName string `json:"SubMerchantName"`
	TransactionId   string `json:"TransactionId"`
	TraceId         string `json:"TraceId"`
	CustomerId      string `json:"CustomerId"`
	Signature       string `json:"Signature"`
	TimeStamp       int    `json:"TimeStamp"`
}

// Cleaner model
type Cleaner struct {
	PurgePeriod int `json:"purge_period"`
	MaxAge      int `json:"max_age"`
}

// SqsDestination model
type SqsDestination struct {
	Url string `json:"url"`
}

// MssqlConfig MSSQL DBConfig model
type MssqlConfig struct {
	Port              int    `json:"port"`
	Host              string `json:"host"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Database          string `json:"database"`
	ConnectionTimeout int    `json:"connectionTimeout"`
	RequestTimeout    int    `json:"requestTimeout"`
}

// DBConfig Treasury DB model
type DBConfig struct {
	Africainv *MssqlConfig `json:"africainv"`
}

type Services struct {
	MailmailSenderQueueUrl string `json:"mailSenderQueueUrl,omitempty"`
	SumoPusherUrl          string `json:"sumoPusherUrl,omitempty"`
}

// Cache RedisConfig model
type Cache struct {
	Type     *string `json:"type,omitempty"`
	Host     *string `json:"host,omitempty"`
	Port     *int    `json:"port,omitempty"`
	Password *string `json:"password,omitempty"`
	Database *int    `json:"db,omitempty"`
}

// SecretModel model
type SecretModel struct {
	Secrets    []string  `json:"secrets,omitempty"`
	Database   *DBConfig `json:"database,omitempty"`
	Services   *Services `json:"services,omitempty"`
	Cache      *Cache    `json:"cache,omitempty"`
	DpoPygwUrl string    `json:"dpo_pygw_url,omitempty"`
}

// Event Request model
type Event struct {
	Action      string    `json:"action"`
	RequestedBy int       `json:"requested_by"`
	RequestedTS time.Time `json:"requested_at"`
}

func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	return json.Marshal(&struct {
		RequestedTS int64 `json:"requested_at"`
		*Alias
	}{
		RequestedTS: e.RequestedTS.Unix(),
		Alias:       (*Alias)(e),
	})
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		RequestedTS int64 `json:"requested_at"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.RequestedTS = time.Unix(aux.RequestedTS, 0)
	return nil
}

func (sm *SecretModel) Merge(src *string) *SecretModel {
	var secretModel = SecretModel{}
	err := json.Unmarshal([]byte(*src), &secretModel)
	if err != nil {
		log.Error("SYSTEM", "Inner secret parse error: "+err.Error())
		return nil
	}

	if secretModel.Services != nil {
		if sm.Services == nil {
			sm.Services = secretModel.Services
		}
	}

	if secretModel.Database != nil {
		if sm.Database == nil {
			sm.Database = secretModel.Database
		} else {
			sm.Database.Merge(secretModel.Database)
		}
	}

	if secretModel.Cache != nil {
		if sm.Cache == nil {
			sm.Cache = secretModel.Cache
		} else {
			sm.Cache.Merge(secretModel.Cache)
		}
	}

	return sm
}

func (tdb *DBConfig) Merge(dbConfig *DBConfig) {
	if tdb.Africainv == nil {
		tdb.Africainv = dbConfig.Africainv
	} else {
		tdb.Africainv.Merge(dbConfig.Africainv)
	}
}

func (dbc *MssqlConfig) Merge(databaseCfg *MssqlConfig) {
	if dbc.Database == "" {
		dbc.Database = databaseCfg.Database
	}

	if dbc.Host == "" {
		dbc.Host = databaseCfg.Host
	}
	if dbc.Password == "" {
		dbc.Password = databaseCfg.Password
	}
	if dbc.Port == 0 {
		dbc.Port = databaseCfg.Port
	}
	if dbc.Username == "" {
		dbc.Username = databaseCfg.Username
	}

	if dbc.ConnectionTimeout == 0 {
		dbc.ConnectionTimeout = databaseCfg.ConnectionTimeout
	}

	if dbc.RequestTimeout == 0 {
		dbc.RequestTimeout = databaseCfg.RequestTimeout
	}
}

func (c *Cache) Merge(cache *Cache) {
	if c.Database == nil {
		c.Database = cache.Database
	}
	if c.Type == nil {
		c.Type = cache.Type
	}
	if c.Host == nil {
		c.Host = cache.Host
	}
	if c.Password == nil {
		c.Password = cache.Password
	}
	if c.Port == nil {
		c.Port = cache.Port
	}
}
