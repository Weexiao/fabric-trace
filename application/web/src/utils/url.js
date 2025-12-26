// 统一的 URL 构建与编码工具
// 使用方法：buildImgUrl(baseApi, fileName)

export function buildImgUrl(baseApi, file) {
  if (!file) return ''
  try {
    const encoded = encodeURIComponent(String(file))
    // 确保 baseApi 末尾有斜杠或拼接正确
    const prefix = String(baseApi || '')
    return `${prefix}getImg/${encoded}`
  } catch (e) {
    return ''
  }
}

export default { buildImgUrl }

