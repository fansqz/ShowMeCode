<template>
  <el-card class="card">
    <div class="user-message">
      <el-avatar
        v-if="accountInfo.avatar != ''"
        :src="accountInfo.avatar"
        class="avatar"
      ></el-avatar>
      <el-avatar v-if="accountInfo.avatar == ''" class="avatar">S</el-avatar>
      <div class="base-info">
        <div class="base-info-item1">
          {{ accountInfo.username }}
        </div>
        <div class="base-info-item"> 性别：{{ accountInfo.sex == '1' ? '男' : '女' }} </div>
        <div class="base-info-item"> 生日：{{ accountInfo.birthDay }} </div>
        <div class="login-name"> 用户id: {{ accountInfo.loginName }} </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
  import { reactive, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { reqAccountInfo } from '@/api/account'
  import { reqGetURL } from '@/api/common'
  let accountInfo = reactive({
    avatar: '',
    sex: '',
    loginName: '',
    username: '',
    birthDay: '',
    introduction: '',
    codingAge: 0,
  })

  onMounted(async () => {
    getAccountInfo()
  })

  const getAccountInfo = async () => {
    let result = await reqAccountInfo()
    if (result.code == 200) {
      let data = result.data
      if (data.avatar != '') {
        let result2 = await reqGetURL(data.avatar)
        if (result2.code == 200) {
          accountInfo.avatar = result2.data
        }
      }
      accountInfo.sex = data.sex
      accountInfo.loginName = data.loginName
      accountInfo.username = data.username
      accountInfo.birthDay = data.birthDay
      accountInfo.introduction = data.introduction
      accountInfo.codingAge = data.codingAge
    } else {
      ElMessage({
        type: 'error',
        message: result.message,
      })
    }
  }
</script>

<style scoped lang="scss">
  .card {
    width: 800px;

    .user-message {
      position: relative;
      box-sizing: border-box;
      height: 110px;
      border-radius: 5px;

      .avatar {
        position: absolute;
        left: 50px;
        box-sizing: border-box;
        width: 110px;
        height: 110px;
        border: 1px solid $border-color;
      }

      .base-info {
        position: absolute;
        left: 180px;
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        width: 500px;
        height: 110px;

        .base-info-item1 {
          margin-top: 5px;
          font-size: xx-large;
          font-weight: 900;
        }

        .base-info-item {
          margin-top: 6px;
          font-size: small;
        }

        .login-name {
          margin-top: 6px;
          font-size: small;
          color: $text-color-secondary;
        }
      }
    }
  }
</style>
