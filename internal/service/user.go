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

func NewUserService() *UserService {
	return &UserService{}
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, userData map[string]interface{}) error {
	// 不允许更新用户名和密码
	delete(userData, "username")
	delete(userData, "password")

	result := database.DB.Model(&model.User{}).Where("id = ?", id).Updates(userData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

//// UpdatePassword 更新用户密码
//func (s *UserService) UpdatePassword(id uint, oldPassword, newPassword string) error {
//	var user model.User
//	if err := database.DB.First(&user, id).Error; err != nil {
//		return err
//	}
//
//	// 验证旧密码
//	if !user.CheckPassword(oldPassword) {
//		return errors.New("旧密码不正确")
//	}
//
//	// 设置新密码
//	if err := user.SetPassword(newPassword); err != nil {
//		return err
//	}
//
//	// 更新密码
//	return database.DB.Model(&user).Update("password", user.Password).Error
//}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	result := database.DB.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
