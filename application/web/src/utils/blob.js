// Centralized helpers for creating and revoking object URLs safely
// Provides a tiny registry to allow bulk revoke as a safety net.

const registry = new Set()

export function createObjectURLSafe(file, track = true) {
  try {
    const url = URL.createObjectURL(file)
    if (track && url) registry.add(url)
    return url
  } catch (e) {
    if (process.env.NODE_ENV !== 'production') console.warn('createObjectURL failed', e)
    return ''
  }
}

export function revokeObjectURLSafe(url) {
  if (!url) return
  try {
    URL.revokeObjectURL(url)
  } catch (e) {
    if (process.env.NODE_ENV !== 'production') console.warn('revokeObjectURL failed', e)
  }
  registry.delete(url)
}

export function revokeAllObjectURLs() {
  registry.forEach((u) => {
    try { URL.revokeObjectURL(u) } catch (e) { if (process.env.NODE_ENV !== 'production') console.warn('revokeObjectURL failed (bulk)', e) }
  })
  registry.clear()
}
