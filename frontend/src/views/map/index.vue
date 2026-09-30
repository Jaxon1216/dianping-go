<template>
  <div class="map-page">
    <header class="map-header">
      <div class="header-main">
        <div>
          <h1>地图找店</h1>
          <button class="location-status" type="button" @click="requestLocation">
            <el-icon><Location /></el-icon>
            <span>{{ locationText }}</span>
          </button>
        </div>
        <div class="shop-total">
          <strong>{{ filteredShops.length }}</strong>
          <span>家商户</span>
        </div>
      </div>

      <div class="category-tabs" role="tablist" aria-label="商户类型">
        <button
          v-for="type in availableTypes"
          :key="type.id"
          class="category-tab"
          :class="{ active: selectedTypeId === type.id }"
          type="button"
          role="tab"
          :aria-selected="selectedTypeId === type.id"
          @click="selectType(type.id)"
        >
          {{ type.name }}
        </button>
      </div>
    </header>

    <main class="map-stage">
      <div ref="mapContainer" class="map-canvas" aria-label="附近商户地图" />

      <div class="map-summary">
        <span class="summary-dot" />
        {{ selectedCategoryName }}
      </div>

      <div class="map-actions">
        <el-tooltip content="定位到当前位置" placement="left">
          <button
            class="map-action"
            type="button"
            :disabled="locationState === 'locating'"
            aria-label="定位到当前位置"
            @click="requestLocation"
          >
            <el-icon :class="{ 'is-loading': locationState === 'locating' }">
              <Loading v-if="locationState === 'locating'" />
              <Aim v-else />
            </el-icon>
          </button>
        </el-tooltip>
        <el-tooltip content="显示全部商户" placement="left">
          <button
            class="map-action"
            type="button"
            aria-label="显示全部商户"
            @click="fitVisibleShops"
          >
            <el-icon><FullScreen /></el-icon>
          </button>
        </el-tooltip>
      </div>

      <div v-if="loading" class="map-state">
        <el-icon class="is-loading" size="22"><Loading /></el-icon>
        <span>正在加载商户</span>
      </div>

      <div v-else-if="errorText" class="map-state map-state-error">
        <el-icon size="22"><WarningFilled /></el-icon>
        <span>{{ errorText }}</span>
        <button type="button" @click="loadShops">重试</button>
      </div>

      <div v-else-if="filteredShops.length === 0" class="map-state">
        <el-icon size="22"><MapLocation /></el-icon>
        <span>该类型暂无商户</span>
      </div>

      <button
        v-if="selectedShop"
        class="selected-shop"
        type="button"
        @click="openShop(selectedShop.id)"
      >
        <img :src="selectedShop.cover" :alt="selectedShop.name">
        <span class="selected-shop-info">
          <span class="selected-shop-heading">
            <strong>{{ selectedShop.name }}</strong>
            <span>{{ formatDistance(selectedShop.distance) }}</span>
          </span>
          <span class="selected-shop-meta">
            {{ getTypeName(selectedShop.typeId) }}
            <i />
            {{ selectedShop.score / 10 }} 分
            <i />
            ￥{{ selectedShop.avgPrice || '--' }}/人
          </span>
          <span class="selected-shop-address">{{ selectedShop.address }}</span>
        </span>
        <el-icon class="detail-arrow"><ArrowRight /></el-icon>
      </button>

      <div v-if="tileError" class="tile-warning">
        底图加载失败，商户坐标仍可使用
      </div>
    </main>

    <FootBar :active-btn="2" />
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import {
  Aim,
  ArrowRight,
  FullScreen,
  Loading,
  Location,
  MapLocation,
  WarningFilled
} from '@element-plus/icons-vue'
import { getShopList, getShopTypes } from '@/api/shop'
import { formatDistance } from '@/utils/format'
import FootBar from '@/components/FootBar.vue'

const DEFAULT_LOCATION = {
  lat: 30.334229,
  lng: 120.149993
}

const MARKER_COLORS = ['#e85d35', '#147d76', '#d19a2a', '#3b6692', '#a64e65']

const router = useRouter()
const mapContainer = ref(null)
const shopTypes = ref([])
const shops = ref([])
const selectedTypeId = ref(0)
const selectedShopId = ref(null)
const loading = ref(true)
const errorText = ref('')
const tileError = ref(false)
const location = ref({ ...DEFAULT_LOCATION })
const locationState = ref('fallback')

let map
let markerLayer
let resizeObserver
let markerByShopId = new Map()

const availableTypes = computed(() => [
  { id: 0, name: '全部' },
  ...shopTypes.value
])

const selectedCategoryName = computed(() => {
  if (selectedTypeId.value === 0) return '全部商户'
  return getTypeName(selectedTypeId.value)
})

const locationText = computed(() => {
  const labels = {
    locating: '正在获取位置',
    located: '已定位到当前位置',
    denied: '定位未授权 · 杭州',
    unavailable: '定位不可用 · 杭州',
    fallback: '杭州 · 拱墅区'
  }
  return labels[locationState.value]
})

const filteredShops = computed(() => {
  const visible = selectedTypeId.value === 0
    ? shops.value
    : shops.value.filter(shop => shop.typeId === selectedTypeId.value)

  return visible
    .map(shop => ({
      ...shop,
      distance: getDistance(location.value, {
        lat: Number(shop.y),
        lng: Number(shop.x)
      })
    }))
    .sort((a, b) => a.distance - b.distance)
})

const selectedShop = computed(() => (
  filteredShops.value.find(shop => shop.id === selectedShopId.value) || null
))

const getTypeName = (typeId) => (
  shopTypes.value.find(type => type.id === typeId)?.name || '商户'
)

const getMarkerColor = (shop) => {
  const typeIndex = shopTypes.value.findIndex(type => type.id === shop.typeId)
  return MARKER_COLORS[Math.max(typeIndex, 0) % MARKER_COLORS.length]
}

const getDistance = (from, to) => {
  const earthRadius = 6371000
  const toRadians = degrees => degrees * Math.PI / 180
  const deltaLat = toRadians(to.lat - from.lat)
  const deltaLng = toRadians(to.lng - from.lng)
  const lat1 = toRadians(from.lat)
  const lat2 = toRadians(to.lat)
  const value = Math.sin(deltaLat / 2) ** 2
    + Math.cos(lat1) * Math.cos(lat2) * Math.sin(deltaLng / 2) ** 2
  const boundedValue = Math.min(1, Math.max(0, value))
  return earthRadius * 2 * Math.atan2(
    Math.sqrt(boundedValue),
    Math.sqrt(1 - boundedValue)
  )
}

const hasCoordinates = shop => (
  Number.isFinite(Number(shop.x)) && Number.isFinite(Number(shop.y))
)

const normalizeShop = (shop) => ({
  ...shop,
  cover: shop.images?.split(',')[0] || ''
})

const initMap = () => {
  map = L.map(mapContainer.value, {
    attributionControl: true,
    zoomControl: false
  }).setView([location.value.lat, location.value.lng], 14)

  const tileLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap',
    maxZoom: 19
  })
    .on('tileerror', () => {
      tileError.value = true
    })
    .on('tileload', () => {
      tileError.value = false
    })
    .addTo(map)

  L.control.zoom({ position: 'topright' }).addTo(map)
  markerLayer = L.layerGroup().addTo(map)

  resizeObserver = new ResizeObserver(() => {
    map?.invalidateSize({ pan: false })
  })
  resizeObserver.observe(mapContainer.value)
}

const createPopup = (shop) => {
  const button = document.createElement('button')
  button.className = 'shop-popup-button'
  button.type = 'button'

  const name = document.createElement('strong')
  name.textContent = shop.name
  const address = document.createElement('span')
  address.textContent = shop.address

  button.append(name, address)
  button.addEventListener('click', () => openShop(shop.id))
  return button
}

const updateMarkerStyles = () => {
  markerByShopId.forEach((marker, shopId) => {
    const shop = filteredShops.value.find(item => item.id === shopId)
    if (!shop) return
    const isSelected = selectedShopId.value === shopId
    marker.setStyle({
      radius: isSelected ? 10 : 7,
      weight: isSelected ? 4 : 2,
      fillOpacity: isSelected ? 1 : 0.82,
      color: isSelected ? '#ffffff' : getMarkerColor(shop),
      fillColor: getMarkerColor(shop)
    })
    if (isSelected) marker.bringToFront()
  })
}

const renderMarkers = () => {
  if (!map || !markerLayer) return

  markerLayer.clearLayers()
  markerByShopId = new Map()

  L.circleMarker([location.value.lat, location.value.lng], {
    radius: 8,
    color: '#ffffff',
    weight: 3,
    fillColor: '#2369d8',
    fillOpacity: 1
  })
    .bindTooltip('当前位置', { direction: 'top', offset: [0, -8] })
    .addTo(markerLayer)

  filteredShops.value.filter(hasCoordinates).forEach(shop => {
    const color = getMarkerColor(shop)
    const marker = L.circleMarker([Number(shop.y), Number(shop.x)], {
      radius: 7,
      color,
      weight: 2,
      fillColor: color,
      fillOpacity: 0.82
    })
      .bindPopup(createPopup(shop), {
        closeButton: false,
        offset: [0, -4],
        minWidth: 190
      })
      .on('click', () => {
        selectedShopId.value = shop.id
      })
      .addTo(markerLayer)

    markerByShopId.set(shop.id, marker)
  })

  updateMarkerStyles()
}

const fitVisibleShops = () => {
  if (!map) return
  const coordinates = filteredShops.value
    .filter(hasCoordinates)
    .map(shop => [Number(shop.y), Number(shop.x)])

  if (coordinates.length === 0) {
    map.setView([location.value.lat, location.value.lng], 14)
    return
  }

  map.fitBounds(L.latLngBounds(coordinates).pad(0.18), {
    animate: true,
    maxZoom: 15
  })
}

const loadShops = async () => {
  loading.value = true
  errorText.value = ''

  try {
    if (shopTypes.value.length === 0) {
      shopTypes.value = await getShopTypes()
    }

    const results = await Promise.allSettled(
      shopTypes.value.map(type => getShopList({
        typeId: type.id,
        current: 1,
        x: location.value.lng,
        y: location.value.lat
      }))
    )

    const uniqueShops = new Map()
    results.forEach(result => {
      if (result.status !== 'fulfilled' || !Array.isArray(result.value)) return
      result.value.forEach(shop => {
        if (hasCoordinates(shop)) uniqueShops.set(shop.id, normalizeShop(shop))
      })
    })
    shops.value = [...uniqueShops.values()]

    if (shops.value.length === 0) {
      errorText.value = '暂时没有可展示的商户'
      selectedShopId.value = null
    } else {
      selectedShopId.value = filteredShops.value[0]?.id || null
    }

    await nextTick()
    renderMarkers()
    fitVisibleShops()
  } catch (error) {
    console.error('加载地图商户失败:', error)
    errorText.value = '商户加载失败'
  } finally {
    loading.value = false
  }
}

const selectType = async (typeId) => {
  selectedTypeId.value = typeId
  await nextTick()
  selectedShopId.value = filteredShops.value[0]?.id || null
  renderMarkers()
  fitVisibleShops()
}

const requestLocation = () => {
  if (!navigator.geolocation) {
    locationState.value = 'unavailable'
    return
  }

  locationState.value = 'locating'
  navigator.geolocation.getCurrentPosition(
    position => {
      location.value = {
        lat: position.coords.latitude,
        lng: position.coords.longitude
      }
      locationState.value = 'located'
      renderMarkers()
      map?.setView([location.value.lat, location.value.lng], 15, { animate: true })
    },
    error => {
      locationState.value = error.code === error.PERMISSION_DENIED ? 'denied' : 'unavailable'
      location.value = { ...DEFAULT_LOCATION }
      renderMarkers()
    },
    {
      enableHighAccuracy: true,
      timeout: 8000,
      maximumAge: 60000
    }
  )
}

const openShop = (shopId) => {
  router.push(`/shop-detail/${shopId}`)
}

watch(selectedShopId, updateMarkerStyles)

onMounted(async () => {
  await nextTick()
  initMap()
  await loadShops()
  requestLocation()
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  map?.remove()
  map = null
  markerLayer = null
  markerByShopId.clear()
})
</script>

<style scoped>
.map-page {
  --map-orange: #e85d35;
  --map-ink: #202426;
  --map-muted: #70767a;
  height: 100vh;
  height: 100dvh;
  padding-bottom: 50px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f2f4f2;
  color: var(--map-ink);
}

.map-header {
  position: relative;
  z-index: 20;
  flex: 0 0 auto;
  background: #ffffff;
  border-bottom: 1px solid #e7e9e7;
}

.header-main {
  height: 62px;
  padding: 10px 16px 6px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-main h1 {
  font-size: 20px;
  line-height: 1.15;
  font-weight: 700;
  letter-spacing: 0;
}

.location-status {
  height: 22px;
  margin-top: 2px;
  padding: 0;
  border: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  color: var(--map-muted);
  font-size: 12px;
  cursor: pointer;
}

.shop-total {
  min-width: 68px;
  padding-left: 14px;
  border-left: 1px solid #e6e8e6;
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.shop-total strong {
  color: var(--map-orange);
  font-size: 24px;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.shop-total span {
  color: var(--map-muted);
  font-size: 11px;
}

.category-tabs {
  height: 44px;
  padding: 4px 12px 8px;
  display: flex;
  align-items: center;
  gap: 4px;
  overflow-x: auto;
  scrollbar-width: none;
}

.category-tabs::-webkit-scrollbar {
  display: none;
}

.category-tab {
  position: relative;
  flex: 0 0 auto;
  height: 32px;
  padding: 0 11px;
  border: 0;
  border-radius: 4px;
  background: transparent;
  color: #666c70;
  font-size: 13px;
  white-space: nowrap;
  cursor: pointer;
}

.category-tab.active {
  background: #fff0e9;
  color: #bf4420;
  font-weight: 600;
}

.map-stage {
  position: relative;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.map-canvas {
  width: 100%;
  height: 100%;
  background: #dfe4df;
}

.map-summary,
.map-actions,
.map-state,
.selected-shop,
.tile-warning {
  position: absolute;
  z-index: 500;
}

.map-summary {
  top: 12px;
  left: 12px;
  height: 30px;
  padding: 0 10px;
  border: 1px solid rgb(255 255 255 / 78%);
  border-radius: 5px;
  display: flex;
  align-items: center;
  gap: 7px;
  background: rgb(255 255 255 / 94%);
  box-shadow: 0 4px 12px rgb(27 39 34 / 12%);
  color: #343a3d;
  font-size: 12px;
  font-weight: 600;
}

.summary-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--map-orange);
}

.map-actions {
  top: 88px;
  right: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.map-action {
  width: 36px;
  height: 36px;
  border: 1px solid rgb(214 219 215 / 90%);
  border-radius: 6px;
  display: grid;
  place-items: center;
  background: rgb(255 255 255 / 96%);
  box-shadow: 0 3px 10px rgb(27 39 34 / 14%);
  color: #3e474a;
  font-size: 18px;
  cursor: pointer;
}

.map-action:disabled {
  cursor: wait;
  opacity: 0.7;
}

.map-state {
  top: 50%;
  left: 50%;
  min-width: 150px;
  padding: 14px 18px;
  border: 1px solid rgb(225 229 225 / 90%);
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transform: translate(-50%, -50%);
  background: rgb(255 255 255 / 96%);
  box-shadow: 0 8px 26px rgb(36 45 40 / 15%);
  color: #596064;
  font-size: 13px;
}

.map-state-error {
  flex-direction: column;
  color: #9e3c2b;
}

.map-state-error button {
  padding: 5px 14px;
  border: 1px solid #da694a;
  border-radius: 4px;
  background: #ffffff;
  color: #b94429;
  cursor: pointer;
}

.selected-shop {
  right: 12px;
  bottom: 14px;
  left: 12px;
  width: calc(100% - 24px);
  min-height: 88px;
  padding: 10px;
  border: 1px solid rgb(229 231 228 / 92%);
  border-radius: 7px;
  display: grid;
  grid-template-columns: 68px minmax(0, 1fr) 18px;
  align-items: center;
  gap: 10px;
  background: rgb(255 255 255 / 97%);
  box-shadow: 0 8px 24px rgb(34 43 39 / 20%);
  color: inherit;
  text-align: left;
  cursor: pointer;
}

.selected-shop img {
  width: 68px;
  height: 68px;
  border-radius: 5px;
  object-fit: cover;
  background: #e9ece9;
}

.selected-shop-info {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.selected-shop-heading {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
}

.selected-shop-heading strong {
  min-width: 0;
  overflow: hidden;
  color: #202426;
  font-size: 15px;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.selected-shop-heading > span {
  flex: 0 0 auto;
  color: #bc4b2a;
  font-size: 12px;
  font-weight: 600;
}

.selected-shop-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #646b6e;
  font-size: 11px;
  white-space: nowrap;
}

.selected-shop-meta i {
  width: 2px;
  height: 2px;
  border-radius: 50%;
  background: #a8ada9;
}

.selected-shop-address {
  overflow: hidden;
  color: #858b8d;
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.detail-arrow {
  color: #8a9092;
}

.tile-warning {
  right: 12px;
  bottom: 112px;
  left: 12px;
  padding: 7px 10px;
  border: 1px solid #ead6a0;
  border-radius: 5px;
  background: rgb(255 249 226 / 96%);
  color: #7a6528;
  font-size: 11px;
  text-align: center;
}

:global(.leaflet-top.leaflet-right) {
  top: 46px;
  right: 2px;
}

:global(.leaflet-control-zoom) {
  border: 1px solid #d6dbd7 !important;
  border-radius: 6px !important;
  overflow: hidden;
  box-shadow: 0 3px 10px rgb(27 39 34 / 14%) !important;
}

:global(.leaflet-control-zoom a) {
  width: 34px !important;
  height: 34px !important;
  border-color: #e5e8e5 !important;
  color: #3e474a !important;
  font: 18px/34px Arial, sans-serif !important;
}

:global(.leaflet-control-attribution) {
  margin-bottom: 108px !important;
  font-size: 9px !important;
}

:global(.leaflet-popup-content-wrapper) {
  border-radius: 6px;
  box-shadow: 0 6px 20px rgb(24 35 30 / 18%);
}

:global(.leaflet-popup-content) {
  margin: 0;
}

:global(.shop-popup-button) {
  width: 100%;
  padding: 10px 12px;
  border: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
  background: #ffffff;
  color: #2a3032;
  text-align: left;
  cursor: pointer;
}

:global(.shop-popup-button strong) {
  max-width: 210px;
  overflow: hidden;
  font-size: 13px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:global(.shop-popup-button span) {
  max-width: 210px;
  overflow: hidden;
  color: #747a7c;
  font-size: 11px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 360px) {
  .header-main {
    padding-right: 12px;
    padding-left: 12px;
  }

  .shop-total {
    min-width: 62px;
    padding-left: 10px;
  }

  .selected-shop {
    grid-template-columns: 60px minmax(0, 1fr) 16px;
    gap: 8px;
  }

  .selected-shop img {
    width: 60px;
    height: 60px;
  }
}
</style>
