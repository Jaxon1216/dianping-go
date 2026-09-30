import request from '@/utils/request'

// 获取商户的优惠券列表
export const getVoucherList = (shopId) => {
  return request.get(`/voucher/list/${shopId}`)
}

// 秒杀抢购
export const seckillVoucher = (id) => {
  return request.post(`/voucher-order/seckill/${id}`)
}
