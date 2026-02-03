import { login, logout, getInfo, register } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { resetRouter } from '@/router'

const getDefaultState = () => {
  return {
    token: getToken(),
    name: '',
    avatar: '',
    userType: '',
    dynamicAttributes: {}
  }
}

const state = getDefaultState()

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
    // 清除 localStorage 中的 userType
    localStorage.removeItem('userType')
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_USERTYPE: (state, userType) => {
    state.userType = userType
    // 保存到 localStorage，供路由重定向使用
    localStorage.setItem('userType', userType)
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_DYNAMIC_ATTRIBUTES: (state, dynamicAttributes) => {
    state.dynamicAttributes = dynamicAttributes
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    const formData = new FormData()
    formData.append('username', username.trim())
    formData.append('password', password)
    return new Promise((resolve, reject) => {
      login(formData).then(response => {
        commit('SET_TOKEN', response.jwt)
        setToken(response.jwt)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user register
  register({ commit }, userInfo) {
    const { username, password, userType } = userInfo
    const formData = new FormData()
    formData.append('username', username.trim())
    formData.append('password', password)
    formData.append('userType', userType)
    return new Promise((resolve, reject) => {
      register(formData).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo(state.token).then(response => {
        const { username } = response
        const { userType } = response
        const { dynamicAttributes } = response

        if (!username) {
          return reject('Verification failed, please Login again.')
        }

        commit('SET_NAME', username)
        commit('SET_USERTYPE', userType)
        // Parse JSON string to object if needed
        const attrs = typeof dynamicAttributes === 'string' ? JSON.parse(dynamicAttributes || '{}') : (dynamicAttributes || {})
        commit('SET_DYNAMIC_ATTRIBUTES', attrs)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        removeToken() // must remove  token  first
        resetRouter()
        commit('RESET_STATE')
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      removeToken() // must remove  token  first
      commit('RESET_STATE')
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}

