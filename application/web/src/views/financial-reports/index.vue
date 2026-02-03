<template>
  <div class="app-container">
    <h2>📊 财务报表 - 数据级别限制示例</h2>

    <el-alert
      title="权限说明"
      :description="`当前用户: ${username} | 角色: ${userType} | 数据级别: ${userDataLevel || '未设置'}`"
      type="info"
      show-icon
      style="margin-bottom: 20px"
    />

    <div v-if="!hasPermission" class="permission-denied">
      <el-empty
        description="您没有权限访问此页面"
        :image-size="200"
      >
        <div style="margin-top: 20px; color: #606266;">
          <p>❌ 权限不足</p>
          <p v-if="userDataLevel">您的数据级别为: <strong>{{ userDataLevel }}</strong></p>
          <p v-else>您没有设置数据级别属性</p>
          <p>此页面需要数据级别为: <strong>Internal</strong></p>
          <p style="font-size: 12px; color: #909399; margin-top: 10px;">
            💡 提示: 请联系管理员为您提升数据访问权限
          </p>
          <el-button type="primary" @click="goBack" style="margin-top: 20px">返回首页</el-button>
        </div>
      </el-empty>
    </div>

    <div v-else class="permission-granted">
      <el-card class="box-card" style="margin-bottom: 20px;">
        <div slot="header" class="clearfix">
          <span>✅ 权限检查通过</span>
          <el-tag type="success" style="float: right;">Internal 级别数据</el-tag>
        </div>
        <p>您具有访问内部机密财务报表的权限</p>
        <p style="color: #909399; font-size: 12px; margin-top: 10px;">⚠️ 此数据为机密信息，请勿外泄</p>
      </el-card>

      <el-row :gutter="20">
        <el-col :span="6">
          <el-statistic
            title="总收入"
            :value="1234567"
            :value-style="{ color: '#27ae60', fontSize: '24px' }"
            suffix="元"
          />
        </el-col>
        <el-col :span="6">
          <el-statistic
            title="总支出"
            :value="987654"
            :value-style="{ color: '#e74c3c', fontSize: '24px' }"
            suffix="元"
          />
        </el-col>
        <el-col :span="6">
          <el-statistic
            title="净利润"
            :value="246913"
            :value-style="{ color: '#3498db', fontSize: '24px' }"
            suffix="元"
          />
        </el-col>
        <el-col :span="6">
          <el-statistic
            title="利润率"
            :value="20"
            :value-style="{ color: '#f39c12', fontSize: '24px' }"
            suffix="%"
          />
        </el-col>
      </el-row>

      <el-divider></el-divider>

      <h4>📈 月度数据统计</h4>
      <el-table
        :data="financialData"
        border
        style="width: 100%; margin-top: 20px"
      >
        <el-table-column prop="month" label="月份" width="100" />
        <el-table-column prop="revenue" label="收入" width="120" />
        <el-table-column prop="expense" label="支出" width="120" />
        <el-table-column prop="profit" label="利润" width="120">
          <template slot-scope="scope">
            <span style="color: #27ae60; font-weight: bold;">{{ scope.row.profit }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="rate" label="利润率" width="100">
          <template slot-scope="scope">
            <el-progress :percentage="parseInt(scope.row.rate)" />
          </template>
        </el-table-column>
        <el-table-column prop="level" label="数据级别" width="100">
          <template slot-scope="scope">
            <el-tag type="danger">{{ scope.row.level }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { hasPermission } from '@/utils/checkPermission'

export default {
  name: 'FinancialReports',
  data() {
    return {
      financialData: [
        { month: '2026-01', revenue: '¥500,000', expense: '¥400,000', profit: '¥100,000', rate: '20', level: 'Internal' },
        { month: '2026-02', revenue: '¥550,000', expense: '¥420,000', profit: '¥130,000', rate: '24', level: 'Internal' },
        { month: '2026-03', revenue: '¥600,000', expense: '¥450,000', profit: '¥150,000', rate: '25', level: 'Internal' }
      ]
    }
  },
  computed: {
    ...mapGetters(['name', 'userType', 'dynamicAttributes']),
    username() {
      return this.name || '未知用户'
    },
    userDataLevel() {
      return this.dynamicAttributes?.data_level || null
    },
    hasPermission() {
      const user = {
        userType: this.userType,
        dynamicAttributes: this.dynamicAttributes || {}
      }
      const resource = {
        requiredAttributes: { data_level: 'Internal' }
      }
      return hasPermission(user, resource)
    }
  },
  methods: {
    goBack() {
      this.$router.push('/')
    }
  },
  mounted() {
    if (!this.hasPermission) {
      this.$message.warning('您没有权限访问此页面')
    }
  }
}
</script>

<style scoped>
.permission-denied {
  padding: 40px 20px;
  text-align: center;
}

.permission-granted {
  animation: slideIn 0.3s ease-in;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.clearfix:after {
  content: "";
  display: table;
  clear: both;
}
</style>

