// Minimal, dependency-free sanitization helpers
// 1) escapeHTML: escape < > & " ' to avoid injecting HTML
// 2) stripDangerous: remove control chars and zero-widths
// 3) clampLen: trim string to a max length

const HTML_ESCAPE_REG = /[&<>"]|\u00A0|\u2028|\u2029|\u00AD|\u200B|\u200C|\u200D/g
const HTML_REPLACER = (ch) => {
  switch (ch) {
    case '&': return '&amp;'
    case '<': return '&lt;'
    case '>': return '&gt;'
    case '"': return '&quot;'
    // some whitespace/control chars that can break contexts
    case '\u00A0': return ' '
    case '\u2028':
    case '\u2029': return '\n'
    case '\u00AD': // soft hyphen
    case '\u200B': // zero width space
    case '\u200C': // zero width non-joiner
    case '\u200D': // zero width joiner
      return ''
    default: return ch
  }
}

export function escapeHTML(input) {
  if (input == null) return ''
  return String(input).replace(HTML_ESCAPE_REG, HTML_REPLACER)
}

export function stripDangerous(input) {
  if (input == null) return ''
  // remove control characters except CR/LF/TAB
  return String(input).replace(/[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]/g, '')
}

export function clampLen(input, max) {
  if (input == null) return ''
  const s = String(input)
  return s.length > max ? s.slice(0, max) : s
}

export function sanitize(input, max) {
  return clampLen(stripDangerous(input), max)
}

// remove default export to avoid unused default lint warning
// export default { escapeHTML, stripDangerous, clampLen, sanitize }
