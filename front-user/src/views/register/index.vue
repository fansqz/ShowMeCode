<template>
  <div class="register-container">
    <el-card class="card">
      <el-form class="register-form" :rules="rules" :model="registerForm" ref="registerElFrom">
        <Logo class="logo"></Logo>
        <el-form-item prop="username">
          <el-input
            :prefix-icon="User"
            placeholder="请设置用户名称"
            v-model="registerForm.username"
          />
        </el-form-item>
        <el-form-item prop="email">
          <el-input :prefix-icon="Message" placeholder="请输入邮箱" v-model="registerForm.email" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            placeholder="请设置密码"
            type="password"
            show-password="true"
            v-model="registerForm.password"
          />
        </el-form-item>
        <el-form-item prop="code">
          <el-input placeholder="验证码" v-model="registerForm.code">
            <template #suffix>
              <el-button
                v-if="!downTimeState.countDownIng"
                link
                type="primary"
                @click="sendEmailCode"
              >
                发送验证码
              </el-button>
              <span v-if="downTimeState.countDownIng">{{ downTimeState.countDownTime }}s</span>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            class="register-button"
            type="primary"
            @click="register"
            :loading="loading"
            size="default"
          >
            注册
          </el-button>
        </el-form-item>
        <el-from-item class="go-login">
          已有账户？
          <div
            class="login-text"
            @mouseover="goLoginTextStyle = 'color:blue'"
            @mouseleave="goLoginTextStyle = 'color:black'"
            @click="changeRoute('login')"
            :style="goLoginTextStyle"
          >
            点击登录
          </div>
        </el-from-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
  import Logo from '@/components/logo/index.vue'
  import { User, Message } from '@element-plus/icons-vue'
  import { reactive, ref, onMounted } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { reqSendCode, reqRegister } from '@/api/auth'

  let $router = useRouter()
  let $route = useRoute()
  // 登录按钮是否加载
  let loading = ref(false)
  let registerForm = reactive({
    username: '',
    email: '',
    password: '',
    code: '',
  })
  let goLoginTextStyle = ref()
  // 是否在验证码倒计时中
  const downTimeState = reactive({
    // 倒计时
    countDownTime: 60,
    timer: null,
    countDownIng: false,
  })

  let rules = {
    username: [
      { required: true, message: '用户名称不能为空', trigger: 'change' },
      { required: true, min: 2, max: 30, message: '用户名称过长或过短' },
    ],
    password: [
      { required: true, message: '密码不能为空', trigger: 'change' },
      { required: true, min: 6, max: 20, message: '密码长度需6-20位' },
    ],
    email: [{ required: true, message: '邮箱不能为空', trigger: 'change' }],
    code: [
      { required: true, message: '验证码不能为空', trigger: 'change' },
      { required: true, min: 6, max: 6, message: '验证码需要6位' },
    ],
  }
  let registerElFrom = ref()

  const sendEmailCode = async () => {
    let result = await reqSendCode({
      email: registerForm.email,
      type: 'register',
    })
    if (result.code == 200) {
      ElMessage({
        showClose: true,
        message: result.message,
        type: 'success',
      })
      countDown()
    } else {
      ElMessage({
        showClose: true,
        message: result.message,
        type: 'warning',
      })
    }
  }

  const register = async () => {
    // 判断表单校验是否通过
    try {
      await registerElFrom.value.validate()
      loading.value = true
      let result = await reqRegister(registerForm)
      // 判断是否有redirect参数
      if (result.code == 200) {
        ElMessage({
          type: 'success',
          message: '注册成功，请登录',
        })
        $router.push('login')
      }
      loading.value = false
    } catch (error) {
      loading.value = false
    }
  }

  const countDown = () => {
    let startTime = parseInt(localStorage.getItem('registerStartTimeLogin'))
    let nowTime = new Date().getTime()
    if (startTime) {
      let surplus = 60 - Math.round((nowTime - startTime) / 1000)
      downTimeState.countDownTime = surplus <= 0 ? 0 : surplus
    } else {
      downTimeState.countDownTime = 60
      localStorage.setItem('registerStartTimeLogin', nowTime.toString())
    }

    downTimeState.timer = setInterval(() => {
      downTimeState.countDownTime--
      downTimeState.countDownIng = true
      if (downTimeState.countDownTime <= 0) {
        localStorage.removeItem('registerStartTimeLogin')
        clearInterval(downTimeState.timer)
        downTimeState.countDownTime = 60
        downTimeState.countDownIng = false
      }
    }, 1000)
  }

  const changeRoute = (routeName: string, params = {}) => {
    if ($route.name === routeName) {
      return
    }
    $router.push({ name: routeName, params })
  }

  onMounted(() => {
    let sendEndTime = localStorage.getItem('registerStartTimeLogin')
    if (sendEndTime) {
      downTimeState.countDownIng = true
      countDown()
    }
  })
</script>

<style scoped lang="scss">
  .register-container {
    position: absolute;
    top: 0;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    width: 100%;
    height: calc(100vh - $header-height);
    background-color: $base-background-color;

    .card {
      margin-top: calc(50vh - 280px); /* 调整偏移量 */
      .register-form {
        width: 280px;
        height: 360px;
        padding: 10px;
        text-align: center;

        .logo {
          height: 50px;
        }

        .go-login {
          font-size: small;

          .login-text {
            display: inline;
            cursor: pointer;
          }
        }

        .register-button {
          width: 100%;
        }
      }
    }
  }
</style>
