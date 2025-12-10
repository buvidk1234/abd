import CryptoJS from 'crypto-js'
import JSEncrypt from 'jsencrypt'
import type { AxiosResponse, InternalAxiosRequestConfig } from 'axios'

// RSA公钥（后端提供的）
const rsaPublicKey = import.meta.env.VITE_RSA_PUBLIC_KEY

// 初始化RSA加密器
const rsa = new JSEncrypt()
rsa.setPublicKey(rsaPublicKey as string)

// 用于存储AES密钥，方便解密响应时使用
const aesKeyStore = new Map<string, CryptoJS.lib.WordArray>()

// 生成请求ID用于匹配请求和响应
const generateRequestId = () => {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

// 标记需要加密的请求头
export const reqHeader = {
  'X-AES-Key': '',
}

export const encryptData = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
  // 1、判断是否需要加密（根据header中包含X-AES-Key判断）
  const { data, headers } = config
  if (headers?.['X-AES-Key'] === undefined) {
    return config
  }

  // 2、生成AES密钥
  const aesKey = CryptoJS.lib.WordArray.random(128 / 8)

  // 3、生成请求ID并存储AES密钥
  const requestId = generateRequestId()
  aesKeyStore.set(requestId, aesKey)

  // 清理过期的密钥（保留最近20个）
  if (aesKeyStore.size > 20) {
    const firstKey = aesKeyStore.keys().next().value as string
    aesKeyStore.delete(firstKey)
  }

  // 4、使用RSA公钥加密AES密钥
  // 将WordArray转为Base64字符串再加密（后端用Base64解码）
  const aesKeyBase64 = CryptoJS.enc.Base64.stringify(aesKey)
  const encryptedAesKey = rsa.encrypt(aesKeyBase64)
  if (!encryptedAesKey) {
    console.error('RSA加密AES密钥失败')
    throw new Error('加密失败')
  }

  // 5、使用AES加密数据
  // 后端使用的是Hutool的AES，默认是ECB模式，PKCS7填充
  const dataStr = JSON.stringify(data)
  const encrypted = CryptoJS.AES.encrypt(dataStr, aesKey, {
    mode: CryptoJS.mode.ECB,
    padding: CryptoJS.pad.Pkcs7,
  })
  // 6、返回加密后的数据
  config.data = encrypted.toString()
  config.headers['X-AES-Key'] = encryptedAesKey
  config.headers['X-Request-Id'] = requestId
  config.headers['Content-Type'] = 'application/json'
  return config
}

export const decryptData = (response: AxiosResponse) => {
  const { data, headers } = response
  // 1、判断是否需要解密（根据header中包含X-AES-Key判断）
  const requestId = headers?.['x-request-id'] || headers?.['X-Request-Id']
  if (!requestId) {
    // 不需要解密
    return response
  }

  try {
    // 2、获取之前保存的AES密钥
    const aesKey = aesKeyStore.get(requestId)
    if (!aesKey) {
      console.warn('未找到对应的AES密钥，可能是响应超时或密钥已过期')
      // 如果找不到密钥，返回原数据，糟糕啦
      return response
    }

    // 3、解密响应数据
    // 后端返回的是Base64编码的密文
    const decrypted = CryptoJS.AES.decrypt(data, aesKey, {
      mode: CryptoJS.mode.ECB,
      padding: CryptoJS.pad.Pkcs7,
    })

    // 4、转换为字符串并解析JSON
    const decryptedStr = decrypted.toString(CryptoJS.enc.Utf8)
    // 处理后端返回的JSON字符串
    const decryptedData = JSON.parse(decryptedStr)
    // 5、清理使用过的密钥
    aesKeyStore.delete(requestId)
    response.data = decryptedData
    response.headers = headers
    return response
  } catch (error) {
    console.error('解密响应数据失败:', error)
    // 解密失败时返回原数据，又糟糕啦
    return response
  }
}

export default {
  reqHeader,
  encryptData,
  decryptData,
}
