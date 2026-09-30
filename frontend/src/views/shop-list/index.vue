<template>
  <div class="shop-list-page">
    <HeaderBar :title="typeName">
      <template #right>
        <el-icon size="20"><Search /></el-icon>
      </template>
    </HeaderBar>

    <!-- 排序栏 -->
    <div class="sort-bar">
      <el-dropdown trigger="click" @command="handleTypeChange">
        <span class="sort-item">
          {{ typeName }}<el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item 
              v-for="t in shopTypes" 
              :key="t.id" 
              :command="t"
            >
              {{ t.name }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <span class="sort-item" @click="sortBy('')">
        距离<el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </span>
      <span class="sort-item" @click="sortBy('comments')">
        人气<el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </span>
      <span class="sort-item" @click="sortBy('score')">
        评分<el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </span>
    </div>

    <!-- 商户列表 -->
    <div 
      ref="scrollContainer"
      class="shop-container"
      @scroll="handleScroll"
    >
      <ShopCard 
        v-for="shop in shops" 
        :key="shop.id" 
        :shop="shop" 
      />
      <div v-if="loading" class="loading">加载中...</div>
      <div v-if="finished" class="finished">没有更多了</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown, Search } from '@element-plus/icons-vue'
import { getShopTypes, getShopList } from '@/api/shop'
import { useScrollLoad } from '@/utils/useScrollLoad'
import HeaderBar from '@/components/HeaderBar.vue'
import ShopCard from '@/components/ShopCard.vue'

const route = useRoute()
const router = useRouter()

const typeName = ref('')
const shopTypes = ref([])
const shops = ref([])

const params = ref({
  typeId: 0,
  current: 1,
  sortBy: '',
  x: 120.149993,
  y: 30.334229
})

const loadMore = async () => {
  const data = await getShopList(params.value)
  if (data && data.length > 0) {
    data.forEach(shop => {
      shop.images = shop.images.split(',')[0]
    })
    shops.value.push(...data)
    params.value.current++
    return true
  }
  return false
}

const { loading, finished, scrollContainer, handleScroll, reset } = useScrollLoad(loadMore, { immediate: false })

const fetchShopTypes = async () => {
  try {
    shopTypes.value = await getShopTypes()
  } catch (error) {
    console.error('获取商户类型失败:', error)
  }
}

const handleTypeChange = (type) => {
  params.value.typeId = type.id
  typeName.value = type.name
  params.value.current = 1
  shops.value = []
  reset()
  loadMore()
}

const sortBy = (sortType) => {
  params.value.sortBy = sortType
  params.value.current = 1
  shops.value = []
  reset()
  loadMore()
}

onMounted(() => {
  params.value.typeId = parseInt(route.query.type) || 0
  typeName.value = route.query.name || '全部'
  fetchShopTypes()
  loadMore()
})
</script>

<style scoped>
.shop-list-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.sort-bar {
  display: flex;
  justify-content: space-around;
  align-items: center;
  height: 44px;
  background: #fff;
  border-bottom: 1px solid #f1f1f1;
}

.sort-item {
  font-size: 14px;
  color: #666;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.shop-container {
  flex: 1;
  overflow-y: auto;
  background: #f5f5f5;
}

.loading,
.finished {
  text-align: center;
  padding: 15px;
  color: #999;
  font-size: 14px;
  background: #fff;
}
</style>
