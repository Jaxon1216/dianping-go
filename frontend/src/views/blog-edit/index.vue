<template>
  <div class="blog-edit-page">
    <!-- 头部 -->
    <div class="header">
      <span class="cancel-btn" @click="goBack">取消</span>
      <span class="title">发笔记</span>
      <button class="publish-btn" @click="publish">发布</button>
    </div>

    <!-- 图片上传 -->
    <div class="upload-section">
      <input 
        ref="fileInput"
        type="file" 
        accept="image/*"
        style="display: none"
        @change="handleFileChange"
      >
      <div class="upload-btn" @click="openFileDialog">
        <el-icon size="24"><Camera /></el-icon>
        <span>上传照片</span>
      </div>
      <div class="image-list">
        <div 
          v-for="(img, index) in fileList" 
          :key="index"
          class="image-item"
        >
          <img :src="img" alt="">
          <el-icon class="delete-btn" @click="deleteImage(index)"><Close /></el-icon>
        </div>
      </div>
    </div>

    <!-- 标题输入 -->
    <div class="input-section">
      <input 
        v-model="form.title"
        type="text"
        placeholder="填写标题更容易上首页哦~"
        class="title-input"
      >
    </div>

    <!-- 内容输入 -->
    <div class="input-section">
      <textarea 
        v-model="form.content"
        placeholder="最近打卡了什么地方，有什么新奇体验呢？"
        class="content-input"
        rows="6"
      ></textarea>
    </div>

    <!-- 关联商户 -->
    <div class="shop-section" @click="showShopDialog = true">
      <span class="label">关联商户</span>
      <span v-if="selectedShop.name" class="value">{{ selectedShop.name }}</span>
      <span v-else class="placeholder">去选择</span>
      <el-icon><ArrowRight /></el-icon>
    </div>

    <!-- 商户选择弹窗 -->
    <div v-if="showShopDialog" class="dialog-mask" @click="showShopDialog = false"></div>
    <transition name="slide-up">
      <div v-if="showShopDialog" class="shop-dialog">
        <div class="dialog-header">
          <span>关联商户</span>
          <el-icon @click="showShopDialog = false"><Close /></el-icon>
        </div>
        <div class="search-bar">
          <span class="city">杭州 <el-icon><ArrowDown /></el-icon></span>
          <div class="search-input">
            <el-icon @click="searchShops"><Search /></el-icon>
            <input 
              v-model="shopName"
              type="text"
              placeholder="搜索商户名称"
              @keyup.enter="searchShops"
            >
          </div>
        </div>
        <div class="shop-list">
          <div 
            v-for="shop in shops" 
            :key="shop.id"
            class="shop-item"
            @click="selectShop(shop)"
          >
            <div class="shop-name">{{ shop.name }}</div>
            <div class="shop-area">{{ shop.area }}</div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Camera, Close, ArrowRight, ArrowDown, Search } from '@element-plus/icons-vue'
import { publishBlog } from '@/api/blog'
import { searchShops as searchShopsApi } from '@/api/shop'
import { uploadBlogImage, deleteBlogImage } from '@/api/upload'
import { ElMessage } from 'element-plus'

const router = useRouter()

const fileInput = ref(null)
const fileList = ref([])
const form = ref({
  title: '',
  content: ''
})
const selectedShop = ref({})
const showShopDialog = ref(false)
const shopName = ref('')
const shops = ref([])

const openFileDialog = () => {
  fileInput.value?.click()
}

const handleFileChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  
  try {
    const path = await uploadBlogImage(file)
    fileList.value.push('/imgs' + path)
    // 清空input，允许重复选择同一文件
    e.target.value = ''
  } catch (error) {
    ElMessage.error('上传失败')
  }
}

const deleteImage = async (index) => {
  try {
    await deleteBlogImage(fileList.value[index])
    fileList.value.splice(index, 1)
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const searchShops = async () => {
  try {
    shops.value = await searchShopsApi(shopName.value)
  } catch (error) {
    console.error('搜索商户失败:', error)
  }
}

const selectShop = (shop) => {
  selectedShop.value = shop
  showShopDialog.value = false
}

const publish = async () => {
  if (!form.value.title && !form.value.content) {
    ElMessage.warning('请填写标题或内容')
    return
  }
  
  if (fileList.value.length === 0) {
    ElMessage.warning('请至少上传一张图片')
    return
  }
  
  try {
    await publishBlog({
      ...form.value,
      images: fileList.value.join(','),
      shopId: selectedShop.value.id
    })
    ElMessage.success('发布成功')
    router.push('/profile')
  } catch (error) {
    ElMessage.error('发布失败')
  }
}

const goBack = () => {
  router.back()
}

onMounted(() => {
  searchShops()
})
</script>

<style scoped>
.blog-edit-page {
  min-height: 100vh;
  background: #fff;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 44px;
  padding: 0 15px;
  border-bottom: 1px solid #f1f1f1;
}

.cancel-btn {
  font-size: 14px;
  color: #666;
  cursor: pointer;
}

.title {
  font-size: 16px;
  font-weight: 500;
}

.publish-btn {
  padding: 6px 16px;
  border: none;
  border-radius: 16px;
  background: #ff6633;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
}

.publish-btn:disabled {
  background: #ccc;
}

.upload-section {
  display: flex;
  gap: 10px;
  padding: 15px;
  overflow-x: auto;
}

.upload-btn {
  width: 80px;
  height: 80px;
  border: 1px dashed #ddd;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
  cursor: pointer;
  flex-shrink: 0;
}

.upload-btn span {
  font-size: 12px;
  margin-top: 4px;
}

.image-list {
  display: flex;
  gap: 10px;
}

.image-item {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
}

.image-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.delete-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 18px;
  height: 18px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.input-section {
  padding: 0 15px;
  margin-bottom: 15px;
}

.title-input,
.content-input {
  width: 100%;
  border: none;
  outline: none;
  font-size: 14px;
  resize: none;
}

.title-input {
  font-size: 16px;
  font-weight: 500;
  padding: 10px 0;
  border-bottom: 1px solid #f1f1f1;
}

.content-input {
  line-height: 1.8;
}

.shop-section {
  display: flex;
  align-items: center;
  padding: 15px;
  margin: 0 15px;
  background: #f9f9f9;
  border-radius: 8px;
  cursor: pointer;
}

.shop-section .label {
  font-size: 14px;
  color: #333;
}

.shop-section .value {
  flex: 1;
  text-align: right;
  font-size: 14px;
  color: #333;
  margin-right: 8px;
}

.shop-section .placeholder {
  flex: 1;
  text-align: right;
  font-size: 14px;
  color: #999;
  margin-right: 8px;
}

.dialog-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 100;
}

.shop-dialog {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 540px;
  height: 70vh;
  background: #fff;
  border-radius: 16px 16px 0 0;
  z-index: 101;
  display: flex;
  flex-direction: column;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f1f1f1;
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 15px;
  border-bottom: 1px solid #f1f1f1;
}

.city {
  font-size: 14px;
  color: #333;
  white-space: nowrap;
}

.search-input {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f5f5f5;
  padding: 8px 12px;
  border-radius: 20px;
}

.search-input input {
  flex: 1;
  border: none;
  background: transparent;
  outline: none;
  font-size: 14px;
}

.shop-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 15px;
}

.shop-item {
  padding: 15px 0;
  border-bottom: 1px solid #f1f1f1;
  cursor: pointer;
}

.shop-name {
  font-size: 14px;
  color: #333;
  margin-bottom: 4px;
}

.shop-area {
  font-size: 12px;
  color: #999;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translate(-50%, 100%);
}
</style>
