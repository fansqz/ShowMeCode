<template>
  <div class="container">
    <el-tabs v-model="actionLanguageValue" type="card" editable class="tabs" @edit="handleTabsEdit">
      <el-tab-pane
        v-for="item in codeList"
        :key="item.language"
        :label="item.language"
        :name="item.language"
        style="padding: 2px"
      >
      </el-tab-pane>
    </el-tabs>
    <div v-for="(item, index) in codeList" class="code-editor-item">
      <div class="code-editor-item-div" v-show="item.language == actionLanguageValue">
        <CodeEditorItem v-model="codeList[index]" />
      </div>
    </div>
  </div>

  <el-dialog v-model="addVisualCode.visible" :title="addVisualCode.title">
    <el-radio-group v-model="addVisualCode.selectedLanguage">
      <el-radio
        v-for="item in addVisualCode.optionalLanguage"
        :key="item"
        :value="item"
        :label="item"
        size="large"
        >{{ item }}</el-radio
      >
    </el-radio-group>
    <template #footer>
      <div class="dialog-footer">
        <el-button
          @click="
            () => {
              addVisualCode.visible = false;
            }
          "
          >取消</el-button
        >
        <el-button @click="handleAddVisualCode()" type="primary">提交</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue';
  import type { TabPaneName } from 'element-plus';
  import useVisualDocumentStore from '@/store/modules/visual-document.ts';
  import CodeEditorItem from './code-editor-item.vue';
  import { defaultCodeMap } from './code/default-code';
  import { languageConstants, supportedLanguages } from '@/constants/languages.ts';

  const visualDocumentStore = useVisualDocumentStore();
  let codeList = visualDocumentStore.codeList;

  const addVisualCode = reactive({
    title: '添加语言',
    visible: false,
    optionalLanguage: [languageConstants.C],
    selectedLanguage: languageConstants.C,
  });

  const actionLanguageValue = ref('');
  if (codeList.length >= 1) {
    actionLanguageValue.value = codeList[0].language;
  }

  const handleAddVisualCode = () => {
    if (addVisualCode.optionalLanguage.length == 0) {
      // 可选语言为空，则不执行添加代码任务
      return;
    }
    let defaultCode = defaultCodeMap.get(addVisualCode.selectedLanguage);
    if (defaultCode == undefined) {
      defaultCode = '';
    }
    codeList.push({
      code: defaultCode,
      language: addVisualCode.selectedLanguage,
      breakpoints: [],
    });
    actionLanguageValue.value = addVisualCode.selectedLanguage;
    addVisualCode.visible = false;
  };

  const handleTabsEdit = (targetName: TabPaneName | undefined, action: 'remove' | 'add') => {
    if (action === 'add') {
      // 处理可选的语言
      let optionalLanguage: languageConstants[] = [];
      for (let i = 0; i < supportedLanguages.length; i++) {
        let isSelected = false;
        for (let j = 0; j < codeList.length; j++) {
          if (codeList[j].language == supportedLanguages[i]) {
            isSelected = true;
          }
        }
        if (!isSelected) {
          optionalLanguage.push(supportedLanguages[i]);
        }
      }
      addVisualCode.optionalLanguage = optionalLanguage;
      if (optionalLanguage.length > 0) {
        addVisualCode.selectedLanguage = optionalLanguage[0];
      }
      addVisualCode.visible = true;
    } else if (action === 'remove') {
      // 删除语言
      const tabs = codeList;
      let activeName = actionLanguageValue.value;
      // 如果删除的语言等于目标语言，那么需要改变当前选择的语言
      if (activeName === targetName) {
        tabs.forEach((tab, index) => {
          if (tab.language === targetName) {
            const nextTab = tabs[index + 1] || tabs[index - 1];
            if (nextTab) {
              activeName = nextTab.language;
            } else {
              activeName = '';
            }
          }
        });
      }

      actionLanguageValue.value = activeName;
      // 删除目标语言
      visualDocumentStore.codeList = tabs.filter((tab) => tab.language !== targetName);
      codeList = visualDocumentStore.codeList;
    }
  };

  onMounted(() => {});
</script>

<style scoped lang="scss">
  .container {
    position: relative;
    width: 100%;
    height: 100%;
    margin: 0px;
    :deep(.el-tabs__header) {
      padding: 0 10px;
      margin: 0px;
    }

    .code-editor-item-div {
      position: absolute;
      top: 42px;
      height: calc(100% - 42px);
      width: 100%;
      padding: 0px;
    }
  }
</style>
@/constants/description@/constants/languages
