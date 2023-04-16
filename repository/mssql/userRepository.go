package mssql

import (
  "malawi-callback/models"
  "context"
  "time"
)

// FindByID attaches the user repository and find data based on id
func (r *repository) FindUserByID(id int) (*models.UserEntity, error) {
  user := new(models.UserEntity)

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  err := r.db.QueryRowContext(ctx, "SELECT id, name, email, phone FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
  if err != nil {
    return nil, err
  }
  return user, nil
}

// Create attaches the user repository and creating the data
func (r *repository) CreateUser(user *models.UserEntity) error {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  query := "INSERT INTO users (id, name, email, phone) VALUES (?, ?, ?, ?)"
  stmt, err := r.db.PrepareContext(ctx, query)
  if err != nil {
    return err
  }
  defer stmt.Close()

  _, err = stmt.ExecContext(ctx, user.ID, user.Name, user.Email, user.Phone)
  return err
}

// Update attaches the user repository and update data based on id
func (r *repository) UpdateUser(user *models.UserEntity) error {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  query := "UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?"
  stmt, err := r.db.PrepareContext(ctx, query)
  if err != nil {
    return err
  }
  defer stmt.Close()

  _, err = stmt.ExecContext(ctx, user.Name, user.Email, user.Phone, user.ID)
  return err
}
