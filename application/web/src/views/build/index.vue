<template>
  <div>
    <h1 style="color: #1f2f3d; text-align: center;">{{ $t('build.title') }}</h1>
    <p style="color: #5e6d82; text-align: center;">{{ $t('build.p1') }}</p>
    <p style="color: #5e6d82; text-align: center;">{{ $t('build.p2') }}</p>
    <div style="text-align: center; margin-bottom: 20px;">
      <a href="https://www.bilibili.com/video/BV1Ar421H7TK" style="color: #409EFF; text-decoration: underline;">
        {{ $t('build.tutorial') }}
      </a>
    </div>
    <el-button type="text" @click="dialog2Visible = true">{{ $t('build.buyCode') }}</el-button>
    <el-dialog
      :title="$t('build.deletePage')"
      :visible.sync="dialog1Visible"
      width="30%"
      :before-close="handleClose"
    >
      <span>{{ $t('build.deleteStep1') }}</span>
      <br>
      <span style="display: block;margin-top: 20px;">{{ $t('build.deleteStep2') }}</span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialog1Visible = false">{{ $t('build.confirm') }}</el-button>
      </span>
    </el-dialog>

    <el-dialog
      :title="$t('build.supportTitle')"
      :visible.sync="dialog2Visible"
      width="30%"
      :before-close="handleClose"
    >
      <span>{{ $t('build.price') }}</span>

      <br>
      <span style="display: block; margin-top: 20px;">{{ $t('build.contact') }}</span>
      <br>
      <span style="display: block; margin-top: 20px;">{{ $t('build.discount') }}</span>
      <br>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialog2Visible = false">{{ $t('build.confirm') }}</el-button>
      </span>
    </el-dialog>

    <el-dialog
      :title="$t('build.buildSuccess')"
      :visible.sync="dialog3Visible"
      width="30%"
      :before-close="handleClose"
    >
      <span>{{ $t('build.download') }}</span>
      <br>
      <span style="display: block;margin-top: 20px;">{{ downloadUrl }} </span>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialog3Visible = false">{{ $t('build.confirm') }}</el-button>
      </span>
    </el-dialog>

    <div class="form-container">
      <el-form ref="form" :model="form" label-width="150px" class="form">

        <el-form-item v-for="(v, key) in form" :key="key" :label="key">
          <el-input v-model="form[key]" :placeholder="$t('validate.pleaseEnter') + ' ' + key + ' ' + '值'" />
        </el-form-item>
        <div style="text-align: center; margin-top: 20px;">
          <el-button :loading="loading" type="primary" @click="submitForm">{{ $t('build.startBuild') }}</el-button>
        </div>
      </el-form>
      <el-button type="text" style="display: block;margin-top: 20px;" @click="dialog1Visible = true">{{ $t('build.howToDelete') }}</el-button>
    </div>
    <div style="height: 30px;" />
    <el-card shadow="hover" class="card-item">
      <div slot="header" class="clearfix">
        <span>系统介绍</span>
      </div>
      <div class="text item">
        <p>parm1: 基于区块链的工业产品溯源系统</p>
        <p>parm2: 工业产品信息</p>
        <p>parm3: 工业产品名称</p>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Build',
  data() {
    return {
      loading: false,
      dialog1Visible: false,
      dialog2Visible: false,
      form: {
        parm1: '基于区块链的农产品溯源系统',
        parm2: '农产品信息',
        parm3: '农产品名称',
        parm4: '产地',
        parm5: '种植时间',
        parm6: '采摘时间',
        parm7: '种植户名称',
        parm8: '种植户',
        parm9: '工厂信息',
        parm10: '商品名称',
        parm11: '生产批次',
        parm12: '生产时间',
        parm13: '工厂名称与厂址',
        parm14: '工厂电话',
        parm15: '工厂',
        parm16: '物流轨迹信息',
        parm17: '姓名',
        parm18: '年龄',
        parm19: '联系电话',
        parm20: '车牌号',
        parm21: '运输记录',
        parm22: '物流司机',
        parm23: '商店信息',
        parm24: '入库时间',
        parm25: '销售时间',
        parm26: '商店名称',
        parm27: '商店位置',
        parm28: '商店电话',
        parm29: '商店',
        parm30: '消费者',
        parm31: '云服务器IP,例如：32.12.243.30/192.168.1.20',
        activatecode: '激活码'
      },
      downloadUrl: '',
      dialog3Visible: false
    }
  },
  methods: {
    submitForm() {
      this.loading = true
      const params = new URLSearchParams()
      for (const key in this.form) {
        params.append(key, this.form[key])
      }
      fetch('http://realcool.top:8088/activate', {
        method: 'POST',
        body: params,
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
      })
        .then(response => response.json())
        .then(data => {
          this.loading = false
          this.downloadUrl = data.msg
          // eslint-disable-next-line eqeqeq
          if (data.code == 0) {
            this.dialog3Visible = true
          } else {
            this.$message.error(this.$t('build.buildFailed', { msg: data.msg }))
          }
        })
        .catch(error => {
          this.$message.error(this.$t('build.submitFailed'), { msg: error.toString() })
          this.loading = false
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.form-container {
  max-width: 700px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f9f9f9;
  border: 1px solid #ebebeb;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.el-form-item {
  margin-bottom: 15px;
}

.el-button {
  width: 100%;
}
</style>
