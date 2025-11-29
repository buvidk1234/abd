<script lang="ts" setup>
import { reactive, ref } from 'vue'
import { loginUser, fetchUsers } from '@/api/user'
import { useSessionStore } from '@/store'

definePage({
  style: {
    navigationBarTitleText: '登录',
  },
})

const sessionStore = useSessionStore()

const form = reactive({
  username: '',
  password: '',
})
const loading = ref(false)

const handleLogin = async () => {
  if (!form.username || !form.password) {
    uni.showToast({ icon: 'none', title: '请输入用户名和密码' })
    return
  }
  loading.value = true
  try {
    await loginUser(form)
    const users = await fetchUsers()
    const matched = users.find(user => user.username === form.username)
    if (matched) {
      sessionStore.setCurrentUser(matched)
    }
    uni.showToast({ icon: 'none', title: '登录成功' })
    setTimeout(() => {
      uni.navigateBack()
    }, 300)
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '登录失败' })
  }
  finally {
    loading.value = false
  }
}

const goRegister = () => {
  uni.navigateTo({ url: '/pages-fg/login/register' })
}
</script>

<template>
  <view class="min-h-screen bg-gray-50 px-4 py-10">
    <view class="rounded-2xl bg-white p-5 shadow">
      <view class="text-xl font-semibold text-gray-800">
        登录账号
      </view>
      <view class="mt-6 space-y-4">
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">用户名</text>
          <input
            v-model="form.username"
            class="mt-1 text-sm"
            placeholder="请输入用户名"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">密码</text>
          <input
            v-model="form.password"
            class="mt-1 text-sm"
            type="password"
            placeholder="请输入密码"
          >
        </view>
        <button type="primary" :loading="loading" @click="handleLogin">
          登录
        </button>
      </view>
      <view class="mt-4 text-center text-xs text-gray-500">
        还没有账号？
        <text class="text-emerald-600" @click="goRegister">
          立即注册
        </text>
      </view>
    </view>
  </view>
</template>
