package logic

import (
	"context"
	"double-token-example/internal/model"
	"double-token-example/internal/repository"
	"double-token-example/pkg/utils"
	"errors"
)

type UserLogic struct {
	userRepo *repository.UserRepository
}

func NewUserLogic() *UserLogic {
	return &UserLogic{
		userRepo: repository.NewUserRepository(),
	}
}

// Register 用户注册
func (l *UserLogic) Register(ctx context.Context, username, password, phone, email string) error {
	// 检查用户名是否已存在
	_, err := l.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 检查手机号是否已存在
	_, err = l.userRepo.GetByPhone(ctx, phone)
	if err == nil {
		return errors.New("手机号已存在")
	}

	// 检查邮箱是否已存在
	if email != "" {
		_, err = l.userRepo.GetByEmail(ctx, email)
		if err == nil {
			return errors.New("邮箱已存在")
		}
	}

	// 生成密码盐值
	salt, err := utils.GenerateSalt()
	if err != nil {
		return err
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(password + salt)
	if err != nil {
		return err
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Salt:     salt,
		Phone:    phone,
		Email:    email,
		Status:   1,
	}

	return l.userRepo.Create(ctx, user)
}

// Login 用户登录
func (l *UserLogic) Login(ctx context.Context, username, password string) (string, error) {
	// 获取用户信息
	user, err := l.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if !utils.CheckPassword(password+user.Salt, user.Password) {
		return "", errors.New("用户名或密码错误")
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserInfo 获取用户信息
func (l *UserLogic) GetUserInfo(ctx context.Context, userID int64) (*model.User, error) {
	return l.userRepo.GetByID(ctx, userID)
}
