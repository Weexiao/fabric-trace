<template>
  <el-dialog :title="$t('dialog.uplinkSuccess')" :visible.sync="visibleProxy" width="480px">
    <div style="line-height: 1.8;">
      <div v-if="code">
        <strong>{{ $t('dialog.traceCode') }}:</strong>
        <span>{{ code }}</span>
        <el-button type="text" size="mini" @click="copyText(code)">{{ $t('actions.copy') }}</el-button>
      </div>
      <div v-if="txid">
        <strong>{{ $t('dialog.txid') }}:</strong>
        <span>{{ txid }}</span>
        <el-button type="text" size="mini" @click="copyText(txid)">{{ $t('actions.copy') }}</el-button>
      </div>
    </div>
    <span slot="footer" class="dialog-footer">
      <el-button @click="onContinue">{{ $t('actions.continue') }}</el-button>
      <el-button type="primary" plain @click="onViewTrace">{{ $t('actions.viewTrace') }}</el-button>
      <el-button v-if="canOpenTx" type="success" plain @click="onViewTx">{{ $t('actions.viewTx') }}</el-button>
    </span>
  </el-dialog>
</template>

<script>
export default {
  name: 'SuccessDialog',
  props: {
    value: { type: Boolean, default: false },
    code: { type: String, default: '' },
    txid: { type: String, default: '' },
    txExplorer: { type: String, default: '' }
  },
  computed: {
    visibleProxy: {
      get() { return this.value },
      set(v) { this.$emit('input', v) }
    },
    canOpenTx() {
      return !!(this.txExplorer && this.txid)
    }
  },
  methods: {
    copyText(text) {
      const val = String(text || '')
      if (!val) return
      if (navigator && navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(val).then(() => {
          this.$message.success(this.$t('actions.copied'))
        }).catch(() => {
          this._legacyCopy(val)
        })
      } else {
        this._legacyCopy(val)
      }
    },
    _legacyCopy(val) {
      const ta = document.createElement('textarea')
      ta.value = val
      ta.style.position = 'fixed'
      ta.style.opacity = '0'
      document.body.appendChild(ta)
      ta.focus()
      ta.select()
      try { document.execCommand('copy'); this.$message.success(this.$t('actions.copied')) } catch (e) { this.$message.error('复制失败') }
      document.body.removeChild(ta)
    },
    onViewTrace() { this.$emit('view-trace') },
    onViewTx() { if (this.canOpenTx) this.$emit('view-tx', this.txid) },
    onContinue() { this.$emit('continue') }
  }
}
</script>

<style scoped>
/* no-op */
</style>

