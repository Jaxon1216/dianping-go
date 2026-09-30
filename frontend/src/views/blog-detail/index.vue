<template>
  <div class="blog-detail-page">
    <HeaderBar :show-share="true">
      <template #right>
        <span>...</span>
      </template>
    </HeaderBar>

    <div class="content" ref="contentRef">
      <!-- 图片轮播 -->
      <ImageSwiper 
        v-if="blog.images?.length" 
        :images="blog.images"
        @change="onSwiperChange"
      />

      <!-- 作者信息 -->
      <div class="author-section">
        <div class="author-info" @click="toUserProfile">
          <img 
            class="avatar" 
            :src="blog.icon || '/imgs/icons/default-icon.png'" 
            alt=""
          >
          <div class="author-meta">
            <div class="name">{{ blog.name }}</div>
            <div class="time">{{ formatTime(blog.createTime) }}</div>
          </div>
        </div>
        <button 
          v-if="!isSelf"
          class="follow-btn"
          :class="{ followed }"
          @click="handleFollow"
        >
          {{ followed ? '取消关注' : '关注' }}
        </button>
      </div>

      <!-- 博客内容 -->
      <div class="blog-content" v-html="blog.content"></div>

      <!-- 关联商户 -->
      <div v-if="shop.id" class="shop-section" @click="toShopDetail">
        <img class="shop-img" :src="shop.image" alt="">
        <div class="shop-info">
          <div class="shop-name">{{ shop.name }}</div>
          <el-rate :model-value="shop.score / 10" disabled />
          <div class="shop-price">￥{{ shop.avgPrice }}/人</div>
        </div>
      </div>

      <!-- 点赞区域 -->
      <div class="likes-section">
        <div class="like-icon">
          <svg viewBox="0 0 1024 1024" width="24" height="24">
            <path 
              d="M160 944c0 8.8-7.2 16-16 16h-32c-26.5 0-48-21.5-48-48V528c0-26.5 21.5-48 48-48h32c8.8 0 16 7.2 16 16v448zM96 416c-53 0-96 43-96 96v416c0 53 43 96 96 96h96c17.7 0 32-14.3 32-32V448c0-17.7-14.3-32-32-32H96zM505.6 64c16.2 0 26.4 8.7 31 13.9 4.6 5.2 12.1 16.3 10.3 32.4l-23.5 203.4c-4.9 42.2 8.6 84.6 36.8 116.4 28.3 31.7 68.9 49.9 111.4 49.9h271.2c6.6 0 10.8 3.3 13.2 6.1s5 7.5 4 14l-48 303.4c-6.9 43.6-29.1 83.4-62.7 112C815.8 944.2 773 960 728.9 960h-317c-33.1 0-59.9-26.8-59.9-59.9v-455c0-6.1 1.7-12 5-17.1 69.5-109 106.4-234.2 107-364h41.6z m0-64h-44.9C427.2 0 400 27.2 400 60.7c0 127.1-39.1 251.2-112 355.3v484.1c0 68.4 55.5 123.9 123.9 123.9h317c122.7 0 227.2-89.3 246.3-210.5l47.9-303.4c7.8-49.4-30.4-94.1-80.4-94.1H671.6c-50.9 0-90.5-44.4-84.6-95l23.5-203.4C617.7 55 568.7 0 505.6 0z" 
              :fill="blog.isLike ? '#ff6633' : '#82848a'"
            />
          </svg>
        </div>
        <div class="like-users">
          <img 
            v-for="user in likes" 
            :key="user.id"
            :src="user.icon || '/imgs/icons/default-icon.png'" 
            alt=""
            class="like-avatar"
          >
          <span class="like-count">{{ blog.liked }}人点赞</span>
        </div>
      </div>

      <!-- 评论区 -->
      <div class="comments-section">
        <div class="section-title">
          网友评价 <span class="count">(119)</span>
        </div>
        <div class="comment-list">
          <div v-for="i in 3" :key="i" class="comment-item">
            <img 
              class="comment-avatar" 
              src="https://p0.meituan.net/userheadpicbackend/57e44d6eba01aad0d8d711788f30a126549507.jpg" 
              alt=""
            >
            <div class="comment-content">
              <div class="comment-user">叶小乙 <span class="level">Lv5</span></div>
              <el-rate :model-value="4.5" disabled />
              <p class="comment-text">某平台上买的券，价格可以当工作餐吃，虽然价格便宜，但是这家店一点都没有...</p>
              <div class="comment-images">
                <img v-for="j in 4" :key="j" src="https://qcloud.dpfile.com/pc/6T7MfXzx7USPIkSy7jzm40qZSmlHUF2jd-FZUL6WpjE9byagjLlrseWxnl1LcbuSGybIjx5eX6WNgCPvcASYAw.jpg" alt="">
              </div>
              <div class="comment-meta">浏览641 &nbsp;&nbsp; 评论5</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div class="footer-bar">
      <div class="action-btn" @click="handleLike">
        <svg viewBox="0 0 1024 1024" width="26" height="26">
          <path 
            d="M160 944c0 8.8-7.2 16-16 16h-32c-26.5 0-48-21.5-48-48V528c0-26.5 21.5-48 48-48h32c8.8 0 16 7.2 16 16v448zM96 416c-53 0-96 43-96 96v416c0 53 43 96 96 96h96c17.7 0 32-14.3 32-32V448c0-17.7-14.3-32-32-32H96zM505.6 64c16.2 0 26.4 8.7 31 13.9 4.6 5.2 12.1 16.3 10.3 32.4l-23.5 203.4c-4.9 42.2 8.6 84.6 36.8 116.4 28.3 31.7 68.9 49.9 111.4 49.9h271.2c6.6 0 10.8 3.3 13.2 6.1s5 7.5 4 14l-48 303.4c-6.9 43.6-29.1 83.4-62.7 112C815.8 944.2 773 960 728.9 960h-317c-33.1 0-59.9-26.8-59.9-59.9v-455c0-6.1 1.7-12 5-17.1 69.5-109 106.4-234.2 107-364h41.6z m0-64h-44.9C427.2 0 400 27.2 400 60.7c0 127.1-39.1 251.2-112 355.3v484.1c0 68.4 55.5 123.9 123.9 123.9h317c122.7 0 227.2-89.3 246.3-210.5l47.9-303.4c7.8-49.4-30.4-94.1-80.4-94.1H671.6c-50.9 0-90.5-44.4-84.6-95l23.5-203.4C617.7 55 568.7 0 505.6 0z" 
            :fill="blog.isLike ? '#ff6633' : '#82848a'"
          />
        </svg>
        <span :class="{ liked: blog.isLike }">{{ blog.liked }}</span>
      </div>
      <div class="action-btn">
        <el-icon size="26"><ChatDotSquare /></el-icon>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ChatDotSquare } from '@element-plus/icons-vue'
import { getBlogById, likeBlog, getBlogLikes } from '@/api/blog'
import { getShopById } from '@/api/shop'
import { isFollowed, followUser } from '@/api/follow'
import { getUserInfo } from '@/api/user'
import { formatTime } from '@/utils/format'
import HeaderBar from '@/components/HeaderBar.vue'
import ImageSwiper from '@/components/ImageSwiper.vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const blog = ref({})
const shop = ref({})
const likes = ref([])
const followed = ref(false)
const currentUser = ref(null)

const isSelf = computed(() => {
  return currentUser.value?.id === blog.value.userId
})

const fetchBlogDetail = async (id) => {
  try {
    const data = await getBlogById(id)
    data.images = data.images.split(',')
    blog.value = data
    
    // 获取关联商户
    if (data.shopId) {
      fetchShopDetail(data.shopId)
    }
    
    // 获取点赞列表
    fetchLikes(id)
    
    // 获取当前登录用户
    fetchCurrentUser()
  } catch (error) {
    ElMessage.error('获取博客详情失败')
  }
}

const fetchShopDetail = async (shopId) => {
  try {
    const data = await getShopById(shopId)
    data.image = data.images.split(',')[0]
    shop.value = data
  } catch (error) {
    console.error('获取商户详情失败:', error)
  }
}

const fetchLikes = async (blogId) => {
  try {
    likes.value = await getBlogLikes(blogId)
  } catch (error) {
    console.error('获取点赞列表失败:', error)
  }
}

const fetchCurrentUser = async () => {
  try {
    currentUser.value = await getUserInfo()
    if (!isSelf.value) {
      checkFollowStatus()
    }
  } catch {
    // 未登录不处理
  }
}

const checkFollowStatus = async () => {
  try {
    followed.value = await isFollowed(blog.value.userId)
  } catch (error) {
    console.error('获取关注状态失败:', error)
  }
}

const handleFollow = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  try {
    await followUser(blog.value.userId, !followed.value)
    followed.value = !followed.value
    ElMessage.success(followed.value ? '已关注' : '已取消关注')
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleLike = async () => {
  try {
    await likeBlog(blog.value.id)
    // 刷新博客数据
    const data = await getBlogById(blog.value.id)
    data.images = data.images.split(',')
    blog.value = data
    // 刷新点赞列表
    fetchLikes(blog.value.id)
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const toUserProfile = () => {
  if (isSelf.value) {
    router.push('/profile')
  } else {
    router.push(`/user/${blog.value.userId}`)
  }
}

const toShopDetail = () => {
  router.push(`/shop-detail/${shop.value.id}`)
}

const onSwiperChange = (index) => {
  console.log('当前图片索引:', index)
}

onMounted(() => {
  const blogId = route.params.id
  if (blogId) {
    fetchBlogDetail(blogId)
  }
})
</script>

<style scoped>
.blog-detail-page {
  min-height: 100vh;
  background: #fff;
  padding-bottom: 60px;
}

.content {
  height: calc(100vh - 104px);
  overflow-y: auto;
}

.author-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
}

.author-meta .name {
  font-size: 16px;
  font-weight: 500;
}

.author-meta .time {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.follow-btn {
  padding: 6px 16px;
  border: 1px solid #ff6633;
  border-radius: 16px;
  background: #ff6633;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
}

.follow-btn.followed {
  background: #fff;
  color: #666;
  border-color: #ddd;
}

.blog-content {
  padding: 0 15px 15px;
  font-size: 14px;
  line-height: 1.8;
  color: #333;
}

.shop-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0 15px 15px;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 8px;
  cursor: pointer;
}

.shop-img {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
}

.shop-info {
  flex: 1;
}

.shop-name {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.shop-price {
  font-size: 12px;
  color: #ff6633;
  margin-top: 4px;
}

.likes-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px;
  background: #f9f9f9;
  margin: 0 15px 15px;
  border-radius: 8px;
}

.like-users {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
}

.like-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid #fff;
  margin-left: -8px;
}

.like-avatar:first-child {
  margin-left: 0;
}

.like-count {
  font-size: 12px;
  color: #666;
  margin-left: 8px;
}

.comments-section {
  padding: 15px;
  border-top: 10px solid #f5f5f5;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 15px;
}

.count {
  font-size: 14px;
  color: #999;
}

.comment-item {
  display: flex;
  gap: 12px;
  padding: 15px 0;
  border-bottom: 1px solid #f1f1f1;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
}

.comment-content {
  flex: 1;
}

.comment-user {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.level {
  font-size: 12px;
  color: #ff6633;
  background: #ffe4d9;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 6px;
}

.comment-text {
  font-size: 14px;
  color: #333;
  line-height: 1.6;
  margin: 8px 0;
}

.comment-images {
  display: flex;
  gap: 8px;
  margin: 10px 0;
}

.comment-images img {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  object-fit: cover;
}

.comment-meta {
  font-size: 12px;
  color: #999;
}

.footer-bar {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 540px;
  height: 50px;
  background: #fff;
  border-top: 1px solid #f1f1f1;
  display: flex;
  align-items: center;
  padding: 0 30px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  cursor: pointer;
}

.action-btn span {
  font-size: 14px;
}

.action-btn span.liked {
  color: #ff6633;
}
</style>
