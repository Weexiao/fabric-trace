// Normalize trace data rows into a unified shape used by views
// Export both per-row and list-level helpers

export function normalizeRow(item) {
  const o = item && typeof item === 'object' ? { ...item } : {}
  // top-level
  o.traceability_code = o.traceability_code || ''
  // source branches (support legacy casing)
  const farmer = o.farmer_input || o.Farmer_input || {}
  const factory = o.factory_input || o.Factory_input || {}
  const driver = o.driver_input || o.Driver_input || {}
  const shop = o.shop_input || o.Shop_input || {}

  o.farmer_input = {
    fa_fruitName: farmer.fa_fruitName || farmer.Fa_fruitName || '',
    fa_origin: farmer.fa_origin || farmer.Fa_origin || '',
    fa_originCodes: farmer.fa_originCodes || farmer.Fa_originCodes || [],
    fa_plantTime: farmer.fa_plantTime || farmer.Fa_plantTime || '',
    fa_pickingTime: farmer.fa_pickingTime || farmer.Fa_pickingTime || '',
    fa_farmerName: farmer.fa_farmerName || farmer.Fa_farmerName || '',
    fa_img: farmer.fa_img || farmer.Fa_img || '',
    fa_txid: farmer.fa_txid || farmer.Fa_txid || '',
    fa_timestamp: farmer.fa_timestamp || farmer.Fa_timestamp || ''
  }

  o.factory_input = {
    fac_productName: factory.fac_productName || factory.Fac_productName || '',
    fac_productionbatch: factory.fac_productionbatch || factory.Fac_productionbatch || '',
    fac_productionTime: factory.fac_productionTime || factory.Fac_productionTime || '',
    fac_factoryName: factory.fac_factoryName || factory.Fac_factoryName || '',
    fac_contactNumber: factory.fac_contactNumber || factory.Fac_contactNumber || '',
    fac_img: factory.fac_img || factory.Fac_img || '',
    fac_txid: factory.fac_txid || factory.Fac_txid || '',
    fac_timestamp: factory.fac_timestamp || factory.Fac_timestamp || ''
  }

  o.driver_input = {
    dr_name: driver.dr_name || driver.Dr_name || '',
    dr_age: driver.dr_age || driver.Dr_age || '',
    dr_phone: driver.dr_phone || driver.Dr_phone || '',
    dr_carNumber: driver.dr_carNumber || driver.Dr_carNumber || '',
    dr_transport: driver.dr_transport || driver.Dr_transport || '',
    dr_img: driver.dr_img || driver.Dr_img || '',
    dr_txid: driver.dr_txid || driver.Dr_txid || '',
    dr_timestamp: driver.dr_timestamp || driver.Dr_timestamp || ''
  }

  o.shop_input = {
    sh_storeTime: shop.sh_storeTime || shop.Sh_storeTime || '',
    sh_sellTime: shop.sh_sellTime || shop.Sh_sellTime || '',
    sh_shopName: shop.sh_shopName || shop.Sh_shopName || '',
    sh_shopAddress: shop.sh_shopAddress || shop.Sh_shopAddress || '',
    sh_shopAddressCodes: shop.sh_shopAddressCodes || shop.Sh_shopAddressCodes || [],
    sh_shopPhone: shop.sh_shopPhone || shop.Sh_shopPhone || '',
    sh_img: shop.sh_img || shop.Sh_img || '',
    sh_txid: shop.sh_txid || shop.Sh_txid || '',
    sh_timestamp: shop.sh_timestamp || shop.Sh_timestamp || ''
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

