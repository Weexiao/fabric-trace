<template>
  <div class="trace-container">
    <div class="trace-toolbar">
      <el-input
        v-model="traceCode"
        :placeholder="$t ? $t('form.inputTraceCode') : '请输入溯源码查询'"
        :maxlength="LENGTHS.traceCode"
        clearable
        class="trace-input"
        @input="onInputChange"
        @keyup.enter.native="fruitInfo"
      />
      <el-button type="primary" plain :disabled="loading" :loading="loading" @click="fruitInfo"> {{ $t('actions.viewTrace') || '查询' }} </el-button>
      <el-button type="success" plain :disabled="loading" :loading="loading" @click="allFruitInfo"> {{ $t('trace.getAll') || '获取所有产品信息' }} </el-button>
      <el-button type="warning" plain :disabled="loading" @click="resetQuery"> {{ $t('common.reset') || '重 置' }} </el-button>
      <el-select
        v-if="recentQueries.length"
        v-model="recentSelected"
        :placeholder="$t('trace.recentQueries') || '最近查询'"
        clearable
        class="trace-select"
        @change="onSelectRecent"
      >
        <el-option v-for="item in recentQueries" :key="item" :label="item" :value="item" />
      </el-select>
      <el-button v-if="recentQueries.length" type="danger" plain @click="clearRecent"> {{ $t('trace.clearHistory') || '清空历史' }} </el-button>
    </div>
    <div>
      <div v-if="loading" class="loading-placeholder" aria-live="polite">
        <i class="el-icon-loading" /> {{ $t('common.uploading') || '正在加载...' }}
      </div>
      <el-card v-else-if="errorMessage" class="box-card" shadow="never" body-style="padding: 16px;" :aria-label="$t('trace.requestError') || '请求出错'">
        <div slot="header" class="clearfix">
          <span>{{ $t('trace.requestError') || '请求出错' }}</span>
        </div>
        <div>
          <p style="margin:0 0 8px;">{{ errorMessage }}</p>
          <el-button type="primary" plain :disabled="loading" @click="retryLast">{{ $t('actions.continue') || '重试' }}</el-button>
        </div>
      </el-card>
      <el-card v-else-if="!tracedata || tracedata.length === 0" class="box-card" shadow="never" body-style="padding: 16px; text-align:center; color:#606266;" :aria-label="$t('trace.empty') || '暂无数据'">
        {{ $t('trace.empty') || '暂无数据' }}
      </el-card>
      <el-table v-else :data="filteredSortedData" style="width: 100%" :row-key="rowKey" @expand-change="onExpandChange" @sort-change="onSortChange">
        <el-table-column type="expand">
          <template v-slot="props">
            <div v-if="isExpanded(props.row)">
              <el-collapse accordion>
                <el-collapse-item name="farmer">
                  <template v-slot:title>
                    <i class="el-icon-apple" aria-hidden="true" />
                    <span class="section-title">{{ $t('trace.sections.farmer') }}</span>
                    <el-tag size="mini" type="success" :aria-label="$t('trace.tag.rawMaterial')">{{ $t('trace.tag.rawMaterial') }}</el-tag>
                  </template>
                  <el-form label-position="left" inline class="demo-table-expand">
                    <el-form-item :label="($t('form.farmer.fruitName') || '原料名称') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'rawSupplierInput.productName', ''), { max: 100 }) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.farmer.origin') || '原料产地') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'rawSupplierInput.rawOrigin', ''), { max: 100 }) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.farmer.plantTime') || '原料生产时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'rawSupplierInput.arrivalTime', '')) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.farmer.pickTime') || '原料到货时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'rawSupplierInput.productionTime', '')) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.farmer.supplier') || '原料供应商名称') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'rawSupplierInput.supplierName', ''), { max: 100 }) }}</span>
                    </el-form-item>
                    <el-form-item v-if="safeGet(props.row, 'rawSupplierInput.img', '')" :label="($t('trace.relatedImage') || '相关图片（点击下载）') + ':'" class="image-item">
                      <a :href="imgHref(props.row, 'rawSupplierInput.img')" target="_blank" rel="noopener noreferrer" :download="safeGet(props.row, 'rawSupplierInput.productName', 'image')">
                        <el-image
                          style="width: 100px; height: 100px;"
                          :src="imgSrc(props.row, 'rawSupplierInput.img')"
                          :preview-src-list="imgPreviewList(props.row, 'rawSupplierInput.img')"
                          :alt="`${safeGet(props.row, 'rawSupplierInput.productName', $t('trace.alt.farmerImage') || '原料图片')}`"
                          :title="`${safeGet(props.row, 'rawSupplierInput.productName', $t('trace.alt.farmerImage') || '原料图片')}`"
                          fit="cover"
                          lazy
                        >
                          <div slot="placeholder" class="image-placeholder">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.loading') || '加载中' }}
                          </div>
                          <div slot="error" class="image-error">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.error') || '加载失败' }}
                          </div>
                        </el-image>
                      </a>
                      <el-button type="text" icon="el-icon-download" :href="imgHref(props.row, 'rawSupplierInput.img')" :aria-label="$t('actions.download')" @click.prevent="openDownload(imgHref(props.row, 'rawSupplierInput.img'), safeGet(props.row, 'rawSupplierInput.productName', 'image'))">{{ $t('actions.download') || '下载' }}</el-button>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txid') || '区块链交易ID') + ':'">
                      <span>{{ safeGet(props.row, 'rawSupplierInput.txid', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txTime') || '区块链交易时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'rawSupplierInput.timestamp', '')) }}</span>
                    </el-form-item>
                  </el-form>
                </el-collapse-item>
                <el-collapse-item name="factory">
                  <template v-slot:title>
                    <i class="el-icon-office-building" aria-hidden="true" />
                    <span class="section-title">{{ $t('trace.sections.factory') }}</span>
                    <el-tag size="mini" type="info" :aria-label="$t('trace.tag.manufacturer')">{{ $t('trace.tag.manufacturer') }}</el-tag>
                  </template>
                  <el-form label-position="left" inline class="demo-table-expand">
                    <el-form-item :label="($t('form.factory.productName') || '产品名称') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'manufacturerInput.productName', ''), { max: 100 }) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.factory.batch') || '生产批次') + ':'">
                      <span>{{ safeGet(props.row, 'manufacturerInput.productionBatch', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.factory.prodTime') || '生产时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'manufacturerInput.factoryTime', '')) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.factory.factoryName') || '制造商名称') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'manufacturerInput.factoryNameAddress', ''), { max: 100 }) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.factory.phone') || '制造商电话') + ':'">
                      <span>{{ safeGet(props.row, 'manufacturerInput.contactPhone', '') }}</span>
                    </el-form-item>
                    <el-form-item v-if="safeGet(props.row, 'manufacturerInput.img', '')" :label="($t('trace.relatedImage') || '相关图片（点击下载）') + ':'" class="image-item">
                      <a :href="imgHref(props.row, 'manufacturerInput.img')" target="_blank" rel="noopener noreferrer" :download="safeGet(props.row, 'manufacturerInput.productName', 'image')">
                        <el-image
                          style="width: 100px; height: 100px;"
                          :src="imgSrc(props.row, 'manufacturerInput.img')"
                          :preview-src-list="imgPreviewList(props.row, 'manufacturerInput.img')"
                          :alt="`${safeGet(props.row, 'manufacturerInput.productName', $t('trace.alt.factoryImage') || '制造商图片')}`"
                          :title="`${safeGet(props.row, 'manufacturerInput.productName', $t('trace.alt.factoryImage') || '制造商图片')}`"
                          fit="cover"
                          lazy
                        >
                          <div slot="placeholder" class="image-placeholder">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.loading') || '加载中' }}
                          </div>
                          <div slot="error" class="image-error">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.error') || '加载失败' }}
                          </div>
                        </el-image>
                      </a>
                      <el-button type="text" icon="el-icon-download" :href="imgHref(props.row, 'manufacturerInput.img')" :aria-label="$t('actions.download')" @click.prevent="openDownload(imgHref(props.row, 'manufacturerInput.img'), safeGet(props.row, 'manufacturerInput.productName', 'image'))">{{ $t('actions.download') || '下载' }}</el-button>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txid') || '区块链交易ID') + ':'">
                      <span>{{ safeGet(props.row, 'manufacturerInput.txid', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txTime') || '区块链交易时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'manufacturerInput.timestamp', '')) }}</span>
                    </el-form-item>
                  </el-form>
                </el-collapse-item>
                <el-collapse-item name="driver">
                  <template v-slot:title>
                    <i class="el-icon-truck" aria-hidden="true" />
                    <span class="section-title">{{ $t('trace.sections.driver') }}</span>
                    <el-tag size="mini" type="warning" :aria-label="$t('trace.tag.logistics')">{{ $t('trace.tag.logistics') }}</el-tag>
                  </template>
                  <el-form label-position="left" inline class="demo-table-expand">
                    <el-form-item :label="($t('form.driver.name') || '运输司机姓名') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'carrierInput.name', ''), { max: 100 }) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.driver.age') || '运输司机年龄') + ':'">
                      <span>{{ safeGet(props.row, 'carrierInput.age', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.driver.phone') || '运输司机联系电话') + ':'">
                      <span>{{ safeGet(props.row, 'carrierInput.phone', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.driver.carNumber') || '车牌号') + ':'">
                      <span>{{ safeGet(props.row, 'carrierInput.plateNumber', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.driver.transport') || '运输记录') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'carrierInput.transportRecord', ''), { max: 200 }) }}</span>
                    </el-form-item>
                    <el-form-item v-if="safeGet(props.row, 'carrierInput.img', '')" :label="($t('trace.relatedImage') || '相关图片（点击下载）') + ':'" class="image-item">
                      <a :href="imgHref(props.row, 'carrierInput.img')" target="_blank" rel="noopener noreferrer" :download="safeGet(props.row, 'carrierInput.name', 'image')">
                        <el-image
                          style="width: 100px; height: 100px;"
                          :src="imgSrc(props.row, 'carrierInput.img')"
                          :preview-src-list="imgPreviewList(props.row, 'carrierInput.img')"
                          :alt="`${safeGet(props.row, 'carrierInput.name', $t('trace.alt.driverImage') || '物流图片')}`"
                          :title="`${safeGet(props.row, 'carrierInput.name', $t('trace.alt.driverImage') || '物流图片')}`"
                          fit="cover"
                          lazy
                        >
                          <div slot="placeholder" class="image-placeholder">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.loading') || '加载中' }}
                          </div>
                          <div slot="error" class="image-error">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.error') || '加载失败' }}
                          </div>
                        </el-image>
                      </a>
                      <el-button type="text" icon="el-icon-download" :href="imgHref(props.row, 'carrierInput.img')" :aria-label="$t('actions.download')" @click.prevent="openDownload(imgHref(props.row, 'carrierInput.img'), safeGet(props.row, 'carrierInput.name', 'image'))">{{ $t('actions.download') || '下载' }}</el-button>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txid') || '区块链交易ID') + ':'">
                      <span>{{ safeGet(props.row, 'carrierInput.txid', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txTime') || '区块链交易时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'carrierInput.timestamp', '')) }}</span>
                    </el-form-item>
                  </el-form>
                </el-collapse-item>
                <el-collapse-item name="shop">
                  <template v-slot:title>
                    <i class="el-icon-shopping-cart-full" aria-hidden="true" />
                    <span class="section-title">{{ $t('trace.sections.shop') }}</span>
                    <el-tag size="mini" type="primary" :aria-label="$t('trace.tag.dealer')">{{ $t('trace.tag.dealer') }}</el-tag>
                  </template>
                  <el-form label-position="left" inline class="demo-table-expand">
                    <el-form-item :label="($t('form.shop.storeTime') || '入库时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'dealerInput.storeTime', '')) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.shop.sellTime') || '销售时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'dealerInput.sellTime', '')) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.shop.name') || '经销商名称') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'dealerInput.dealerName', ''), { max: 100 }) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.shop.address') || '经销商地址') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'dealerInput.dealerLocation', ''), { max: 200 }) }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('form.shop.phone') || '经销商电话') + ':'">
                      <span>{{ normalizeDisplay(safeGet(props.row, 'dealerInput.dealerPhone', ''), { max: 30 }) }}</span>
                    </el-form-item>
                    <el-form-item v-if="safeGet(props.row, 'dealerInput.img', '')" :label="($t('trace.relatedImage') || '相关图片（点击下载）') + ':'" class="image-item">
                      <a :href="imgHref(props.row, 'dealerInput.img')" target="_blank" rel="noopener noreferrer" :download="safeGet(props.row, 'dealerInput.dealerName', 'image')">
                        <el-image
                          style="width: 100px; height: 100px;"
                          :src="imgSrc(props.row, 'dealerInput.img')"
                          :preview-src-list="imgPreviewList(props.row, 'dealerInput.img')"
                          :alt="`${safeGet(props.row, 'dealerInput.dealerName', $t('trace.alt.shopImage') || '经销商图片')}`"
                          :title="`${safeGet(props.row, 'dealerInput.dealerName', $t('trace.alt.shopImage') || '经销商图片')}`"
                          fit="cover"
                          lazy
                        >
                          <div slot="placeholder" class="image-placeholder">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.loading') || '加载中' }}
                          </div>
                          <div slot="error" class="image-error">
                            <i class="el-icon-picture-outline" /> {{ $t('trace.image.error') || '加载失败' }}
                          </div>
                        </el-image>
                      </a>
                      <el-button type="text" icon="el-icon-download" :href="imgHref(props.row, 'dealerInput.img')" :aria-label="$t('actions.download')" @click.prevent="openDownload(imgHref(props.row, 'dealerInput.img'), safeGet(props.row, 'dealerInput.dealerName', 'image'))">{{ $t('actions.download') || '下载' }}</el-button>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txid') || '区块链交易ID') + ':'">
                      <span>{{ safeGet(props.row, 'dealerInput.txid', '') }}</span>
                    </el-form-item>
                    <el-form-item :label="($t('trace.txTime') || '区块链交易时间') + ':'">
                      <span>{{ formatDateTime(safeGet(props.row, 'dealerInput.timestamp', '')) }}</span>
                    </el-form-item>
                  </el-form>
                </el-collapse-item>
              </el-collapse>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('form.traceCode') || '溯源码'"
          prop="traceabilityCode"
          sortable="custom"
        />
        <el-table-column
          :label="$t('form.factory.productName') || '产品名称'"
          prop="rawSupplierInput.productName"
          :formatter="(row)=>normalizeDisplay(safeGet(row, 'rawSupplierInput.productName', ''), { max: 100 })"
          :filters="nameFilters"
          :filter-method="filterByName"
          sortable="custom"
        />
        <el-table-column
          :label="$t('form.farmer.origin') || '产品产地'"
          prop="rawSupplierInput.rawOrigin"
          :formatter="(row)=>normalizeDisplay(safeGet(row, 'rawSupplierInput.rawOrigin', ''), { max: 100 })"
          :filters="originFilters"
          :filter-method="filterByOrigin"
          sortable="custom"
        />
        <el-table-column
          :label="$t('form.farmer.pickTime') || '产品生产时间'"
          prop="rawSupplierInput.productionTime"
          :formatter="(row)=>formatDateTime(safeGet(row, 'rawSupplierInput.productionTime', ''))"
          sortable="custom"
        />
      </el-table>
      <div v-if="!loading && filteredSortedData && filteredSortedData.length" class="table-footer">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          :page-size="pageSize"
          :current-page="currentPage"
          @size-change="onPageSizeChange"
          @current-change="onPageChange"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { getIndustrialProductInfo, getIndustrialProductList, getAllIndustrialProductInfo } from '@/api/trace'
import { sanitize } from '@/utils/sanitize'
import { LENGTHS } from '@/utils/limits'
import { apiWrap, retryLast } from '@/utils/error'
import safeGetUtil from '@/utils/safeGet'
import { buildImgUrl } from '@/utils/url'
import { normalizeResults, toList } from '@/utils/normalize'

const RECENT_KEY = 'trace_recent'
const MAX_RECENT = 10

export default {
  name: 'Trace',
  data() {
    return {
      tracedata: [],
      loading: false,
      errorMessage: null,
      traceCode: '',
      baseApi: process.env.VUE_APP_BASE_API,
      LENGTHS,
      recentQueries: [],
      recentSelected: null,
      expandedSet: new Set(),
      sortProp: null,
      sortOrder: null, // ascending | descending | null
      activeNameFilters: [],
      activeOriginFilters: [],
      currentPage: 1,
      pageSize: 20,
      total: 0
    }
  },
  computed: {
    nameFilters() {
      const set = new Set();
      (this.tracedata || []).forEach((r) => {
        const v = this.safeGet(r, 'rawSupplierInput.productName', '')
        if (v) set.add(v)
      })
      return Array.from(set).slice(0, 20).map(v => ({ text: v, value: v }))
    },
    originFilters() {
      const set = new Set();
      (this.tracedata || []).forEach((r) => {
        const v = this.safeGet(r, 'rawSupplierInput.rawOrigin', '')
        if (v) set.add(v)
      })
      return Array.from(set).slice(0, 20).map(v => ({ text: v, value: v }))
    },
    filteredSortedData() {
      let arr = Array.isArray(this.tracedata) ? this.tracedata.slice() : []
      if (this.activeNameFilters && this.activeNameFilters.length) {
        arr = arr.filter(r => this.activeNameFilters.includes(this.safeGet(r, 'rawSupplierInput.productName', '')))
      }
      if (this.activeOriginFilters && this.activeOriginFilters.length) {
        arr = arr.filter(r => this.activeOriginFilters.includes(this.safeGet(r, 'rawSupplierInput.rawOrigin', '')))
      }
      // 排序
      const prop = this.sortProp
      const order = this.sortOrder
      if (prop && order) {
        const getVal = (row) => this.safeGet(row, prop, '')
        const normalizeStr = (v) => String(v == null ? '' : v)
        const cmpDate = (a, b) => {
          const ta = this.parseTimeValue(getVal(a))
          const tb = this.parseTimeValue(getVal(b))
          return ta - tb
        }
        const cmpStr = (a, b) => normalizeStr(getVal(a)).localeCompare(normalizeStr(getVal(b)))
        const isTimeProp = prop === 'rawSupplierInput.productionTime'
        const comparator = isTimeProp ? cmpDate : cmpStr
        arr.sort((a, b) => {
          const r = comparator(a, b)
          return order === 'ascending' ? r : -r
        })
      }
      return arr
    }
  },
  watch: {
    $route(to) {
      const code = to && to.params && to.params.traceability_code
      if (code) {
        const next = (this.LENGTHS ? String(code).replace(/[^\d]/g, '').slice(0, this.LENGTHS.traceCode) : String(code).replace(/[^\d]/g, ''))
        if (next !== this.traceCode) {
          this.traceCode = next
          this.fruitInfo()
        }
      } else if (this.tracedata && this.tracedata.length) {
        this.resetQuery()
      }
    }
  },
  created() {
    this.loadRecent()
    const code = this.$route.params.traceability_code
    if (code) {
      this.traceCode = code
      this.fruitInfo()
    } else {
      apiWrap(this, () => getIndustrialProductList(), (res) => {
        const raw = Array.isArray(res.data) ? res.data : []
        this.tracedata = normalizeResults(raw)
      }, '获取初始列表失败，请稍后重试')
    }
  },
  beforeRouteUpdate(to, from, next) {
    const code = to && to.params && to.params.traceability_code
    if (code) {
      this.traceCode = (this.LENGTHS ? String(code).replace(/[^\d]/g, '').slice(0, this.LENGTHS.traceCode) : String(code).replace(/[^\d]/g, ''))
      this.fruitInfo()
    } else {
      this.resetQuery()
    }
    next()
  },
  methods: {
    onInputChange(val) {
      const v = sanitize(String(val || '').replace(/[^\d]/g, ''), LENGTHS.traceCode)
      if (v !== this.traceCode) this.traceCode = v
    },
    resetQuery() {
      this.traceCode = ''
      // keep current data until next successful fetch to avoid flicker
      this.tracedata = []
      this.recentSelected = null
      this.errorMessage = null
    },
    loadRecent() {
      try {
        const raw = localStorage.getItem(RECENT_KEY)
        this.recentQueries = raw ? JSON.parse(raw) : []
      } catch (e) {
        // ignore storage failures
        this.recentQueries = []
      }
    },
    saveRecent(code) {
      if (!code) return
      const list = this.recentQueries.slice()
      const idx = list.indexOf(code)
      if (idx !== -1) list.splice(idx, 1)
      list.unshift(code)
      if (list.length > MAX_RECENT) list.length = MAX_RECENT
      this.recentQueries = list
      try {
        localStorage.setItem(RECENT_KEY, JSON.stringify(list))
      } catch (e) {
        // storage may fail; ignore
      }
    },
    onSelectRecent(val) {
      if (!val) return
      this.traceCode = val
      this.fruitInfo()
    },
    clearRecent() {
      this.recentQueries = []
      this.recentSelected = null
      try {
        localStorage.removeItem(RECENT_KEY)
      } catch (e) {
        // storage may fail; ignore
      }
    },
    allFruitInfo() {
      const payload = { page: this.currentPage, pageSize: this.pageSize }
      apiWrap(this, () => getAllIndustrialProductInfo(payload), (res) => {
        const raw = Array.isArray(res.data && res.data.items) ? res.data.items : (Array.isArray(res.data) ? res.data : [])
        this.tracedata = normalizeResults(raw)
        this.total = typeof res.data === 'object' && res.data && typeof res.data.total === 'number' ? res.data.total : raw.length
      }, '获取所有产品信息失败，请稍后重试')
    },
    fruitInfo() {
      const code = sanitize(String(this.traceCode || '').replace(/[^\d]/g, ''), LENGTHS.traceCode)
      if (!code) {
        this.$message && this.$message.error(this.$t ? this.$t('validate.pleaseEnterValidTraceCode') : '请输入有效的溯源码')
        return
      }
      const payload = { traceabilityCode: code }
      apiWrap(this, () => getIndustrialProductInfo(payload), (res) => {
        const item = res && res.data ? res.data : null
        this.applyNewResults(item ? toList(item) : [])
        if (item) this.saveRecent(code)
      }, this.$t ? this.$t('error.fetchTraceFailed') : '查询接口请求失败，请稍后重试')
    },
    retryLast() {
      retryLast(this)
    },
    formatDateTime(val) {
      if (!val) return ''
      try {
        let d
        if (typeof val === 'number') {
          d = new Date(val > 1e12 ? val : val * 1000)
        } else if (/^\d+$/.test(val)) {
          const num = Number(val)
          d = new Date(num > 1e12 ? num : num * 1000)
        } else {
          d = new Date(val)
        }
        if (isNaN(d.getTime())) return String(val)
        const locale = (this.$i18n && this.$i18n.locale) || (typeof navigator !== 'undefined' && navigator.language) || 'zh-CN'
        const options = { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' }
        try {
          return new Intl.DateTimeFormat(locale, options).format(d)
        } catch (e) {
          // ignore and fallback
        }
        const pad = n => String(n).padStart(2, '0')
        const Y = d.getFullYear()
        const M = pad(d.getMonth() + 1)
        const D = pad(d.getDate())
        const h = pad(d.getHours())
        const m = pad(d.getMinutes())
        const s = pad(d.getSeconds())
        return `${Y}-${M}-${D} ${h}:${m}:${s}`
      } catch (e) {
        return String(val)
      }
    },
    safeGet: safeGetUtil,
    rowKey(row) {
      return row && (row.traceabilityCode || JSON.stringify(row))
    },
    isExpanded(row) {
      const key = this.rowKey(row)
      return this.expandedSet.has(key)
    },
    onExpandChange(row, expandedRows) {
      // Element-UI 的 @expand-change 传 (row, expandedRows)
      const key = this.rowKey(row)
      const isOpen = Array.isArray(expandedRows) && expandedRows.some(r => this.rowKey(r) === key)
      if (isOpen) this.expandedSet.add(key); else this.expandedSet.delete(key)
    },
    onSortChange({ prop, order }) {
      this.sortProp = prop
      this.sortOrder = order
    },
    filterByName(value, row) {
      const v = this.safeGet(row, 'rawSupplierInput.productName', '')
      return value ? v === value : true
    },
    filterByOrigin(value, row) {
      const v = this.safeGet(row, 'rawSupplierInput.rawOrigin', '')
      return value ? v === value : true
    },
    parseTimeValue(val) {
      if (!val) return 0
      if (typeof val === 'number') return val > 1e12 ? val : val * 1000
      if (/^\d+$/.test(val)) {
        const num = Number(val)
        return num > 1e12 ? num : num * 1000
      }
      const t = Date.parse(val)
      return isNaN(t) ? 0 : t
    },
    applyNewResults(arr) {
      const raw = Array.isArray(arr) ? arr : []
      this.tracedata = normalizeResults(raw)
      this.loading = false
      this.errorMessage = null
      this.expandedSet.clear()
    },
    imgSrc(row, path) {
      try {
        const file = this.safeGet(row, path, '')
        return buildImgUrl(this.baseApi, file)
      } catch (e) {
        return ''
      }
    },
    imgPreviewList(row, path) {
      const src = this.imgSrc(row, path)
      return src ? [src] : []
    },
    imgHref(row, path) {
      return this.imgSrc(row, path)
    },
    normalizeDisplay(value, opts) {
      // 延迟加载，避免未使用导入报错；从工具中动态引用
      const { normalizeDisplay } = require('@/utils/sanitize')
      return normalizeDisplay(value, opts)
    },
    openDownload(url, name) {
      try {
        const a = document.createElement('a')
        a.href = url
        a.download = name || 'image'
        a.rel = 'noopener noreferrer'
        a.target = '_blank'
        a.setAttribute('aria-label', (this.$t && this.$t('actions.download')) || '下载')
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
      } catch (e) {
        window.open(url, '_blank', 'noopener,noreferrer')
      }
    }
  },
  errorCaptured(err) {
    // 捕获渲染/子组件异常，避免白屏，直接展示错误卡片
    this.loading = false
    this.errorMessage = (err && err.message) ? err.message : '页面渲染异常，请稍后重试'
    return false // 不中断默认的错误处理
  }
}
</script>

<style lang="scss" scoped>
.demo-descriptions { margin-bottom: 12px; }

// 去除重复，合并 demo-table-expand 样式
.demo-table-expand {
  font-size: 0;
}

.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}

.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
  display: inline-block;
  vertical-align: top;
}

.demo-table-expand .image-item {
  width: 100%;
  margin-top: 10px;
  margin-bottom: 10px;
}

.demo-table-expand .image-item .el-form-item__content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.trace-container {
  margin: 30px;
  // 局部CSS变量，统一控制宽度/间距
  --toolbar-gap: 10px;
  --input-width: 300px;
  --input-margin-right: 15px;
  --select-width: 200px;
}

.trace-text {
  font-size: 30px;
  line-height: 46px;
}

// 顶部工具栏与表单控件统一样式
.trace-toolbar {
  display: flex;
  align-items: center;
  gap: var(--toolbar-gap);
  flex-wrap: wrap;
}

.trace-input {
  width: var(--input-width);
  margin-right: var(--input-margin-right);
}

.trace-select {
  width: var(--select-width);
}

.loading-placeholder {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 8px;
  color: #606266;
}

.table-footer { display:flex; justify-content:flex-end; padding: 12px 0; }

.section-title { margin-left: 6px; margin-right: 6px; font-weight: 600; }

// 基础响应式：在窄屏下控件全宽，减小间距
@media (max-width: 768px) {
  .trace-container {
    --toolbar-gap: 8px;
    --input-width: 100%;
    --input-margin-right: 0px;
    --select-width: 100%;
  }
  .trace-input,
  .trace-select {
    width: 100%;
    margin-right: 0;
  }
  .trace-toolbar {
    gap: var(--toolbar-gap);
  }
}

.image-placeholder { width:100px;height:100px;display:flex;align-items:center;justify-content:center;color:#909399;background:#f5f7fa; }
.image-error { width:100px;height:100px;display:flex;align-items:center;justify-content:center;color:#f56c6c;background:#fff; }
</style>
