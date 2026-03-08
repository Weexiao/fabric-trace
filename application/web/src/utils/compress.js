/**
 * compress.js - 前端数据压缩工具
 *
 * 使用 pako (Gzip) 对 JSON 业务数据进行压缩 + Base64 编码，
 * 适配后端 /uplink/compressed 接口的 { compressedPayload } 传输协议。
 *
 * 技术说明：
 * pako.gzip() 生成符合 RFC 1952 的 Gzip 字节流，与后端 Go compress/gzip 完全兼容。
 * Base64 编码确保二进制数据在 JSON 中安全传输，无需 multipart/form-data。
 */
import pako from 'pako'

/**
 * 将对象压缩为 Base64 字符串
 * @param {Object|string} data - 待压缩的数据（对象将自动 JSON.stringify）
 * @returns {{ compressedB64: string, originalSize: number, compressedSize: number }}
 */
export function compressPayload(data) {
  const jsonStr = typeof data === 'string' ? data : JSON.stringify(data)
  const originalBytes = new TextEncoder().encode(jsonStr)
  const compressed = pako.gzip(originalBytes)

  // 将 Uint8Array 转为 Base64
  const compressedB64 = uint8ArrayToBase64(compressed)

  return {
    compressedB64,
    originalSize: originalBytes.length,
    compressedSize: compressed.length
  }
}

/**
 * 从 Base64 字符串解压还原为对象
 * @param {string} b64 - Base64 编码的 Gzip 数据
 * @returns {any} 解压后 JSON.parse 的对象
 */
export function decompressPayload(b64) {
  const compressed = base64ToUint8Array(b64)
  const decompressed = pako.ungzip(compressed)
  const jsonStr = new TextDecoder().decode(decompressed)
  return JSON.parse(jsonStr)
}

/**
 * Uint8Array 转 Base64（兼容大数据量，避免 btoa + String.fromCharCode 的栈溢出）
 */
function uint8ArrayToBase64(uint8Array) {
  let binary = ''
  const chunkSize = 8192
  for (let i = 0; i < uint8Array.length; i += chunkSize) {
    const chunk = uint8Array.subarray(i, i + chunkSize)
    binary += String.fromCharCode.apply(null, chunk)
  }
  return btoa(binary)
}

/**
 * Base64 转 Uint8Array
 */
function base64ToUint8Array(b64) {
  const binary = atob(b64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i)
  }
  return bytes
}

export default {
  compressPayload,
  decompressPayload
}

