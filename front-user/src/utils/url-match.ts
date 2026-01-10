export const matchPattern = (pattern: string, string: string): boolean => {
  const patternParts = pattern.split('/')
  const stringParts = string.split('/')

  if (patternParts.length !== stringParts.length) {
    return false
  }

  for (let i = 0; i < patternParts.length; i++) {
    if (patternParts[i] !== stringParts[i] && patternParts[i] !== ':id') {
      return false
    }
  }

  return true
}
