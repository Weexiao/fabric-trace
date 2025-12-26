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

// getFruitInfo 使用 JSON
export function getFruitInfo(data) {
  return request({
    url: '/getFruitInfo',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data
  })
}

// getFruitList
export function getFruitList(data) {
  return request({
    url: '/getFruitList',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}

// 服务端分页的“获取所有产品信息”
// 约定 data: { page: number, pageSize: number }
export function getAllFruitInfo(data) {
  const payload = data && typeof data === 'object' ? data : {}
  return request({
    url: '/getAllFruitInfo',
    method: 'post',
    headers: {
      'Content-Type': 'application/json'
    },
    data: payload
  })
}

// getFruitHistory
export function getFruitHistory(data) {
  return request({
    url: '/getFruitHistory',
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data
  })
}
