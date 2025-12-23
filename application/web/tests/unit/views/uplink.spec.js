import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import ElementUI from 'element-ui'
import Uplink from '@/views/uplink/index.vue'

jest.mock('@/api/trace', () => ({
  uplink: jest.fn(() => Promise.resolve({ code: 200, txid: 'TX123', traceability_code: '123456789012345678' }))
}))

const localVue = createLocalVue()
localVue.use(Vuex)
localVue.use(ElementUI)

function makeStore(userType) {
  return new Vuex.Store({
    getters: {
      name: () => 'tester',
      userType: () => userType
    }
  })
}

function mountUplink(userType = '制造商') {
  const store = makeStore(userType)
  const $router = { push: jest.fn() }
  const wrapper = shallowMount(Uplink, {
    localVue,
    store,
    mocks: {
      $t: (k) => k,
      $router,
      $message: { success: jest.fn(), error: jest.fn() },
      $loading: () => ({ close: jest.fn() })
    },
    stubs: {
      'success-dialog': true,
      'el-upload': true,
      'el-date-picker': true,
      'el-cascader': true
    }
  })
  // Patch shortcuts method to stable no-op
  wrapper.vm.commonShortcuts = () => []
  wrapper.vm.$forceUpdate()
  return { wrapper, store, $router }
}

describe('Uplink view', () => {
  test('trace code validation requires 18 digits for manufacturer', async() => {
    const { wrapper } = mountUplink('制造商')
    wrapper.setData({ tracedata: { ...wrapper.vm.tracedata, traceability_code: '123' }})
    await wrapper.vm.$nextTick()
    const res = await new Promise((resolve) => wrapper.vm.$refs.form.validate((v) => resolve(v)))
    expect(res).toBe(false)
    wrapper.setData({ tracedata: { ...wrapper.vm.tracedata, traceability_code: '123456789012345678' }})
    await wrapper.vm.$nextTick()
    const res2 = await new Promise((resolve) => wrapper.vm.$refs.form.validate((v) => resolve(v)))
    expect(res2).toBe(false) // still invalid until required fields filled
  })

  test('phone validator accepts +86 mobile and ext formats', () => {
    const { wrapper } = mountUplink('制造商')
    expect(wrapper.vm.isValidPhone('+86 13312345678')).toBe(true)
    expect(wrapper.vm.isValidPhone('010-88886666 转123')).toBe(true)
    expect(wrapper.vm.isValidPhone('+123')).toBe(false)
    expect(wrapper.vm.isValidPhone('abc')).toBe(false)
  })

  test('submit success opens success dialog and sets successInfo', async() => {
    const { wrapper } = mountUplink('制造商')
    // fill minimal required
    wrapper.setData({
      tracedata: {
        ...wrapper.vm.tracedata,
        traceability_code: '123456789012345678',
        Factory_input: {
          Fac_productName: 'p',
          Fac_productionbatch: 'b1',
          Fac_productionTime: '2020-01-01 00:00:00',
          Fac_factoryName: 'factory',
          Fac_contactNumber: '+86 13312345678'
        }
      }
    })
    await wrapper.vm.submittracedata()
    expect(wrapper.vm.successDialogVisible).toBe(true)
    expect(wrapper.vm.successInfo.code).toBe('123456789012345678')
    expect(wrapper.vm.successInfo.txid).toBe('TX123')
  })

  test('role switch resets unrelated branch data', async() => {
    const { wrapper } = mountUplink('制造商')
    wrapper.setData({
      tracedata: {
        ...wrapper.vm.tracedata,
        Factory_input: {
          Fac_productName: 'p',
          Fac_productionbatch: 'b1',
          Fac_productionTime: '2020-01-01 00:00:00',
          Fac_factoryName: 'factory',
          Fac_contactNumber: '+86 13312345678'
        },
        Driver_input: {
          Dr_name: 'n', Dr_age: 30, Dr_phone: '+1 2222222', Dr_carNumber: 'ABC123', Dr_transport: 't'
        }
      }
    })
    // change role to 经销商
    wrapper.vm.$options.computed.userType.get = () => '经销商'
    wrapper.vm.$forceUpdate()
    await wrapper.vm.$nextTick()
    // trigger watch
    wrapper.vm.userType
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.tracedata.Driver_input.Dr_name).toBe('')
    expect(wrapper.vm.tracedata.Factory_input.Fac_productName).toBe('')
  })
})
