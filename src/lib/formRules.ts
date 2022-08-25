export const testPath = (path: string) => {
  return /.*[^\\/]$/.test(path);
};
