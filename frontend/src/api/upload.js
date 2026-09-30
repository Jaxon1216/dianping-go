import request from '@/utils/request'

// 上传博客图片
export const uploadBlogImage = (file) => {
  const formData = new FormData()
  formData.append('file', file)
  
  return request.post('/upload/blog', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 删除博客图片
export const deleteBlogImage = (name) => {
  return request.get('/upload/blog/delete', { params: { name } })
}
