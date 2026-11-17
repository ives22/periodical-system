package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"periodical/conf"
	"periodical/internal/apps/role"
	roleImpl "periodical/internal/apps/role/impl"
	"periodical/internal/apps/user"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserImpl struct {
	db      *gorm.DB
	roleSvc role.Service
}

func NewUserImpl(roleSvc ...role.Service) *UserImpl {
	var rSvc role.Service
	if len(roleSvc) > 0 && roleSvc[0] != nil {
		rSvc = roleSvc[0]
	} else {
		rSvc = roleImpl.NewRoleImpl()
	}
	return &UserImpl{
		db:      conf.C().MySQL.GetConn().Debug(),
		roleSvc: rSvc,
	}
}

// revokeUserTokens 将指定用户的所有有效Token置为失效
func (i *UserImpl) revokeUserTokens(ctx context.Context, username string) {
	_ = i.db.WithContext(ctx).Table("tokens").
		Where("username = ? AND is_enable = 1", username).
		Update("is_enable", 0).Error
}

// DeleteUser 删除用户并清理关联角色与Token
func (i *UserImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) error {
	uObj, err := i.DescribeUser(ctx, user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		log.Printf("删除用户失败，用户%s不存在", req.Username)
		return err
	}

	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", uObj.Id).Delete(&role.UserRole{}).Error; err != nil {
			return err
		}
		if err := tx.Table("tokens").Where("username = ?", uObj.Username).Update("is_enable", 0).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", uObj.Id).Delete(&user.User{}).Error
	})
}

// CreateUser 创建用户并绑定多角色
func (i *UserImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	var count int64
	if err := i.db.WithContext(ctx).Model(&user.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	ins := user.NewUser(req)
	ins.UpdatedAt = ins.CreatedAt

	if err := i.db.WithContext(ctx).Create(ins).Error; err != nil {
		return nil, err
	}

	if err := i.roleSvc.SetUserRoles(ctx, ins.Id, req.RoleIds); err != nil {
		return nil, err
	}

	return i.DescribeUser(ctx, user.NewDescribeUserRequestById(fmt.Sprintf("%d", ins.Id)))
}

// DescribeUser 查询单个用户及其角色、权限集
func (i *UserImpl) DescribeUser(ctx context.Context, req *user.DescribeUserRequest) (*user.User, error) {
	query := i.db.WithContext(ctx)

	switch req.DescribeBy {
	case user.DESCRIBE_BY_ID:
		query = query.Where("id = ?", req.DescribeValue)
	case user.DESCRIBE_BY_USERNAME:
		query = query.Where("username = ?", req.DescribeValue)
	}

	ins := user.NewUser(user.NewCreateUserRequest())
	if err := query.First(ins).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user %s not found", req.DescribeValue)
		}
		return nil, err
	}

	ins.SetIsHashed()

	// 加载用户的多角色与权限并集
	roles, perms, isSuper, err := i.roleSvc.GetUserRolesAndPermissions(ctx, ins.Id)
	if err == nil {
		ins.Roles = roles
		ins.Permissions = perms
		ins.IsSuperAdmin = isSuper
		if len(roles) > 0 {
			ids := make([]int64, 0, len(roles))
			for _, r := range roles {
				ids = append(ids, r.Id)
			}
			ins.RoleIds = ids
		}
	}

	return ins, nil
}

// QueryUsers 分页查询用户列表
func (i *UserImpl) QueryUsers(ctx context.Context, req *user.QueryUserRequest) (*user.UserSet, error) {
	query := i.db.WithContext(ctx).Model(&user.User{})

	if req.Username != "" {
		query = query.Where("users.username LIKE ?", "%"+req.Username+"%")
	}
	if req.RoleId > 0 {
		subQuery := i.db.Table("user_roles").Select("user_id").Where("role_id = ?", req.RoleId)
		query = query.Where("users.id IN (?)", subQuery)
	}

	set := user.NewUserSet()
	if err := query.Count(&set.Total).Error; err != nil {
		return nil, err
	}

	if err := query.Preload("Roles").Order("users.id DESC").Offset(req.Offset()).Limit(req.PageSize).Find(&set.Items).Error; err != nil {
		return nil, err
	}

	for _, item := range set.Items {
		item.Desensitize()
	}

	return set, nil
}

// UpdatePassword 修改用户个人密码
func (i *UserImpl) UpdatePassword(ctx context.Context, req *user.UpdatePasswordRequest) error {
	if req.NewPassword == "" || len(req.NewPassword) < 6 {
		return errors.New("新密码长度不能少于6位")
	}
	uObj, err := i.DescribeUser(ctx, user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		return err
	}

	if err = uObj.CheckPassword(req.OldPassword); err != nil {
		log.Printf("修改密码失败，%s 旧密码错误", req.Username)
		return fmt.Errorf("旧密码错误")
	}

	uObj.Password = req.NewPassword
	uObj.SetNoHashed()
	uObj.PasswordHash()

	if err := i.db.WithContext(ctx).Model(&user.User{}).Where("username = ?", uObj.Username).Updates(map[string]interface{}{
		"password":   uObj.Password,
		"updated_at": time.Now().Unix(),
	}).Error; err != nil {
		return err
	}

	i.revokeUserTokens(ctx, uObj.Username)
	return nil
}

// UpdateUser 更新用户信息（多角色分配、重置密码、标签）
func (i *UserImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	uObj, err := i.DescribeUser(ctx, user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	needRevokeToken := false
	updates := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}

	// 如果需要重置密码
	if req.NewPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
		needRevokeToken = true
	}

	// 更新多角色绑定
	if len(req.RoleIds) > 0 {
		if err := i.roleSvc.SetUserRoles(ctx, uObj.Id, req.RoleIds); err != nil {
			return nil, err
		}
		updates["role"] = user.Role(req.RoleIds[0])
		needRevokeToken = true
	}

	// 更新标签
	if req.Label != nil {
		uObj.Label = req.Label
		updates["label"] = uObj.Label
	}

	if err := i.db.WithContext(ctx).Model(&user.User{}).Where("id = ?", uObj.Id).Updates(updates).Error; err != nil {
		return nil, err
	}

	if needRevokeToken {
		i.revokeUserTokens(ctx, uObj.Username)
	}

	updatedUser, err := i.DescribeUser(ctx, user.NewDescribeUserRequestById(fmt.Sprintf("%d", uObj.Id)))
	if err != nil {
		return nil, err
	}
	updatedUser.Desensitize()
	return updatedUser, nil
}
