import request from "@/utils/request"

enum API {
    LOGIN_URL = '/v1/auth/login',
    LOGOUT_URL = '/v1/auth/logout',
    USER_PROFILE_URL = '/v1/user/profile',
    UPDATE_PASSWORD_URL = '/v1/user/update',
    USER_LIST_URL = '/v1/user/list',
    USER_BASE_URL = '/v1/user',
    USER_DELETE_URL = '/v1/user/delete',
}

// reqLogin 登录接口
export const reqLogin = (data: any) => request.post<any, any>(API.LOGIN_URL, data)

// reqLogout 退出登录接口
export const reqLogout = (data: any) => request.post<any, any>(API.LOGOUT_URL, data)

// reqUserProfile 获取当前登录用户信息、绑定的多角色与权限集
export const reqUserProfile = () => request.get<any, any>(API.USER_PROFILE_URL)

// reqUpdatePassword 修改当前登录人密码
export const reqUpdatePassword = (data: any) => request.post<any, any>(API.UPDATE_PASSWORD_URL, data)

// reqUserList 分页查询用户列表
export const reqUserList = (data: any) => request.post<any, any>(API.USER_LIST_URL, data)

// reqCreateUser 创建新用户（支持绑定多个角色）
export const reqCreateUser = (data: any) => request.post<any, any>(API.USER_BASE_URL, data)

// reqUpdateUser 更新用户信息（多角色、标签、重置密码）
export const reqUpdateUser = (data: any) => request.put<any, any>(API.USER_BASE_URL, data)

// reqDeleteUser 删除用户
export const reqDeleteUser = (data: { username: string }) => request.post<any, any>(API.USER_DELETE_URL, data)
