import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
import { getToken } from '@/utils/auth' // get token from cookie
import getPageTitle from '@/utils/get-page-title'
import { hasPermission } from '@/utils/checkPermission' // dynamic attribute permission checker

NProgress.configure({ showSpinner: false }) // NProgress Configuration

// 免登录白名单，支持子路径匹配
const whiteList = ['/login', '/trace']

router.beforeEach(async(to, from, next) => {
  // start progress bar
  NProgress.start()

  // set page title
  document.title = getPageTitle(to.meta.title)

  // determine whether the user has logged in
  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      // if logged in, redirect to home page
      next({ path: '/' })
      NProgress.done()
    } else {
      const hasGetUserInfo = store.getters.name
      if (hasGetUserInfo) {
        // 检查角色权限
        if (to.meta.roles && to.meta.roles.length > 0) {
          const userRole = store.getters.userType
          if (!to.meta.roles.includes(userRole)) {
            Message.error('你没有权限访问此页面（角色不符）')
            NProgress.done()
            return
          }
        }

        // 检查动态属性权限（Admin 用户豁免）
        if (to.meta.requiredAttributes) {
          const userRole = store.getters.userType
          // Admin 用户跳过属性检查
          if (userRole !== 'Admin') {
            const user = {
              userType: store.getters.userType,
              dynamicAttributes: store.getters.dynamicAttributes || {}
            }
            if (!hasPermission(user, { requiredAttributes: to.meta.requiredAttributes })) {
              Message.error('你没有权限访问此资源（属性不匹配）')
              NProgress.done()
              return
            }
          }
        }
        next()
      } else {
        try {
          // get user info
          await store.dispatch('user/getInfo')

          // 检查角色权限
          if (to.meta.roles && to.meta.roles.length > 0) {
            const userRole = store.getters.userType
            if (!to.meta.roles.includes(userRole)) {
              Message.error('你没有权限访问此页面（角色不符）')
              NProgress.done()
              return
            }
          }

          // 检查动态属性权限（Admin 用户豁免）
          if (to.meta.requiredAttributes) {
            const userRole = store.getters.userType
            // Admin 用户跳过属性检查
            if (userRole !== 'Admin') {
              const user = {
                userType: store.getters.userType,
                dynamicAttributes: store.getters.dynamicAttributes || {}
              }
              if (!hasPermission(user, { requiredAttributes: to.meta.requiredAttributes })) {
                Message.error('你没有权限访问此资源（属性不匹配）')
                NProgress.done()
                return
              }
            }
          }
          next()
        } catch (error) {
          // remove token and go to login page to re-login
          await store.dispatch('user/resetToken')
          Message.error(error || 'Has Error')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    // no token
    // 检查是否在白名单内（支持子路径匹配）
    if (whiteList.some(path => to.path.startsWith(path))) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})
