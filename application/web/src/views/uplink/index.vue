<template>
  <div class="uplink-container">
    <div style="color:#909399;margin-bottom: 30px">
      当前用户：{{ name }};
      用户角色: {{ userType }}
    </div>
    <div>
      <el-form ref="form" :model="tracedata" :rules="rules" label-width="80px" size="mini" style="">
        <el-form-item v-show="userType!=='原料供应商'&&userType!=='零售商'" label="溯源码:" style="width: 500px" label-width="200px" prop="traceability_code">
          <el-input v-model.trim="tracedata.traceability_code" placeholder="请输入18位数字溯源码" clearable maxlength="18" show-word-limit />
        </el-form-item>

        <div v-show="userType==='原料供应商'">
          <el-form-item label="原料名称:" style="width: 500px" label-width="200px" prop="Farmer_input.Fa_fruitName">
            <el-input v-model="tracedata.Farmer_input.Fa_fruitName" />
          </el-form-item>
          <el-form-item label="原料产地:" style="width: 500px" label-width="200px" prop="Farmer_input.Fa_originCodes">
            <el-cascader
              v-model="tracedata.Farmer_input.Fa_originCodes"
              :options="regionOptions"
              clearable
              filterable
              placeholder="请选择省/市/区"
              style="width: 100%;"
              @change="onRegionChange"
            />
          </el-form-item>
          <el-form-item label="原料生产时间:" style="width: 500px" label-width="200px" prop="Farmer_input.Fa_plantTime">
            <el-date-picker v-model="tracedata.Farmer_input.Fa_plantTime" type="datetime" value-format="yyyy-MM-dd HH:mm:ss" placeholder="选择日期时间" style="width: 100%;" :picker-options="noFuturePickerOptions" />
          </el-form-item>
          <el-form-item label="原料到货时间:" style="width: 500px" label-width="200px" prop="Farmer_input.Fa_pickingTime">
            <el-date-picker v-model="tracedata.Farmer_input.Fa_pickingTime" type="datetime" value-format="yyyy-MM-dd HH:mm:ss" placeholder="选择日期时间" style="width: 100%;" :picker-options="noFuturePickerOptions" />
          </el-form-item>
          <el-form-item label="原料供应商名称:" style="width: 500px" label-width="200px" prop="Farmer_input.Fa_farmerName">
            <el-input v-model="tracedata.Farmer_input.Fa_farmerName" />
          </el-form-item>
        </div>
        <div v-show="userType==='制造商'">
          <el-form-item label="产品名称:" style="width: 500px" label-width="200px" prop="Factory_input.Fac_productName">
            <el-input v-model="tracedata.Factory_input.Fac_productName" />
          </el-form-item>
          <el-form-item label="生产批次:" style="width: 500px" label-width="200px" prop="Factory_input.Fac_productionbatch">
            <el-input v-model="tracedata.Factory_input.Fac_productionbatch" />
          </el-form-item>
          <el-form-item label="生产时间:" style="width: 500px" label-width="200px" prop="Factory_input.Fac_productionTime">
            <el-date-picker v-model="tracedata.Factory_input.Fac_productionTime" type="datetime" value-format="yyyy-MM-dd HH:mm:ss" placeholder="选择日期时间" style="width: 100%;" :picker-options="noFuturePickerOptions" />
          </el-form-item>
          <el-form-item label="制造商名称:" style="width: 500px" label-width="200px" prop="Factory_input.Fac_factoryName">
            <el-input v-model="tracedata.Factory_input.Fac_factoryName" />
          </el-form-item>
          <el-form-item label="制造商电话:" style="width: 500px" label-width="200px" prop="Factory_input.Fac_contactNumber">
            <el-input v-model="tracedata.Factory_input.Fac_contactNumber" />
          </el-form-item>
        </div>
        <div v-show="userType==='物流承运商'">
          <el-form-item label="姓名:" style="width: 500px" label-width="200px" prop="Driver_input.Dr_name">
            <el-input v-model="tracedata.Driver_input.Dr_name" />
          </el-form-item>
          <el-form-item label="年龄:" style="width: 500px" label-width="200px" prop="Driver_input.Dr_age">
            <el-input v-model="tracedata.Driver_input.Dr_age" />
          </el-form-item>
          <el-form-item label="联系电话:" style="width: 500px" label-width="200px" prop="Driver_input.Dr_phone">
            <el-input v-model="tracedata.Driver_input.Dr_phone" />
          </el-form-item>
          <el-form-item label="车牌号:" style="width: 500px" label-width="200px" prop="Driver_input.Dr_carNumber">
            <el-input v-model="tracedata.Driver_input.Dr_carNumber" clearable maxlength="10" show-word-limit placeholder="请输入车牌号" style="text-transform: uppercase;" @input="onCarNumberInput" />
          </el-form-item>
          <el-form-item label="运输记录：" style="width: 500px" label-width="200px" prop="Driver_input.Dr_transport">
            <el-input v-model="tracedata.Driver_input.Dr_transport" />
          </el-form-item>
        </div>
        <div v-show="userType==='经销商'">
          <el-form-item label="存入时间:" style="width: 500px" label-width="200px" prop="Shop_input.Sh_storeTime">
            <el-date-picker v-model="tracedata.Shop_input.Sh_storeTime" type="datetime" value-format="yyyy-MM-dd HH:mm:ss" placeholder="选择日期时间" style="width: 100%;" :picker-options="noFuturePickerOptions" />
          </el-form-item>
          <el-form-item label="销售时间:" style="width: 500px" label-width="200px" prop="Shop_input.Sh_sellTime">
            <el-date-picker v-model="tracedata.Shop_input.Sh_sellTime" type="datetime" value-format="yyyy-MM-dd HH:mm:ss" placeholder="选择日期时间" style="width: 100%;" :picker-options="sellPickerOptions" />
          </el-form-item>
          <el-form-item label="经销商名称:" style="width: 500px" label-width="200px" prop="Shop_input.Sh_shopName">
            <el-input v-model="tracedata.Shop_input.Sh_shopName" />
          </el-form-item>
          <el-form-item label="经销商位置:" style="width: 500px" label-width="200px" prop="Shop_input.Sh_shopAddressCodes">
            <el-cascader
              v-model="tracedata.Shop_input.Sh_shopAddressCodes"
              :options="regionOptions"
              clearable
              filterable
              placeholder="请选择省/市/区"
              style="width: 100%;"
              @change="onShopRegionChange"
            />
          </el-form-item>
          <el-form-item label="经销商电话:" style="width: 500px" label-width="200px" prop="Shop_input.Sh_shopPhone">
            <el-input v-model="tracedata.Shop_input.Sh_shopPhone" />
          </el-form-item>
        </div>
        <el-form-item v-show="userType !== '零售商'" label="上传图片:" style="width: 500px" label-width="200px">
          <el-upload
            action="#"
            class="upload-demo"
            :show-file-list="false"
            :on-change="onImageSelected"
            :before-upload="beforeUpload"
            accept="image/*"
          >
            <el-button type="primary" size="mini">选择图片</el-button>
          </el-upload>

          <div v-if="imagePreview" style="margin-top: 10px;">
            <img :src="imagePreview" alt="预览图" style="max-width: 100%; max-height: 150px; border: 1px solid #dcdfe6;">
          </div>
        </el-form-item>
      </el-form>
      <span slot="footer" style="color: gray;" class="dialog-footer">
        <el-button v-show="userType !== '零售商'" type="primary" plain style="margin-left: 220px;" @click="submittracedata()">提 交</el-button>
      </span>
      <span v-show="userType === '零售商'" slot="footer" style="color: gray;" class="dialog-footer">
        零售商没有权限录入！请使用溯源功能!
      </span>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { uplink } from '@/api/trace'
import { regionData as regionOptions } from 'element-china-area-data'

export default {
  name: 'Uplink',
  data() {
    return {
      tracedata: {
        traceability_code: '',
        Farmer_input: {
          Fa_fruitName: '',
          Fa_origin: '',
          Fa_originCodes: [],
          Fa_plantTime: null,
          Fa_pickingTime: null,
          Fa_farmerName: '',
          Fa_img: null
        },
        Factory_input: {
          Fac_productName: '',
          Fac_productionbatch: '',
          Fac_productionTime: null,
          Fac_factoryName: '',
          Fac_contactNumber: ''
        },
        Driver_input: {
          Dr_name: '',
          Dr_age: '',
          Dr_phone: '',
          Dr_carNumber: '',
          Dr_transport: ''
        },
        Shop_input: {
          Sh_storeTime: null,
          Sh_sellTime: null,
          Sh_shopName: '',
          Sh_shopAddress: '',
          Sh_shopAddressCodes: [],
          Sh_shopPhone: ''
        }
      },
      regionOptions,
      imageFile: null,
      imagePreview: null,
      loading: false,
      // 表单校验规则（在 created 中初始化，便于使用 this）
      rules: {}
    }
  },
  computed: {
    ...mapGetters([
      'name',
      'userType'
    ]),
    noFuturePickerOptions() {
      return {
        disabledDate: (time) => time.getTime() > Date.now(),
        shortcuts: this.commonShortcuts()
      }
    },
    sellPickerOptions() {
      return {
        disabledDate: (time) => {
          const t = time.getTime()
          if (t > Date.now()) return true
          const store = this.tracedata.Shop_input.Sh_storeTime
          if (!store) return false
          const sd = this.parseToLocalDate(store)
          if (!sd) return false
          const sod = new Date(sd.getFullYear(), sd.getMonth(), sd.getDate()).getTime()
          return t < sod
        },
        shortcuts: this.commonShortcuts()
      }
    }
  },
  watch: {
    userType() {
      // 角色切换时清除校验状态，避免显示无关错误
      this.$nextTick(() => {
        this.$refs.form && this.$refs.form.clearValidate()
      })
    }
  },
  created() {
    this.initRules()
  },
  methods: {
    // 初始化各字段校验规则
    initRules() {
      this.rules = {
        traceability_code: [
          { validator: this.validateTraceCode, trigger: 'blur' }
        ],
        'Farmer_input.Fa_fruitName': [
          { validator: this.validateFaFruitName, trigger: 'blur' }
        ],
        'Farmer_input.Fa_originCodes': [
          { validator: this.validateFaOriginCodes, trigger: 'change' }
        ],
        'Farmer_input.Fa_plantTime': [
          { validator: this.validateFaPlantTime, trigger: 'change' }
        ],
        'Farmer_input.Fa_pickingTime': [
          { validator: this.validateFaPickingTime, trigger: 'change' }
        ],
        'Farmer_input.Fa_farmerName': [
          { validator: this.validateFaFarmerName, trigger: 'blur' }
        ],
        'Factory_input.Fac_productName': [
          { validator: this.validateFacProductName, trigger: 'blur' }
        ],
        'Factory_input.Fac_productionbatch': [
          { validator: this.validateFacBatch, trigger: 'blur' }
        ],
        'Factory_input.Fac_productionTime': [
          { validator: this.validateFacProductionTime, trigger: 'change' }
        ],
        'Factory_input.Fac_factoryName': [
          { validator: this.validateFacFactoryName, trigger: 'blur' }
        ],
        'Factory_input.Fac_contactNumber': [
          { validator: this.validatePhoneForFactory, trigger: 'blur' }
        ],
        'Driver_input.Dr_name': [
          { validator: this.validateDrName, trigger: 'blur' }
        ],
        'Driver_input.Dr_age': [
          { validator: this.validateDrAge, trigger: 'blur' }
        ],
        'Driver_input.Dr_phone': [
          { validator: this.validateDrPhone, trigger: 'blur' }
        ],
        'Driver_input.Dr_carNumber': [
          { validator: this.validateDrCarNumber, trigger: 'blur' }
        ],
        'Driver_input.Dr_transport': [
          { validator: this.validateDrTransport, trigger: 'blur' }
        ],
        'Shop_input.Sh_storeTime': [
          { validator: this.validateShStoreTime, trigger: 'change' }
        ],
        'Shop_input.Sh_sellTime': [
          { validator: this.validateShSellTime, trigger: 'change' }
        ],
        'Shop_input.Sh_shopName': [
          { validator: this.validateShShopName, trigger: 'blur' }
        ],
        'Shop_input.Sh_shopAddressCodes': [
          { validator: this.validateShShopAddressCodes, trigger: 'change' }
        ],
        'Shop_input.Sh_shopPhone': [
          { validator: this.validateShShopPhone, trigger: 'blur' }
        ]
      }
    },

    // 将 codes 数组转成 "省-市-区" 文本，独立于 CodeToText，避免运行时不兼容
    codesToText(codes) {
      if (!Array.isArray(codes) || codes.length === 0) return ''
      const names = []
      let nodes = this.regionOptions
      for (const code of codes) {
        if (!Array.isArray(nodes)) { names.push(String(code)); break }
        const found = nodes.find(n => n.value === code)
        if (!found) { names.push(String(code)); break }
        names.push(found.label)
        nodes = found.children || []
      }
      return names.join('-')
    },

    // 基础校验工具
    isValidPhone(value) {
      if (!value) return false
      const simple = /^[0-9\-+()\s]{6,20}$/
      const cnMobile = /^1[3-9]\d{9}$/
      return simple.test(value) || cnMobile.test(value)
    },
    isValidCarNumber(value) {
      if (!value) return false
      // 宽松校验：5-10 位字母数字（兼容多地区车牌）
      return /^[A-Za-z0-9]{5,10}$/.test(value)
    },

    // 具体字段校验器（按角色动态启用）
    validateTraceCode(rule, value, callback) {
      if (this.userType === '原料供应商' || this.userType === '零售商') return callback()
      const v = (value || '').trim()
      if (!v) return callback(new Error('请输入溯源码'))
      const ok = /^\d{18}$/.test(v)
      return ok ? callback() : callback(new Error('溯源码必须为18位数字'))
    },

    // 农户
    validateFaFruitName(rule, value, callback) {
      if (this.userType !== '原料供应商') return callback()
      return value ? callback() : callback(new Error('请输入原料名称'))
    },
    validateFaOriginCodes(rule, value, callback) {
      if (this.userType !== '原料供应商') return callback()
      if (!Array.isArray(value) || value.length === 0) return callback(new Error('请选择原料产地'))
      return callback()
    },
    validateFaPlantTime(rule, value, callback) {
      if (this.userType !== '原料供应商') return callback()
      if (!value) return callback(new Error('请选择原料生产时间'))
      if (this.isFutureStringDate(value)) return callback(new Error('原料生产时间不能晚于当前时间'))
      return callback()
    },
    validateFaPickingTime(rule, value, callback) {
      if (this.userType !== '原料供应商') return callback()
      if (!value) return callback(new Error('请选择原料到货时间'))
      if (this.isFutureStringDate(value)) return callback(new Error('原料到货时间不能晚于当前时间'))
      return callback()
    },
    validateFaFarmerName(rule, value, callback) {
      if (this.userType !== '原料供应商') return callback()
      return value ? callback() : callback(new Error('请输入供应商名称'))
    },

    // 制造商
    validateFacProductName(rule, value, callback) {
      if (this.userType !== '制造商') return callback()
      return value ? callback() : callback(new Error('请输入产品名称'))
    },
    validateFacBatch(rule, value, callback) {
      if (this.userType !== '制造商') return callback()
      if (!value) return callback(new Error('请输入生产批次'))
      const ok = /^[A-Za-z0-9_-]{1,32}$/.test(value)
      return ok ? callback() : callback(new Error('批次格式不正确(1-32位字母/数字/_-)'))
    },
    validateFacProductionTime(rule, value, callback) {
      if (this.userType !== '制造商') return callback()
      if (!value) return callback(new Error('请选择生产时间'))
      if (this.isFutureStringDate(value)) return callback(new Error('生产时间不能晚于当前时间'))
      return callback()
    },
    validateFacFactoryName(rule, value, callback) {
      if (this.userType !== '制造商') return callback()
      return value ? callback() : callback(new Error('请输入制造商名称与厂址'))
    },
    validatePhoneForFactory(rule, value, callback) {
      if (this.userType !== '制造商') return callback()
      if (!value) return callback(new Error('请输入制造商电话'))
      return this.isValidPhone(value) ? callback() : callback(new Error('请输入有效的电话号码'))
    },

    // 司机
    validateDrName(rule, value, callback) {
      if (this.userType !== '物流承运商') return callback()
      return value ? callback() : callback(new Error('请输入姓名'))
    },
    validateDrAge(rule, value, callback) {
      if (this.userType !== '物流承运商') return callback()
      if (value === '' || value === null || value === undefined) return callback(new Error('请输入年龄'))
      const n = Number(value)
      if (!Number.isInteger(n)) return callback(new Error('年龄必须为整数'))
      if (n < 18 || n > 70) return callback(new Error('年龄需在18-70之间'))
      return callback()
    },
    validateDrPhone(rule, value, callback) {
      if (this.userType !== '物流承运商') return callback()
      if (!value) return callback(new Error('请输入联系电话'))
      return this.isValidPhone(value) ? callback() : callback(new Error('请输入有效的联系电话'))
    },
    validateDrCarNumber(rule, value, callback) {
      if (this.userType !== '物流承运商') return callback()
      if (!value) return callback(new Error('请输入车牌号'))
      return this.isValidCarNumber(value) ? callback() : callback(new Error('车牌号格式不正确'))
    },
    validateDrTransport(rule, value, callback) {
      if (this.userType !== '物流承运商') return callback()
      return value ? callback() : callback(new Error('请输入运输记录'))
    },

    // 经销商
    validateShStoreTime(rule, value, callback) {
      if (this.userType !== '经销商') return callback()
      if (!value) return callback(new Error('请选择存入时间'))
      if (this.isFutureStringDate(value)) return callback(new Error('存入时间不能晚于当前时间'))
      return callback()
    },
    validateShSellTime(rule, value, callback) {
      if (this.userType !== '经销商') return callback()
      if (!value) return callback(new Error('请选择销售时间'))
      if (this.isFutureStringDate(value)) return callback(new Error('销售时间不能晚于当前时间'))
      const store = this.tracedata.Shop_input.Sh_storeTime
      if (store) {
        const sellD = this.parseToLocalDate(value)
        const storeD = this.parseToLocalDate(store)
        if (sellD && storeD && sellD.getTime() < storeD.getTime()) {
          return callback(new Error('销售时间不能早于存入时间'))
        }
      }
      return callback()
    },
    validateShShopName(rule, value, callback) {
      if (this.userType !== '经销商') return callback()
      return value ? callback() : callback(new Error('请输入经销商名称'))
    },
    validateShShopAddressCodes(rule, value, callback) {
      if (this.userType !== '经销商') return callback()
      if (!Array.isArray(value) || value.length === 0) return callback(new Error('请选择经销商位置'))
      return callback()
    },
    validateShShopPhone(rule, value, callback) {
      if (this.userType !== '经销商') return callback()
      if (!value) return callback(new Error('请输入经销商电话'))
      return this.isValidPhone(value) ? callback() : callback(new Error('请输入有效的经销商电话'))
    },

    onRegionChange(codes) {
      this.tracedata.Farmer_input.Fa_origin = this.codesToText(codes)
    },
    onShopRegionChange(codes) {
      this.tracedata.Shop_input.Sh_shopAddress = this.codesToText(codes)
    },

    submittracedata() {
      // 先触发表单校验
      this.$refs.form.validate((valid) => {
        if (!valid) {
          this.$message({ message: '请先修正表单校验错误', type: 'error' })
          return
        }
        // 兜底：若未触发 change，也根据 codes 生成一次文本
        if (this.userType === '原料供应商' && (!this.tracedata.Farmer_input.Fa_origin) && Array.isArray(this.tracedata.Farmer_input.Fa_originCodes)) {
          this.tracedata.Farmer_input.Fa_origin = this.codesToText(this.tracedata.Farmer_input.Fa_originCodes)
        }
        if (this.userType === '经销商' && (!this.tracedata.Shop_input.Sh_shopAddress) && Array.isArray(this.tracedata.Shop_input.Sh_shopAddressCodes)) {
          this.tracedata.Shop_input.Sh_shopAddress = this.codesToText(this.tracedata.Shop_input.Sh_shopAddressCodes)
        }
        console.log(this.tracedata)
        const loading = this.$loading({
          lock: true,
          text: '数据上链中...',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        })
        var formData = new FormData()
        formData.append('traceability_code', this.tracedata.traceability_code)
        formData.append('file', this.imageFile)
        // 根据不同的用户给arg1、arg2、arg3..赋值,
        switch (this.userType) {
          case '原料供应商':
            formData.append('arg1', this.tracedata.Farmer_input.Fa_fruitName)
            formData.append('arg2', this.tracedata.Farmer_input.Fa_origin)
            formData.append('arg3', this.tracedata.Farmer_input.Fa_plantTime)
            formData.append('arg4', this.tracedata.Farmer_input.Fa_pickingTime)
            formData.append('arg5', this.tracedata.Farmer_input.Fa_farmerName)
            break
          case '制造商':
            formData.append('arg1', this.tracedata.Factory_input.Fac_productName)
            formData.append('arg2', this.tracedata.Factory_input.Fac_productionbatch)
            formData.append('arg3', this.tracedata.Factory_input.Fac_productionTime)
            formData.append('arg4', this.tracedata.Factory_input.Fac_factoryName)
            formData.append('arg5', this.tracedata.Factory_input.Fac_contactNumber)
            break
          case '物流承运商':
            formData.append('arg1', this.tracedata.Driver_input.Dr_name)
            formData.append('arg2', this.tracedata.Driver_input.Dr_age)
            formData.append('arg3', this.tracedata.Driver_input.Dr_phone)
            formData.append('arg4', this.tracedata.Driver_input.Dr_carNumber)
            formData.append('arg5', this.tracedata.Driver_input.Dr_transport)
            break
          case '经销商':
            formData.append('arg1', this.tracedata.Shop_input.Sh_storeTime)
            formData.append('arg2', this.tracedata.Shop_input.Sh_sellTime)
            formData.append('arg3', this.tracedata.Shop_input.Sh_shopName)
            formData.append('arg4', this.tracedata.Shop_input.Sh_shopAddress)
            formData.append('arg5', this.tracedata.Shop_input.Sh_shopPhone)
            break
        }
        console.log(formData)
        for (const [k, v] of formData.entries()) { console.log(k, v) }
        uplink(formData).then(res => {
          if (res.code === 200) {
            loading.close()
            this.$message({
              message: '上链成功，交易ID：' + res.txid + '\n溯源码：' + res.traceability_code,
              type: 'success'
            })
          } else {
            loading.close()
            this.$message({
              message: '上链失败',
              type: 'error'
            })
          }
        }).catch(err => {
          loading.close()
          console.log(err)
        })
      })
    },
    onImageSelected(file) {
      // 生成预览图
      this.imagePreview = URL.createObjectURL(file.raw)
      this.imageFile = file.raw
    },
    beforeUpload() {
      // 禁止自动上传
      return false
    },
    commonShortcuts() {
      return [
        {
          text: '此刻',
          onClick: (picker) => picker.$emit('pick', new Date())
        },
        {
          text: '30分钟前',
          onClick: (picker) => picker.$emit('pick', new Date(Date.now() - 30 * 60 * 1000))
        },
        {
          text: '1小时前',
          onClick: (picker) => picker.$emit('pick', new Date(Date.now() - 60 * 60 * 1000))
        },
        {
          text: '昨天此刻',
          onClick: (picker) => picker.$emit('pick', new Date(Date.now() - 24 * 60 * 60 * 1000))
        },
        {
          text: '一周前此刻',
          onClick: (picker) => picker.$emit('pick', new Date(Date.now() - 7 * 24 * 60 * 60 * 1000))
        }
      ]
    },
    parseToLocalDate(str) {
      if (!str || typeof str !== 'string') return null
      const s = str.replace(' ', 'T')
      const d = new Date(s)
      return isNaN(d.getTime()) ? null : d
    },
    isFutureStringDate(str) {
      const d = this.parseToLocalDate(str)
      if (!d) return false
      return d.getTime() > Date.now()
    }
  }
}

</script>

<style lang="scss" scoped>
.uplink {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
