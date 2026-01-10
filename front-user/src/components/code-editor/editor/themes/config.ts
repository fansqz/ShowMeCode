// 语言id到作用域名称的映射
export const scopeNameMap = {
  go: 'source.go',
  c: 'source.c',
  java: 'source.java',
}

// 作用域名称到语法文件的映射
export const tmGrammarJsonMap = {
  'source.go': 'go.tmLanguage.json',
  'source.c': 'c.tmLanguage.json',
  'source.asm': 'c.tmLanguage.json',
  'source.x86': 'c.tmLanguage.json',
  'source.x86_64': 'c.tmLanguage.json',
  'source.arm': 'c.tmLanguage.json',
  'source.java': 'java.tmLanguage.json',
}

// 主题配置
export const codeThemeList = ['BlulocoLight', 'GitHubLight', 'Monokai', 'OneDarkPro']
