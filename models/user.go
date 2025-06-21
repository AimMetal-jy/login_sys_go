package models

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // 不在JSON中显示密码哈希
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `json:"is_active"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// HashPassword 对密码进行哈希加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser 创建新用户
func (ur *UserRepository) CreateUser(username, password string) (*User, error) {
	// 检查用户名是否已存在
	exists, err := ur.UserExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// 对密码进行哈希加密
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// 插入用户到数据库
	query := `INSERT INTO users (username, password_hash) VALUES (?, ?)`
	result, err := ur.DB.Exec(query, username, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %v", err)
	}

	// 返回创建的用户信息
	return ur.GetUserByID(int(id))
}

// GetUserByUsername 根据用户名获取用户
func (ur *UserRepository) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password_hash, created_at, updated_at, is_active FROM users WHERE username = ?`
	row := ur.DB.QueryRow(query, username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (ur *UserRepository) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, password_hash, created_at, updated_at, is_active FROM users WHERE id = ?`
	row := ur.DB.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// UserExists 检查用户名是否已存在
func (ur *UserRepository) UserExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ?`
	var count int
	err := ur.DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %v", err)
	}
	return count > 0, nil
}

// ValidateUser 验证用户登录
func (ur *UserRepository) ValidateUser(username, password string) (*User, error) {
	user, err := ur.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}