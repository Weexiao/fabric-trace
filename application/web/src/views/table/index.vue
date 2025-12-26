<template>
  <div class="app-container">
    <el-table
      :data="list"
      :element-loading-text="$t('table.loading')"
      border
      fit
      highlight-current-row
    >
      <el-table-column align="center" :label="$t('table.id')" width="95">
        <template v-slot="scope">
          {{ scope.$index }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.title')">
        <template v-slot="scope">
          {{ safeGet(scope.row, 'title', '') }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.author')" width="110" align="center">
        <template v-slot="scope">
          <span>{{ safeGet(scope.row, 'author', '') }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.pageviews')" width="110" align="center">
        <template v-slot="scope">
          {{ safeGet(scope.row, 'pageviews', 0) }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" :label="$t('table.status')" width="110" align="center">
        <template v-slot="scope">
          <el-tag :type="(safeGet(scope.row, 'status', '') | statusFilter)">{{ safeGet(scope.row, 'status', '') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" :label="$t('table.displayTime')" width="200">
        <template v-slot="scope">
          <i class="el-icon-time" />
          <span>{{ safeGet(scope.row, 'display_time', '') }}</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { getList } from '@/api/table'
import { safeGet as safeGetUtil } from '@/utils'
import { apiWrap } from '@/utils/error'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'gray',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      list: [],
      listLoading: false,
      errorMessage: null
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    safeGet: safeGetUtil,
    fetchData() {
      apiWrap(this, () => getList(), (res) => {
        this.list = this.safeGet(res, 'data.items', [])
      }, this.$t('table.fetchFailed'))
    }
  }
}
</script>
