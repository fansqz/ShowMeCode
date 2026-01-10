<template>
  <div class="console-container">
    <div class="console-header">
      <el-text>调试终端</el-text>
      <div class="debug-buttons">
        <StepOrContinueButton type="continue" class="debug-btn continue-button" />
        <StepOrContinueButton type="step-over" class="debug-btn step-over-button" />
        <StepOrContinueButton type="step-in" class="debug-btn step-in-button" />
        <StepOrContinueButton type="step-out" class="debug-btn step-out-button" />
        <TerminateButton class="debug-btn terminal-button" />
      </div>
    </div>
    <div class="console-terminal">
      <DebugTerminal />
    </div>
    <div class="console-bottom">
      <CodeButtonBar :debug="true" :execute="true" class="code-button-bar" />
    </div>
  </div>
</template>

<script setup lang="ts">
  import DebugTerminal from './debug-terminal/index.vue'
  import StepOrContinueButton from './debug-button/debug-button-step-continue.vue'
  import TerminateButton from './debug-button/debug-button-terminate.vue'
  import CodeButtonBar from './coding-button/index.vue'
</script>

<style scoped lang="scss">
  .console-container {
    display: flex;
    flex-direction: column;
    height: 100%;

    .console-header {
      display: flex;
      flex-shrink: 0; // 固定30px
      align-items: center;
      justify-content: space-between;
      height: 30px;
      padding: 0 20px;
      background-color: $base-background-color;

      .debug-buttons {
        display: flex;
        align-items: center;
        justify-content: space-around;
        width: 140px;
        height: 25px;
      }
    }

    .console-terminal {
      flex: 1; // 占据剩余空间
      min-height: 0; // 允许收缩
    }

    .console-bottom {
      flex-shrink: 0; // 固定40px
      height: 36px;
      padding: 2px 0;
      background-color: $base-background-color;

      .code-button-bar {
        display: flex;
        justify-content: flex-end;
        width: 100%;
        padding: 0 10px;
      }
    }
  }
</style>
