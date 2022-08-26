export const testPath = (path: string) => {
  return /.*[^\\/]$/.test(path);
};

export const testUrl = (url: string) => {
  return /^(http|https):\/\/.*?[^\\/]$/.test(url);
};
