import request from "@/utils/request";


enum API {
    LOGIN_URL = '/v1/auth/login',
    LOGOUT_URL = '/v1/auth/logout',
}

// reqLogin 登录接口
export const reqLogin = (data: any) => request.post<any, any>(API.LOGIN_URL, data)

// reqLogout 退出登录接口
export const reqLogout = (data: any) => request.post<any, any>(API.LOGOUT_URL, data)

