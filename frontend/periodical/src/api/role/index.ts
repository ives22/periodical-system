import request from "@/utils/request"

enum API {
    ALL_ENABLED_ROLES_URL = '/v1/role/all',
    ROLE_LIST_URL = '/v1/role/list',
    ROLE_BASE_URL = '/v1/role',
    PERMISSION_TREE_URL = '/v1/permission/tree',
}

// 获取所有已启用的角色列表（用于用户分配角色下拉框）
export const reqAllEnabledRoles = () => request.get<any, any>(API.ALL_ENABLED_ROLES_URL)

// 分页获取角色列表
export const reqRoleList = (data: any) => request.post<any, any>(API.ROLE_LIST_URL, data)

// 创建角色
export const reqCreateRole = (data: any) => request.post<any, any>(API.ROLE_BASE_URL, data)

// 更新角色信息
export const reqUpdateRole = (id: number, data: any) => request.put<any, any>(`${API.ROLE_BASE_URL}/${id}`, data)

// 删除角色
export const reqDeleteRole = (id: number) => request.delete<any, any>(`${API.ROLE_BASE_URL}/${id}`)

// 获取完整权限树
export const reqPermissionTree = () => request.get<any, any>(API.PERMISSION_TREE_URL)

// 获取指定角色拥有的权限ID列表
export const reqRolePermissionIds = (id: number) => request.get<any, any>(`${API.ROLE_BASE_URL}/${id}/permissions`)

// 为角色分配权限
export const reqAssignRolePermissions = (id: number, data: { permission_ids: number[] }) =>
    request.put<any, any>(`${API.ROLE_BASE_URL}/${id}/permissions`, data)
