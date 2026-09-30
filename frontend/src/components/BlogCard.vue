<template>
  <div class="blog-card" @click="toDetail">
    <div class="blog-img">
      <img :src="blog.img || blog.images?.split(',')[0]" alt="">
    </div>
    <div class="blog-title">{{ blog.title }}</div>
    <div class="blog-footer">
      <div class="blog-user">
        <img :src="blog.icon || '/imgs/icons/default-icon.png'" alt="">
        <span>{{ blog.name }}</span>
      </div>
      <div class="blog-stats">
        <div class="blog-liked" @click.stop="handleLike">
          <svg 
            class="like-icon" 
            viewBox="0 0 1024 1024" 
            width="14" 
            height="14"
          >
            <path 
              d="M160 944c0 8.8-7.2 16-16 16h-32c-26.5 0-48-21.5-48-48V528c0-26.5 21.5-48 48-48h32c8.8 0 16 7.2 16 16v448zM96 416c-53 0-96 43-96 96v416c0 53 43 96 96 96h96c17.7 0 32-14.3 32-32V448c0-17.7-14.3-32-32-32H96zM505.6 64c16.2 0 26.4 8.7 31 13.9 4.6 5.2 12.1 16.3 10.3 32.4l-23.5 203.4c-4.9 42.2 8.6 84.6 36.8 116.4 28.3 31.7 68.9 49.9 111.4 49.9h271.2c6.6 0 10.8 3.3 13.2 6.1s5 7.5 4 14l-48 303.4c-6.9 43.6-29.1 83.4-62.7 112C815.8 944.2 773 960 728.9 960h-317c-33.1 0-59.9-26.8-59.9-59.9v-455c0-6.1 1.7-12 5-17.1 69.5-109 106.4-234.2 107-364h41.6z m0-64h-44.9C427.2 0 400 27.2 400 60.7c0 127.1-39.1 251.2-112 355.3v484.1c0 68.4 55.5 123.9 123.9 123.9h317c122.7 0 227.2-89.3 246.3-210.5l47.9-303.4c7.8-49.4-30.4-94.1-80.4-94.1H671.6c-50.9 0-90.5-44.4-84.6-95l23.5-203.4C617.7 55 568.7 0 505.6 0z" 
              :fill="blog.isLike ? '#ff6633' : '#82848a'"
            />
          </svg>
          {{ blog.liked }}
        </div>
        <div class="blog-comments">
          <svg 
            class="comment-icon" 
            viewBox="0 0 1024 1024" 
            width="14" 
            height="14"
          >
            <path 
              d="M928 160c8.8 0 16-7.2 16-16V80c0-8.8-7.2-16-16-16H96c-8.8 0-16 7.2-16 16v80c0 8.8 7.2 16 16 16h16v688c0 35.3 28.7 64 64 64h640c35.3 0 64-28.7 64-64V160h16zM208 224h608v576H272c-17.7 0-32-14.3-32-32V256c0-17.7 14.3-32 32-32z" 
              fill="#82848a"
            />
          </svg>
          {{ blog.comments }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { likeBlog, getBlogById } from '@/api/blog'
import { ElMessage } from 'element-plus'

const props = defineProps({
  blog: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update'])
const router = useRouter()

const toDetail = () => {
  router.push(`/blog-detail/${props.blog.id}`)
}

const handleLike = async () => {
  try {
    await likeBlog(props.blog.id)
    const data = await getBlogById(props.blog.id)
    emit('update', data)
  } catch (error) {
    ElMessage.error('操作失败')
  }
}
</script>

<style scoped>
.blog-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.blog-img {
  width: 100%;
  aspect-ratio: 16/9;
  overflow: hidden;
}

.blog-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.blog-title {
  padding: 8px;
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.blog-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 8px 8px;
}

.blog-user {
  display: flex;
  align-items: center;
  gap: 6px;
}

.blog-user img {
  width: 20px;
  height: 20px;
  border-radius: 50%;
}

.blog-user span {
  font-size: 12px;
  color: #666;
}

.blog-stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.blog-liked {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
  cursor: pointer;
}

.blog-comments {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}

.like-icon,
.comment-icon {
  transition: transform 0.2s;
}

.blog-liked:active .like-icon {
  transform: scale(1.2);
}
</style>
