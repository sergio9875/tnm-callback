package models

import (
  log "malawi-callback/logger"
  "encoding/json"
  "time"
)

// Cleaner model
type Cleaner struct {
  PurgePeriod		int  `json:"purge_period"`
  MaxAge				int  `json:"max_age"`
}

// SqsDestination model
type SqsDestination struct {
  Url   string  `json:"url"`
}

// DBConfig model
type DBConfig struct {
  Dialect   string  `json:"dialect"`
  Database  string  `json:"database"`
  Host      string  `json:"host"`
  Port      int     `json:"port"`
  User      string  `json:"user"`
  Password  string  `json:"password"`
}

// Treasury DB model
type MyDB struct {
  Config  *DBConfig `json:"treasury,omitempty"`
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
  Secrets   []string          `json:"secrets,omitempty"`
  Port      *int              `json:"port,omitempty"`
  Cleaner   *Cleaner          `json:"cleaner"`
  Sender    *SqsDestination   `json:"result_queue,omitempty"`
  Database  *MyDB             `json:"db,omitempty"`
}

// Request model
type Event struct {
  Action        string      `json:"action"`
  RequestedBy   int         `json:"requested_by"`
  RequestedTS   time.Time   `json:"requested_at"`
}

func (e *Event) MarshalJSON() ([]byte, error) {
  type Alias Event
  return json.Marshal(&struct {
    RequestedTS int64 `json:"requested_at"`
    *Alias
  }{
    RequestedTS: e.RequestedTS.Unix(),
    Alias:    (*Alias)(e),
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
    log.Error("SYSTEM", "Inner secret parse error: " + err.Error())
    return nil
  }
  if secretModel.Database != nil {
    if sm.Database == nil {
      sm.Database = secretModel.Database
    } else {
      sm.Database.Merge(secretModel.Database)
    }
  }
  return sm
}

func (mdb *MyDB) Merge(MyDB *MyDB) {
  if mdb.Config == nil {
    mdb.Config = MyDB.Config
  } else {
    mdb.Config.Merge(MyDB.Config)
  }
}

func (dbc *DBConfig) Merge( databaseCfg *DBConfig) {
  if dbc.Database == "" {
    dbc.Database = databaseCfg.Database
  }
  if dbc.Dialect == "" {
    dbc.Dialect = databaseCfg.Dialect
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
  if dbc.User == "" {
    dbc.User = databaseCfg.User
  }
}
