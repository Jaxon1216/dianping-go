import request from '@/utils/request'

const normalizeBlog = (blog) => {
  if (!blog) return blog
  return {
    ...blog,
    images: blog.images ?? blog.Images ?? ''
  }
}

const normalizeBlogList = (blogs) => {
  return Array.isArray(blogs) ? blogs.map(normalizeBlog) : []
}

// 获取热门博客列表
export const getHotBlogs = (current = 1) => {
  return request
    .get('/blog/hot', { params: { id: current } })
    .then(normalizeBlogList)
}

// 获取博客详情
export const getBlogById = (id) => {
  return request.get(`/blog/${id}`).then(normalizeBlog)
}

// 获取用户的博客列表
export const getUserBlogs = (id, current = 1) => {
  return request
    .get('/blog/of/user', { params: { id, current } })
    .then(normalizeBlogList)
}

// 获取当前用户的博客列表
export const getMyBlogs = (current = 1) => {
  return request
    .get('/blog/of/me', { params: { id: current } })
    .then(normalizeBlogList)
}

// 获取关注的人的博客
export const getFollowBlogs = (params) => {
  return request.get('/blog/of/follow', { params }).then((result) => {
    if (!result) return result
    return {
      ...result,
      minTime: result.minTime ?? result.min_time,
      list: normalizeBlogList(result.list)
    }
  })
}

// 点赞/取消点赞
export const likeBlog = (id) => {
  return request.put(`/blog/like/${id}`)
}

// 获取点赞列表
export const getBlogLikes = (id) => {
  return request.get(`/blog/likes/${id}`)
}

// 发布博客
export const publishBlog = (data) => {
  return request.post('/blog', data)
}
