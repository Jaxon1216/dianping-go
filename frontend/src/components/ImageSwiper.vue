<template>
  <div ref="container" class="image-swiper" @touchstart="onTouchStart" @touchmove="onTouchMove" @touchend="onTouchEnd">
    <div v-for="(img, index) in images" :key="index" class="swiper-item" :style="getItemStyle(index)">
      <img :src="img" alt="">
    </div>
    <div v-if="images.length > 1" class="swiper-indicator">
      <span v-for="(_, index) in images" :key="index" :class="{ active: currentIndex === index }" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  images: {
    type: Array,
    required: true
  },
  duration: {
    type: Number,
    default: 300
  }
})

const emit = defineEmits(['change'])

const container = ref(null)
const currentIndex = ref(0)
const containerWidth = ref(0)
const isDragging = ref(false)

// 触摸相关
const touchStart = ref({ x: 0, y: 0 })
const touchMove = ref({ x: 0, y: 0 })
const translateX = ref(0)

const getItemStyle = (index) => {
  const offset = (index - currentIndex.value) * containerWidth.value + translateX.value
  return {
    transform: `translate3d(${offset}px, 0, 0)`,
    transition: isDragging.value ? 'none' : `transform ${props.duration}ms`
  }
}

const updateWidth = () => {
  if (container.value) {
    containerWidth.value = container.value.offsetWidth
  }
}

const onTouchStart = (e) => {
  isDragging.value = true
  touchStart.value = {
    x: e.touches[0].pageX,
    y: e.touches[0].pageY
  }
}

const onTouchMove = (e) => {
  if (!isDragging.value) return

  const deltaX = e.touches[0].pageX - touchStart.value.x
  const deltaY = e.touches[0].pageY - touchStart.value.y

  // 水平滑动才处理
  if (Math.abs(deltaX) > Math.abs(deltaY)) {
    e.preventDefault()

    // 边界阻力效果
    if ((currentIndex.value === 0 && deltaX > 0) ||
      (currentIndex.value === props.images.length - 1 && deltaX < 0)) {
      translateX.value = deltaX * 0.3
    } else {
      translateX.value = deltaX
    }
  }
}

const onTouchEnd = (e) => {
  if (!isDragging.value) return

  isDragging.value = false
  const deltaX = touchMove.value.x - touchStart.value.x
  const threshold = containerWidth.value * 0.2

  if (Math.abs(translateX.value) > threshold) {
    if (translateX.value < 0 && currentIndex.value < props.images.length - 1) {
      currentIndex.value++
    } else if (translateX.value > 0 && currentIndex.value > 0) {
      currentIndex.value--
    }
  }

  translateX.value = 0
  emit('change', currentIndex.value)
}

onMounted(() => {
  updateWidth()
  window.addEventListener('resize', updateWidth)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateWidth)
})
</script>

<style scoped>
.image-swiper {
  position: relative;
  width: 100%;
  height: 300px;
  overflow: hidden;
  background: #f5f5f5;
}

.swiper-item {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  will-change: transform;
}

.swiper-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.swiper-indicator {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 6px;
}

.swiper-indicator span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  transition: all 0.3s;
}

.swiper-indicator span.active {
  width: 12px;
  border-radius: 3px;
  background: #fff;
}
</style>
