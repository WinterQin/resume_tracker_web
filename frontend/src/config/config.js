export const API_BASE_URL = import.meta.env.PROD 
  ? 'http://8.147.235.248:8080'  // 生产环境 API 地址，使用相对路径
  : 'http://GWDesktop:8080'  // 开发环境 API 地址 