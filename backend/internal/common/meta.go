package common

import "time"

// 通用信息

type Meta struct {
	CreatedAt int64 `json:"created_at"` // 创建时间，这里存储unix时间戳
	UpdatedAt int64 `json:"updated_at"` // 更新时间
}

func NewMeta() *Meta {
	return &Meta{
		CreatedAt: time.Now().Unix(),
	}
}
