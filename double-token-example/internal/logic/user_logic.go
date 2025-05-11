package logic

import (
	"context"
	"double-token-example/internal/config"
	"double-token-example/internal/model"
	"double-token-example/internal/repository"
	"double-token-example/pkg/utils"
	"errors"
	"time"
)

type UserLogic struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
}

func NewUserLogic() *UserLogic {
	return &UserLogic{
		userRepo:         repository.NewUserRepository(),
		refreshTokenRepo: repository.NewRefreshTokenRepository(),
	}
}

// Register 用户注册
func (l *UserLogic) Register(ctx context.Context, username, password, phone, email string) error {
	// 检查用户名是否已存在
	var (
		err                            error
		existByUsernameFromBloomFilter bool
		existByUsernameFromRedis       bool
	)
	existByUsernameFromBloomFilter, err = l.userRepo.ExistByUsernameFromBloomFilter(ctx, username)
	if err != nil {
		return err
	}
	if existByUsernameFromBloomFilter {
		existByUsernameFromRedis, err = l.userRepo.ExistByUsernameFromRedis(ctx, username)
		if err != nil {
			return err
		}
		if existByUsernameFromRedis {
			return errors.New("redis : 用户名已存在")
		}

		_, err = l.userRepo.GetByUsername(ctx, username)
		if err == nil {
			return errors.New("mysql : 用户名已存在")
		}
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
func (l *UserLogic) Login(ctx context.Context, username, password string) (string, string, error) {
	// 获取用户信息
	user, err := l.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if !utils.CheckPassword(password+user.Salt, user.Password) {
		return "", "", errors.New("用户名或密码错误")
	}

	// 生成token
	accessToken, err := utils.GenerateToken(
		utils.GenerateUUID(),
		user.ID,
		user.Username,
		config.Cfg.JWT.AccessTokenExpireTime,
		config.Cfg.JWT.AccessTokenType,
	)
	if err != nil {
		return "", "", err
	}
	refreshId := utils.GenerateUUID()
	refreshToken, err := utils.GenerateToken(
		refreshId,
		user.ID,
		user.Username,
		config.Cfg.JWT.RefreshTokenExpireTime,
		config.Cfg.JWT.RefreshTokenType,
	)
	if err != nil {
		return "", "", err
	}
	token := &model.RefreshToken{
		UserID:    user.ID,
		JTI:       refreshId,
		ExpiresAt: time.Now().Add(time.Duration(config.Cfg.JWT.RefreshTokenExpireTime) * time.Second),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	// 持久化refreshToken
	err = l.refreshTokenRepo.Create(ctx, token)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// RefreshToken 刷新token
func (l *UserLogic) RefreshToken(ctx context.Context, userID int64, username string) (string, string, error) {

	// 生成新的token
	accessToken, err := utils.GenerateToken(
		utils.GenerateUUID(),
		userID,
		username,
		config.Cfg.JWT.AccessTokenExpireTime,
		config.Cfg.JWT.AccessTokenType,
	)
	if err != nil {
		return "", "", err
	}
	refreshId := utils.GenerateUUID()
	refreshToken, err := utils.GenerateToken(
		refreshId,
		userID,
		username,
		config.Cfg.JWT.RefreshTokenExpireTime,
		config.Cfg.JWT.RefreshTokenType,
	)
	if err != nil {
		return "", "", err
	}
	token := &model.RefreshToken{
		UserID:    userID,
		JTI:       refreshId,
		ExpiresAt: time.Now().Add(time.Duration(config.Cfg.JWT.RefreshTokenExpireTime) * time.Second),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	err = l.refreshTokenRepo.Refresh(ctx, token)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// GetUserInfo 获取用户信息
func (l *UserLogic) GetUserInfo(ctx context.Context, userID int64) (*model.User, error) {
	return l.userRepo.GetByID(ctx, userID)
}

func (l *UserLogic) Logout(ctx context.Context, jti string, expiresAt time.Time, userID int64) error {
	err := l.refreshTokenRepo.DeleteByUserID(ctx, jti, expiresAt, userID)
	if err != nil {
		return err
	}
	return nil
}
