import request from '@/utils/request'

export function uploadFile(formData) {
  return request({
    url: '/file/upload',
    method: 'post',
    headers: { 'Content-Type': 'multipart/form-data' },
    data: formData
  })
}

export function listManifests(traceabilityCode) {
  const data = new FormData()
  data.append('traceabilityCode', traceabilityCode)
  return request({
    url: '/file/list',
    method: 'post',
    data
  })
}

export function downloadFile(fileID) {
  return request({
    url: `/file/download/${fileID}`,
    method: 'get',
    responseType: 'blob'
  })
}

