<script lang="ts" setup>
import { computed, onLoad, reactive, ref, watch } from 'vue'
import { fetchUsers, updateUserInfo } from '@/api/user'
import type { UserProfile } from '@/types/backend'
import { useSessionStore } from '@/store'

definePage({
  style: {
    navigationBarTitleText: '我的',
  },
})

const sessionStore = useSessionStore()
const currentUser = computed(() => sessionStore.currentUser)

const profileForm = reactive({
  nickname: '',
  avatarURL: '',
  signature: '',
  gender: 0,
  phone: '',
  email: '',
})

const genderOptions = [
  { label: '保密', value: 0 },
  { label: '男', value: 1 },
  { label: '女', value: 2 },
]

const allUsers = ref<UserProfile[]>([])
const loadingUsers = ref(false)

const fillForm = (user: UserProfile | null) => {
  profileForm.nickname = user?.nickname || ''
  profileForm.avatarURL = user?.avatarURL || ''
  profileForm.signature = user?.signature || ''
  profileForm.gender = user?.gender || 0
  profileForm.phone = user?.phone || ''
  profileForm.email = user?.email || ''
}

watch(
  () => currentUser.value,
  (val) => {
    fillForm(val || null)
  },
  { immediate: true },
)

const ensureLogin = () => {
  if (!currentUser.value?.userID) {
    uni.showToast({ icon: 'none', title: '请先选择/登录用户' })
    return false
  }
  return true
}

const handleSaveProfile = async () => {
  if (!ensureLogin()) {
    return
  }
  try {
    await updateUserInfo({
      userID: currentUser.value!.userID,
      ...profileForm,
    })
    sessionStore.updateCurrentProfile(profileForm)
    uni.showToast({ icon: 'none', title: '资料已更新' })
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '更新失败' })
  }
}

const loadUsers = async () => {
  loadingUsers.value = true
  try {
    allUsers.value = await fetchUsers()
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '获取用户失败' })
  }
  finally {
    loadingUsers.value = false
  }
}

const pickUser = (user: UserProfile) => {
  sessionStore.setCurrentUser(user)
  uni.showToast({ icon: 'none', title: '已切换用户' })
}

const navigateToLogin = () => {
  uni.navigateTo({ url: '/pages-fg/login/login' })
}
const navigateToRegister = () => {
  uni.navigateTo({ url: '/pages-fg/login/register' })
}

const handleGenderChange = (e: any) => {
  const idx = Number(e.detail.value)
  profileForm.gender = genderOptions[idx].value
}

onLoad(() => {
  loadUsers()
})
</script>

<template>
  <view class="min-h-screen bg-gray-50 px-4 py-5">
    <view class="rounded-3xl bg-white p-4 shadow">
      <view class="flex items-center justify-between">
        <view>
          <view class="text-lg font-semibold text-gray-800">
            {{ currentUser?.nickname || currentUser?.username || '未登录' }}
          </view>
          <view class="mt-1 text-xs text-gray-500">
            {{ currentUser?.userID || '选择一个账号后再进行好友操作' }}
          </view>
        </view>
        <view class="flex gap-2">
          <button size="mini" @click="navigateToLogin">
            去登录
          </button>
          <button type="primary" size="mini" @click="navigateToRegister">
            去注册
          </button>
        </view>
      </view>
      <view v-if="currentUser?.signature" class="mt-2 text-xs text-gray-500">
        {{ currentUser.signature }}
      </view>
    </view>

    <view class="mt-4 rounded-2xl bg-white p-4 shadow">
      <view class="flex items-center justify-between">
        <text class="text-base font-semibold text-gray-800">选择一个现有用户</text>
        <button size="mini" :loading="loadingUsers" @click="loadUsers">
          刷新
        </button>
      </view>
      <view class="mt-3 space-y-3">
        <view
          v-for="item in allUsers"
          :key="item.userID"
          class="rounded-xl border border-gray-100 p-3"
        >
          <view class="flex items-center justify-between">
            <view>
              <view class="text-sm font-semibold text-gray-800">
                {{ item.nickname || item.username || '未命名用户' }}
              </view>
              <view class="text-xs text-gray-500">
                {{ item.userID }}
              </view>
            </view>
            <button size="mini" @click="pickUser(item)">
              切换
            </button>
          </view>
        </view>
        <view v-if="!allUsers.length && !loadingUsers" class="text-center text-xs text-gray-400">
          暂无用户，请先注册
        </view>
      </view>
    </view>

    <view class="my-4 rounded-2xl bg-white p-4 shadow">
      <view class="text-base font-semibold text-gray-800">
        编辑资料
      </view>
      <view class="mt-3 space-y-3">
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">昵称</text>
          <input
            v-model="profileForm.nickname"
            class="mt-1 text-sm"
            placeholder="给自己起个名字"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">头像 URL</text>
          <input
            v-model="profileForm.avatarURL"
            class="mt-1 text-sm"
            placeholder="可不填"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">签名</text>
          <textarea
            v-model="profileForm.signature"
            class="mt-1 h-16 text-sm"
            placeholder="简单介绍一下自己"
            auto-height
          ></textarea>
        </view>
        <view class="flex items-center justify-between rounded-xl bg-gray-50 px-3 py-2">
          <view>
            <text class="text-xs text-gray-500">性别</text>
            <view class="text-sm text-gray-800">
              {{ genderOptions.find(item => item.value === profileForm.gender)?.label || '保密' }}
            </view>
          </view>
          <picker :range="genderOptions" range-key="label" @change="handleGenderChange">
            <view class="text-xs text-emerald-600">
              选择
            </view>
          </picker>
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">手机号</text>
          <input
            v-model="profileForm.phone"
            class="mt-1 text-sm"
            type="number"
            placeholder="可选"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">邮箱</text>
          <input
            v-model="profileForm.email"
            class="mt-1 text-sm"
            placeholder="可选"
          >
        </view>
        <button type="primary" @click="handleSaveProfile">
          保存资料
        </button>
      </view>
    </view>
  </view>
</template>
