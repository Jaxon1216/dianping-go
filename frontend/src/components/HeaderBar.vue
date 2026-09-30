<template>
  <div class="header-bar">
    <div class="header-back" @click="goBack">
      <el-icon size="20"><ArrowLeft /></el-icon>
    </div>
    <div class="header-title">{{ title }}</div>
    <div class="header-right">
      <slot name="right">
        <span v-if="showShare">...</span>
      </slot>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  showBack: {
    type: Boolean,
    default: true
  },
  showShare: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()

const goBack = () => {
  if (props.showBack) {
    router.back()
  }
}
</script>

<style scoped>
.header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  padding: 0 15px;
  background: #fff;
  border-bottom: 1px solid #f1f1f1;
}

.header-back {
  width: 40px;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.header-title {
  flex: 1;
  text-align: center;
  font-size: 16px;
  font-weight: 500;
}

.header-right {
  width: 40px;
  text-align: right;
  font-size: 14px;
  color: #666;
}
</style>
