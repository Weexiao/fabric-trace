import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  {
    path: '/',
    component: Layout,
    redirect: '/uplink',
    children: [{
      path: 'uplink',
      name: 'Uplink',
      component: () => import('@/views/uplink/index'),
      meta: {
        title: '溯源信息录入',
        icon: 'el-icon-edit-outline',
        roles: ['原料供应商', '制造商', '物流承运商', '经销商'] // 排除 Admin
      }
    }]
  },
  {
    path: '/trace',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Trace',
        component: () => import('@/views/trace/index.vue'),
        meta: { title: '溯源查询', icon: 'el-icon-search' }
      },
      {
        path: ':traceability_code',
        name: 'TraceDetail',
        component: () => import('@/views/trace/index.vue'),
        meta: { title: '溯源查询详情' },
        hidden: true
      }
    ]
  },
  // 动态属性权限控制示例页面
  // {
  //   path: '/attribute-demo',
  //   component: Layout,
  //   meta: { title: '权限演示', icon: 'el-icon-key' },
  //   children: [
  //     {
  //       path: '',
  //       name: 'AttributeDemo',
  //       component: () => import('@/views/attribute-demo/index.vue'),
  //       meta: { title: '权限演示中心' }
  //     },
  //     {
  //       path: 'orders-region',
  //       name: 'OrdersRegion',
  //       component: () => import('@/views/orders-region/index.vue'),
  //       meta: {
  //         title: '订单管理（地域限制）',
  //         requiredAttributes: { region: 'Sichuan' }
  //       }
  //     },
  //     {
  //       path: 'financial-reports',
  //       name: 'FinancialReports',
  //       component: () => import('@/views/financial-reports/index.vue'),
  //       meta: {
  //         title: '财务报表（数据级别限制）',
  //         requiredAttributes: { data_level: 'Internal' }
  //       }
  //     }
  //   ]
  // },
  {
    path: '/user-management',
    component: Layout,
    children: [{
      path: '',
      name: 'UserManagement',
      component: () => import('@/views/user-management/index.vue'),
      meta: {
        title: '用户管理',
        icon: 'el-icon-user',
        roles: ['Admin'] // 只有 Admin 可以访问
      }
    }]
  },
  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: 'http://192.168.1.12:8080',
        meta: { title: '区块链浏览器', icon: 'el-icon-discover' }
      }
    ]
  },
  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
