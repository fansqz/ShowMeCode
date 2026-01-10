module.exports = {
    "env": {
        browser: true,
        es2021: true,
        node: true,
        jest: true,
    },
    /* 接信息语法 */
    parser: "vue-eslint-parser",
    parserOptions: {
        ecmaVersion: "latest",
        sourceType: "module",
        parser: "@typescript-eslint/parser",
        jsxPragma: "React",
        ecmaFeatures: {
            jsx: true,
        },
    },
    /* 继承已有规则 */
    extends: [
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:vue/vue3-essential"
    ],
    overrides: [
        {
            "env": {
                "node": true
            },
            "files": [
                ".eslintrc.{js,cjs}"
            ],
            "parserOptions": {
                "sourceType": "script"
            }
        }
    ],
    plugins: [
        "@typescript-eslint",
        "vue",
        "import"
    ],
    rules: {
        "no-var": "error", // 不使用var
        "no-multiple-empty-lines": ["warn", { max: 1}], //不允许多个空格
        "no-console": process.env.NODE_ENV == "production" ? "error" : "off",
        "no-debugger": process.env.NODE_ENV == "production" ? "error" : "off",
        "no-unexpected-multiline": "error", //禁止空余的多行
        "no-useless-escape": "off", //禁止不必要的转义字符

        "@typescript-eslint/no-unused-vars": ["error", {
            "argsIgnorePattern": "^_",
            "varsIgnorePattern": "^_",
            "caughtErrorsIgnorePattern": "^_"
        }], //禁止定义未使用的变量，但忽略下划线前缀的变量
        "@typescript-eslint/prefer-ts-expect-error": "off", //禁止使用@ts-ignore
        "@typescript-eslint/no-explicit-any": "off", //禁止使用any
        "@typescript-eslint/no-non-null-assertion": "off",
        "@typescript-eslint/no-namespace": "off", //禁止使用自定义ts模块
        "@typescript-eslint/ban-types": "off", // 关闭对{}类型的严格限制
        "@typescript-eslint/semi": ["error", "never"],

        // Import 规范 - 合理的顺序检查
        "import/order": [
          "error",
          {
            "groups": [
              "builtin",
              "external",
              "internal",
              "parent",
              "sibling",
              "index"
            ],
            "pathGroups": [
              {
                "pattern": "vue",
                "group": "external",
                "position": "before"
              },
              {
                "pattern": "@/**",
                "group": "internal",
                "position": "after"
              }
            ],
            "pathGroupsExcludedImportTypes": ["vue"],
            "newlines-between": "always",
            "alphabetize": {
              "order": "asc",
              "caseInsensitive": true
            }
          }
        ],
        "import/no-duplicates": "error", // 保留重复导入检查
        "import/no-unresolved": "off", // 关闭，因为 TypeScript 会处理
        "import/extensions": "off", // 关闭文件扩展名检查，Vue/TS 项目通常不需要

        "vue/multi-word-component-names": "off", //要求组件以"-"链接单词
        "vue/script-setup-uses-vars": "error",
        "vue/no-mutation-props": "off",
        "vue/attribute-hyphenation": "off",

    }
};
