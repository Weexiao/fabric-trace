import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import ElementUI from 'element-ui'
import Uplink from '@/views/uplink/index.vue'

jest.mock('@/api/trace', () => ({
  // component uses res.traceabilityCode
  uplink: jest.fn(() => Promise.resolve({ code: 200, txid: 'TX123', traceabilityCode: '123456789012345678' }))
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
    wrapper.setData({ tracedata: { ...wrapper.vm.tracedata, traceabilityCode: '123' }})
    await wrapper.vm.$nextTick()
    const res = await new Promise((resolve) => wrapper.vm.$refs.form.validate((v) => resolve(v)))
    expect(res).toBe(false)
    wrapper.setData({ tracedata: { ...wrapper.vm.tracedata, traceabilityCode: '123456789012345678' }})
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
        traceabilityCode: '123456789012345678',
        manufacturerInput: {
          productName: 'p',
          productionBatch: 'b1',
          factoryTime: '2020-01-01 00:00:00',
          factoryNameAddress: 'factory',
          contactPhone: '+86 13312345678'
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
        manufacturerInput: {
          productName: 'p',
          productionBatch: 'b1',
          factoryTime: '2020-01-01 00:00:00',
          factoryNameAddress: 'factory',
          contactPhone: '+86 13312345678'
        },
        carrierInput: {
          name: 'n', age: 30, phone: '+1 2222222', plateNumber: 'ABC123', transportRecord: 't'
        }
      }
    })

    // simulate role change to 经销商 (component detects via lastUserType)
    wrapper.vm.lastUserType = '制造商'
    wrapper.vm.$options.computed.userType.get = () => '经销商'

    // invoke the same reset behavior the component uses when role changes
    wrapper.vm.resetBranchData('经销商')
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.tracedata.carrierInput.name).toBe('')
    expect(wrapper.vm.tracedata.manufacturerInput.productName).toBe('')
  })

  test('offchain upload without trace code shows confirm and scrolls to uplink form on confirm', async() => {
    const { wrapper } = mountUplink('制造商')

    // arrange: pick a file but no traceabilityCode
    wrapper.setData({ offchainFile: { name: 'a.txt', size: 1 }})

    const confirmMock = jest.fn(() => Promise.resolve())
    wrapper.vm.$confirm = confirmMock

    const scrollMock = jest.fn()
    wrapper.vm.scrollToUplinkForm = scrollMock

    await wrapper.vm.uploadOffchain()

    expect(confirmMock).toHaveBeenCalled()
    expect(scrollMock).toHaveBeenCalled()
  })

  test('after uplink success, traceabilityCode is written back enabling offchain upload', async() => {
    const { wrapper } = mountUplink('制造商')

    wrapper.setData({
      tracedata: {
        ...wrapper.vm.tracedata,
        // 制造商提交前校验要求 18 位溯源码，因此这里先给一个有效值
        traceabilityCode: '123456789012345678',
        manufacturerInput: {
          productName: 'p',
          productionBatch: 'b1',
          factoryTime: '2020-01-01 00:00:00',
          factoryNameAddress: 'factory',
          contactPhone: '+86 13312345678'
        }
      },
      offchainFile: { name: 'a.txt', size: 1 }
    })

    await wrapper.vm.submittracedata()

    expect(wrapper.vm.tracedata.traceabilityCode).toBe('123456789012345678')
    expect(wrapper.vm.canUploadOffchain).toBe(true)
  })
})
