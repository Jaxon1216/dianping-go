/**
 * 格式化价格
 * @param {string|number} val - 价格值
 * @returns {string|null} 格式化后的价格
 */
export const formatPrice = (val) => {
  if (typeof val === 'string') {
    if (isNaN(val)) return null
    
    const index = val.lastIndexOf('.')
    let p = ''
    
    if (index < 0) {
      p = val + '00'
    } else if (index === val.length - 2) {
      p = val.replace('.', '') + '0'
    } else {
      p = val.replace('.', '')
    }
    return (parseInt(p) / 100).toFixed(2)
  }
  
  if (typeof val === 'number') {
    if (!val) return '0.00'
    return (val / 100).toFixed(2)
  }
  
  return null
}

/**
 * 格式化距离
 * @param {number} distance - 距离（米）
 * @returns {string} 格式化后的距离
 */
export const formatDistance = (distance) => {
  if (!distance && distance !== 0) return ''
  if (distance < 1000) {
    return distance.toFixed(1) + 'm'
  }
  return (distance / 1000).toFixed(1) + 'km'
}

/**
 * 格式化时间
 * @param {string|Date} time - 时间
 * @returns {string} 格式化后的时间
 */
export const formatTime = (time) => {
  const date = new Date(time)
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`
}

/**
 * 格式化日期时间（用于秒杀）
 * @param {string|Date} beginTime - 开始时间
 * @param {string|Date} endTime - 结束时间
 * @returns {string} 格式化后的时间范围
 */
export const formatSeckillTime = (beginTime, endTime) => {
  const b = new Date(beginTime)
  const e = new Date(endTime)
  
  const formatMinutes = (m) => m < 10 ? '0' + m : m
  
  return `${b.getMonth() + 1}月${b.getDate()}日 ${b.getHours()}:${formatMinutes(b.getMinutes())} ~ ${e.getHours()}:${formatMinutes(e.getMinutes())}`
}

/**
 * 检查秒杀是否未开始
 * @param {string|Date} beginTime - 开始时间
 * @returns {boolean}
 */
export const isSeckillNotBegin = (beginTime) => {
  return new Date(beginTime).getTime() > Date.now()
}

/**
 * 检查秒杀是否已结束
 * @param {string|Date} endTime - 结束时间
 * @returns {boolean}
 */
export const isSeckillEnd = (endTime) => {
  return new Date(endTime).getTime() < Date.now()
}
