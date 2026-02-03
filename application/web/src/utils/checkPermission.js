/**
 * 权限检查工具 - 支持RBAC + 动态属性权限
 */

/**
 * 检查用户的动态属性是否与资源属性匹配
 * @param {Object} userAttributes - 用户的动态属性 e.g., {region: "Sichuan", data_level: "Internal"}
 * @param {Object} resourceAttributes - 资源的属性限制 e.g., {region: "Sichuan"}
 * @returns {boolean} - 是否有权限访问
 */
export function checkAttributeMatch(userAttributes, resourceAttributes) {
  if (!resourceAttributes || Object.keys(resourceAttributes).length === 0) {
    // 如果资源没有属性限制，则允许访问
    return true
  }

  if (!userAttributes || Object.keys(userAttributes).length === 0) {
    // 如果用户没有属性但资源有限制，拒绝访问
    return false
  }

  // 检查用户属性是否包含资源所需的所有属性
  for (const key in resourceAttributes) {
    if (userAttributes[key] !== resourceAttributes[key]) {
      return false
    }
  }

  return true
}

/**
 * 检查用户是否有权限访问指定的资源
 * @param {Object} user - 用户对象 {userType, dynamicAttributes}
 * @param {Object} resource - 资源对象 {requiredRole, requiredAttributes}
 * @returns {boolean} - 是否有权限
 */
export function hasPermission(user, resource) {
  // 先检查角色权限
  if (resource.requiredRole && user.userType !== resource.requiredRole) {
    return false
  }

  // 再检查属性权限
  if (resource.requiredAttributes) {
    return checkAttributeMatch(user.dynamicAttributes, resource.requiredAttributes)
  }

  return true
}

/**
 * 在路由中标记需要的属性权限
 * 使用示例：
 * {
 *   path: '/orders',
 *   component: () => import('@/views/orders'),
 *   meta: {
 *     requiresAuth: true,
 *     requiredAttributes: { region: 'Sichuan' }  // 只有region为Sichuan的用户可以访问
 *   }
 * }
 */
export default {
  checkAttributeMatch,
  hasPermission
}

