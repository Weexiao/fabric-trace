import Vue from 'vue'
import VueI18n from 'vue-i18n'

Vue.use(VueI18n)

const messages = {
  zh: {
    app: {
      title: '价值网溯源系统'
    },
    nav: {
      logout: '登出'
    },
    common: {
      currentUser: '当前用户',
      userRole: '用户角色',
      submit: '提 交',
      reset: '重 置',
      edit: '更 改',
      noPermission: '零售商没有权限录入！请使用溯源功能!',
      uploading: '数据上链中...',
      uploadChooseImg: '选择图片',
      deleteImage: '删除图片',
      errorSubmitFlow: '提交流程异常，请重试'
    },
    actions: {
      copy: '复制',
      copied: '已复制',
      viewTrace: '查看详情',
      viewTx: '查看交易',
      continue: '继续录入'
    },
    dialog: {
      uplinkSuccess: '上链成功',
      traceCode: '溯源码',
      txid: '交易ID'
    },
    form: {
      traceCode: '溯源码',
      inputTraceCode: '请输入18位数字溯源码',
      traceCodeHelp: '原料供应商无需填写，系统会自动生成溯源码；制造商/物流/经销商需填写与该批次关联的18位溯源码。支持从链接自动填充。',
      farmer: {
        fruitName: '原料名称',
        inputFruitName: '请输入原料名称',
        origin: '原料产地',
        pickTime: '原料到货时间',
        plantTime: '原料生产时间',
        supplier: '原料供应商名称',
        inputSupplier: '请输入原料供应商名称',
        phone: '供应商电话',
        inputPhone: '请输入供应商电话'
      },
      factory: {
        productName: '产品名称',
        inputProductName: '请输入产品名称',
        batch: '生产批次',
        inputBatch: '请输入生产批次',
        prodTime: '生产时间',
        factoryName: '制造商名称',
        inputFactoryName: '请输入制造商名称与厂址',
        phone: '制造商电话',
        inputPhone: '请输入电话号码'
      },
      driver: {
        name: '姓名',
        inputName: '请输入姓名',
        age: '年龄',
        inputAge: '请输入年龄',
        phone: '联系电话',
        inputPhone: '请输入联系电话',
        carNumber: '车牌号',
        inputCarNumber: '请输入车牌号',
        transport: '运输记录',
        inputTransport: '请输入运输记录'
      },
      shop: {
        storeTime: '存入时间',
        sellTime: '销售时间',
        name: '经销商名称',
        inputName: '请输入经销商名称',
        address: '经销商位置',
        phone: '经销商电话',
        inputPhone: '请输入经销商电话'
      },
      regionPlaceholder: '请选择省/市/区',
      datetimePlaceholder: '选择日期时间'
    },
    validate: {
      pleaseFixErrors: '请先修正表单校验错误',
      pleaseEnter: '请输入',
      pleaseSelect: '请选择',
      codeMustBe18Digits: '溯源码必须为18位数字',
      notInFuture: '不能晚于当前时间',
      sellNotEarlierThanStore: '销售时间不能早于存入时间',
      invalidPhone: '请输入有效的电话号码',
      invalidCar: '车牌号格式不正确',
      batchInvalid: '批次格式不正确(1-32位字母/数字/_-)'
    },
    result: {
      success: '上链成功，交易ID：{txid}\n溯源码：{code}',
      fail: '上链失败',
      exception: '上链异常{detail}'
    },
    login: {
      title: '基于区块链的价值网溯源系统',
      username: '账号',
      password: '密码',
      usernamePlaceholder: '请输入账号',
      passwordPlaceholder: '请输入密码',
      passwordAgainPlaceholder: '请再次输入密码',
      register: '注册',
      login: '登录',
      back: '返回',
      submitRegister: '提交注册',
      selectRole: '请选择角色',
      rules: {
        usernameRequired: '请输入账号',
        usernameMin: '账号至少3个字符',
        passwordRequired: '请输入密码',
        passwordMin: '密码至少6个字符',
        passwordComplex: '密码需包含字母和数字',
        passwordAgainRequired: '请再次输入密码',
        passwordNotEqual: '两次密码不一致'
      },
      messages: {
        loginFailed: '登录失败',
        registerFailed: '注册失败',
        registering: '注册中...'
      }
    },
    build: {
      title: '5分钟构建任意溯源系统',
      p1: '请对比字段填写，生成个性化的溯源系统。激活码仅可使用一次，',
      p2: '提交前请认真对比，生成后请尽快下载并备份，源码在服务器保留一周后删除。',
      tutorial: 'B站：使用教程',
      buyCode: '购买激活码',
      deletePage: '删除此页面',
      deleteStep1: '1. 在fabric-trace/application/web目录下，运行./rmbuildsyspage.sh',
      deleteStep2: '2.重新启动前端，npm run dev',
      confirm: '确 定',
      supportTitle: '开发不易，感谢支持！',
      price: '激活码售价：269元',
      contact: '请加QQ群776873343联系群主购买 （补差价可包搭建）',
      discount: '购买课程后可99元购买激活码',
      buildSuccess: '构建成功！',
      download: '下载地址：',
      building: '构建中，稍等1分钟',
      startBuild: '开始构建',
      howToDelete: '如何删除此页面？',
      submitFailed: '提交失败，请重试！',
      buildFailed: '构建失败：{msg}'
    }
  },
  en: {
    app: {
      title: 'Value Network Traceability System'
    },
    nav: {
      logout: 'Logout'
    },
    common: {
      currentUser: 'Current user',
      userRole: 'User role',
      submit: 'Submit',
      reset: 'Reset',
      edit: 'Edit',
      noPermission: 'Retailer has no permission to input, please use tracing feature!',
      uploading: 'Submitting to blockchain...',
      uploadChooseImg: 'Choose Image',
      deleteImage: 'Remove Image',
      errorSubmitFlow: 'Submission flow error, please retry'
    },
    actions: {
      copy: 'Copy',
      copied: 'Copied',
      viewTrace: 'View Details',
      viewTx: 'View Transaction',
      continue: 'Continue Input'
    },
    dialog: {
      uplinkSuccess: 'Uplink Succeeded',
      traceCode: 'Trace Code',
      txid: 'TxID'
    },
    form: {
      traceCode: 'Trace Code',
      inputTraceCode: 'Please enter 18-digit numeric trace code',
      traceCodeHelp: 'Supplier doesn\'t need to input; system generates code automatically. Manufacturer/Carrier/Dealer must input the 18-digit code. Auto-fill via link is supported.',
      farmer: {
        fruitName: 'Ingredient Name',
        inputFruitName: 'Enter ingredient name',
        origin: 'Origin',
        pickTime: 'Arrival Time',
        plantTime: 'Production Time',
        supplier: 'Supplier Name',
        inputSupplier: 'Enter supplier name',
        phone: 'Supplier Phone',
        inputPhone: 'Enter supplier phone'
      },
      factory: {
        productName: 'Product Name',
        inputProductName: 'Enter product name',
        batch: 'Batch',
        inputBatch: 'Enter batch',
        prodTime: 'Production Time',
        factoryName: 'Manufacturer',
        inputFactoryName: 'Enter manufacturer and address',
        phone: 'Phone',
        inputPhone: 'Enter phone number'
      },
      driver: {
        name: 'Name',
        inputName: 'Enter name',
        age: 'Age',
        inputAge: 'Enter age',
        phone: 'Phone',
        inputPhone: 'Enter phone',
        carNumber: 'Plate Number',
        inputCarNumber: 'Enter plate number',
        transport: 'Transport Record',
        inputTransport: 'Enter transport record'
      },
      shop: {
        storeTime: 'Store Time',
        sellTime: 'Sell Time',
        name: 'Dealer Name',
        inputName: 'Enter dealer name',
        address: 'Dealer Address',
        phone: 'Dealer Phone',
        inputPhone: 'Enter dealer phone'
      },
      regionPlaceholder: 'Select province/city/district',
      datetimePlaceholder: 'Select date time'
    },
    validate: {
      pleaseFixErrors: 'Please fix form errors first',
      pleaseEnter: 'Please enter',
      pleaseSelect: 'Please select',
      codeMustBe18Digits: 'Trace code must be 18 digits',
      notInFuture: 'Cannot be later than now',
      sellNotEarlierThanStore: 'Sell time cannot be earlier than store time',
      invalidPhone: 'Please enter a valid phone number',
      invalidCar: 'Invalid plate number format',
      batchInvalid: 'Invalid batch format (1-32 letters/numbers/_-)'
    },
    result: {
      success: 'Success, Tx: {txid}\nTrace code: {code}',
      fail: 'Failed',
      exception: 'Exception{detail}'
    },
    login: {
      title: 'Blockchain-based Value Network Traceability',
      username: 'Username',
      password: 'Password',
      usernamePlaceholder: 'Enter username',
      passwordPlaceholder: 'Enter password',
      passwordAgainPlaceholder: 'Re-enter password',
      register: 'Register',
      login: 'Login',
      back: 'Back',
      submitRegister: 'Submit Registration',
      selectRole: 'Select role',
      rules: {
        usernameRequired: 'Please enter username',
        usernameMin: 'At least 3 characters',
        passwordRequired: 'Please enter password',
        passwordMin: 'At least 6 characters',
        passwordComplex: 'Password must include letters and numbers',
        passwordAgainRequired: 'Please re-enter password',
        passwordNotEqual: 'Passwords do not match'
      },
      messages: {
        loginFailed: 'Login failed',
        registerFailed: 'Registration failed',
        registering: 'Registering...'
      }
    },
    build: {
      title: 'Build any traceability system in 5 minutes',
      p1: 'Fill the fields carefully. Activation code can be used once only.',
      p2: 'Review before submit; please download and back up after generated. Code kept for one week on server.',
      tutorial: 'Bilibili: Tutorial',
      buyCode: 'Buy activation code',
      deletePage: 'Delete this page',
      deleteStep1: '1. In fabric-trace/application/web, run ./rmbuildsyspage.sh',
      deleteStep2: '2. Restart frontend: npm run dev',
      confirm: 'Confirm',
      supportTitle: 'Support the project, thanks!',
      price: 'Activation code: ¥269',
      contact: 'Join QQ group 776873343 to contact admin (deployment available with extra fee)',
      discount: 'Buy the course to get the code at ¥99',
      buildSuccess: 'Build succeeded!',
      download: 'Download URL:',
      building: 'Building, about 1 minute...',
      startBuild: 'Start Build',
      howToDelete: 'How to delete this page?',
      submitFailed: 'Submit failed, please retry!',
      buildFailed: 'Build failed: {msg}'
    }
  }
}

const i18n = new VueI18n({
  locale: 'zh',
  fallbackLocale: 'zh',
  messages
})

export default i18n
