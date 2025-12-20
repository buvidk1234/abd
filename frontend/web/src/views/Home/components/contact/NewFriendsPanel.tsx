import { ArrowLeft, ArrowUpRight, ArrowDownLeft, CheckCircle, Clock, Ban } from 'lucide-react'
import { toast } from 'sonner'

import { useUserStore } from '@/store/userStore'
import type { FriendRequest } from '@/modules'
import { Badge, IconButton } from '../common'

interface NewFriendsPanelProps {
  themeColor: string
  applies: FriendRequest[]
  loading: boolean
  isMobile: boolean
  onBack: () => void
  onAccept?: (applyId: string) => Promise<void>
  onReject?: (applyId: string) => Promise<void>
}

export function NewFriendsPanel({
  themeColor,
  applies,
  loading,
  isMobile,
  onBack,
  onAccept,
  onReject,
}: NewFriendsPanelProps) {
  const currentUserId = useUserStore((state) => state.user?.id)

  return (
    <div className="flex min-w-0 flex-1 flex-col bg-slate-50">
      {/* Header */}
      <div className="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-4 shadow-sm md:px-6">
        <div className="flex items-center gap-3">
          {isMobile && (
            <IconButton icon={<ArrowLeft className="size-5" />} onClick={onBack} ariaLabel="返回" />
          )}
          <div className="text-base font-semibold text-slate-900">新的朋友</div>
        </div>
        <Badge variant="primary">{applies.filter((a) => a.handleResult === 0).length} 待处理</Badge>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto px-4 py-6 md:px-8">
        {loading && (
          <div className="flex items-center justify-center py-12">
            <div className="text-center">
              <div className="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-b-2 border-[#E46342]"></div>
              <p className="text-sm text-slate-500">加载中...</p>
            </div>
          </div>
        )}

        {!loading && applies.length === 0 && (
          <div className="flex items-center justify-center py-12">
            <p className="text-sm text-slate-500">暂无好友申请</p>
          </div>
        )}

        {!loading && applies.length > 0 && (
          <div className="space-y-3">
            {applies.map((apply) => {
              // 判断是发送还是接收
              const isSent = apply.fromUser.id === currentUserId
              const otherUser = isSent ? apply.toUser : apply.fromUser
              const direction = isSent ? 'sent' : 'received'

              return (
                <div
                  key={apply.id}
                  className="flex items-start gap-4 rounded-2xl border border-slate-100 bg-white px-4 py-4 shadow-sm transition hover:shadow-md"
                >
                  {/* 头像 */}
                  <div className="relative">
                    <div
                      className="flex size-12 items-center justify-center rounded-xl text-white text-sm font-bold shadow-sm"
                      style={{
                        background: `linear-gradient(135deg, ${themeColor}, ${themeColor}dd)`,
                      }}
                    >
                      {otherUser.nickname?.[0]?.toUpperCase() ||
                        otherUser.username?.[0]?.toUpperCase() ||
                        'U'}
                    </div>
                    {/* 方向图标 */}
                    <div
                      className="absolute -bottom-1 -right-1 flex size-5 items-center justify-center rounded-full bg-white shadow-md"
                      title={direction === 'sent' ? '我发送的请求' : '收到的请求'}
                    >
                      {direction === 'sent' ? (
                        <ArrowUpRight className="size-3 text-blue-500" />
                      ) : (
                        <ArrowDownLeft className="size-3 text-green-500" />
                      )}
                    </div>
                  </div>

                  {/* 信息 */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <div className="font-semibold text-slate-900 truncate">
                        {otherUser.nickname || otherUser.username}
                      </div>
                      {direction === 'sent' && (
                        <span className="text-xs text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full">
                          我发送的
                        </span>
                      )}
                    </div>
                    <div className="text-xs text-slate-500 mb-1">ID: {otherUser.id}</div>
                    {apply.message && (
                      <div className="text-sm text-slate-600 mb-2 line-clamp-2">
                        {apply.message}
                      </div>
                    )}
                    <div className="text-[11px] text-slate-400">{apply.createdAt}</div>
                  </div>

                  {/* 状态/操作按钮 */}
                  <div className="flex flex-col items-end gap-2">
                    {apply.handleResult === 0 ? (
                      // 待处理
                      direction === 'received' ? (
                        // 收到的请求，显示同意/拒绝按钮
                        <div className="flex gap-2">
                          <button
                            onClick={async () => {
                              if (onReject) {
                                try {
                                  await onReject(apply.id)
                                  toast.success('已拒绝')
                                } catch (error) {
                                  console.error('拒绝失败', error)
                                  toast.error('拒绝失败')
                                }
                              }
                            }}
                            className="rounded-lg bg-slate-100 px-3 py-1.5 text-xs font-medium text-slate-700 hover:bg-slate-200 transition"
                          >
                            拒绝
                          </button>
                          <button
                            onClick={async () => {
                              if (onAccept) {
                                try {
                                  await onAccept(apply.id)
                                  toast.success('已同意，成为好友')
                                } catch (error) {
                                  console.error('同意失败', error)
                                  toast.error('同意失败')
                                }
                              }
                            }}
                            className="rounded-lg px-3 py-1.5 text-xs font-medium text-white hover:shadow-lg transition"
                            style={{
                              background: `linear-gradient(to right, ${themeColor}, ${themeColor}dd)`,
                            }}
                          >
                            同意
                          </button>
                        </div>
                      ) : (
                        // 发送的请求，显示等待状态
                        <StatusPill status="pending" text="等待对方确认" themeColor={themeColor} />
                      )
                    ) : apply.handleResult === 1 ? (
                      <StatusPill status="accepted" text="已通过" themeColor={themeColor} />
                    ) : (
                      <StatusPill status="rejected" text="已拒绝" themeColor={themeColor} />
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}

function StatusPill({
  status,
  text,
  themeColor,
}: {
  status: 'pending' | 'accepted' | 'rejected'
  text: string
  themeColor: string
}) {
  const icon =
    status === 'pending' ? (
      <Clock className="size-3" />
    ) : status === 'accepted' ? (
      <CheckCircle className="size-3" style={{ color: themeColor }} />
    ) : (
      <Ban className="size-3" />
    )

  return (
    <span
      className={`flex items-center gap-1 rounded-full px-3 py-1 text-xs font-semibold ${
        status === 'pending'
          ? 'bg-amber-50 text-amber-600'
          : status === 'accepted'
            ? 'bg-emerald-50 text-emerald-600'
            : 'bg-slate-100 text-slate-600'
      }`}
    >
      {icon}
      {text}
    </span>
  )
}
