package models

import "time"

const (
	RoleEmployee = "EMPLOYEE"
	RoleAdmin    = "ADMIN"
	RoleApprover = "APPROVER"
	RoleFinance  = "FINANCE"
)

type User struct {
	ID             int64     `json:"id"`
	EmployeeNumber string    `json:"employee_number"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Position       string    `json:"position"`
	Division       string    `json:"division"`
	Role           string    `json:"role"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
