<template>
  <div class="app-container">
    <h2>🔐 动态属性权限控制 - 演示中心</h2>

    <el-card class="box-card" style="margin-bottom: 20px;">
      <div slot="header" class="clearfix">
        <span>👤 当前用户信息</span>
      </div>
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="user-info-item">
            <p class="label">用户名</p>
            <p class="value">{{ username }}</p>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="user-info-item">
            <p class="label">用户类型</p>
            <p class="value"><el-tag>{{ userType }}</el-tag></p>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="user-info-item">
            <p class="label">地域属性</p>
            <p class="value">
              <el-tag v-if="userRegion" type="success">{{ userRegion }}</el-tag>
              <el-tag v-else type="info">未设置</el-tag>
            </p>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="user-info-item">
            <p class="label">数据级别</p>
            <p class="value">
              <el-tag v-if="userDataLevel" type="danger">{{ userDataLevel }}</el-tag>
              <el-tag v-else type="info">未设置</el-tag>
            </p>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <el-divider>权限演示</el-divider>

    <el-row :gutter="20">
      <!-- 演示1：地域限制 -->
      <el-col :span="12">
        <el-card class="demo-card">
          <div slot="header" class="clearfix">
            <span>🌍 演示1：地域限制</span>
            <span style="float: right; font-size: 12px;">订单管理（Sichuan）</span>
          </div>

          <p style="margin-bottom: 10px;">
            <strong>需求：</strong> 仅 Sichuan 地域用户可访问订单
          </p>

          <el-alert
            v-if="canAccessOrders"
            title="✅ 权限检查通过"
            type="success"
            :closable="false"
            style="margin-bottom: 15px"
          >
            您具有 Sichuan 地域属性，可以访问订单管理
          </el-alert>

          <el-alert
            v-else
            title="❌ 权限检查未通过"
            type="error"
            :closable="false"
            style="margin-bottom: 15px"
          >
            <div>
              <p v-if="!userRegion">您未设置地域属性</p>
              <p v-else>您的地域属性为 {{ userRegion }}，无权访问 Sichuan 订单</p>
            </div>
          </el-alert>

          <el-button
            :type="canAccessOrders ? 'primary' : 'info'"
            :disabled="!canAccessOrders"
            @click="goToOrders"
            style="width: 100%"
          >
            {{ canAccessOrders ? '进入订单管理' : '权限不足' }}
          </el-button>
        </el-card>
      </el-col>

      <!-- 演示2：数据级别限制 -->
      <el-col :span="12">
        <el-card class="demo-card">
          <div slot="header" class="clearfix">
            <span>🔒 演示2：数据级别限制</span>
            <span style="float: right; font-size: 12px;">财务报表（Internal）</span>
          </div>

          <p style="margin-bottom: 10px;">
            <strong>需求：</strong> 仅 Internal 级别用户可访问财务报表
          </p>

          <el-alert
            v-if="canAccessFinance"
            title="✅ 权限检查通过"
            type="success"
            :closable="false"
            style="margin-bottom: 15px"
          >
            您具有 Internal 数据级别属性，可以访问财务报表
          </el-alert>

          <el-alert
            v-else
            title="❌ 权限检查未通过"
            type="error"
            :closable="false"
            style="margin-bottom: 15px"
          >
            <div>
              <p v-if="!userDataLevel">您未设置数据级别属性</p>
              <p v-else>您的数据级别为 {{ userDataLevel }}，无权访问 Internal 数据</p>
            </div>
          </el-alert>

          <el-button
            :type="canAccessFinance ? 'primary' : 'info'"
            :disabled="!canAccessFinance"
            @click="goToFinance"
            style="width: 100%"
          >
            {{ canAccessFinance ? '进入财务报表' : '权限不足' }}
          </el-button>
        </el-card>
      </el-col>
    </el-row>

    <el-divider>如何测试</el-divider>

    <el-steps :active="activeStep" finish-status="success" style="margin-top: 20px;">
      <el-step title="登录系统" />
      <el-step title="进入用户管理" />
      <el-step title="编辑用户属性" />
      <el-step title="添加属性值" />
      <el-step title="保存并登出" />
      <el-step title="重新登录测试" />
    </el-steps>

    <div class="test-guide">
      <h4>📝 详细步骤:</h4>
      <ol>
        <li>
          <strong>编辑用户属性</strong>
          <ul>
            <li>点击左侧菜单 "用户管理"</li>
            <li>找到当前用户，点击"编辑属性"</li>
            <li>添加属性: <code>region</code> = <code>Sichuan</code></li>
            <li>添加属性: <code>data_level</code> = <code>Internal</code></li>
            <li>点击"确定"保存</li>
          </ul>
        </li>
        <li>
          <strong>退出登录</strong>
          <ul>
            <li>点击右上角用户菜单</li>
            <li>选择"登出"</li>
          </ul>
        </li>
        <li>
          <strong>重新登录</strong>
          <ul>
            <li>使用相同账号重新登录</li>
            <li>系统会自动加载新的属性</li>
          </ul>
        </li>
        <li>
          <strong>查看权限更新</strong>
          <ul>
            <li>返回此页面刷新</li>
            <li>应该能看到权限已开放</li>
            <li>可以进入对应的演示页面</li>
          </ul>
        </li>
      </ol>
    </div>

    <el-divider>权限匹配规则</el-divider>

    <el-collapse>
      <el-collapse-item title="🔍 查看详细规则" name="1">
        <el-table :data="permissionRules" border style="width: 100%">
          <el-table-column prop="scenario" label="场景" width="200" />
          <el-table-column prop="userAttrs" label="用户属性" width="250" />
          <el-table-column prop="required" label="资源要求" width="200" />
          <el-table-column prop="result" label="结果" width="100">
            <template slot-scope="scope">
              <el-tag :type="scope.row.result === '✅ 允许' ? 'success' : 'danger'">
                {{ scope.row.result }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { hasPermission } from '@/utils/checkPermission'

export default {
  name: 'AttributeDemo',
  data() {
    return {
      activeStep: 4,
      permissionRules: [
        {
          scenario: '用户具有所有需要的属性',
          userAttrs: '{ region: "Sichuan", level: "admin" }',
          required: '{ region: "Sichuan" }',
          result: '✅ 允许'
        },
        {
          scenario: '用户属性过多（包含所需属性）',
          userAttrs: '{ region: "Sichuan", level: "admin", dept: "IT" }',
          required: '{ region: "Sichuan" }',
          result: '✅ 允许'
        },
        {
          scenario: '用户缺少必需属性',
          userAttrs: '{ level: "admin" }',
          required: '{ region: "Sichuan" }',
          result: '❌ 拒绝'
        },
        {
          scenario: '用户属性值不匹配',
          userAttrs: '{ region: "Shanghai" }',
          required: '{ region: "Sichuan" }',
          result: '❌ 拒绝'
        },
        {
          scenario: '资源无权限要求',
          userAttrs: '{}',
          required: '无',
          result: '✅ 允许'
        }
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
    userDataLevel() {
      return this.dynamicAttributes?.data_level || null
    },
    canAccessOrders() {
      const user = {
        userType: this.userType,
        dynamicAttributes: this.dynamicAttributes || {}
      }
      return hasPermission(user, { requiredAttributes: { region: 'Sichuan' } })
    },
    canAccessFinance() {
      const user = {
        userType: this.userType,
        dynamicAttributes: this.dynamicAttributes || {}
      }
      return hasPermission(user, { requiredAttributes: { data_level: 'Internal' } })
    }
  },
  methods: {
    goToOrders() {
      this.$router.push('/attribute-demo/orders-region')
    },
    goToFinance() {
      this.$router.push('/attribute-demo/financial-reports')
    }
  }
}
</script>

<style scoped>
.user-info-item {
  text-align: center;
}

.user-info-item .label {
  color: #909399;
  font-size: 12px;
  margin-bottom: 5px;
}

.user-info-item .value {
  font-weight: bold;
  margin: 0;
}

.demo-card {
  height: 100%;
}

.test-guide {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 4px;
  margin-top: 20px;
}

.test-guide ol {
  margin-left: 20px;
}

.test-guide li {
  margin-bottom: 15px;
  line-height: 1.6;
}

.test-guide code {
  background: #eff2f7;
  padding: 2px 6px;
  border-radius: 2px;
  font-family: 'Courier New', monospace;
  color: #d63031;
}
</style>

