<template>
  <div class="app-container">
    <h2>📋 订单管理 - 地域限制示例</h2>

    <el-alert
      title="权限说明"
      :description="`当前用户: ${username} | 角色: ${userType} | 地域: ${userRegion || '未设置'}`"
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
          <p v-if="userRegion">您的地域属性为: <strong>{{ userRegion }}</strong></p>
          <p v-else>您没有设置地域属性</p>
          <p>此页面需要地域属性为: <strong>Sichuan</strong></p>
          <el-button type="primary" @click="goBack" style="margin-top: 20px">返回首页</el-button>
        </div>
      </el-empty>
    </div>

    <div v-else class="permission-granted">
      <el-card class="box-card" style="margin-bottom: 20px;">
        <div slot="header" class="clearfix">
          <span>✅ 权限检查通过</span>
          <el-tag type="success" style="float: right;">Sichuan 地域</el-tag>
        </div>
        <p>您具有访问四川地域订单的权限</p>
      </el-card>

      <el-table
        :data="orders"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="订单ID" width="100" />
        <el-table-column prop="orderNo" label="订单号" width="150" />
        <el-table-column prop="region" label="地域" width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.region === 'Sichuan'" type="success">{{ scope.row.region }}</el-tag>
            <el-tag v-else type="info">{{ scope.row.region }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === '已完成' ? 'success' : 'info'">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="100" />
        <el-table-column prop="createTime" label="创建时间" width="150" />
      </el-table>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { hasPermission } from '@/utils/checkPermission'

export default {
  name: 'OrdersRegion',
  data() {
    return {
      orders: [
        { id: 1, orderNo: 'ORD001', region: 'Sichuan', status: '已完成', amount: '¥1000', createTime: '2026-02-01' },
        { id: 2, orderNo: 'ORD002', region: 'Sichuan', status: '处理中', amount: '¥2000', createTime: '2026-02-02' },
        { id: 3, orderNo: 'ORD003', region: 'Sichuan', status: '已完成', amount: '¥1500', createTime: '2026-02-03' }
      ]
    }
  },
  computed: {
    ...mapGetters(['name', 'userType', 'dynamicAttributes']),
    username() {
      return this.name || '未知用户'
    },
    userRegion() {
      return this.dynamicAttributes?.region || null
    },
    hasPermission() {
      const user = {
        userType: this.userType,
        dynamicAttributes: this.dynamicAttributes || {}
      }
      const resource = {
        requiredAttributes: { region: 'Sichuan' }
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

