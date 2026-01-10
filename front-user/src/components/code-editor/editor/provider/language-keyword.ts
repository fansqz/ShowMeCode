// 封装获取不同语言关键字的逻辑
export const getLanguageKeywords = (language) => {
  if (language === 'go') {
    return [
      'break',
      'default',
      'func',
      'interface',
      'select',
      'case',
      'defer',
      'go',
      'map',
      'struct',
      'chan',
      'else',
      'goto',
      'package',
      'switch',
      'const',
      'fallthrough',
      'if',
      'range',
      'type',
      'continue',
      'for',
      'import',
      'return',
      'var',
    ]
  } else if (language === 'c') {
    return [
      'auto',
      'break',
      'case',
      'char',
      'const',
      'continue',
      'default',
      'do',
      'double',
      'else',
      'enum',
      'extern',
      'float',
      'for',
      'goto',
      'if',
      'int',
      'long',
      'register',
      'return',
      'short',
      'signed',
      'sizeof',
      'static',
      'struct',
      'switch',
      'typedef',
      'union',
      'unsigned',
      'void',
      'volatile',
      'while',
    ]
  }
  return []
}

// 封装获取不同语言基本数据类型的逻辑
export const getLanguageBaseDataTypes = (language) => {
  if (language === 'go') {
    return [
      'bool',
      'int',
      'int8',
      'int16',
      'int32',
      'int64',
      'uint',
      'uint8',
      'uint16',
      'uint32',
      'uint64',
      'float32',
      'float64',
      'complex64',
      'complex128',
      'string',
      'byte',
      'rune',
    ]
  } else if (language === 'c') {
    return [
      'char',
      'int',
      'float',
      'double',
      'void',
      'signed',
      'unsigned',
      'short',
      'long',
      'long long',
    ]
  }
  return []
}
