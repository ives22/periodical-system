package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"periodical/conf"
	"periodical/internal/apps/role"
	"periodical/internal/common"
	"sort"
	"time"

	"gorm.io/gorm"
)

type RoleImpl struct {
	db *gorm.DB
}

func NewRoleImpl() *RoleImpl {
	impl := &RoleImpl{
		db: conf.C().MySQL.GetConn().Debug(),
	}
	if err := impl.bootstrap(); err != nil {
		log.Printf("初始化RBAC角色权限表警告: %v", err)
	}
	return impl
}

// bootstrap 自动检查并初始化RBAC表、预设角色与权限树，并迁移存量用户角色
func (i *RoleImpl) bootstrap() error {
	// 1. 使用幂等SQL确保四张RBAC核心表存在
	tableDDLs := []string{
		`CREATE TABLE IF NOT EXISTS roles (
			id int unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
			name varchar(64) NOT NULL COMMENT '角色名称',
			code varchar(64) NOT NULL COMMENT '角色唯一编码',
			description varchar(255) DEFAULT '' COMMENT '角色描述',
			status tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
			is_system tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否内置角色：1-是，0-否',
			sort int NOT NULL DEFAULT 0 COMMENT '排序',
			created_at int NOT NULL DEFAULT 0 COMMENT '创建时间',
			updated_at int NOT NULL DEFAULT 0 COMMENT '更新时间',
			PRIMARY KEY (id),
			UNIQUE KEY idx_role_code (code)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色表';`,
		`CREATE TABLE IF NOT EXISTS permissions (
			id int unsigned NOT NULL AUTO_INCREMENT COMMENT '权限ID',
			parent_id int unsigned NOT NULL DEFAULT 0 COMMENT '父级权限ID，0为顶级',
			name varchar(64) NOT NULL COMMENT '权限名称',
			code varchar(100) NOT NULL COMMENT '权限唯一标识码',
			type tinyint(1) NOT NULL DEFAULT 1 COMMENT '类型：1-菜单/目录，2-按钮/接口操作',
			sort int NOT NULL DEFAULT 0 COMMENT '排序',
			created_at int NOT NULL DEFAULT 0 COMMENT '创建时间',
			updated_at int NOT NULL DEFAULT 0 COMMENT '更新时间',
			PRIMARY KEY (id),
			UNIQUE KEY idx_perm_code (code)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='权限资源表';`,
		`CREATE TABLE IF NOT EXISTS role_permissions (
			role_id int unsigned NOT NULL COMMENT '角色ID',
			permission_id int unsigned NOT NULL COMMENT '权限ID',
			PRIMARY KEY (role_id, permission_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色权限关联表';`,
		`CREATE TABLE IF NOT EXISTS user_roles (
			user_id int unsigned NOT NULL COMMENT '用户ID',
			role_id int unsigned NOT NULL COMMENT '角色ID',
			PRIMARY KEY (user_id, role_id),
			KEY idx_user_roles_role_id (role_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户角色关联表';`,
	}
	for _, ddl := range tableDDLs {
		if err := i.db.Exec(ddl).Error; err != nil {
			return err
		}
	}

	now := time.Now().Unix()

	// 2. 初始化预设权限树（幂等插入或补充）
	defaultPerms := []*role.Permission{
		{Id: 1, ParentId: 0, Name: "首页看板", Code: "dashboard:view", Type: role.PERM_TYPE_MENU, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 2, ParentId: 0, Name: "期刊管理", Code: "periodical:view", Type: role.PERM_TYPE_MENU, Sort: 2, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 3, ParentId: 2, Name: "查询与查看详情", Code: "periodical:query", Type: role.PERM_TYPE_BUTTON, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 4, ParentId: 2, Name: "新增期刊", Code: "periodical:create", Type: role.PERM_TYPE_BUTTON, Sort: 2, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 5, ParentId: 2, Name: "编辑期刊", Code: "periodical:update", Type: role.PERM_TYPE_BUTTON, Sort: 3, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 6, ParentId: 2, Name: "删除期刊", Code: "periodical:delete", Type: role.PERM_TYPE_BUTTON, Sort: 4, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 7, ParentId: 2, Name: "复制期刊描述", Code: "periodical:copy", Type: role.PERM_TYPE_BUTTON, Sort: 5, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 8, ParentId: 0, Name: "系统管理", Code: "system:view", Type: role.PERM_TYPE_MENU, Sort: 3, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 9, ParentId: 8, Name: "用户管理", Code: "user:view", Type: role.PERM_TYPE_MENU, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 10, ParentId: 9, Name: "查询用户列表", Code: "user:query", Type: role.PERM_TYPE_BUTTON, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 11, ParentId: 9, Name: "新增用户", Code: "user:create", Type: role.PERM_TYPE_BUTTON, Sort: 2, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 12, ParentId: 9, Name: "编辑用户/重置密码", Code: "user:update", Type: role.PERM_TYPE_BUTTON, Sort: 3, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 13, ParentId: 9, Name: "删除用户", Code: "user:delete", Type: role.PERM_TYPE_BUTTON, Sort: 4, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 14, ParentId: 8, Name: "角色与权限管理", Code: "role:view", Type: role.PERM_TYPE_MENU, Sort: 2, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 15, ParentId: 14, Name: "查询角色与权限树", Code: "role:query", Type: role.PERM_TYPE_BUTTON, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 16, ParentId: 14, Name: "新增角色", Code: "role:create", Type: role.PERM_TYPE_BUTTON, Sort: 2, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 17, ParentId: 14, Name: "编辑角色与分配权限", Code: "role:update", Type: role.PERM_TYPE_BUTTON, Sort: 3, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 18, ParentId: 14, Name: "删除角色", Code: "role:delete", Type: role.PERM_TYPE_BUTTON, Sort: 4, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 19, ParentId: 0, Name: "AI 智能工作台", Code: "ai:workspace", Type: role.PERM_TYPE_MENU, Sort: 3, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 20, ParentId: 19, Name: "AI 对话与智能助手", Code: "ai:chat", Type: role.PERM_TYPE_BUTTON, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 21, ParentId: 8, Name: "AI 模型与审计管理", Code: "ai:config", Type: role.PERM_TYPE_MENU, Sort: 3, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		{Id: 22, ParentId: 21, Name: "全员对话审计查看", Code: "ai:audit", Type: role.PERM_TYPE_BUTTON, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
	}

	for _, p := range defaultPerms {
		var count int64
		i.db.Model(&role.Permission{}).Where("code = ?", p.Code).Count(&count)
		if count == 0 {
			_ = i.db.Create(p).Error
			// 新增的AI工作台权限默认关联给管理员(1)、普通用户(2)与超级管理员(3)
			if p.Id == 19 || p.Id == 20 {
				_ = i.db.Exec("INSERT IGNORE INTO role_permissions (role_id, permission_id) VALUES (1, ?), (2, ?), (3, ?)", p.Id, p.Id, p.Id).Error
			} else if p.Id == 21 || p.Id == 22 {
				_ = i.db.Exec("INSERT IGNORE INTO role_permissions (role_id, permission_id) VALUES (3, ?)", p.Id).Error
			}
		}
	}

	// 3. 初始化默认角色（若角色表为空则初始化）
	var roleCount int64
	i.db.Model(&role.Role{}).Count(&roleCount)
	if roleCount == 0 {
		defaultRoles := []*role.Role{
			{Id: 1, Name: "管理员", Code: role.ADMIN_CODE, Description: "业务管理员，负责期刊管理及用户日常维护", Status: 1, IsSystem: 0, Sort: 2, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
			{Id: 2, Name: "普通用户", Code: role.AUTHOR_CODE, Description: "普通用户，可查询、浏览期刊详情及复制期刊描述", Status: 1, IsSystem: 0, Sort: 3, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
			{Id: 3, Name: "超级管理员", Code: role.SUPER_ADMIN_CODE, Description: "系统最高权限管理员，拥有系统全部权限", Status: 1, IsSystem: 1, Sort: 1, Meta: &common.Meta{CreatedAt: now, UpdatedAt: now}},
		}
		for _, r := range defaultRoles {
			_ = i.db.Create(r).Error
		}

		// 默认角色权限关联
		authorPerms := []int64{1, 2, 3, 7}
		for _, pid := range authorPerms {
			_ = i.db.Create(&role.RolePermission{RoleId: 2, PermissionId: pid}).Error
		}
		adminPerms := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
		for _, pid := range adminPerms {
			_ = i.db.Create(&role.RolePermission{RoleId: 1, PermissionId: pid}).Error
		}
		for pid := int64(1); pid <= 18; pid++ {
			_ = i.db.Create(&role.RolePermission{RoleId: 3, PermissionId: pid}).Error
		}
	}

	// 4. 自动迁移存量用户的单角色字段至 user_roles 表
	_ = i.db.Exec(`
		INSERT IGNORE INTO user_roles (user_id, role_id)
		SELECT id, CASE WHEN role IN (1, 2, 3) THEN role ELSE 2 END
		FROM users
		WHERE id NOT IN (SELECT DISTINCT user_id FROM user_roles)
	`).Error

	return nil
}

// CreateRole 创建角色
func (i *RoleImpl) CreateRole(ctx context.Context, req *role.CreateRoleRequest) (*role.Role, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var existCount int64
	if err := i.db.WithContext(ctx).Model(&role.Role{}).Where("code = ? OR name = ?", req.Code, req.Name).Count(&existCount).Error; err != nil {
		return nil, err
	}
	if existCount > 0 {
		return nil, errors.New("角色名称或角色编码已存在")
	}

	now := time.Now().Unix()
	ins := &role.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
		IsSystem:    0,
		Sort:        req.Sort,
		Meta: &common.Meta{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ins).Error; err != nil {
			return err
		}
		if len(req.PermissionIds) > 0 {
			var rps []role.RolePermission
			for _, pid := range req.PermissionIds {
				rps = append(rps, role.RolePermission{
					RoleId:       ins.Id,
					PermissionId: pid,
				})
			}
			if err := tx.Create(&rps).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	ins.PermissionIds = req.PermissionIds
	return ins, nil
}

// UpdateRole 更新角色基本信息
func (i *RoleImpl) UpdateRole(ctx context.Context, req *role.UpdateRoleRequest) (*role.Role, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	ins, err := i.DescribeRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 内置超级管理员禁止禁用或修改编码
	if ins.IsSuperAdmin() {
		if req.Status != nil && *req.Status == 0 {
			return nil, errors.New("系统内置超级管理员角色不可禁用")
		}
		if req.Code != "" && req.Code != ins.Code {
			return nil, errors.New("系统内置超级管理员编码不可修改")
		}
	}

	// 检查名称或编码冲突
	var existCount int64
	checkCode := ins.Code
	if req.Code != "" && ins.IsSystem == 0 {
		checkCode = req.Code
	}
	if err := i.db.WithContext(ctx).Model(&role.Role{}).
		Where("id != ? AND (name = ? OR code = ?)", ins.Id, req.Name, checkCode).
		Count(&existCount).Error; err != nil {
		return nil, err
	}
	if existCount > 0 {
		return nil, errors.New("角色名称或编码与其他角色冲突")
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"updated_at":  time.Now().Unix(),
	}
	if req.Code != "" && ins.IsSystem == 0 {
		updates["code"] = req.Code
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}

	if err := i.db.WithContext(ctx).Model(&role.Role{}).Where("id = ?", ins.Id).Updates(updates).Error; err != nil {
		return nil, err
	}

	return i.DescribeRole(ctx, ins.Id)
}

// DeleteRole 删除角色
func (i *RoleImpl) DeleteRole(ctx context.Context, id int64) error {
	ins, err := i.DescribeRole(ctx, id)
	if err != nil {
		return err
	}
	if ins.IsSystem == 1 || ins.IsSuperAdmin() {
		return errors.New("系统内置角色不可删除")
	}

	// 检查是否有用户关联了该角色
	var userCount int64
	if err := i.db.WithContext(ctx).Model(&role.UserRole{}).Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return fmt.Errorf("该角色下仍有关联的 %d 名用户，请先解除用户角色绑定后再删除", userCount)
	}

	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&role.RolePermission{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", id).Delete(&role.Role{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// DescribeRole 查询单个角色详情
func (i *RoleImpl) DescribeRole(ctx context.Context, id int64) (*role.Role, error) {
	ins := &role.Role{}
	if err := i.db.WithContext(ctx).Where("id = ?", id).First(ins).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	pids, _ := i.GetRolePermissionIds(ctx, ins.Id)
	ins.PermissionIds = pids
	return ins, nil
}

// QueryRoles 分页查询角色列表
func (i *RoleImpl) QueryRoles(ctx context.Context, req *role.QueryRoleRequest) (*role.RoleSet, error) {
	query := i.db.WithContext(ctx).Model(&role.Role{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?", like, like, like)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	set := role.NewRoleSet()
	if err := query.Count(&set.Total).Error; err != nil {
		return nil, err
	}

	if err := query.Order("sort ASC, id ASC").Offset(req.Offset()).Limit(req.PageSize).Find(&set.Items).Error; err != nil {
		return nil, err
	}

	for _, item := range set.Items {
		pids, _ := i.GetRolePermissionIds(ctx, item.Id)
		item.PermissionIds = pids
	}

	return set, nil
}

// QueryAllEnabledRoles 获取所有启用状态的角色列表
func (i *RoleImpl) QueryAllEnabledRoles(ctx context.Context) ([]*role.Role, error) {
	var items []*role.Role
	if err := i.db.WithContext(ctx).Where("status = 1").Order("sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// QueryPermissionTree 获取权限资源树
func (i *RoleImpl) QueryPermissionTree(ctx context.Context) ([]*role.Permission, error) {
	var allPerms []*role.Permission
	if err := i.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&allPerms).Error; err != nil {
		return nil, err
	}

	permMap := make(map[int64]*role.Permission, len(allPerms))
	for _, p := range allPerms {
		p.Children = []*role.Permission{}
		permMap[p.Id] = p
	}

	var roots []*role.Permission
	for _, p := range allPerms {
		if p.ParentId == 0 {
			roots = append(roots, p)
		} else if parent, ok := permMap[p.ParentId]; ok {
			parent.Children = append(parent.Children, p)
		} else {
			roots = append(roots, p)
		}
	}

	return roots, nil
}

// GetRolePermissionIds 获取指定角色的权限ID列表
func (i *RoleImpl) GetRolePermissionIds(ctx context.Context, roleId int64) ([]int64, error) {
	r := &role.Role{}
	if err := i.db.WithContext(ctx).Where("id = ?", roleId).First(r).Error; err == nil && r.IsSuperAdmin() {
		var allIds []int64
		_ = i.db.WithContext(ctx).Model(&role.Permission{}).Pluck("id", &allIds).Error
		return allIds, nil
	}

	var ids []int64
	if err := i.db.WithContext(ctx).Model(&role.RolePermission{}).
		Where("role_id = ?", roleId).
		Pluck("permission_id", &ids).Error; err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []int64{}
	}
	return ids, nil
}

// AssignRolePermissions 为角色分配权限
func (i *RoleImpl) AssignRolePermissions(ctx context.Context, req *role.AssignPermissionsRequest) error {
	ins, err := i.DescribeRole(ctx, req.RoleId)
	if err != nil {
		return err
	}
	if ins.IsSuperAdmin() {
		return errors.New("系统内置超级管理员默认拥有全部权限，无需修改")
	}

	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", req.RoleId).Delete(&role.RolePermission{}).Error; err != nil {
			return err
		}
		if len(req.PermissionIds) > 0 {
			// 去重
			uniqueMap := make(map[int64]struct{})
			var rps []role.RolePermission
			for _, pid := range req.PermissionIds {
				if _, exists := uniqueMap[pid]; !exists && pid > 0 {
					uniqueMap[pid] = struct{}{}
					rps = append(rps, role.RolePermission{
						RoleId:       req.RoleId,
						PermissionId: pid,
					})
				}
			}
			if len(rps) > 0 {
				if err := tx.Create(&rps).Error; err != nil {
					return err
				}
			}
		}
		return tx.Model(&role.Role{}).Where("id = ?", req.RoleId).Update("updated_at", time.Now().Unix()).Error
	})
}

// GetUserRolesAndPermissions 获取用户绑定的所有角色及合并去重后的权限码列表
func (i *RoleImpl) GetUserRolesAndPermissions(ctx context.Context, userId int64) ([]*role.Role, []string, bool, error) {
	var roles []*role.Role
	err := i.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userId).
		Order("roles.sort ASC, roles.id ASC").
		Find(&roles).Error
	if err != nil {
		return nil, nil, false, err
	}

	var enabledRoleIds []int64
	isSuperAdmin := false
	for _, r := range roles {
		if r.Status == 1 {
			enabledRoleIds = append(enabledRoleIds, r.Id)
			if r.IsSuperAdmin() {
				isSuperAdmin = true
			}
		}
	}

	var permCodes []string
	if isSuperAdmin {
		// 超级管理员拥有全部权限码
		if err := i.db.WithContext(ctx).Model(&role.Permission{}).Pluck("code", &permCodes).Error; err != nil {
			return roles, []string{}, true, err
		}
	} else if len(enabledRoleIds) > 0 {
		// 取所有已启用角色的权限码并集
		if err := i.db.WithContext(ctx).
			Table("permissions").
			Distinct("permissions.code").
			Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
			Where("role_permissions.role_id IN ?", enabledRoleIds).
			Pluck("permissions.code", &permCodes).Error; err != nil {
			return roles, []string{}, false, err
		}
	}

	if permCodes == nil {
		permCodes = []string{}
	}
	sort.Strings(permCodes)
	return roles, permCodes, isSuperAdmin, nil
}

// SetUserRoles 设置用户的多角色绑定关系
func (i *RoleImpl) SetUserRoles(ctx context.Context, userId int64, roleIds []int64) error {
	if len(roleIds) == 0 {
		return errors.New("请至少为用户分配一个角色")
	}

	// 校验角色ID是否存在
	var validCount int64
	if err := i.db.WithContext(ctx).Model(&role.Role{}).Where("id IN ?", roleIds).Count(&validCount).Error; err != nil {
		return err
	}
	if validCount == 0 {
		return errors.New("所选角色不存在")
	}

	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userId).Delete(&role.UserRole{}).Error; err != nil {
			return err
		}
		unique := make(map[int64]struct{})
		var urs []role.UserRole
		for _, rid := range roleIds {
			if _, ok := unique[rid]; !ok && rid > 0 {
				unique[rid] = struct{}{}
				urs = append(urs, role.UserRole{
					UserId: userId,
					RoleId: rid,
				})
			}
		}
		return tx.Create(&urs).Error
	})
}
