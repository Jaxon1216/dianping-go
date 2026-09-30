import request from '@/utils/request'

// 判断是否已关注
export const isFollowed = (id) => {
  return request.get(`/follow/or/not/${id}`)
}

// 关注/取消关注
export const followUser = (id, isFollow) => {
  return request.put(`/follow/${id}/${isFollow}`)
}

// 获取共同关注
export const getCommonFollows = (id) => {
  return request.get(`/follow/common/${id}`)
}
