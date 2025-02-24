package periodical

type UpdateMode int

const (
	UPDATE_MODE_PUT   UpdateMode = iota // 全量更新
	UPDATE_MODE_PATCH                   // 部分更新
)
