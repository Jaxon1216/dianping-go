import request from '@/utils/request'

// 获取商户类型列表
export const getShopTypes = () => {
  return request.get('/shop-type/list')
}

// 获取商户列表
export const getShopList = (params) => {
  return request.get('/shop/of/type', { params })
}

// 根据名称搜索商户
export const searchShops = (name) => {
  return request.get('/shop/of/name', { params: { name } })
}

// 获取商户详情
export const getShopById = (id) => {
  return request.get(`/shop/${id}`)
}
