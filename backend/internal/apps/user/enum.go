package user

type Role int

const (
	ROLE_UNKNOWN     Role = iota // 未知角色
	ROLE_ADMIN                   // 管理员
	ROLE_AUTHOR                  // 普通用户
	ROLE_SUPER_ADMIN             // 超级管理员
)

type DescribeBy int

const (
	DESCRIBE_BY_ID DescribeBy = iota
	DESCRIBE_BY_USERNAME
)

// IsValid 验证角色是否有效
func (r Role) IsValid() bool {
	return r == ROLE_ADMIN || r == ROLE_SUPER_ADMIN || r == ROLE_AUTHOR
}
