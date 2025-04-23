
import request from "@/utils/request"

enum API {
    // 查询期刊列表
    GET_PERIODICAL_LIST_URL = "/v1/periodical/list",
    // 查看期刊详情
    GET_PERIODICAL_DETAIL_URL = '/v1/periodical/',
    // 创建期刊
    CREATE_PERIODICAL_URL = '/v1/periodical',
    // 删除期刊
    DELETE_PERIODICAL_URL = '/v1/periodical/',
    // 编辑期刊
    UPDATE__PERIODICAL_URL = '/v1/periodical/',
    // 查询通用分类列表
    GET_CATEGORIZE_LIST_URL = "/v1/periodical/base/list"
}


// reqPeriodicalList 查询期刊列表
export const reqPeriodicalList = (data: any) => request.post<any, any>(API.GET_PERIODICAL_LIST_URL, data)

// reqCategorizeList 获取通用分类列表
export const reqCategorizeList = (data: any) => request.post<any, any>(API.GET_CATEGORIZE_LIST_URL, data)

// reqPeriodicalDetail 获取期刊详情
export const reqPeriodicalDetail = (id: number) => request.get<any, any>(API.GET_PERIODICAL_DETAIL_URL + `${id}`)

// reqCreatePeriodical 创建期刊
export const reqCreatePeriodical = (data: any) => request.post<any, any>(API.CREATE_PERIODICAL_URL, data)

// reqUpdatePeriodical 编辑期刊
export const reqUpdatePeriodical = (id: number, data: any) => request.put<any, any>(API.UPDATE__PERIODICAL_URL + `${id}`, data)

// 删除期刊
export const reqDeletePeriodical = (id: number) => request.delete<any, any>(API.DELETE_PERIODICAL_URL + `${id}`)