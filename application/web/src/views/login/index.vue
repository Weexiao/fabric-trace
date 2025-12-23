<template>
  <div class="login-container">
    <el-form ref="loginForm" :model="isLoginPage ? loginForm : registerForm" :rules="isLoginPage ? loginRules : registerRules" class="login-form" autocomplete="on" label-position="left">

      <div class="title-container">
        <h3 class="title" style="color: white;">{{ $t('login.title') }}</h3>
      </div>
      <div v-show="isLoginPage">
        <el-form-item prop="username">
          <span class="svg-container">
            <svg-icon icon-class="user" />
          </span>
          <el-input
            ref="username"
            v-model.trim="loginForm.username"
            v-sanitize="{ max: LENGTHS.login.username }"
            :placeholder="$t('login.usernamePlaceholder')"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="on"
          />
        </el-form-item>

        <el-form-item prop="password">
          <span class="svg-container">
            <svg-icon icon-class="password" />
          </span>
          <el-input
            :key="passwordTypes.login"
            ref="loginPassword"
            v-model="loginForm.password"
            v-sanitize="{ max: LENGTHS.login.password }"
            :type="passwordTypes.login"
            :placeholder="$t('login.passwordPlaceholder')"
            name="password"
            tabindex="2"
            autocomplete="on"
            @keyup.enter.native="handleLogin"
          />
          <span class="show-pwd" @click="showPwd('login')">
            <svg-icon :icon-class="passwordTypes.login === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>
        <el-button :loading="loading" :disabled="loading" type="info" style="width:20%;margin-bottom:30px;" @click="handleRegister">{{ $t('login.register') }}</el-button>
        <el-button :loading="loading" :disabled="loading" type="primary" style="width:30%;margin-bottom:30px; float: right" @click.prevent="handleLogin">{{ $t('login.login') }}</el-button>
      </div>
      <div v-show="!isLoginPage">
        <el-form-item prop="username">
          <span class="svg-container">
            <svg-icon icon-class="user" />
          </span>
          <el-input
            v-model.trim="registerForm.username"
            v-sanitize="{ max: LENGTHS.login.username }"
            :placeholder="$t('login.usernamePlaceholder')"
            name="username"
            type="text"
            autocomplete="on"
          />
        </el-form-item>
        <el-form-item prop="password">
          <span class="svg-container">
            <svg-icon icon-class="password" />
          </span>
          <el-input
            :key="passwordTypes.register1"
            ref="registerPassword1"
            v-model="registerForm.password"
            v-sanitize="{ max: LENGTHS.login.password }"
            :type="passwordTypes.register1"
            :placeholder="$t('login.passwordPlaceholder')"
            name="password"
            autocomplete="on"
            style="color: white !important;"
          />
          <span class="show-pwd" @click="showPwd('register1')">
            <svg-icon :icon-class="passwordTypes.register1 === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>
        <el-form-item prop="password2">
          <span class="svg-container">
            <svg-icon icon-class="password" />
          </span>
          <el-input
            :key="passwordTypes.register2"
            ref="registerPassword2"
            v-model="registerForm.password2"
            v-sanitize="{ max: LENGTHS.login.password }"
            :placeholder="$t('login.passwordAgainPlaceholder')"
            name="password"
            autocomplete="on"
            :type="passwordTypes.register2"
          />
          <span class="show-pwd" @click="showPwd('register2')">
            <svg-icon :icon-class="passwordTypes.register2 === 'password' ? 'eye' : 'eye-open'" />
          </span>
        </el-form-item>
        <el-form-item prop="userType" style="width: 200px">
          <el-select v-model="registerForm.userType" :placeholder="$t('login.selectRole')">
            <el-option
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-button :loading="loading" :disabled="loading" type="info" style="width:20%;margin-bottom:30px;" @click="handleRegister">{{ $t('login.back') }}</el-button>
        <el-button :loading="loading" :disabled="loading" type="primary" style="width:30%;margin-bottom:30px; float: right" @click.prevent="submitRegister">{{ $t('login.submitRegister') }}</el-button>
      </div>
      <!-- <div class="tips">
        <span style="margin-right:20px;">提示：可以放一些提示</span>
      </div> -->

    </el-form>
  </div>
</template>

<script>
import { LENGTHS } from '@/utils/limits'
import { sanitize } from '@/utils/sanitize'

export default {
  name: 'Login',
  data() {
    // 仅字母/数字/下划线
    const usernamePattern = /^[a-zA-Z0-9_]+$/
    // 注册密码复杂度（>=8 且包含字母和数字）
    const validatePasswordComplex = (rule, value, callback) => {
      if (!value) return callback(new Error(this.$t('login.rules.passwordRequired')))
      if (String(value).length < 8) return callback(new Error(this.$t('login.rules.passwordMin')))
      if (!/[A-Za-z]/.test(value) || !/\d/.test(value)) return callback(new Error(this.$t('login.rules.passwordComplex')))
      return callback()
    }
    // 二次密码一致性校验
    const equalTo = (rule, value, callback) => {
      if (!value) return callback(new Error(this.$t('login.rules.passwordAgainRequired')))
      if (value !== this.registerForm.password) return callback(new Error(this.$t('login.rules.passwordNotEqual')))
      return callback()
    }

    return {
      loginForm: {
        username: '',
        password: ''
      },
      loginRules: {
        username: [
          { required: true, message: this.$t('login.rules.usernameRequired'), trigger: 'blur' },
          { type: 'string', min: 3, message: this.$t('login.rules.usernameMin'), trigger: 'blur' }
        ],
        password: [
          { required: true, message: this.$t('login.rules.passwordRequired'), trigger: 'blur' },
          { type: 'string', min: 6, message: this.$t('login.rules.passwordMin'), trigger: 'blur' }
        ]
      },
      registerRules: {
        username: [
          { required: true, message: this.$t('login.rules.usernameRequired'), trigger: 'blur' },
          { type: 'string', min: 3, message: this.$t('login.rules.usernameMin'), trigger: 'blur' },
          { pattern: usernamePattern, message: '仅限字母/数字/下划线', trigger: 'blur' }
        ],
        password: [
          { validator: validatePasswordComplex, trigger: 'blur' }
        ],
        password2: [
          { validator: equalTo, trigger: 'blur' }
        ],
        userType: [
          { required: true, message: this.$t('login.selectRole'), trigger: 'change' }
        ]
      },
      loading: false,
      // 独立的密码可见性控制
      passwordTypes: {
        login: 'password',
        register1: 'password',
        register2: 'password'
      },
      redirect: undefined,
      isLoginPage: true,
      registerForm: {
        username: '',
        password: '',
        password2: '',
        userType: ''
      },
      options: [{
        value: '原料供应商',
        label: '原料供应商'
      }, {
        value: '制造商',
        label: '制造商'
      }, {
        value: '物流承运商',
        label: '物流承运商'
      }, {
        value: '经销商',
        label: '经销商'
      }, {
        value: '零售商',
        label: '零售商'
      }],
      LENGTHS: Object.assign({ login: { username: 50, password: 100 }}, LENGTHS)
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  methods: {
    showPwd(which) {
      // 切换指定输入框的显示/隐藏
      const nextType = this.passwordTypes[which] === 'password' ? 'text' : 'password'
      this.$set(this.passwordTypes, which, nextType)
      // 聚焦到对应输入框
      const refMap = {
        login: 'loginPassword',
        register1: 'registerPassword1',
        register2: 'registerPassword2'
      }
      this.$nextTick(() => {
        const refName = refMap[which]
        if (refName && this.$refs[refName]) {
          this.$refs[refName].focus()
        }
      })
    },
    handleLogin() {
      // fallback sanitize before submit
      this.loginForm.username = sanitize(this.loginForm.username, this.LENGTHS.login.username)
      this.loginForm.password = sanitize(this.loginForm.password, this.LENGTHS.login.password)
      this.$refs.loginForm.validate(valid => {
        if (!valid) {
          return
        }
        this.loading = true
        this.$store.dispatch('user/login', this.loginForm).then(() => {
          this.$router.push({ path: this.redirect || '/' })
          // 登录成功后清空密码
          this.loginForm.password = ''
          // 登录成功后重置可见性
          this.passwordTypes.login = 'password'
        }).catch(err => {
          this.$message.error((err && err.message) || this.$t('login.messages.loginFailed'))
        }).finally(() => {
          this.loading = false
        })
      })
    },
    handleRegister() {
      // 切换登录注册
      this.isLoginPage = !this.isLoginPage
      // 切换时统一重置所有密码为隐藏
      this.passwordTypes = { login: 'password', register1: 'password', register2: 'password' }
      // 切换后清空并清理校验提示
      this.$nextTick(() => {
        if (this.isLoginPage) {
          this.$refs.loginForm.clearValidate()
        }
      })
    },
    submitRegister() {
      // fallback sanitize before submit
      this.registerForm.username = sanitize(this.registerForm.username, this.LENGTHS.login.username)
      this.registerForm.password = sanitize(this.registerForm.password, this.LENGTHS.login.password)
      this.registerForm.password2 = sanitize(this.registerForm.password2, this.LENGTHS.login.password)
      this.$refs.loginForm.validate(valid => {
        if (!valid) {
          return
        }
        if (this.registerForm.password !== this.registerForm.password2) {
          this.$message.error(this.$t('login.rules.passwordNotEqual'))
          return
        }
        const overlay = this.$loading({
          lock: true,
          text: this.$t('login.messages.registering'),
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        })
        this.loading = true
        this.$store.dispatch('user/register', this.registerForm).then(response => {
          this.$router.push({ path: this.redirect || '/' })
          this.$message({
            message: '注册成功，链上交易ID：' + response.txid,
            type: 'success'
          })
          // 注册成功后清空注册表单
          this.registerForm = { username: '', password: '', password2: '', userType: '' }
          // 重置密码可见性并切回登录页
          this.passwordTypes = { login: 'password', register1: 'password', register2: 'password' }
          this.handleRegister() // 切回登录页
        }).catch(err => {
          this.$message.error((err && err.message) || this.$t('login.messages.registerFailed'))
        }).finally(() => {
          overlay.close()
          this.loading = false
        })
      })
    }
  }
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg:#283443;
$light_gray:#fff;
$cursor: #fff;

@supports (-webkit-mask: none) and (not (caret-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  // 添加背景图片
  background-image: url("../../assets/login_images/nature.jpg");
  background-size: cover;
  background-position: center;

  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      background: transparent;
      border: 0px;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $light_gray;
      height: 47px;
      caret-color: $cursor;
      // -webkit-appearance: none;
      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $cursor !important;
      }
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    color: #454545;
  }
}
</style>

<style lang="scss" scoped>
$bg:#2d3a4b;
$dark_gray:#889aa4;
$light_gray:#eee;

.login-container {
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  overflow: hidden;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }

  ::v-deep .el-input__inner::placeholder {
  color: white !important;
  }
}
</style>
