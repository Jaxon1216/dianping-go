<template>
  <div class="shop-detail-page">
    <HeaderBar :show-share="true">
      <template #right>
        <span>...</span>
      </template>
    </HeaderBar>

    <div class="content">
      <!-- 商户基本信息 -->
      <div class="shop-header">
        <h1 class="shop-name">{{ shop.name }}</h1>
        <div class="shop-rate">
          <el-rate
            :model-value="shop.score / 10"
            disabled
            text-color="#F63"
            show-score
          />
          <span class="comments">{{ shop.comments }}条</span>
        </div>
        <div class="shop-rank">
          <img src="/imgs/bd.png" width="63" height="20" alt="榜单">
          <span>拱墅区好评榜第3名</span>
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>

      <!-- 图片展示 -->
      <div v-if="shop.images?.length" class="shop-images">
        <img 
          v-for="(img, index) in shop.images.slice(0, 3)" 
          :key="index"
          :src="img" 
          alt=""
        >
      </div>

      <!-- 地址信息 -->
      <div class="shop-address">
        <el-icon><MapLocation /></el-icon>
        <span>{{ shop.address }}</span>
        <span class="divider">|</span>
        <el-icon><Phone /></el-icon>
        <el-icon><Position /></el-icon>
      </div>

      <!-- 营业时间 -->
      <div class="shop-hours">
        <el-icon><Clock /></el-icon>
        <span>营业时间</span>
        <span>{{ shop.openHours }}</span>
        <el-icon class="arrow"><ArrowRight /></el-icon>
      </div>

      <!-- 优惠券 -->
      <div v-if="vouchers.length" class="voucher-section">
        <div class="section-title">
          <span class="tag">券</span>
          <span>代金券</span>
        </div>
        <div 
          v-for="voucher in validVouchers" 
          :key="voucher.id"
          class="voucher-card"
        >
          <div class="voucher-left">
            <div class="voucher-title">{{ voucher.title }}</div>
            <div class="voucher-subtitle">{{ voucher.subTitle }}</div>
            <div class="voucher-price">
              <span class="price">￥{{ formatPrice(voucher.payValue) }}</span>
              <span class="discount">{{ ((voucher.payValue * 10) / voucher.actualValue).toFixed(1) }}折</span>
            </div>
          </div>
          <div class="voucher-right">
            <template v-if="voucher.type === 1">
              <button 
                class="seckill-btn"
                :disabled="isSeckillNotBegin(voucher.beginTime) || isSeckillEnd(voucher.endTime) || voucher.stock < 1"
                @click="handleSeckill(voucher)"
              >
                {{ getSeckillBtnText(voucher) }}
              </button>
              <div class="seckill-info">
                <span v-if="voucher.stock > 0">剩余 {{ voucher.stock }} 张</span>
                <span v-else>已售罄</span>
              </div>
              <div class="seckill-time">{{ formatSeckillTime(voucher.beginTime, voucher.endTime) }}</div>
            </template>
            <button v-else class="buy-btn">抢购</button>
          </div>
        </div>
      </div>

      <!-- 评价区域 -->
      <div class="comments-section">
        <div class="section-header">
          <span>网友评价 <span class="count">(119)</span></span>
          <el-icon><ArrowRight /></el-icon>
        </div>
        <div class="comment-tags">
          <span v-for="tag in commentTags" :key="tag" class="tag">{{ tag }}</span>
        </div>
      </div>

      <div class="copyright">copyright ©2024 hmdp.com</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowRight, MapLocation, Phone, Position, Clock } from '@element-plus/icons-vue'
import { getShopById } from '@/api/shop'
import { getVoucherList, seckillVoucher } from '@/api/voucher'
import { formatPrice, formatSeckillTime, isSeckillNotBegin, isSeckillEnd } from '@/utils/format'
import HeaderBar from '@/components/HeaderBar.vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const shop = ref({})
const vouchers = ref([])

const commentTags = ['味道赞(19)', '牛肉赞(16)', '菜品不错(11)', '回头客(4)', '分量足(4)']

const validVouchers = computed(() => {
  return vouchers.value.filter(v => !isSeckillEnd(v.endTime))
})

const getSeckillBtnText = (voucher) => {
  if (isSeckillNotBegin(voucher.beginTime)) return '即将开始'
  if (isSeckillEnd(voucher.endTime)) return '已结束'
  if (voucher.stock < 1) return '已售罄'
  return '限时抢购'
}

const fetchShopDetail = async (id) => {
  try {
    const data = await getShopById(id)
    data.images = data.images.split(',')
    shop.value = data
  } catch (error) {
    ElMessage.error('获取商户详情失败')
  }
}

const fetchVouchers = async (shopId) => {
  try {
    vouchers.value = await getVoucherList(shopId)
  } catch (error) {
    console.error('获取优惠券失败:', error)
  }
}

const handleSeckill = async (voucher) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  if (isSeckillNotBegin(voucher.beginTime)) {
    ElMessage.warning('优惠券抢购尚未开始')
    return
  }
  
  if (voucher.stock < 1) {
    ElMessage.warning('库存不足')
    return
  }
  
  try {
    const orderId = await seckillVoucher(voucher.id)
    ElMessage.success(`抢购成功，订单id：${orderId}`)
    // 刷新优惠券列表
    fetchVouchers(shop.value.id)
  } catch (error) {
    // 错误已在拦截器处理
  }
}

onMounted(() => {
  const shopId = route.params.id
  if (shopId) {
    fetchShopDetail(shopId)
    fetchVouchers(shopId)
  }
})
</script>

<style scoped>
.shop-detail-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.content {
  padding-bottom: 20px;
}

.shop-header {
  padding: 15px;
  background: #fff;
}

.shop-name {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 10px;
}

.shop-rate {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.comments {
  font-size: 12px;
  color: #999;
}

.shop-rank {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.shop-images {
  display: flex;
  gap: 8px;
  padding: 0 15px 15px;
  background: #fff;
  overflow-x: auto;
}

.shop-images img {
  width: 120px;
  height: 120px;
  border-radius: 8px;
  object-fit: cover;
}

.shop-address,
.shop-hours {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 15px;
  background: #fff;
  margin-top: 10px;
  font-size: 14px;
  color: #333;
}

.divider {
  color: #e1e2e3;
  margin: 0 8px;
}

.arrow {
  margin-left: auto;
  color: #999;
}

.voucher-section {
  background: #fff;
  margin-top: 10px;
  padding: 15px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 15px;
  font-weight: 500;
}

.section-title .tag {
  background: #ff6633;
  color: #fff;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.voucher-card {
  display: flex;
  border: 1px solid #ffe4d9;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 10px;
}

.voucher-left {
  flex: 1;
  padding: 12px;
  background: linear-gradient(to right, #fff5f0, #fff);
}

.voucher-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 4px;
}

.voucher-subtitle {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
}

.voucher-price {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.price {
  font-size: 20px;
  color: #ff6633;
  font-weight: 600;
}

.discount {
  font-size: 12px;
  color: #ff6633;
  background: #ffe4d9;
  padding: 2px 6px;
  border-radius: 4px;
}

.voucher-right {
  width: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fff;
  border-left: 1px dashed #ffe4d9;
  padding: 10px;
}

.seckill-btn,
.buy-btn {
  width: 80px;
  height: 32px;
  border: none;
  border-radius: 16px;
  background: #ff6633;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
}

.seckill-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.seckill-info {
  font-size: 12px;
  color: #999;
  margin-top: 6px;
}

.seckill-info span {
  color: #ff6633;
}

.seckill-time {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
}

.comments-section {
  background: #fff;
  margin-top: 10px;
  padding: 15px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.count {
  color: #999;
  font-size: 14px;
}

.comment-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.comment-tags .tag {
  padding: 6px 12px;
  background: #f5f5f5;
  border-radius: 12px;
  font-size: 12px;
  color: #666;
}

.copyright {
  text-align: center;
  padding: 20px;
  font-size: 12px;
  color: #999;
}
</style>
