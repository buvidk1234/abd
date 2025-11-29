<script lang="ts" setup>
import { reactive, ref } from 'vue'
import { registerUser } from '@/api/user'
import { useSessionStore } from '@/store'

definePage({
  style: {
    navigationBarTitleText: '注册',
  },
})

const sessionStore = useSessionStore()

const form = reactive({
  username: '',
  password: '',
  nickname: '',
  phone: '',
  email: '',
  avatarURL: '',
})
const loading = ref(false)

const handleRegister = async () => {
  if (!form.username || !form.password) {
    uni.showToast({ icon: 'none', title: '用户名和密码必填' })
    return
  }
  loading.value = true
  try {
    const res = await registerUser(form)
    const userID = (res as any)?.userID || ''
    sessionStore.setCurrentUser({
      userID,
      username: form.username,
      nickname: form.nickname || form.username,
      avatarURL: form.avatarURL,
      phone: form.phone,
      email: form.email,
    })
    uni.showToast({ icon: 'none', title: '注册成功' })
    setTimeout(() => {
      uni.navigateBack()
    }, 400)
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '注册失败' })
  }
  finally {
    loading.value = false
  }
}
</script>

<template>
  <view class="min-h-screen bg-gray-50 px-4 py-10">
    <view class="rounded-2xl bg-white p-5 shadow">
      <view class="text-xl font-semibold text-gray-800">
        注册新用户
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
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">昵称</text>
          <input
            v-model="form.nickname"
            class="mt-1 text-sm"
            placeholder="展示用昵称，可选"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">手机号</text>
          <input
            v-model="form.phone"
            class="mt-1 text-sm"
            type="number"
            placeholder="可选"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">邮箱</text>
          <input
            v-model="form.email"
            class="mt-1 text-sm"
            placeholder="可选"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">头像 URL</text>
          <input
            v-model="form.avatarURL"
            class="mt-1 text-sm"
            placeholder="可选"
          >
        </view>
        <button type="primary" :loading="loading" @click="handleRegister">
          注册
        </button>
      </view>
      <view class="mt-4 text-center text-xs text-gray-500">
        注册后会自动切换为当前用户
      </view>
    </view>
  </view>
</template>
