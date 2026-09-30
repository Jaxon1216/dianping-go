<template>
  <div class="shop-card" @click="toDetail">
    <div class="shop-img">
      <img :src="shop.images" alt="">
    </div>
    <div class="shop-info">
      <div class="shop-name">{{ shop.name }}</div>
      <div class="shop-rate">
        <el-rate
          :model-value="shop.score / 10"
          disabled
          text-color="#F63"
          show-score
          score-template="{value}分"
        />
        <span class="comments">{{ shop.comments }}条</span>
      </div>
      <div class="shop-area">
        <span>{{ shop.area }}</span>
        <span v-if="shop.distance">{{ formatDistance(shop.distance) }}</span>
      </div>
      <div class="shop-price">￥{{ shop.avgPrice }}/人</div>
      <div class="shop-address">
        <el-icon><MapLocation /></el-icon>
        <span>{{ shop.address }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { formatDistance } from '@/utils/format'

const props = defineProps({
  shop: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const toDetail = () => {
  router.push(`/shop-detail/${props.shop.id}`)
}
</script>

<style scoped>
.shop-card {
  display: flex;
  padding: 15px;
  background: #fff;
  border-bottom: 1px solid #f1f1f1;
}

.shop-img {
  width: 100px;
  height: 100px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
}

.shop-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.shop-info {
  flex: 1;
  margin-left: 12px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.shop-name {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.shop-rate {
  display: flex;
  align-items: center;
  gap: 10px;
}

.comments {
  font-size: 12px;
  color: #999;
}

.shop-area {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #666;
}

.shop-price {
  font-size: 14px;
  color: #ff6633;
  font-weight: 500;
}

.shop-address {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}

.shop-address span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
