<template>
  <div class="login-page">
    <HeaderBar title="密码登录" />

    <div class="login-form">
      <el-input
        v-model="form.phone"
        placeholder="请输入手机号"
        size="large"
      />

      <el-input
        v-model="form.password"
        type="password"
        placeholder="请输入密码"
        size="large"
        show-password
      />

      <div class="forgot-link">
        <a href="javascript:void(0)">忘记密码</a>
      </div>

      <el-button
        type="primary"
        size="large"
        class="login-btn"
        @click="login"
      >
        登录
      </el-button>

      <div class="other-link">
        <router-link to="/login">验证码登录</router-link>
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
import { loginByPassword } from '@/api/user'
import { useUserStore } from '@/stores/user'
import HeaderBar from '@/components/HeaderBar.vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const form = reactive({
  phone: '',
  password: ''
})

const agreed = ref(false)

const login = async () => {
  if (!agreed.value) {
    ElMessage.warning('请先同意用户协议')
    return
  }
  
  if (!form.phone || !form.password) {
    ElMessage.warning('请输入手机号和密码')
    return
  }
  
  try {
    const token = await loginByPassword(form)
    userStore.setToken(token)
    await userStore.fetchUserInfo()
    ElMessage.success('登录成功')
    router.push('/profile')
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

.login-form .el-input {
  margin-bottom: 15px;
}

.forgot-link {
  text-align: center;
  margin: 10px 0 20px;
}

.forgot-link a {
  color: #999;
  font-size: 14px;
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
}

.agreement :deep(.el-radio__label) {
  font-size: 12px;
  color: #666;
  line-height: 1.6;
}

.agreement a {
  color: #ff6633;
}
</style>
