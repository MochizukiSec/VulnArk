// 这是一个vue-i18n的替代模块，返回一个空函数以避免导入错误
export const useI18n = () => {
  return {
    t: (key) => {
      // 提取最后一个部分作为中文文本的最佳猜测
      const parts = key.split('.');
      return parts[parts.length - 1] || key;
    },
    d: (value) => value,
    n: (value) => value
  }
}

export const createI18n = () => {
  return {
    global: {
      t: (key) => {
        const parts = key.split('.');
        return parts[parts.length - 1] || key;
      }
    }
  }
} 