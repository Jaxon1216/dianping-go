<template>
  <div class="foot-bar">
    <div 
      class="foot-item" 
      :class="{ active: activeBtn === 1 }" 
      @click="toPage(1)"
    >
      <el-icon size="20"><HomeFilled /></el-icon>
      <span class="foot-text">首页</span>
    </div>
    <div class="foot-item" @click="toPage(2)">
      <el-icon size="20"><MapLocation /></el-icon>
      <span class="foot-text">地图</span>
    </div>
    <div class="foot-item" @click="toPage(0)">
      <img class="add-btn" src="/imgs/add.png" alt="发布">
    </div>
    <div class="foot-item" @click="toPage(3)">
      <el-icon size="20"><ChatDotRound /></el-icon>
      <span class="foot-text">消息</span>
    </div>
    <div 
      class="foot-item" 
      :class="{ active: activeBtn === 4 }" 
      @click="toPage(4)"
    >
      <el-icon size="20"><UserFilled /></el-icon>
      <span class="foot-text">我的</span>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const props = defineProps({
  activeBtn: {
    type: Number,
    default: 1
  }
})

const router = useRouter()
const userStore = useUserStore()

const toPage = (index) => {
  switch (index) {
    case 0:
      if (!userStore.isLoggedIn) {
        ElMessage.warning('请先登录')
        router.push('/login')
        return
      }
      router.push('/blog-edit')
      break
    case 1:
      router.push('/')
      break
    case 4:
      if (!userStore.isLoggedIn) {
        router.push('/login')
        return
      }
      router.push('/profile')
      break
    default:
      ElMessage.info('功能开发中')
  }
}
</script>

<style scoped>
.foot-bar {
  position: fixed;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  max-width: 540px;
  height: 50px;
  background: #fff;
  border-top: 1px solid #e1e2e3;
  display: flex;
  justify-content: space-around;
  align-items: center;
  z-index: 100;
}

.foot-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
  cursor: pointer;
  flex: 1;
}

.foot-item.active {
  color: #ff6633;
}

.foot-text {
  font-size: 10px;
  margin-top: 2px;
}

.add-btn {
  width: 40px;
  height: 40px;
}
</style>
