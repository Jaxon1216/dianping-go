import { ref, onMounted, onUnmounted } from 'vue'

/**
 * 滚动加载组合式函数
 * @param {Function} loadMore - 加载更多数据的函数
 * @param {Object} options - 配置选项
 * @returns {Object} 滚动相关状态和方法
 */
export function useScrollLoad(loadMore, options = {}) {
  const { threshold = 50, immediate = true } = options
  
  const loading = ref(false)
  const finished = ref(false)
  const error = ref(false)
  const scrollContainer = ref(null)
  
  const handleScroll = async (e) => {
    if (loading.value || finished.value || error.value) return
    
    const target = e.target
    const scrollTop = target.scrollTop
    const offsetHeight = target.offsetHeight
    const scrollHeight = target.scrollHeight
    
    if (scrollTop + offsetHeight + threshold >= scrollHeight) {
      loading.value = true
      try {
        const hasMore = await loadMore()
        if (!hasMore) {
          finished.value = true
        }
      } catch (err) {
        error.value = true
      } finally {
        loading.value = false
      }
    }
  }
  
  const reset = () => {
    finished.value = false
    error.value = false
    loading.value = false
  }
  
  const retry = () => {
    error.value = false
    handleScroll({ target: scrollContainer.value })
  }
  
  onMounted(() => {
    if (scrollContainer.value && immediate) {
      // 如果内容不足一屏，自动加载
      const target = scrollContainer.value
      if (target.scrollHeight <= target.offsetHeight) {
        handleScroll({ target })
      }
    }
  })
  
  return {
    loading,
    finished,
    error,
    scrollContainer,
    handleScroll,
    reset,
    retry
  }
}
