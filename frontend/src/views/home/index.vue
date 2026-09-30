<template>
  <div class="home-page">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="city-btn">
        杭州 <el-icon><ArrowDown /></el-icon>
      </div>
      <div class="search-input">
        <el-input 
          v-model="searchKey"
          size="small"
          placeholder="请输入商户名、地点"
          :prefix-icon="Search"
        />
      </div>
      <div class="user-icon" @click="toProfile">
        <el-icon size="20"><User /></el-icon>
      </div>
    </div>

    <!-- 类型列表 -->
    <div class="type-list">
      <div 
        v-for="type in shopTypes" 
        :key="type.id" 
        class="type-item"
        @click="toShopList(type)"
      >
        <div class="type-icon">
          <img :src="`/imgs/${type.icon}`" alt="">
        </div>
        <div class="type-name">{{ type.name }}</div>
      </div>
    </div>

    <!-- 博客列表 -->
    <div 
      ref="scrollContainer"
      class="blog-list" 
      @scroll="handleScroll"
    >
      <div class="blog-grid">
        <BlogCard 
          v-for="blog in blogs" 
          :key="blog.id"
          :blog="blog"
          @update="updateBlog"
        />
      </div>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-if="finished" class="finished">没有更多了</div>
    </div>

    <FootBar :active-btn="1" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, ArrowDown, User } from '@element-plus/icons-vue'
import { getShopTypes } from '@/api/shop'
import { getHotBlogs } from '@/api/blog'
import { useScrollLoad } from '@/utils/useScrollLoad'
import BlogCard from '@/components/BlogCard.vue'
import FootBar from '@/components/FootBar.vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const searchKey = ref('')
const shopTypes = ref([])
const blogs = ref([])
const currentPage = ref(1)

const loadMore = async () => {
  const data = await getHotBlogs(currentPage.value)
  if (data && data.length > 0) {
    data.forEach(blog => {
      blog.img = blog.images.split(',')[0]
    })
    blogs.value.push(...data)
    currentPage.value++
    return true
  }
  return false
}

const { loading, finished, scrollContainer, handleScroll } = useScrollLoad(loadMore)

const fetchShopTypes = async () => {
  try {
    shopTypes.value = await getShopTypes()
  } catch (error) {
    console.error('获取商户类型失败:', error)
  }
}

const toShopList = (type) => {
  router.push(`/shop-list?type=${type.id}&name=${encodeURIComponent(type.name)}`)
}

const toProfile = () => {
  if (userStore.isLoggedIn) {
    router.push('/profile')
  } else {
    router.push('/login')
  }
}

const updateBlog = (updatedBlog) => {
  const index = blogs.value.findIndex(b => b.id === updatedBlog.id)
  if (index !== -1) {
    updatedBlog.img = updatedBlog.images.split(',')[0]
    blogs.value[index] = updatedBlog
  }
}

onMounted(() => {
  fetchShopTypes()
})
</script>

<style scoped>
.home-page {
  padding-bottom: 50px;
}

.search-bar {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  background: #fff;
  gap: 10px;
}

.city-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #333;
  white-space: nowrap;
}

.search-input {
  flex: 1;
}

.search-input :deep(.el-input__inner) {
  border-radius: 20px;
  background: #f5f5f5;
  border: none;
}

.user-icon {
  color: #666;
  cursor: pointer;
}

.type-list {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
  padding: 15px;
  background: #fff;
}

.type-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
}

.type-icon {
  width: 44px;
  height: 44px;
}

.type-icon img {
  width: 100%;
  height: 100%;
}

.type-name {
  margin-top: 6px;
  font-size: 12px;
  color: #666;
}

.blog-list {
  height: calc(100vh - 200px);
  overflow-y: auto;
  padding: 10px;
}

.blog-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.loading,
.finished {
  text-align: center;
  padding: 15px;
  color: #999;
  font-size: 14px;
}
</style>
