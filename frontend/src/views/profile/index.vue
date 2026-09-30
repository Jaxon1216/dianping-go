<template>
  <div class="profile-page">
    <HeaderBar title="个人主页" />

    <!-- 用户信息 -->
    <div class="user-section">
      <div class="user-info">
        <img 
          class="avatar" 
          :src="user.icon || '/imgs/icons/default-icon.png'" 
          alt=""
        >
        <div class="user-meta">
          <div class="name">{{ user.nickName }}</div>
          <div class="location">
            <el-icon><Location /></el-icon>
            杭州
          </div>
        </div>
        <button class="edit-btn" @click="toEdit">编辑资料</button>
      </div>
      <button class="logout-btn" @click="logout">退出登录</button>
    </div>

    <!-- 简介 -->
    <div class="intro-section">
      <p v-if="userDetail.introduce">{{ userDetail.introduce }}</p>
      <p v-else class="placeholder">
        添加个人简介，让大家更好的认识你
        <el-icon><Edit /></el-icon>
      </p>
    </div>

    <!-- 内容标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="笔记" name="blogs">
        <div class="blog-list">
          <div 
            v-for="blog in blogs" 
            :key="blog.id"
            class="blog-item"
            @click="toBlogDetail(blog.id)"
          >
            <img :src="blog.images.split(',')[0]" alt="">
            <div class="blog-info">
              <div class="blog-title">{{ blog.title }}</div>
              <div class="blog-stats">
                <span><img src="/imgs/thumbup.png" alt=""> {{ blog.liked }}</span>
                <span><el-icon><ChatDotRound /></el-icon> {{ blog.comments }}</span>
              </div>
            </div>
          </div>
          <div v-if="blogs.length === 0" class="empty">暂无笔记</div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="评价" name="reviews">
        <div class="empty">暂无评价</div>
      </el-tab-pane>

      <el-tab-pane :label="`粉丝(${fansCount})`" name="fans">
        <div class="empty">暂无粉丝</div>
      </el-tab-pane>

      <el-tab-pane :label="`关注(${followCount})`" name="follows">
        <div 
          ref="scrollContainer"
          class="follow-list"
          @scroll="handleFollowScroll"
          @touchstart="handleTouchStart"
          @touchmove="handleTouchMove"
          @touchend="handleTouchEnd"
          @mousedown="handleMouseDown"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
          @mouseleave="handleMouseUp"
        >
          <!-- 下拉刷新提示 -->
          <div class="refresh-tip" :style="{ height: refreshHeight + 'px' }">
            <div v-if="refreshing" class="loading-spinner">
              <el-icon class="is-loading"><Loading /></el-icon>
              <span>刷新中...</span>
            </div>
            <div v-else-if="refreshHeight > 0" class="pull-tip">
              {{ refreshHeight > 50 ? '释放刷新' : '下拉刷新' }}
            </div>
          </div>
          
          <BlogCard 
            v-for="blog in followBlogs" 
            :key="blog.id"
            :blog="blog"
            @update="updateFollowBlog"
          />
          <div v-if="followLoading" class="loading">加载中...</div>
          <div v-if="followFinished" class="finished">没有更多了</div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <FootBar :active-btn="4" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Location, Edit, ChatDotRound, Loading } from '@element-plus/icons-vue'
import { getMyBlogs, getFollowBlogs } from '@/api/blog'
import { getUserDetail } from '@/api/user'
import { useUserStore } from '@/stores/user'
import { useScrollLoad } from '@/utils/useScrollLoad'
import HeaderBar from '@/components/HeaderBar.vue'
import FootBar from '@/components/FootBar.vue'
import BlogCard from '@/components/BlogCard.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const user = ref({})
const userDetail = ref({})
const blogs = ref([])
const activeTab = ref('blogs')
const fansCount = ref(0)
const followCount = ref(0)

// 关注列表相关
const followBlogs = ref([])
const followParams = ref({
  minTime: 0,
  offset: 0
})

// 下拉刷新相关
const refreshHeight = ref(0)
const refreshing = ref(false)
const startY = ref(0)
const maxRefreshHeight = 80

const loadFollowBlogs = async (isRefresh = false) => {
  const params = {
    offset: isRefresh ? 0 : followParams.value.offset,
    lastId: isRefresh ? Date.now() + 1 : (followParams.value.minTime || Date.now() + 1)
  }
  
  const data = await getFollowBlogs(params)
  console.log('getFollowBlogs返回的数据:', data)
  if (!data || !data.list) return false
  
  const { list, minTime, offset } = data
  console.log('处理前的list:', list)
  list.forEach(blog => {
    // 确保所有必要的属性都存在
    blog.img = blog.images?.split(',')[0] || ''
    blog.title = blog.title || '默认标题'
    blog.name = blog.name || '未知用户'
    blog.icon = blog.icon || '/imgs/icons/default-icon.png'
    blog.liked = blog.liked || 0
    blog.comments = blog.comments || 0
    blog.isLike = blog.isLike || false
  })
  console.log('处理后的list:', list)
  
  if (isRefresh) {
    followBlogs.value = list
  } else {
    followBlogs.value.push(...list)
  }
  followParams.value = { minTime, offset }
  return list.length > 0
}

const { 
  loading: followLoading, 
  finished: followFinished, 
  scrollContainer, 
  handleScroll: handleFollowScroll,
  reset: resetFollow
} = useScrollLoad(loadFollowBlogs, { immediate: false })

// 下拉刷新相关方法
const handleTouchStart = (e) => {
  if (scrollContainer.value && scrollContainer.value.scrollTop === 0) {
    startY.value = e.touches[0].clientY
  }
}

const handleTouchMove = (e) => {
  if (startY.value === 0 || refreshing.value) return
  
  const currentY = e.touches[0].clientY
  const diff = currentY - startY.value
  
  if (diff > 0 && scrollContainer.value && scrollContainer.value.scrollTop === 0) {
    e.preventDefault()
    refreshHeight.value = Math.min(diff * 0.5, maxRefreshHeight)
  }
}

const handleTouchEnd = async () => {
  if (refreshHeight.value > 50) {
    // 触发刷新
    refreshing.value = true
    await loadFollowBlogs(true)
    resetFollow()
  }
  // 重置
  refreshHeight.value = 0
  refreshing.value = false
  startY.value = 0
}

// 鼠标事件处理（PC端支持）
const handleMouseDown = (e) => {
  if (scrollContainer.value && scrollContainer.value.scrollTop === 0) {
    startY.value = e.clientY
  }
}

const handleMouseMove = (e) => {
  if (startY.value === 0 || refreshing.value) return
  
  const currentY = e.clientY
  const diff = currentY - startY.value
  
  if (diff > 0 && scrollContainer.value && scrollContainer.value.scrollTop === 0) {
    e.preventDefault()
    refreshHeight.value = Math.min(diff * 0.5, maxRefreshHeight)
  }
}

const handleMouseUp = async () => {
  if (refreshHeight.value > 50) {
    // 触发刷新
    refreshing.value = true
    await loadFollowBlogs(true)
    resetFollow()
  }
  // 重置
  refreshHeight.value = 0
  refreshing.value = false
  startY.value = 0
}

const fetchUserInfo = async () => {
  try {
    user.value = userStore.userInfo || await userStore.fetchUserInfo()
    fetchUserDetail(user.value.id)
  } catch (error) {
    ElMessage.error('获取用户信息失败')
    router.push('/login')
  }
}

const fetchUserDetail = async (id) => {
  try {
    const data = await getUserDetail(id)
    if (data) {
      userDetail.value = data
      userStore.setUserDetail(data)
    }
  } catch (error) {
    console.error('获取用户详情失败:', error)
  }
}

const fetchMyBlogs = async () => {
  try {
    blogs.value = await getMyBlogs()
  } catch (error) {
    console.error('获取笔记失败:', error)
  }
}

const handleTabChange = (tab) => {
  if (tab === 'follows' && followBlogs.value.length === 0) {
    followParams.value.minTime = Date.now() + 1
    followParams.value.offset = 0
    resetFollow()
    loadFollowBlogs()
  }
}

const updateFollowBlog = (updatedBlog) => {
  const index = followBlogs.value.findIndex(b => b.id === updatedBlog.id)
  if (index !== -1) {
    updatedBlog.img = updatedBlog.images.split(',')[0]
    followBlogs.value[index] = updatedBlog
  }
}

const toEdit = () => {
  router.push('/profile-edit')
}

const toBlogDetail = (id) => {
  router.push(`/blog-detail/${id}`)
}

const logout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await userStore.logout()
    router.push('/')
  } catch {
    // 取消退出
  }
}

onMounted(() => {
  fetchUserInfo()
  fetchMyBlogs()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 50px;
}

.user-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 15px;
  background: #fff;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
}

.user-meta .name {
  font-size: 18px;
  font-weight: 600;
}

.user-meta .location {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.edit-btn {
  padding: 6px 16px;
  border: 1px solid #ddd;
  border-radius: 16px;
  background: #fff;
  font-size: 14px;
  cursor: pointer;
}

.logout-btn {
  padding: 6px 16px;
  border: 1px solid #ff6633;
  border-radius: 16px;
  background: #fff;
  color: #ff6633;
  font-size: 14px;
  cursor: pointer;
}

.intro-section {
  padding: 15px;
  background: #fff;
  border-top: 1px solid #f1f1f1;
}

.intro-section p {
  font-size: 14px;
  color: #666;
}

.intro-section .placeholder {
  color: #999;
  display: flex;
  align-items: center;
  gap: 4px;
}

:deep(.el-tabs) {
  background: #fff;
  margin-top: 10px;
}

:deep(.el-tabs__header) {
  margin: 0;
}

:deep(.el-tabs__content) {
  padding: 15px;
}

.blog-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.blog-item {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
}

.blog-item img {
  width: 100%;
  aspect-ratio: 1;
  object-fit: cover;
}

.blog-info {
  padding: 8px;
}

.blog-title {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.blog-stats {
  display: flex;
  gap: 12px;
  margin-top: 6px;
  font-size: 12px;
  color: #999;
}

.blog-stats span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.blog-stats img {
  width: 14px;
  height: 14px;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #999;
  font-size: 14px;
}

.follow-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
  max-height: calc(100vh - 200px);
  overflow-y: auto;
  position: relative;
  padding: 15px;
  box-sizing: border-box;
}

.refresh-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: height 0.3s ease;
  overflow: hidden;
}
.follow-list > * {
  flex-shrink: 0;
}

.pull-tip {
  color: #999;
  font-size: 14px;
  padding: 10px;
}

.loading-spinner {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #999;
  font-size: 14px;
  padding: 10px;
}

.loading,
.finished {
  text-align: center;
  padding: 15px;
  color: #999;
  font-size: 14px;
  width: 100%;
}
</style>
