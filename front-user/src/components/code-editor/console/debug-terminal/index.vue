<template>
  <div class="debug_terminal">
    <el-splitter class="console-splitter">
      <el-splitter-panel size="30" class="left-panel">
        <!--栈帧列表-->
        <div class="frames-card">
          <div class="frames_title">
            <el-text class="frames_text">栈帧</el-text>
          </div>
          <Frames class="frames" />
        </div>
      </el-splitter-panel>
      <el-splitter-panel size="70" class="right-panel">
        <div class="variables-and-input-card">
          <div class="custom-tabs">
            <el-text
              class="tab-item"
              :class="{ active: activeIndex === 'console' }"
              @click="handleSelect('console')"
            >
              控制台
            </el-text>
            <el-text
              class="tab-item"
              :class="{ active: activeIndex === 'variables' }"
              @click="handleSelect('variables')"
            >
              变量
            </el-text>
          </div>
          <div class="console" v-show="activeIndex == 'console'">
            <Console />
          </div>
          <div class="variables" v-show="activeIndex == 'variables'">
            <el-scrollbar class="scrollbar">
              <Variables />
            </el-scrollbar>
          </div>
        </div>
      </el-splitter-panel>
    </el-splitter>
  </div>
</template>

<script setup lang="ts">
  import Console from './console.vue'
  import Variables from './variables.vue'
  import Frames from './frames.vue'
  import { ref } from 'vue'

  const activeIndex = ref('console')

  const handleSelect = (key: string) => {
    activeIndex.value = key
  }
</script>

<style lang="scss" scoped>
  .debug_terminal {
    display: flex;
    flex-flow: row nowrap;
    width: 100%;
    height: 100%;
    background-color: $base-background-color;

    :deep(.left-panel) {
      margin: 6px;
      overflow: hidden;
      border: 1px solid $border-color;
      border-radius: $border-radius-large;
    }

    :deep(.right-panel) {
      margin: 6px;
      overflow: hidden;
      border: 1px solid $border-color;
      border-radius: $border-radius-large;
    }

    .frames-card {
      position: relative;
      width: 100%;
      height: 100%;
      background-color: $base-background-color;

      .frames_title {
        box-sizing: border-box;
        display: flex;
        width: 100%;
        height: 25px;
        padding-left: 15px;
        background-color: $base-background-color;
        border-bottom: 1px solid $border-color;
      }

      .frames {
        width: 100%;
        height: calc(100% - 25px);
        background-color: $fill-color-lighter;
      }
    }

    .variables-and-input-card {
      position: relative;
      width: 100%;
      height: 100%;
      background-color: $base-background-color;

      .custom-tabs {
        display: flex;
        align-items: center;
        width: 100%;
        height: 25px;
        background-color: $base-background-color;
        border-bottom: 1px solid $border-color;

        .tab-item {
          display: flex;
          align-items: center;
          height: 100%;
          padding: 0 15px;
          cursor: pointer;

          &:hover {
            color: $primary-color;
            background-color: rgba($primary-color, 0.1);
          }

          &.active {
            font-weight: 550;
            color: $primary-color;
            background-color: $base-background-color;
          }
        }
      }

      .console {
        width: 100%;
        height: calc(100% - 25px);
        background-color: $fill-color-lighter;
      }

      .variables {
        width: 100%;
        height: calc(100% - 25px);
        background-color: $fill-color-lighter;

        .scrollbar {
          width: 100%;
          height: 100%;
        }
      }
    }
  }
</style>
