export default function download(data: Blob, fileName = 'download') {
  const blobURL = window.URL.createObjectURL(new Blob([data]))
  const link = document.createElement('a')
  link.href = blobURL
  link.setAttribute('download', fileName)
  link.click()
}
