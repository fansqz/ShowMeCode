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
        "vue"
    ],
    rules: {
        "no-var": "error", // 不使用var
        "no-multiple-empty-lines": ["warn", { max: 1}], //不允许多个空格
        "no-console": process.env.NODE_ENV == "production" ? "error" : "off",
        "no-debugger": process.env.NODE_ENV == "production" ? "error" : "off",
        "no-unexpected-multiline": "error", //禁止空余的多行
        "no-useless-escape": "off", //禁止不必要的转义字符

        "@typescript-eslint/no-unused-vars": "error", //禁止定义未使用的变量
        "@typescript-eslint/prefer-ts-expect-error": "off", //禁止使用@ts-ignore
        "@typescript-eslint/no-explicit-any": "off", //禁止使用any
        "@typescript-eslint/no-non-null-assertion": "off",
        "@typescript-eslint/no-namespace": "off", //禁止使用自定义ts模块
        "@typescript-eslint/semi": "error",

        "vue/multi-word-component-names": "off", //要求组件以“-”链接单词
        "vue/script-setup-uses-vars": "error",
        "vue/no-mutation-props": "off",
        "vue/attribute-hyphenation": "off",

    }
};