<template>
  <div class="login-page">
    <HeaderBar title="手机号码快捷登录" />

    <div class="login-form">
      <div class="input-row">
        <el-input
          v-model="form.phone"
          placeholder="请输入手机号"
          size="large"
        />
        <el-button
          type="success"
          :disabled="codeSending || countdown > 0"
          @click="sendCode"
        >
          {{ countdown > 0 ? `${countdown}秒后重发` : '发送验证码' }}
        </el-button>
      </div>

      <el-input
        v-model="form.code"
        placeholder="请输入验证码"
        size="large"
      />

      <p class="tips">未注册的手机号码验证后自动创建账户</p>

      <el-button
        type="primary"
        size="large"
        class="login-btn"
        @click="login"
      >
        登录
      </el-button>

      <div class="other-link">
        <router-link to="/login-password">密码登录</router-link>
      </div>
    </div>

    <div class="agreement">
      <el-radio v-model="agreed" :label="true">
        我已阅读并同意
        <a href="javascript:void(0)">《黑马点评用户服务协议》</a>、
        <a href="javascript:void(0)">《隐私政策》</a>
        等，接受免除或者限制责任、诉讼管辖约定等粗体标示条款
      </el-radio>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { sendCode as sendCodeApi, loginByCode } from '@/api/user'
import { useUserStore } from '@/stores/user'
import HeaderBar from '@/components/HeaderBar.vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const form = reactive({
  phone: '',
  code: ''
})

const agreed = ref(false)
const codeSending = ref(false)
const countdown = ref(0)
let countdownTimer = null

const sendCode = async () => {
  if (!form.phone) {
    ElMessage.warning('请输入手机号')
    return
  }
  
  const phoneReg = /^1[3-9]\d{9}$/
  if (!phoneReg.test(form.phone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }
  
  codeSending.value = true
  try {
    await sendCodeApi(form.phone)
    ElMessage.success('验证码已发送')
    
    // 开始倒计时
    countdown.value = 60
    countdownTimer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(countdownTimer)
      }
    }, 1000)
  } catch (error) {
    ElMessage.error('发送失败')
  } finally {
    codeSending.value = false
  }
}

const login = async () => {
  if (!agreed.value) {
    ElMessage.warning('请先同意用户协议')
    return
  }
  
  if (!form.phone || !form.code) {
    ElMessage.warning('请输入手机号和验证码')
    return
  }
  
  try {
    const token = await loginByCode(form)
    userStore.setToken(token)
    await userStore.fetchUserInfo()
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    // 错误已在拦截器处理
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: #fff;
}

.login-form {
  padding: 30px 20px;
}

.input-row {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.input-row .el-input {
  flex: 1;
}

.input-row .el-button {
  width: 120px;
}

.tips {
  text-align: center;
  color: #999;
  font-size: 12px;
  margin: 15px 0;
}

.login-btn {
  width: 100%;
  background: #ff6633;
  border-color: #ff6633;
}

.login-btn:hover {
  background: #ff5500;
  border-color: #ff5500;
}

.other-link {
  text-align: right;
  margin-top: 15px;
}

.other-link a {
  color: #333;
  text-decoration: none;
}

.agreement {
  padding: 0 20px;
  margin-top: 30px;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.agreement :deep(.el-radio__label) {
  font-size: 12px;
  color: #666;
  line-height: 1.6;
  display: block;
  word-wrap: break-word;
  overflow-wrap: break-word;
  white-space: normal;
}

.agreement a {
  color: #ff6633;
}
</style>
