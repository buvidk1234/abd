import { ArrowLeft, MessageCircle, Phone, ShieldCheck, Star } from 'lucide-react'
import { toast } from 'sonner'

import { ActionButton, Avatar, Badge, IconButton } from '../common'
import type { Friend } from '@/modules'

interface FriendDetailProps {
  themeColor: string
  friend: Friend
  isMobile: boolean
  onBack: () => void
}

export function FriendDetail({ themeColor, friend, isMobile, onBack }: FriendDetailProps) {
  if (!friend) {
    return (
      <div className="flex min-w-0 flex-1 items-center justify-center bg-slate-50">
        <div className="text-center text-sm text-slate-500">选择联系人查看详情</div>
      </div>
    )
  }

  return (
    <div className="flex min-w-0 flex-1 flex-col bg-slate-50">
      <div className="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-4 shadow-sm md:px-6">
        <div className="flex items-center gap-3">
          {isMobile && (
            <IconButton icon={<ArrowLeft className="size-5" />} onClick={onBack} ariaLabel="返回" />
          )}
          <div className="text-base font-semibold text-slate-900">联系人信息</div>
          <Badge variant="primary" size="sm">
            内部
          </Badge>
        </div>
        <div className="flex items-center gap-2">
          <ActionButton
            icon={<MessageCircle className="size-4" />}
            label="发消息"
            themeColor={themeColor}
            onClick={() => toast.info('消息功能稍后接入')}
          />
          <ActionButton
            icon={<Phone className="size-4" />}
            label="语音"
            themeColor={themeColor}
            onClick={() => toast.info('语音功能稍后接入')}
          />
        </div>
      </div>

      <div className="flex flex-1 flex-col gap-6 overflow-y-auto px-4 py-6 md:px-8">
        <div className="flex flex-col gap-4 rounded-3xl bg-white p-5 shadow-sm md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-4">
            <Avatar
              name={friend.friendUser.nickname}
              avatar={friend.friendUser.avatar}
              accent={themeColor}
              status={'online'} // TODO: 根据好友状态显示
              themeColor={themeColor}
              size="xl"
            />
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-2">
                <div className="truncate text-lg font-semibold text-slate-900">
                  {friend.friendUser.nickname}
                </div>
                <Badge variant="default" size="md">
                  {friend.remark}
                </Badge>
              </div>
              <div className="text-sm text-slate-500">{friend.friendUser.signature}</div>
            </div>
          </div>
          <div className="flex flex-wrap gap-2">
            <ActionButton
              icon={<Star className="size-4" />}
              label="设为星标"
              themeColor={themeColor}
              variant="ghost"
              onClick={() => toast.success('已添加星标')}
            />
            <ActionButton
              icon={<ShieldCheck className="size-4" />}
              label="备注/分组"
              themeColor={themeColor}
              variant="ghost"
              onClick={() => toast.info('稍后支持自定义分组')}
            />
          </div>
        </div>

        {/* TODO: other info */}
        {/* <div className="grid gap-4 md:grid-cols-2">
          <InfoCard
            icon={<Mail className="size-4" />}
            label="邮箱"
            value={friend.friendUser.email}
            themeColor={themeColor}
          />
          <InfoCard
            icon={<Phone className="size-4" />}
            label="电话"
            value={friend.phone}
            themeColor={themeColor}
          />
          <InfoCard
            icon={<MapPin className="size-4" />}
            label="地点"
            value={friend.location || '远程/协作'}
            themeColor={themeColor}
          />
          <InfoCard
            icon={<Tag className="size-4" />}
            label="标签"
            value={friend.tags?.join(' / ') || '暂无标签'}
            themeColor={themeColor}
          />
        </div> */}

        <div className="rounded-3xl border border-dashed border-slate-200 bg-white/60 p-5">
          <div className="text-sm font-semibold text-slate-900">备注</div>
          <p className="mt-2 text-sm leading-6 text-slate-600">{friend.remark}</p>
        </div>
      </div>
    </div>
  )
}

function InfoCard({
  icon,
  label,
  value,
  themeColor,
}: {
  icon: React.ReactNode
  label: string
  value: string
  themeColor: string
}) {
  return (
    <div className="flex items-start gap-3 rounded-3xl border border-slate-100 bg-white px-4 py-3 shadow-sm">
      <div
        className="mt-0.5 flex size-9 items-center justify-center rounded-2xl bg-slate-100 text-slate-600"
        style={{ color: themeColor }}
      >
        {icon}
      </div>
      <div className="min-w-0">
        <div className="text-xs font-semibold text-slate-500">{label}</div>
        <div className="truncate text-sm font-semibold text-slate-900">{value}</div>
      </div>
    </div>
  )
}
