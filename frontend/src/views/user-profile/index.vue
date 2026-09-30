<template>
  <div class="user-profile-page">
    <HeaderBar />

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
      </div>
      <button 
        class="follow-btn"
        :class="{ followed }"
        @click="handleFollow"
      >
        {{ followed ? '取消关注' : '关注' }}
      </button>
    </div>

    <!-- 简介 -->
    <div class="intro-section">
      <p v-if="userDetail.introduce">{{ userDetail.introduce }}</p>
      <p v-else class="placeholder">这个人很懒，什么都没有留下</p>
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
              <div class="blog-title" v-html="blog.title"></div>
              <div class="blog-stats">
                <span><img src="/imgs/thumbup.png" alt=""> {{ blog.liked }}</span>
                <span><el-icon><ChatDotRound /></el-icon> {{ blog.comments }}</span>
              </div>
            </div>
          </div>
          <div v-if="blogs.length === 0" class="empty">暂无笔记</div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="共同关注" name="common">
        <div class="common-follows">
          <p class="section-title">你们都关注了：</p>
          <div 
            v-for="u in commonFollows" 
            :key="u.id"
            class="follow-item"
          >
            <img 
              class="avatar" 
              :src="u.icon || '/imgs/icons/default-icon.png'" 
              alt=""
              @click="toUserProfile(u.id)"
            >
            <div class="info">
              <div class="name">{{ u.nickName }}</div>
            </div>
            <button class="view-btn" @click="toUserProfile(u.id)">
              去主页看看
            </button>
          </div>
          <div v-if="commonFollows.length === 0" class="empty">暂无共同关注</div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <FootBar :active-btn="0" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Location, ChatDotRound } from '@element-plus/icons-vue'
import { getUserById, getUserDetail } from '@/api/user'
import { getUserBlogs } from '@/api/blog'
import { isFollowed, followUser, getCommonFollows } from '@/api/follow'
import HeaderBar from '@/components/HeaderBar.vue'
import FootBar from '@/components/FootBar.vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const user = ref({})
const userDetail = ref({})
const blogs = ref([])
const activeTab = ref('blogs')
const followed = ref(false)
const commonFollows = ref([])

const fetchUserInfo = async (id) => {
  try {
    user.value = await getUserById(id)
    fetchUserDetail(id)
    fetchUserBlogs(id)
    checkFollowStatus(id)
  } catch (error) {
    ElMessage.error('获取用户信息失败')
  }
}

const fetchUserDetail = async (id) => {
  try {
    const data = await getUserDetail(id)
    if (data) {
      userDetail.value = data
    }
  } catch (error) {
    console.error('获取用户详情失败:', error)
  }
}

const fetchUserBlogs = async (id) => {
  try {
    blogs.value = await getUserBlogs(id)
  } catch (error) {
    console.error('获取笔记失败:', error)
  }
}

const checkFollowStatus = async (id) => {
  try {
    followed.value = await isFollowed(id)
  } catch (error) {
    console.error('获取关注状态失败:', error)
  }
}

const fetchCommonFollows = async (id) => {
  try {
    commonFollows.value = await getCommonFollows(id)
  } catch (error) {
    console.error('获取共同关注失败:', error)
  }
}

const handleFollow = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  try {
    await followUser(user.value.id, !followed.value)
    followed.value = !followed.value
    ElMessage.success(followed.value ? '已关注' : '已取消关注')
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleTabChange = (tab) => {
  if (tab === 'common') {
    fetchCommonFollows(user.value.id)
  }
}

const toBlogDetail = (id) => {
  router.push(`/blog-detail/${id}`)
}

const toUserProfile = (id) => {
  router.push(`/user/${id}`)
}

onMounted(() => {
  const userId = route.params.id
  if (userId) {
    fetchUserInfo(userId)
  }
})
</script>

<style scoped>
.user-profile-page {
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
  cursor: pointer;
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

.follow-btn {
  padding: 6px 20px;
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

.common-follows {
  padding: 0 15px;
}

.section-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 15px;
}

.follow-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f1f1f1;
}

.follow-item .avatar {
  width: 44px;
  height: 44px;
}

.follow-item .info {
  flex: 1;
}

.follow-item .name {
  font-size: 14px;
  color: #333;
}

.view-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 12px;
  background: #fff;
  font-size: 12px;
  color: #666;
  cursor: pointer;
}
</style>
