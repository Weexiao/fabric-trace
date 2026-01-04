// Normalize industrial product data rows into a unified shape used by views
// Export both per-row and list-level helpers

export function normalizeRow(item) {
  const o = item && typeof item === 'object' ? { ...item } : {}
  o.traceabilityCode = o.traceabilityCode || ''
  // branches (no legacy fallback per requirement)
  const raw = o.rawSupplierInput || {}
  const manuf = o.manufacturerInput || {}
  const carrier = o.carrierInput || {}
  const dealer = o.dealerInput || {}

  o.rawSupplierInput = {
    productName: raw.productName || '',
    rawOrigin: raw.rawOrigin || '',
    arrivalTime: raw.arrivalTime || '',
    productionTime: raw.productionTime || '',
    supplierName: raw.supplierName || '',
    img: raw.img || '',
    txid: raw.txid || '',
    timestamp: raw.timestamp || ''
  }

  o.manufacturerInput = {
    productName: manuf.productName || '',
    productionBatch: manuf.productionBatch || '',
    factoryTime: manuf.factoryTime || '',
    factoryNameAddress: manuf.factoryNameAddress || '',
    contactPhone: manuf.contactPhone || '',
    img: manuf.img || '',
    txid: manuf.txid || '',
    timestamp: manuf.timestamp || ''
  }

  o.carrierInput = {
    name: carrier.name || '',
    age: carrier.age || '',
    phone: carrier.phone || '',
    plateNumber: carrier.plateNumber || '',
    transportRecord: carrier.transportRecord || '',
    img: carrier.img || '',
    txid: carrier.txid || '',
    timestamp: carrier.timestamp || ''
  }

  o.dealerInput = {
    storeTime: dealer.storeTime || '',
    sellTime: dealer.sellTime || '',
    dealerName: dealer.dealerName || '',
    dealerLocation: dealer.dealerLocation || '',
    dealerPhone: dealer.dealerPhone || '',
    img: dealer.img || '',
    txid: dealer.txid || '',
    timestamp: dealer.timestamp || ''
  }

  return o
}

export function normalizeResults(list) {
  return (Array.isArray(list) ? list : []).map(normalizeRow)
}

// Optional: normalize single object response to list for table usage
export function toList(item) {
  if (!item || typeof item !== 'object') return []
  return [normalizeRow(item)]
}
