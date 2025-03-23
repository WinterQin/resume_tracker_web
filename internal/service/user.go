package service

import (
	"errors"
	"internship-manager/internal/model"
	"internship-manager/pkg/database"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

// Register 用户注册
func (s *UserService) Register(username, password, email string) error {
	// 检查用户名是否已存在
	var existingUser model.User
	result := database.DB.Where("username = ?", username).First(&existingUser)
	if result.RowsAffected > 0 {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	result = database.DB.Where("email = ?", email).First(&existingUser)
	if result.RowsAffected > 0 {
		return errors.New("邮箱已被使用")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建新用户
	user := model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	return database.DB.Create(&user).Error
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

// LoginByEmail 通过邮箱登录
func (s *UserService) LoginByEmail(email, password string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
