package repository

import (
	"database/sql"
	"tenant-service/internal/model"
)

type TenantRepo struct {
	db *sql.DB
}

func NewTenantRepo(db *sql.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) Create(t model.Tenant) error {
	query := "INSERT INTO tenants (id, name, email, api_key) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, t.ID, t.Name, t.Email, t.APIKey)
	return err
}

func (r *TenantRepo) ExistsByEmail(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tenants WHERE email = ?", email).Scan(&count)
	return count > 0, err
}

func (r *TenantRepo) GetByAPIKey(apiKey string) (*model.Tenant, error) {
	row := r.db.QueryRow("SELECT id, name, email, api_key FROM tenants WHERE api_key = ?", apiKey)

	var t model.Tenant
	err := row.Scan(&t.ID, &t.Name, &t.Email, &t.APIKey)
	if err != nil {
		return nil, err
	}
	return &t, nil
}