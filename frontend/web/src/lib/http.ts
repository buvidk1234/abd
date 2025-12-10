import axios from 'axios'
import { useUserStore } from '../store/userStore'
import { toast } from 'sonner'
import { decryptData, encryptData } from '@/utils/encrypt'
export const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
})

interface R<T> {
  code: number
  msg: string
  data: T
}

// 添加token
http.interceptors.request.use((config) => {
  const token = useUserStore.getState().getToken()
  if (token) {
    config.headers.Authorization = token
  }
  return config
})

// 加密请求体
http.interceptors.request.use((config) => {
  return encryptData(config)
})

// 统一处理错误
http.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (axios.isAxiosError(error)) {
      if (error.code === 'ECONNABORTED') {
        toast('请求超时，请稍后重试')
      } else if (error.response) {
        const status = error.response.status
        const serverMsg = (error.response.data &&
          (error.response.data.msg || error.response.data.message)) as string | undefined
        const fallbackMsg = error.message || `请求错误（${status}）`
        toast(serverMsg || fallbackMsg)
      } else {
        toast('网络异常，请检查网络连接')
      }
    } else {
      toast(String(error))
    }
    return Promise.reject(error)
  }
)

// 解密响应体
http.interceptors.response.use((response) => {
  return decryptData(response)
})
// 统一解包
http.interceptors.response.use((response) => {
  const { data, msg, code } = response.data as R<unknown>
  if (code !== 0) {
    toast(msg)
    return Promise.reject(msg)
  }
  return {
    ...response,
    data,
  }
})
