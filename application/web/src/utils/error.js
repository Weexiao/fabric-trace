// Unified error handling helpers
export function parseJsonSafe(input, fallback = null) {
  try {
    if (typeof input === 'string') return JSON.parse(input)
    return input
  } catch (e) {
    return fallback
  }
}

export function handleApiError(err, defaultMessage = '请求失败，请稍后重试', vm) {
  const msg = (err && err.message) ? err.message : defaultMessage
  if (vm) {
    // expose to UI error card
    vm.errorMessage = msg
  }
  if (vm && vm.$message) {
    vm.$message.error(msg)
  }
  return msg
}

export function finalizeLoading(vm) {
  if (vm) vm.loading = false
}

/**
 * Wrap a promise with standardized error handling and finally loading cleanup.
 * Supports passing a promise factory for retries.
 */
export function apiWrap(vm, promiseOrFactory, onSuccess, failMessage) {
  vm.loading = true
  vm.errorMessage = null
  // Remember last request info for retry
  vm._lastFactory = typeof promiseOrFactory === 'function' ? promiseOrFactory : () => promiseOrFactory
  vm._lastOnSuccess = onSuccess
  vm._lastFailMessage = failMessage

  const p = typeof promiseOrFactory === 'function' ? promiseOrFactory() : promiseOrFactory
  return p
    .then(res => {
      // clear error on success
      vm.errorMessage = null
      return onSuccess(res)
    })
    .catch(err => handleApiError(err, failMessage, vm))
    .finally(() => finalizeLoading(vm))
}

export function retryLast(vm) {
  if (!vm || typeof vm._lastFactory !== 'function') return
  return apiWrap(vm, vm._lastFactory, vm._lastOnSuccess, vm._lastFailMessage)
}
