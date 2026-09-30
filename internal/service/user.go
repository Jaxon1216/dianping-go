package service

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/regex_utils"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserService interface {
	SendCode(ctx context.Context, req *v1.SendCodeReq) error
	Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRespData, error)
	Me(ctx context.Context) (*v1.SimpleUser, error)
	QueryUserByID(ctx context.Context, userID uint64) (*v1.SimpleUser, error)
	QueryUserInfoByID(ctx context.Context, userID uint64) (*model.UserInfo, error)
	Logout(ctx context.Context, token string) error
	Sign(ctx context.Context) error
	SignCount(ctx context.Context) (int, error)
}

func (s *userService) SignCount(ctx context.Context) (int, error) {
	user := user_holder.GetUser(ctx)
	if user == nil {
		return 0, v1.ErrCanNotGetUser
	}
	now := time.Now()
	keySuffix := now.Format("200601")
	key := constants.RedisUserSignKey + strconv.FormatUint(*user.ID, 10) + ":" + keySuffix
	dayOfMonth := now.Day()
	result, err := s.rdb.BitField(ctx, key,
		"GET", fmt.Sprintf("u%d", dayOfMonth), 0,
	).Result()
	if errors.Is(err, redis.Nil) || len(result) == 0 {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	num := result[0]
	if num == 0 {
		return 0, nil
	}

	var count int
	for {
		if (num & 1) == 0 {
			break
		} else {
			count++
		}
		num >>= 1
	}
	return count, nil
}

type userService struct {
	*Service
}

// NewUserService 演示 Go 的「隐式接口实现」（structural typing / duck typing）。
//
// 关键点：函数返回类型写的是接口 UserService，但 return 的却是具体结构体 *userService。
// Go 没有、也不需要 Java/TS 那样的 `implements` 关键字来声明「谁实现了谁」。
// 只要 *userService 定义齐了 UserService 接口要求的全部方法
// （SendCode / Login / Me / QueryUserByID / Sign / SignCount），
// 编译器就在这里自动检查并认定「*userService 满足 UserService」，检查不过则编译报错。
//
// 好处：上层（handler）只依赖接口 UserService，拿不到也不关心具体结构体 userService，
// 将来替换实现或在测试里塞 mock，上层代码一行都不用改。
//
// 注意 & 的作用：userService 是 struct，用 &userService{...} 取地址得到 *userService（指针）；
// 而返回类型 UserService 是接口，接口本身不加 *。
func NewUserService(service *Service) UserService {
	return &userService{
		Service: service,
	}
}

func (s *userService) SendCode(ctx context.Context, req *v1.SendCodeReq) error {
	// 1. 校验手机号
	if regex_utils.IsPhoneInvalid(req.Phone) {
		// 2. 如果不符合，返回错误信息
		return v1.ErrPhoneIsInvalid
	}
	// 3. 符合，生成验证码
	var code string
	if s.conf.Get("env") == "prod" {
		code = random.RandNumeral(6)
	} else {
		code = "123456"
	}

	// 4. 保存验证码到 redis
	key := constants.RedisLoginCodeKey + req.Phone
	if err := s.rdb.Set(ctx, key, code, constants.RedisLoginCodeTTL).Err(); err != nil {
		return err
	}

	// 5. 发送验证码，接入第三方服务
	s.logger.Info("发送短信验证码成功", zap.String("验证码：", code))
	// 返回 ok
	return nil
}

func (s *userService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRespData, error) {
	// 1. 校验手机号
	if regex_utils.IsPhoneInvalid(req.Phone) {
		// 2. 如果不符合，返回错误信息
		return nil, v1.ErrPhoneIsInvalid
	}

	// 3. 从 redis 获取验证码并校验
	key := constants.RedisLoginCodeKey + req.Phone
	cacheCode, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if cacheCode != req.Code {
		// 不一致，报错
		return nil, v1.ErrCodeIsInvalid
	}

	// 4. 一致，根据手机号查询用户 select * from tb_user where phone = ?
	user, err := s.query.User.Where(s.query.User.Phone.Eq(req.Phone)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 5. 判断用户是否存在
	if user == nil {
		// 6. 不存在，创建新用户并保存
		user, err = s.createUserWithPhone(req.Phone)
		if err != nil {
			return nil, err
		}
	}

	// 7. 保存用户信息到 redis 中
	// 7.1. 随机生成 token，作为登录令牌
	token, err := random.UUIdV4()
	if err != nil {
		return nil, err
	}
	// 7.2. 存储 User 对象
	var simpleUser v1.SimpleUser
	if err := copier.Copy(&simpleUser, &user); err != nil {
		return nil, err
	}

	// 7.3. 存储
	key = constants.RedisLoginUserKey + token
	if err := s.rdb.HSet(ctx, key, simpleUser).Err(); err != nil {
		return nil, err
	}
	// 7.4. 设置 token 有效期
	if err := s.rdb.Expire(ctx, key, constants.RedisLoginUserTTL).Err(); err != nil {
		return nil, err
	}

	// 8. 返回 token
	return &v1.LoginRespData{
		Token: token,
	}, nil
}

func (s *userService) Me(ctx context.Context) (*v1.SimpleUser, error) {
	user := user_holder.GetUser(ctx)
	if user == nil {
		return nil, v1.ErrCanNotGetUser
	}
	return user, nil
}

func (s *userService) QueryUserByID(ctx context.Context, userID uint64) (*v1.SimpleUser, error) {
	user, err := s.query.User.Where(s.query.User.ID.Eq(userID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &v1.SimpleUser{
		ID:       &user.ID,
		NickName: user.NickName,
		Icon:     user.Icon,
	}, nil
}

func (s *userService) QueryUserInfoByID(ctx context.Context, userID uint64) (*model.UserInfo, error) {
	userInfo, err := s.query.UserInfo.Where(s.query.UserInfo.UserID.Eq(userID)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return userInfo, err
}

func (s *userService) Logout(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, constants.RedisLoginUserKey+token).Err()
}

func (s *userService) createUserWithPhone(phone string) (*model.User, error) {
	// 1. 创建用户
	nickname := constants.UserNickNamePrefix + random.RandString(10)
	user := model.User{
		Phone:    phone,
		NickName: &nickname,
	}
	// 2. 保存用户
	if err := s.query.User.Save(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) Sign(ctx context.Context) error {
	user := user_holder.GetUser(ctx)
	if user == nil {
		return v1.ErrCanNotGetUser
	}
	now := time.Now()
	keySuffix := now.Format("200601")
	key := constants.RedisUserSignKey + strconv.FormatUint(*user.ID, 10) + ":" + keySuffix
	dayOfMonth := now.Day()
	return s.rdb.SetBit(ctx, key, int64(dayOfMonth-1), 1).Err()
}
