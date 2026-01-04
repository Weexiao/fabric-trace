import axios from 'axios'
import { Message } from 'element-ui'
import store from '@/store'
import { getToken } from '@/utils/auth'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 10000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    // attach auth token if exists
    if (store.getters.token && !config.headers['Authorization']) {
      config.headers['Authorization'] = getToken()
    }
    // default JSON headers for POST unless explicitly set
    if (config.method === 'post' && !config.headers['Content-Type']) {
      config.headers['Content-Type'] = 'application/json'
    }
    return config
  },
  error => {
    // request error
    console.log(error)
    return Promise.reject(error)
  }
)

// helper: safe JSON parse
function parseJsonSafe(input) {
  try {
    return JSON.parse(input)
  } catch (e) {
    return input
  }
}

// response interceptor
service.interceptors.response.use(
  response => {
    // Some backends return raw string; normalize first
    let res = response && response.data !== undefined ? response.data : response
    if (typeof res === 'string') {
      res = parseJsonSafe(res)
    }
    // If payload has inner data as string, parse once
    if (res && typeof res.data === 'string') {
      res.data = parseJsonSafe(res.data)
    }

    // Normalize error code convention: accept code===200 or success===true
    const code = res && (res.code !== undefined ? res.code : (res.status !== undefined ? res.status : undefined))
    const ok = (code === 200) || (res && res.success === true)

    if (!ok) {
      const message = (res && (res.message || res.msg || res.error)) || '请求异常'
      Message({ message, type: 'error', duration: 5 * 1000 })
      const err = new Error(message)
      err.response = response
      err.payload = res
      err.code = code
      return Promise.reject(err)
    }
    // 成功时直接返回后端payload，避免再次包装导致字段缺失
    return res
  },
  error => {
    // Network or unexpected errors
    const msg = (error && (error.message || (error.response && error.response.statusText))) || '网络异常或服务器错误'
    console.log('err', error)
    Message({ message: msg, type: 'error', duration: 5 * 1000 })
    // Propagate a normalized error
    const err = new Error(msg)
    err.original = error
    return Promise.reject(err)
  }
)

export default service
