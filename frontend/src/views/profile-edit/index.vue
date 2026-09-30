<template>
  <div class="profile-edit-page">
    <HeaderBar title="资料编辑" />

    <div class="edit-list">
      <!-- 基本信息 -->
      <div class="edit-group">
        <div class="edit-item">
          <span class="label">头像</span>
          <div class="value">
            <img :src="user.icon || '/imgs/icons/default-icon.png'" alt="">
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
        <div class="divider"></div>
        <div class="edit-item">
          <span class="label">昵称</span>
          <div class="value">
            <span>{{ user.nickName }}</span>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
        <div class="divider"></div>
        <div class="edit-item">
          <span class="label">个人介绍</span>
          <div class="value">
            <span class="text-ellipsis">{{ userDetail.introduce || '介绍一下自己' }}</span>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </div>

      <!-- 个人信息 -->
      <div class="edit-group">
        <div class="edit-item">
          <span class="label">性别</span>
          <div class="value">
            <span>{{ userDetail.gender || '选择' }}</span>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
        <div class="divider"></div>
        <div class="edit-item">
          <span class="label">城市</span>
          <div class="value">
            <span>{{ userDetail.city || '选择' }}</span>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
        <div class="divider"></div>
        <div class="edit-item">
          <span class="label">生日</span>
          <div class="value">
            <span>{{ userDetail.birthday || '添加' }}</span>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </div>

      <!-- 其他 -->
      <div class="edit-group">
        <div class="edit-item">
          <span class="label">我的积分</span>
          <div class="value">
            <span>查看积分</span>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
        <div class="divider"></div>
        <div class="edit-item">
          <span class="label">会员等级</span>
          <div class="value">
            <span class="vip-text">成为VIP尊享特权</span>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <FootBar :active-btn="4" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ArrowRight } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import HeaderBar from '@/components/HeaderBar.vue'
import FootBar from '@/components/FootBar.vue'

const userStore = useUserStore()

const user = ref({})
const userDetail = ref({})

onMounted(() => {
  user.value = userStore.userInfo || {}
  const stored = sessionStorage.getItem('userInfo')
  if (stored) {
    userDetail.value = JSON.parse(stored)
  }
})
</script>

<style scoped>
.profile-edit-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 50px;
}

.edit-list {
  padding-top: 10px;
}

.edit-group {
  background: #fff;
  margin-bottom: 10px;
  padding: 0 15px;
}

.edit-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
}

.label {
  font-size: 14px;
  color: #333;
}

.value {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.value img {
  width: 35px;
  height: 35px;
  border-radius: 50%;
}

.text-ellipsis {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.vip-text {
  color: #ff6633;
}

.divider {
  height: 1px;
  background: #f1f1f1;
}
</style>
