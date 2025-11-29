<script lang="ts" setup>
import type { FriendInfo, FriendRequestInfo, UserProfile } from '@/types/backend'
import { onLoad } from '@dcloudio/uni-app'
import { computed, reactive, ref, watch } from 'vue'
import { applyToAddFriend, deleteFriend, getFriendList, getIncomingFriendRequests, getOutgoingFriendRequests, respondFriendApply } from '@/api/friend'
import { fetchUsers } from '@/api/user'
import { useSessionStore } from '@/store'

defineOptions({
  name: 'Home',
})
definePage({
  type: 'home',
  style: {
    navigationStyle: 'default',
    navigationBarTitleText: '好友中心',
  },
})

const sessionStore = useSessionStore()
const currentUser = computed(() => sessionStore.currentUser)

const addFriendForm = reactive({
  toUserID: '',
  message: '',
})

const allUsers = ref<UserProfile[]>([])
const loadingUsers = ref(false)

const friendList = ref<FriendInfo[]>([])
const loadingFriends = ref(false)

const incomingRequests = ref<FriendRequestInfo[]>([])
const outgoingRequests = ref<FriendRequestInfo[]>([])
const loadingRequests = ref(false)

function ensureUserSelected() {
  if (!currentUser.value?.userID) {
    uni.showToast({
      icon: 'none',
      title: '请先登录或选择一个用户',
    })
    return false
  }
  return true
}

async function loadUsers() {
  loadingUsers.value = true
  try {
    allUsers.value = await fetchUsers()
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '获取用户列表失败' })
  }
  finally {
    loadingUsers.value = false
  }
}

async function refreshFriends() {
  if (!ensureUserSelected()) {
    return
  }
  loadingFriends.value = true
  try {
    const res = await getFriendList({
      userID: currentUser.value!.userID,
      page: 1,
      pageSize: 50,
    })
    friendList.value = res?.friends || []
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '获取好友失败' })
  }
  finally {
    loadingFriends.value = false
  }
}

async function refreshRequests() {
  if (!ensureUserSelected()) {
    return
  }
  loadingRequests.value = true
  try {
    const incoming = await getIncomingFriendRequests({
      toUserID: currentUser.value!.userID,
      page: 1,
      pageSize: 50,
    })
    incomingRequests.value = incoming?.list || []
    const outgoing = await getOutgoingFriendRequests({
      fromUserID: currentUser.value!.userID,
      page: 1,
      pageSize: 50,
    })
    outgoingRequests.value = outgoing?.list || []
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '获取申请失败' })
  }
  finally {
    loadingRequests.value = false
  }
}

async function handleAddFriend() {
  if (!ensureUserSelected()) {
    return
  }
  if (!addFriendForm.toUserID) {
    uni.showToast({ icon: 'none', title: '请输入好友ID' })
    return
  }
  if (addFriendForm.toUserID === currentUser.value!.userID) {
    uni.showToast({ icon: 'none', title: '不能添加自己为好友' })
    return
  }
  try {
    await applyToAddFriend({
      fromUserID: currentUser.value!.userID,
      toUserID: addFriendForm.toUserID.trim(),
      message: addFriendForm.message,
    })
    uni.showToast({ icon: 'none', title: '好友申请已发送' })
    addFriendForm.message = ''
    refreshRequests()
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '发送失败' })
  }
}

async function handleRespond(id: number, result: number, handleMsg?: string) {
  if (!ensureUserSelected()) {
    return
  }
  try {
    await respondFriendApply({
      id,
      handlerUserID: currentUser.value!.userID,
      handleResult: result,
      handleMsg: handleMsg || '',
    })
    uni.showToast({ icon: 'none', title: result === 1 ? '已同意' : '已拒绝' })
    refreshFriends()
    refreshRequests()
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '操作失败' })
  }
}

async function handleDeleteFriend(friendUserID: string) {
  if (!ensureUserSelected()) {
    return
  }
  try {
    await deleteFriend({
      ownerUserID: currentUser.value!.userID,
      friendUserID,
    })
    uni.showToast({ icon: 'none', title: '好友已删除' })
    refreshFriends()
  }
  catch (error) {
    console.error(error)
    uni.showToast({ icon: 'none', title: '删除失败' })
  }
}

function pickUser(user: UserProfile) {
  sessionStore.setCurrentUser(user)
  uni.showToast({ icon: 'none', title: '已切换当前用户' })
}

function formatStatus(status: number) {
  if (status === 1) {
    return '已同意'
  }
  if (status === 2) {
    return '已拒绝'
  }
  return '待处理'
}

watch(
  () => currentUser.value?.userID,
  (id, oldID) => {
    if (id && id !== oldID) {
      refreshFriends()
      refreshRequests()
    }
    if (!id) {
      friendList.value = []
      incomingRequests.value = []
      outgoingRequests.value = []
    }
  },
)

onLoad(() => {
  if (currentUser.value?.userID) {
    refreshFriends()
    refreshRequests()
  }
  loadUsers()
})
</script>

<template>
  <view class="min-h-screen from-emerald-50 to-white bg-gradient-to-b px-4 py-5">
    <view class="rounded-3xl bg-white/80 p-4 shadow-lg backdrop-blur">
      <view class="flex items-center justify-between">
        <view>
          <view class="text-lg text-gray-800 font-semibold">
            {{ currentUser?.nickname || currentUser?.username || '未登录' }}
          </view>
          <view class="text-sm text-gray-500">
            {{ currentUser?.userID ? `ID: ${currentUser.userID}` : '请先登录或选择一个用户' }}
          </view>
        </view>
        <view class="flex gap-2">
          <button size="mini" @click="uni.navigateTo({ url: '/pages-fg/login/login' })">
            登录
          </button>
          <button type="primary" size="mini" @click="uni.navigateTo({ url: '/pages-fg/login/register' })">
            注册
          </button>
        </view>
      </view>
    </view>

    <view class="mt-4 rounded-2xl bg-white p-4 shadow">
      <view class="flex items-center justify-between">
        <text class="text-base text-gray-800 font-semibold">用户目录</text>
        <button size="mini" :loading="loadingUsers" @click="loadUsers">
          刷新
        </button>
      </view>
      <view class="grid grid-cols-1 mt-3 gap-3">
        <view
          v-for="item in allUsers"
          :key="item.userID"
          class="border border-emerald-50 rounded-xl bg-emerald-50/60 p-3"
        >
          <view class="flex items-center justify-between">
            <view>
              <view class="text-sm text-gray-800 font-semibold">
                {{ item.nickname || item.username || '未命名用户' }}
              </view>
              <view class="mt-1 text-xs text-gray-500">
                {{ item.userID }}
              </view>
            </view>
            <button size="mini" @click="pickUser(item)">
              设为当前
            </button>
          </view>
        </view>
        <view v-if="!allUsers.length && !loadingUsers" class="text-center text-xs text-gray-400">
          暂无用户，先注册一个吧
        </view>
      </view>
    </view>

    <view class="mt-4 rounded-2xl bg-white p-4 shadow">
      <view class="text-base text-gray-800 font-semibold">
        添加好友
      </view>
      <view class="mt-3 space-y-3">
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">好友用户ID</text>
          <input
            v-model="addFriendForm.toUserID"
            class="mt-1 text-sm"
            placeholder="请输入要添加的用户ID"
          >
        </view>
        <view class="rounded-xl bg-gray-50 px-3 py-2">
          <text class="text-xs text-gray-500">备注</text>
          <textarea
            v-model="addFriendForm.message"
            class="mt-1 h-16 text-sm"
            placeholder="简单说明一下你是谁"
            auto-height
          />
        </view>
        <button type="primary" @click="handleAddFriend">
          发送好友申请
        </button>
      </view>
    </view>

    <view class="mt-4 rounded-2xl bg-white p-4 shadow">
      <view class="flex items-center justify-between">
        <text class="text-base text-gray-800 font-semibold">我的好友</text>
        <button size="mini" :loading="loadingFriends" @click="refreshFriends">
          刷新
        </button>
      </view>
      <view class="mt-3 space-y-3">
        <view
          v-for="item in friendList"
          :key="item.friendUser.userID"
          class="border border-gray-100 rounded-xl p-3"
        >
          <view class="flex items-center justify-between">
            <view>
              <view class="text-sm text-gray-800 font-semibold">
                {{ item.friendUser.nickname || item.friendUser.username }}
              </view>
              <view class="text-xs text-gray-500">
                {{ item.friendUser.userID }}
              </view>
            </view>
            <button size="mini" type="warn" plain @click="handleDeleteFriend(item.friendUser.userID)">
              删除
            </button>
          </view>
          <view class="mt-2 text-xs text-gray-400">
            备注：{{ item.remark || '暂无备注' }}
          </view>
        </view>
        <view v-if="!friendList.length && !loadingFriends" class="text-center text-xs text-gray-400">
          暂无好友
        </view>
      </view>
    </view>

    <view class="mt-4 rounded-2xl bg-white p-4 shadow">
      <view class="flex items-center justify-between">
        <text class="text-base text-gray-800 font-semibold">收到的好友申请</text>
        <button size="mini" :loading="loadingRequests" @click="refreshRequests">
          刷新
        </button>
      </view>
      <view class="mt-3 space-y-3">
        <view
          v-for="item in incomingRequests"
          :key="item.id"
          class="border border-gray-100 rounded-xl p-3"
        >
          <view class="flex items-center justify-between">
            <view>
              <view class="text-sm text-gray-800 font-semibold">
                {{ item.fromUser.nickname || item.fromUser.username || '未知用户' }}
              </view>
              <view class="text-xs text-gray-500">
                {{ item.reqMsg || '无附言' }}
              </view>
            </view>
            <view class="flex gap-2">
              <button
                v-if="item.handleResult === 0"
                size="mini"
                type="primary"
                @click="handleRespond(item.id, 1)"
              >
                同意
              </button>
              <button
                v-if="item.handleResult === 0"
                size="mini"
                type="warn"
                plain
                @click="handleRespond(item.id, 2)"
              >
                拒绝
              </button>
              <view v-else class="rounded-full bg-gray-100 px-2 py-1 text-xs text-gray-500">
                {{ formatStatus(item.handleResult) }}
              </view>
            </view>
          </view>
        </view>
        <view v-if="!incomingRequests.length && !loadingRequests" class="text-center text-xs text-gray-400">
          暂无申请
        </view>
      </view>
    </view>

    <view class="my-4 rounded-2xl bg-white p-4 shadow">
      <view class="text-base text-gray-800 font-semibold">
        我发出的申请
      </view>
      <view class="mt-3 space-y-3">
        <view
          v-for="item in outgoingRequests"
          :key="item.id"
          class="border border-gray-100 rounded-xl p-3"
        >
          <view class="flex items-center justify-between">
            <view>
              <view class="text-sm text-gray-800 font-semibold">
                目标：{{ item.toUser.nickname || item.toUser.username || '未知用户' }}
              </view>
              <view class="text-xs text-gray-500">
                状态：{{ formatStatus(item.handleResult) }}
              </view>
            </view>
            <view class="text-xs text-gray-400">
              {{ item.reqMsg || '无备注' }}
            </view>
          </view>
        </view>
        <view v-if="!outgoingRequests.length && !loadingRequests" class="text-center text-xs text-gray-400">
          暂无记录
        </view>
      </view>
    </view>
  </view>
</template>
