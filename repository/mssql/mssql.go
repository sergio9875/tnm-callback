package mssql

import (
  "malawi-callback/models"
  repo "malawi-callback/repository"
  "database/sql"
  "fmt"

  _ "github.com/denisenkom/go-mssqldb"
)

// repository represent the repository model
type repository struct {
  db *sql.DB
}

// NewRepository will create a variable that represent the Repository struct
func NewRepository(dbConfig *models.DBConfig) (repo.Repository, error) {
  connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
   dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Port, dbConfig.Database)
  db, err := sql.Open(dbConfig.Dialect, connString)
  if err != nil {
    return nil, err
  }

  db.SetMaxIdleConns(1)
  db.SetMaxOpenConns(3)

  err = db.Ping()
  if err != nil {
    return nil, err
  }

  return &repository{db}, nil
}

// Close attaches the provider and close the connection
func (r *repository) Close() error {
  return r.db.Close()
}


