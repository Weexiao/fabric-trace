<template>
  <div class="trace-container">
    <el-input v-model="input" :placeholder="$t ? $t('form.inputTraceCode') : '请输入溯源码查询'" style="width: 300px;margin-right: 15px;" :maxlength="LENGTHS.traceCode" @input="onInputChange" />
    <el-button type="primary" plain @click="FruitInfo"> 查询 </el-button>
    <el-button type="success" plain @click="AllFruitInfo"> 获取所有产品信息 </el-button>
    <el-table
      :data="tracedata"
      style="width: 100%"
    >
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-form label-position="left" inline class="demo-table-expand">
            <div><span class="trace-text" style="color: #67C23A;">原料信息</span></div>
            <el-form-item label="原料名称:">
              <span>{{ props.row.farmer_input.fa_fruitName }}</span>
            </el-form-item>
            <el-form-item label="原料产地:">
              <span>{{ props.row.farmer_input.fa_origin }}</span>
            </el-form-item>
            <el-form-item label="原料生产时间:">
              <span>{{ props.row.farmer_input.fa_plantTime }}</span>
            </el-form-item>
            <el-form-item label="原料到货时间:">
              <span>{{ props.row.farmer_input.fa_pickingTime }}</span>
            </el-form-item>
            <el-form-item label="原料供应商名称:">
              <span>{{ props.row.farmer_input.fa_farmerName }}</span>
            </el-form-item>
            <el-form-item v-if="props.row.farmer_input.fa_img" label="相关图片（点击下载）:" class="image-item">
              <a :href="`${baseApi}getImg/${props.row.farmer_input.fa_img}`" target="_blank">
                <el-image
                  style="width: 100px; height: 100px;"
                  :src="`${baseApi}getImg/${props.row.farmer_input.fa_img}`"
                  fit="cover"
                />
              </a>
            </el-form-item>
            <el-form-item label="区块链交易ID:">
              <span>{{ props.row.farmer_input.fa_txid }}</span>
            </el-form-item>
            <el-form-item label="区块链交易时间:">
              <span>{{ props.row.farmer_input.fa_timestamp }}</span>
            </el-form-item>
            <div><span class="trace-text" style="color: #409EFF;">制造商信息</span></div>
            <el-form-item label="产品名称:">
              <span>{{ props.row.factory_input.fac_productName }}</span>
            </el-form-item>
            <el-form-item label="生产批次:">
              <span>{{ props.row.factory_input.fac_productionbatch }}</span>
            </el-form-item>
            <el-form-item label="生产时间:">
              <span>{{ props.row.factory_input.fac_productionTime }}</span>
            </el-form-item>
            <el-form-item label="制造商名称与厂址:">
              <span>{{ props.row.factory_input.fac_factoryName }}</span>
            </el-form-item>
            <el-form-item label="制造商电话:">
              <span>{{ props.row.factory_input.fac_contactNumber }}</span>
            </el-form-item>
            <el-form-item v-if="props.row.factory_input.fac_img" label="相关图片（点击下载）:" class="image-item">
              <a :href="`${baseApi}getImg/${props.row.factory_input.fac_img}`" target="_blank">
                <el-image
                  style="width: 100px; height: 100px;"
                  :src="`${baseApi}getImg/${props.row.factory_input.fac_img}`"
                  fit="cover"
                />
              </a>
            </el-form-item>
            <el-form-item label="区块链交易ID:">
              <span>{{ props.row.factory_input.fac_txid }}</span>
            </el-form-item>
            <el-form-item label="区块链交易时间:">
              <span>{{ props.row.factory_input.fac_timestamp }}</span>
            </el-form-item>
            <div><span class="trace-text" style="color: #E6A23C;">物流承运商信息</span></div>
            <el-form-item label="运输司机姓名:">
              <span>{{ props.row.driver_input.dr_name }}</span>
            </el-form-item>
            <el-form-item label="运输司机年龄:">
              <span>{{ props.row.driver_input.dr_age }}</span>
            </el-form-item>
            <el-form-item label="运输司机联系电话:">
              <span>{{ props.row.driver_input.dr_phone }}</span>
            </el-form-item>
            <el-form-item label="车牌号:">
              <span>{{ props.row.driver_input.dr_carNumber }}</span>
            </el-form-item>
            <el-form-item label="运输记录:">
              <span>{{ props.row.driver_input.dr_transport }}</span>
            </el-form-item>
            <el-form-item v-if="props.row.driver_input.dr_img" label="相关图片（点击下载）:" class="image-item">
              <a :href="`${baseApi}getImg/${props.row.driver_input.dr_img}`" target="_blank">
                <el-image
                  style="width: 100px; height: 100px;"
                  :src="`${baseApi}getImg/${props.row.driver_input.dr_img}`"
                  fit="cover"
                />
              </a>
            </el-form-item>
            <el-form-item label="区块链交易ID:">
              <span>{{ props.row.driver_input.dr_txid }}</span>
            </el-form-item>
            <el-form-item label="区块链交易时间:">
              <span>{{ props.row.driver_input.dr_timestamp }}</span>
            </el-form-item>
            <div><span class="trace-text" style="color: #909399;">经销商信息</span></div>
            <el-form-item label="入库时间:">
              <span>{{ props.row.shop_input.sh_storeTime }}</span>
            </el-form-item>
            <el-form-item label="销售时间:">
              <span>{{ props.row.shop_input.sh_sellTime }}</span>
            </el-form-item>
            <el-form-item label="经销商名称:">
              <span>{{ props.row.shop_input.sh_shopName }}</span>
            </el-form-item>
            <el-form-item label="经销商地址:">
              <span>{{ props.row.shop_input.sh_shopAddress }}</span>
            </el-form-item>
            <el-form-item label="经销商电话:">
              <span>{{ props.row.shop_input.sh_shopPhone }}</span>
            </el-form-item>
            <el-form-item v-if="props.row.shop_input.sh_img" label="相关图片(点击下载）:" class="image-item">
              <a :href="`${baseApi}getImg/${props.row.shop_input.sh_img}`" target="_blank">
                <el-image
                  style="width: 100px; height: 100px;"
                  :src="`${baseApi}getImg/${props.row.shop_input.sh_img}`"
                  fit="cover"
                />
              </a>
            </el-form-item>
            <el-form-item label="区块链交易ID:">
              <span>{{ props.row.shop_input.sh_txid }}</span>
            </el-form-item>
            <el-form-item label="区块链交易时间:">
              <span>{{ props.row.shop_input.sh_timestamp }}</span>
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column
        label="溯源码"
        prop="traceability_code"
      />
      <el-table-column
        label="产品名称"
        prop="farmer_input.fa_fruitName"
      />
      <el-table-column
        label="产品生产时间"
        prop="farmer_input.fa_pickingTime"
      />
    </el-table>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { getFruitInfo, getFruitList, getAllFruitInfo } from '@/api/trace'
import { sanitize } from '@/utils/sanitize'
import { LENGTHS } from '@/utils/limits'

export default {
  name: 'Trace',
  data() {
    return {
      tracedata: [],
      loading: false,
      input: '',
      baseApi: process.env.VUE_APP_BASE_API,
      LENGTHS
    }
  },
  computed: {
    ...mapGetters([
      'name',
      'userType'
    ])
  },
  created() {
    const code = this.$route.params.traceability_code
    if (code) {
      this.input = code
      this.FruitInfo()
    } else {
      getFruitList().then(res => {
        this.tracedata = JSON.parse(res.data)
      })
    }
  },
  methods: {
    onInputChange(val) {
      const v = sanitize(String(val || '').replace(/[^\d]/g, ''), LENGTHS.traceCode)
      if (v !== this.input) this.input = v
    },
    AllFruitInfo() {
      getAllFruitInfo().then(res => {
        this.tracedata = JSON.parse(res.data)
      })
    },
    FruitInfo() {
      const code = sanitize(String(this.input || '').replace(/[^\d]/g, ''), LENGTHS.traceCode)
      if (!code) {
        this.$message.error('请输入有效的溯源码')
        return
      }
      var formData = new FormData()
      formData.append('traceability_code', code)
      getFruitInfo(formData).then(res => {
        if (res.code === 200) {
          this.tracedata = []
          try {
            this.tracedata[0] = JSON.parse(res.data)
          } catch (e) {
            this.$message.error('返回数据解析失败')
          }
        } else {
          this.$message.error(res.message)
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>

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
  }
.trace {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}

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

</style>
