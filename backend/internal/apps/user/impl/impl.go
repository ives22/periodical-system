package impl

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"periodical/conf"
	"periodical/internal/apps/user"
)

type UserImpl struct {
	db *gorm.DB
}

func (i *UserImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) error {
	// 查询用户是否存在
	uObj, err := i.DescribeUser(ctx, user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		log.Printf("删除用户失败，用户%s不存在", req.Username)
		return err
	}

	// 用户存在，则进行删除
	return i.db.WithContext(ctx).Where("username = ?", req.Username).Delete(uObj).Error
}

func NewUserImpl() *UserImpl {
	return &UserImpl{
		db: conf.C().MySQL.GetConn().Debug(),
	}
}

// CreateUser 创建用户
func (i *UserImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	// 1. 校验用户参数
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. 生成一个User对象（ORM对象）
	ins := user.NewUser(req)

	// 3. 保存到数据库
	if err := i.db.WithContext(ctx).Create(ins).Error; err != nil {
		return nil, err
	}

	return ins, nil
}

// DescribeUser 查询用户
func (i *UserImpl) DescribeUser(ctx context.Context, req *user.DescribeUserRequest) (*user.User, error) {
	query := i.db.WithContext(ctx)

	// 根据条件来构造where语句
	switch req.DescribeBy {
	case user.DESCRIBE_BY_ID:
		query = query.Where("id = ?", req.DescribeValue)
	case user.DESCRIBE_BY_USERNAME:
		query = query.Where("username = ?", req.DescribeValue)

	}

	ins := user.NewUser(user.NewCreateUserRequest())
	if err := query.First(ins).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user %s not found", req.DescribeValue)
		}
		return nil, err
	}

	ins.SetIsHashed()

	return ins, nil
}

// UpdatePassword 修改用户密码
func (i *UserImpl) UpdatePassword(ctx context.Context, req *user.UpdatePasswordRequest) error {
	// 查询用户是否存在
	uObj, err := i.DescribeUser(ctx, user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		return err
	}

	// 验证旧密码
	if err = uObj.CheckPassword(req.OldPassword); err != nil {
		log.Printf("修改密码失败，%s 旧密码错误", req.Username)
		return fmt.Errorf("旧密码错误")
	}

	// 生成新密码的hash
	uObj.Password = req.NewPassword
	uObj.SetNoHashed()  // 设置密码为未加密状态
	uObj.PasswordHash() // 进行密码加密

	return i.db.WithContext(ctx).Model(&user.User{}).Where("username = ?", uObj.Username).Update("password", uObj.Password).Error
}

// UpdateUser 更新用户信息
func (i *UserImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	// 查询用户是否存在
	uObj, err := i.DescribeUser(ctx, user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 开始更新用户信息
	if req.NewPassword != "" {
		// 如果需要更新密码，先对密码进行hash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		uObj.Password = string(hashedPassword)
	}

	// 更新角色
	if req.Role != user.ROLE_UNKNOWN {
		if !req.Role.IsValid() {
			return nil, fmt.Errorf("无效的角色值: %d", req.Role)
		}
		uObj.Role = req.Role
	}

	// 更新标签
	if len(req.Label) > 0 {
		uObj.Label = req.Label
	}

	// 保存更新
	if err := i.db.WithContext(ctx).Where("username = ?", req.Username).Updates(uObj).Error; err != nil {
		return nil, err
	}
	// 清空密码后返回
	uObj.Password = ""
	return uObj, err
}
