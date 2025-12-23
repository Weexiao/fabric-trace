import Vue from 'vue'
import { clampLen, stripDangerous } from '@/utils/sanitize'

function getInputEl(el) {
  if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') return el
  return el.querySelector('input, textarea')
}

function applySanitize(inputEl, binding) {
  if (!inputEl) return
  const opts = binding && binding.value ? binding.value : {}
  const max = typeof opts.max === 'number' ? opts.max : undefined
  const digitsOnly = !!(binding && binding.modifiers && binding.modifiers.digits)

  const raw = inputEl.value
  let sanitized = stripDangerous(raw)
  if (digitsOnly) sanitized = sanitized.replace(/[^\d]/g, '')
  if (typeof max === 'number') sanitized = clampLen(sanitized, max)

  if (sanitized !== raw) {
    // avoid re-entrancy
    if (inputEl.__sanitizing) return
    inputEl.__sanitizing = true
    inputEl.value = sanitized
    inputEl.dispatchEvent(new Event('input', { bubbles: true }))
    inputEl.__sanitizing = false
  }
}

Vue.directive('sanitize', {
  bind(el, binding) {
    const inputEl = getInputEl(el)
    // set maxlength on the real input for UX
    if (inputEl && binding && binding.value && typeof binding.value.max === 'number') {
      inputEl.setAttribute('maxlength', String(binding.value.max))
    }
    const handler = (e) => applySanitize(getInputEl(el), binding)
    el.__sanitizeHandler = handler
    if (inputEl) inputEl.addEventListener('input', handler)
  },
  inserted(el, binding) {
    // initial sanitize
    applySanitize(getInputEl(el), binding)
  },
  update(el, binding) {
    // re-apply on updates
    applySanitize(getInputEl(el), binding)
  },
  unbind(el) {
    const inputEl = getInputEl(el)
    if (inputEl && el.__sanitizeHandler) inputEl.removeEventListener('input', el.__sanitizeHandler)
    delete el.__sanitizeHandler
  }
})

