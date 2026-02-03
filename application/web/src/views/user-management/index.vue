<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button type="primary" @click="handleAddUser">添加用户</el-button>
    </div>

    <el-table
      :data="userList"
      border
      fit
      highlight-current-row
      style="width: 100%"
      v-loading="loading"
    >
      <el-table-column align="center" label="用户ID" width="150">
        <template v-slot="scope">
          {{ scope.row.user_id }}
        </template>
      </el-table-column>

      <el-table-column label="用户名" width="150">
        <template v-slot="scope">
          {{ scope.row.username }}
        </template>
      </el-table-column>

      <el-table-column label="用户类型" width="120" align="center">
        <template v-slot="scope">
          <el-tag>{{ scope.row.user_type || scope.row.userType }}</el-tag>
        </template>
      </el-table-column>

      <!-- 约束属性列 -->
      <el-table-column label="约束属性" min-width="200">
        <template v-slot="scope">
          <div v-if="scope.row.dynamicAttributes && Object.keys(scope.row.dynamicAttributes).length > 0">
            <div v-for="(value, key) in scope.row.dynamicAttributes" :key="key" class="attribute-item">
              <span class="attr-key">{{ key }}:</span>
              <span class="attr-value">{{ value }}</span>
            </div>
          </div>
          <span v-else class="no-attr">无约束属性</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="操作" width="200">
        <template v-slot="scope">
          <el-button size="small" type="primary" @click="handleEditAttributes(scope.row)">编辑属性</el-button>
          <el-button size="small" type="danger" @click="handleDeleteUser(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 编辑约束属性弹窗 -->
    <el-dialog
      title="编辑约束属性"
      :visible.sync="dialogVisible"
      width="500px"
      @close="resetForm"
    >
      <div v-if="currentUser">
        <p class="user-info">用户: <strong>{{ currentUser.username }}</strong></p>

        <div class="attributes-container">
          <div v-for="(attr, index) in attributeInputs" :key="index" class="attribute-input-group">
            <el-input
              v-model="attr.key"
              placeholder="属性名 (如: region, data_level)"
              class="attr-key-input"
            />
            <span class="colon">:</span>
            <el-input
              v-model="attr.value"
              placeholder="属性值 (如: Sichuan, Internal)"
              class="attr-value-input"
            />
            <el-button
              icon="el-icon-delete"
              circle
              type="danger"
              size="small"
              @click="removeAttribute(index)"
              v-if="attributeInputs.length > 1"
            />
          </div>
        </div>

        <el-button
          type="dashed"
          icon="el-icon-plus"
          @click="addAttributeInput"
          style="width: 100%; margin-top: 10px"
        >
          添加属性
        </el-button>
      </div>

      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="saveAttributes">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { updateUserDynamicAttributes, getAllUsers } from '@/api/user'
import { Message } from 'element-ui'

export default {
  name: 'UserManagement',
  data() {
    return {
      userList: [],
      loading: false,
      dialogVisible: false,
      currentUser: null,
      attributeInputs: [{ key: '', value: '' }]
    }
  },
  created() {
    this.fetchUserList()
  },
  methods: {
    /**
     * 获取用户列表
     */
    fetchUserList() {
      this.loading = true
      getAllUsers().then(response => {
        if (response.code === 200) {
          this.userList = response.data.map(user => ({
            ...user,
            dynamicAttributes: typeof user.dynamic_attributes === 'string'
              ? JSON.parse(user.dynamic_attributes || '{}')
              : (user.dynamic_attributes || {})
          }))
        } else {
          Message.error(response.message || '获取用户列表失败')
        }
        this.loading = false
      }).catch(error => {
        Message.error('获取用户列表失败: ' + error.message)
        this.loading = false
      })
    },

    /**
     * 返回首页
     */
    goBack() {
      this.$router.push('/')
    },

    /**
     * 打开编辑属性弹窗
     */
    handleEditAttributes(user) {
      this.currentUser = JSON.parse(JSON.stringify(user))
      // 转换dynamicAttributes为输入数组
      if (this.currentUser.dynamicAttributes && Object.keys(this.currentUser.dynamicAttributes).length > 0) {
        this.attributeInputs = Object.entries(this.currentUser.dynamicAttributes).map(([key, value]) => ({
          key,
          value
        }))
      } else {
        this.attributeInputs = [{ key: '', value: '' }]
      }
      this.dialogVisible = true
    },

    /**
     * 添加属性输入框
     */
    addAttributeInput() {
      this.attributeInputs.push({ key: '', value: '' })
    },

    /**
     * 移除属性输入框
     */
    removeAttribute(index) {
      this.attributeInputs.splice(index, 1)
    },

    /**
     * 保存约束属性
     */
    saveAttributes() {
      // 构建属性对象，过滤空的键值对
      const attributes = {}
      for (const attr of this.attributeInputs) {
        if (attr.key && attr.value) {
          attributes[attr.key] = attr.value
        }
      }

      const formData = new FormData()
      formData.append('user_id', this.currentUser.user_id)
      formData.append('dynamic_attributes', JSON.stringify(attributes))

      updateUserDynamicAttributes(formData).then(response => {
        if (response.code === 200) {
          Message.success('属性更新成功')
          // 更新本地列表
          const index = this.userList.findIndex(u => u.user_id === this.currentUser.user_id)
          if (index !== -1) {
            this.userList[index].dynamicAttributes = attributes
          }
          this.dialogVisible = false
        } else {
          Message.error(response.message || '更新失败')
        }
      }).catch(error => {
        Message.error('更新失败: ' + error.message)
      })
    },

    /**
     * 重置表单
     */
    resetForm() {
      this.currentUser = null
      this.attributeInputs = [{ key: '', value: '' }]
    },

    /**
     * 添加用户
     */
    handleAddUser() {
      Message.info('添加用户功能待实现')
    },

    /**
     * 删除用户
     */
    handleDeleteUser(user) {
      this.$confirm(`确定删除用户 ${user.username} 吗?`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        Message.success('用户删除成功')
        this.fetchUserList()
      }).catch(() => {
        Message.info('已取消删除')
      })
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.permission-denied {
  padding: 40px 20px;
  text-align: center;
}

.admin-content {
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

.filter-container {
  margin-bottom: 20px;
}

.attribute-item {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
  font-size: 12px;
}

.attr-key {
  font-weight: bold;
  color: #606266;
  margin-right: 5px;
}

.attr-value {
  color: #909399;
  background: #f5f7fa;
  padding: 2px 8px;
  border-radius: 3px;
}

.no-attr {
  color: #909399;
  font-style: italic;
}

.user-info {
  margin-bottom: 15px;
  font-size: 14px;
  color: #606266;
}

.attributes-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 15px;
  background: #f9fafc;
}

.attribute-input-group {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  gap: 8px;
}

.attribute-input-group:last-child {
  margin-bottom: 0;
}

.attr-key-input {
  flex: 1;
}

.colon {
  margin: 0 5px;
  color: #606266;
}

.attr-value-input {
  flex: 1.5;
}

.dialog-footer {
  text-align: right;
}
</style>

