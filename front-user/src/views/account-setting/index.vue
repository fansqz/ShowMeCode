<template>
  <div class="account-setting-container">
    <el-card class="card">
      <template #header> 账号信息管理 </template>
      <el-form :model="account" label-width="100px" class="account-form">
        <el-form-item label="头像" label-width="100px">
          <!--头像-->
          <el-upload
            class="avatar-uploader"
            :show-file-list="false"
            :http-request="uploadAvatar"
            :before-upload="beforeUpload"
          >
            <el-image v-if="account.avatarURL" :src="account.avatarURL" class="avatar"
              >点击上传</el-image
            >
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="用户名" label-width="100px">
          <el-input v-model="account.username" />
        </el-form-item>
        <el-form-item label="性别" label-width="100px ">
          <el-radio-group v-model="account.sex">
            <el-radio border :label="2">女</el-radio>
            <el-radio border :label="1">男</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="出生日期" label-width="100px">
          <el-date-picker v-model="account.birthDay" type="date" />
        </el-form-item>
        <el-form-item label="个性签名" label-width="100px">
          <el-input v-model="account.introduction" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <div class="account-submit">
        <el-button type="info" @click="cancelAccountSubmit">取消</el-button>
        <el-button type="primary" @click="changeAccountSubmit">提交</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
  import { reactive, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { useRoute, useRouter } from 'vue-router'
  import { reqAccountInfo, reqUpdateAccount, reqUploadAvatar } from '@/api/account'
  import { reqGetURL } from '@/api/common'

  const $route = useRoute()
  const $router = useRouter()
  let account = reactive({
    avatar: '',
    avatarURL: '',
    username: '',
    introduction: '',
    sex: '',
    birthDay: new Date(),
  })

  // 获取题目
  const readAccount = async () => {
    try {
      let result = await reqAccountInfo()
      if (result.code == 200) {
        account.avatar = result.data.avatar
        let result2 = await reqGetURL(result.data.avatar)
        if (account.avatar != '') {
          if (result2.code == 200) {
            account.avatarURL = result2.data
          }
        }
        account.username = result.data.username
        account.introduction = result.data.introduction
        account.sex = result.data.sex
        account.birthDay = new Date(result.data.birthDay)
      }
    } catch (err) {
      ElMessage({
        showClose: true,
        message: '用户数据读取失败',
        type: 'error',
      })
    }
  }

  onMounted(() => {
    readAccount()
  })

  // 提交账号修改
  const changeAccountSubmit = async () => {
    try {
      let result = await reqUpdateAccount({
        avatar: account.avatar,
        username: account.username,
        introduction: account.introduction,
        sex: account.sex,
        birthDay: dateToStr(account.birthDay),
      })
      if (result.code == 200) {
        ElMessage({
          showClose: true,
          message: result.message,
          type: 'success',
        })
      } else {
        ElMessage({
          showClose: true,
          message: '提交失败',
          type: 'error',
        })
      }
    } catch {
      ElMessage({
        showClose: true,
        message: '提交失败',
        type: 'error',
      })
    }
  }

  // 添加文件并校验
  const beforeUpload = (file: any) => {
    // 检验
    const isLt2M = file.size / 1024 / 1024 < 2
    if (!isLt2M) {
      ElMessage.error('上传图片大小不能超过 2MB！')
      return false
    }
    return true
  }

  // 上传头像
  const uploadAvatar = async (params: any) => {
    let result = await reqUploadAvatar({
      avatar: params.file,
    })
    if (result.code == 200) {
      account.avatar = result.data
      let result2 = await reqGetURL(account.avatar)
      if (result2.code == 200) {
        account.avatarURL = result2.data
      }
      ElMessage({
        type: 'success',
        message: '头像上传成功',
      })
    } else {
      ElMessage({
        type: 'error',
        message: result.message,
      })
    }
  }

  const cancelAccountSubmit = () => {
    changeRoute('myprofile')
  }

  const dateToStr = (date: Date): string => {
    const year: number = date.getFullYear()
    const month: number = date.getMonth() + 1
    const day: number = date.getDate()
    const hours: number = date.getHours()
    const minutes: number = date.getMinutes()
    const seconds: number = date.getSeconds()

    // 格式化日期字符串
    const formattedDate: string = `${year}-${padZero(month)}-${padZero(day)} ${padZero(
      hours,
    )}:${padZero(minutes)}:${padZero(seconds)}`
    return formattedDate
  }

  // 辅助函数：在不足两位数时补0
  const padZero = (value: number): string => {
    return value.toString().padStart(2, '0')
  }

  const changeRoute = (routeName: string, params = {}) => {
    if ($route.name === routeName) {
      return
    }
    $router.push({ name: routeName, params: params })
  }
</script>

<style scoped lang="scss">
  .account-setting-container {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    padding: 30px;

    .card {
      width: 70%;
      height: 90%;

      .account-form {
        padding: 10px 20px;

        .avatar-uploader {
          position: relative;
          overflow: hidden;
          cursor: pointer;
          border: 1px dashed $border-color;
          border-radius: 6px;
          transition: var(--el-transition-duration-fast);

          .avatar {
            display: block;
            width: 178px;
            height: 178px;
          }

          .avatar-uploader-icon {
            width: 150px;
            height: 150px;
            font-size: 28px;
            color: $text-color-secondary;
            text-align: center;
          }
        }
      }

      .account-submit {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: auto;
      }
    }
  }
</style>
