export const toFormData = (data: Record<string, any>) => {
  if (!data) {
    return new FormData();
  }
  if (data instanceof FormData) {
    return data;
  }
  const formData = new FormData();
  if (data instanceof File) {
    formData.append('file', data);
    return formData;
  }
  Object.keys(data).forEach((key) => {
    // 后来使用formData出现了传undefined变成字符串的问题，故删去值为undefined的字段
    if (data[key] !== undefined) formData.append(key, data[key]);
  });
  return formData;
};
