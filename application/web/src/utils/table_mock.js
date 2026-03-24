const TABLE_MOCK_KEY = 'uplink_table_mock_records'

function readFileAsText(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result || ''))
    reader.onerror = () => reject(new Error('read file failed'))
    reader.readAsText(file)
  })
}

function detectDelimiter(lines) {
  const sample = lines.slice(0, 3)
  const score = (d) => sample.reduce((acc, line) => acc + (line.split(d).length - 1), 0)
  const candidates = [',', '\t', ';']
  let best = ','
  let bestScore = -1
  candidates.forEach((d) => {
    const s = score(d)
    if (s > bestScore) {
      best = d
      bestScore = s
    }
  })
  return best
}

export function parseCsv(text, delimiter) {
  const rows = []
  let row = []
  let field = ''
  let i = 0
  let inQuotes = false
  const src = String(text || '').replace(/^\uFEFF/, '')

  while (i < src.length) {
    const ch = src[i]
    if (inQuotes) {
      if (ch === '"') {
        if (src[i + 1] === '"') {
          field += '"'
          i += 1
        } else {
          inQuotes = false
        }
      } else {
        field += ch
      }
    } else if (ch === '"') {
      inQuotes = true
    } else if (ch === delimiter) {
      row.push(field)
      field = ''
    } else if (ch === '\n') {
      row.push(field)
      rows.push(row)
      row = []
      field = ''
    } else if (ch === '\r') {
      // ignore CR, LF branch handles end of line
    } else {
      field += ch
    }
    i += 1
  }

  row.push(field)
  const hasAny = row.some(cell => String(cell || '').trim() !== '')
  if (hasAny) rows.push(row)
  return rows
}

export async function parseTableFile(file) {
  const text = await readFileAsText(file)
  const lines = text.split(/\r?\n/).filter(Boolean)
  const delimiter = detectDelimiter(lines)
  const rows = parseCsv(text, delimiter)
  if (!rows.length) {
    return { headers: [], rows: [], rowCount: 0, colCount: 0 }
  }
  const header = rows[0].map(v => String(v || '').trim())
  const body = rows.slice(1)
  return {
    headers: header,
    rows: body,
    rowCount: body.length,
    colCount: header.length,
    delimiter
  }
}

export function loadTableMockMap() {
  try {
    const raw = localStorage.getItem(TABLE_MOCK_KEY)
    const parsed = raw ? JSON.parse(raw) : {}
    return parsed && typeof parsed === 'object' ? parsed : {}
  } catch (e) {
    return {}
  }
}

export function saveTableMockRecord(traceabilityCode, record) {
  if (!traceabilityCode || !record) return
  const code = String(traceabilityCode).trim()
  if (!code) return
  const map = loadTableMockMap()
  const list = Array.isArray(map[code]) ? map[code] : []
  list.unshift(record)
  map[code] = list.slice(0, 20)
  try {
    localStorage.setItem(TABLE_MOCK_KEY, JSON.stringify(map))
  } catch (e) {
    // ignore
  }
}

export function getTableMockRecords(traceabilityCode) {
  const map = loadTableMockMap()
  const code = String(traceabilityCode || '').trim()
  if (!code) return []
  return Array.isArray(map[code]) ? map[code] : []
}

