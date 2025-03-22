/**
 * 时间格式化工具
 */

/**
 * 将时间字符串或时间对象格式化为标准格式
 * @param {Date|string} time 时间对象或ISO时间字符串
 * @param {string} format 格式模板，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns {string} 格式化后的时间字符串
 */
export function formatTime(time, format = 'YYYY-MM-DD HH:mm:ss') {
  if (!time) return '';
  
  // 转换为Date对象
  let date;
  if (typeof time === 'string') {
    // 尝试解析字符串
    date = new Date(time);
  } else if (time instanceof Date) {
    date = time;
  } else {
    return '';
  }
  
  // 检查日期有效性
  if (isNaN(date.getTime())) {
    return '';
  }
  
  // 格式化时间
  const year = date.getFullYear();
  const month = padZero(date.getMonth() + 1);
  const day = padZero(date.getDate());
  const hours = padZero(date.getHours());
  const minutes = padZero(date.getMinutes());
  const seconds = padZero(date.getSeconds());
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds);
}

/**
 * 格式化相对时间（多久之前）
 * @param {Date|string} time 时间对象或ISO时间字符串
 * @returns {string} 相对时间
 */
export function formatRelativeTime(time) {
  if (!time) return '';
  
  const date = typeof time === 'string' ? new Date(time) : time;
  if (isNaN(date.getTime())) return '';
  
  const now = new Date();
  const diffSeconds = Math.floor((now - date) / 1000);
  
  if (diffSeconds < 60) {
    return '刚刚';
  } else if (diffSeconds < 3600) {
    return `${Math.floor(diffSeconds / 60)}分钟前`;
  } else if (diffSeconds < 86400) {
    return `${Math.floor(diffSeconds / 3600)}小时前`;
  } else if (diffSeconds < 604800) {
    return `${Math.floor(diffSeconds / 86400)}天前`;
  } else {
    return formatTime(date, 'YYYY-MM-DD');
  }
}

/**
 * 数字补零
 * @param {number} num 数字
 * @returns {string} 补零后的字符串
 */
function padZero(num) {
  return num < 10 ? `0${num}` : `${num}`;
} 