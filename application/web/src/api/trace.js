import request from '@/utils/request'

export function uplink(data) {
  return request({
    url: '/uplink',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// getIndustrialProductInfo 使用 JSON
export function getIndustrialProductInfo(data) {
  return request({
    url: '/getIndustrialProductInfo',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// getIndustrialProductList
export function getIndustrialProductList(data) {
  return request({
    url: '/getIndustrialProductList',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 服务端分页的“获取所有产品信息”
// 约定 data: { page: number, pageSize: number }
export function getAllIndustrialProductInfo(data) {
  const payload = data && typeof data === 'object' ? data : {}
  return request({
    url: '/getAllIndustrialProductInfo',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data: payload
  })
}

// getIndustrialProductHistory
export function getIndustrialProductHistory(data) {
  return request({
    url: '/getIndustrialProductHistory',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
