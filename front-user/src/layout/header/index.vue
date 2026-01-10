<template>
  <div class="header" :class="{ 'header--dark': isDarkHeader }">
    <div class="header_left">
      <Logo class="logo"></Logo>
      <div class="empty" />
      <div :class="{ nav_item: true }" @click="changeRoute('coding')">
        <el-text size="large" :class="{ 'active-text': isActiveNavItem('coding') }">编程</el-text>
      </div>
      <div :class="{ learn: true, nav_item: true }">
        <el-text
          size="large"
          :class="{ 'active-text': isActiveNavItem('learn') }"
          @click="toggleDropdown"
          >学习</el-text
        >
        <el-dropdown ref="dropdownRef" @command="handleSelectBank">
          <span class="dropdown-trigger">
            <el-icon class="el-icon--right">
              <arrow-down />
            </el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-for="item in visualDocumentBanks" :command="item.id">{{
                item.name
              }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    <div class="header_right">
      <ThemeSwitcher class="theme-switcher" />
      <Setting class="setting" v-show="isLogged()"></Setting>
      <div class="login_button" v-show="!isLogged()">
        <div class="login_content">
          <el-button class="button" type="primary" @click="getoLogin()" link>登陆</el-button>
          <el-text>/</el-text>
          <el-button class="button" type="primary" @click="getoRegister()" link>注册</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import Logo from '@/components/logo/index.vue'
  import Setting from './setting.vue'
  import ThemeSwitcher from '@/components/theme-switcher/index.vue'
  import { useRouter, useRoute } from 'vue-router'
  import useUserStore from '@/store/modules/user'
  import { onMounted, computed } from 'vue'
  import { ref } from 'vue'
  import { reqAllVisualDocumentBank } from '@/api/visual-document-bank/index.ts'
  let userStore = useUserStore()
  let $router = useRouter()
  let $route = useRoute()
  const visualDocumentBanks = ref([])
  const dropdownRef = ref()

  // 判断是否为深色header（coding和learn页面）
  const isDarkHeader = computed(() => {
    const routeName = $route.name as string
    return routeName === 'coding' || routeName === 'learn'
  })

  const changeRoute = (routeName: string, params = {}) => {
    if ($route.name === routeName) {
      return
    }
    $router.push({ name: routeName, params })
  }
  const isActiveNavItem = (routeName: string) => {
    return $route.name === routeName
  }

  const isLogged = (): boolean => {
    return !!userStore.token
  }

  const getoLogin = () => {
    $router.push({ name: 'login' })
  }

  const getoRegister = () => {
    $router.push({ name: 'register' })
  }

  const handleSelectBank = (bankID: number) => {
    $router.push({
      name: 'learn',
      params: {
        bankID: bankID,
      },
    })
  }

  const toggleDropdown = () => {
    dropdownRef.value?.handleOpen()
  }

  onMounted(async () => {
    let result = await reqAllVisualDocumentBank()
    if (result.code == 200) {
      visualDocumentBanks.value = result.data
    }
  })
</script>

<style scoped lang="scss">
  .header {
    position: absolute;
    top: 0;
    box-sizing: border-box;
    width: 100%;
    height: $header-height;
    background-color: $base-background-color;
    border-bottom: 1px solid $border-color;

    &--dark {
      background-color: $deep-background-color;
      border-bottom: none;
    }

    .header_left {
      position: absolute;
      left: 0%;
      display: flex;
      align-items: center;
      width: 50%;
      height: 100%;

      .logo {
        display: flex;
        align-items: center;
        height: calc($header-height - 15px);
        margin: 0 20px;
      }

      .empty {
        width: 20px;
      }

      .nav_item {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 60px;
        min-width: 60px;
        height: $header-height;
        margin: 0 10px;
        cursor: pointer;
      }

      .active-text {
        position: relative;

        &::after {
          position: absolute;
          bottom: -8px;
          left: 50%;
          width: 25px;
          height: 3px !important;
          content: '';
          background-color: $primary-color !important;
          border: none !important;
          box-shadow: none !important;
          transform: translateX(-50%);
        }
      }

      .learn {
        display: flex;
        gap: 4px;
        align-items: center;
      }

      .dropdown-trigger {
        display: flex;
        align-items: center;
        cursor: pointer;
      }
    }

    .header_right {
      position: absolute;
      right: 0%;
      display: flex;
      flex-direction: row-reverse;
      width: 50%;
      height: 100%;

      .login_button {
        position: relative;
        width: 160px;
        height: 100%;

        .login_content {
          position: absolute;
          display: flex;
          align-items: center;
          justify-content: center;
          width: 160px;
          height: 100%;

          .button {
            margin: 10px;
          }
        }
      }

      .theme-switcher {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 100%;
        margin-right: 20px;
      }

      .setting {
        width: 120px;
        height: 100%;
        margin-right: 30px;
      }
    }
  }
</style>
