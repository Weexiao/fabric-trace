import request from '@/utils/request'
import axios from 'axios'
import { getToken } from '@/utils/auth'

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
  // Blob download is not JSON and doesn't contain {code:200}; bypass the JSON interceptor.
  return axios({
    url: `${process.env.VUE_APP_BASE_API}/file/download/${fileID}`,
    method: 'get',
    responseType: 'blob',
    headers: {
      Authorization: getToken()
    },
    timeout: 60000
  })
}
