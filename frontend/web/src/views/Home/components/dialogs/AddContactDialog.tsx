import { useRef, useState } from 'react'
import { useImmer } from 'use-immer'
import { X, UserPlus, Search, Send } from 'lucide-react'
import { toast } from 'sonner'

import { useClickOutside } from '../../hooks/useClickOutside'
import type { SearchFriendParams, User } from '@/modules'

interface AddContactDialogProps {
  open: boolean
  onSearch: (params: SearchFriendParams) => Promise<User>
  onClose: () => void
  onSubmit: (id: string, message?: string) => Promise<void>
}

const ACCENT_COLORS = [
  '#E46342',
  '#0ea5e9',
  '#22c55e',
  '#f59e0b',
  '#8b5cf6',
  '#ec4899',
  '#14b8a6',
  '#f97316',
  '#6366f1',
  '#84cc16',
]

export function AddContactDialog({ open, onSearch, onClose, onSubmit }: AddContactDialogProps) {
  const dialogRef = useRef<HTMLDivElement>(null)
  useClickOutside(dialogRef, onClose, open)
  const [searchParams, setSearchParams] = useImmer<SearchFriendParams>({ id: '' })
  const [message, setMessage] = useImmer<string>('')
  const [searchedUser, setSearchedUser] = useState<User | null>(null)
  const [isSearching, setIsSearching] = useState(false)

  if (!open) return null

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!searchParams.id.trim()) {
      toast.error('请输入用户ID')
      return
    }

    setIsSearching(true)
    try {
      const user = await onSearch(searchParams)
      setSearchedUser(user)
      toast.success('搜索成功')
    } catch {
      toast.error('搜索失败，请检查用户ID是否正确')
      setSearchedUser(null)
    } finally {
      setIsSearching(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!searchedUser) {
      toast.error('请先搜索用户')
      return
    }

    await onSubmit(searchedUser.id, message || undefined)
    toast.success('已发送添加好友请求')

    // 重置状态
    setSearchParams({ id: '' })
    setMessage('')
    setSearchedUser(null)
    onClose()
  }

  const getAccentColor = (id: string) => {
    const index = parseInt(id.slice(-1) || '0', 16) % ACCENT_COLORS.length
    return ACCENT_COLORS[index]
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div
        ref={dialogRef}
        className="relative w-full max-w-2xl rounded-3xl bg-white p-8 shadow-2xl animate-in fade-in zoom-in-95 duration-200"
      >
        {/* Header */}
        <div className="mb-6 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="flex size-12 items-center justify-center rounded-2xl bg-gradient-to-br from-[#E46342] to-[#E46342]/80 text-white shadow-lg">
              <UserPlus className="size-6" />
            </div>
            <div>
              <h2 className="text-2xl font-bold text-slate-900">添加联系人</h2>
              <p className="text-sm text-slate-500">搜索并添加新的联系人</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="flex size-10 items-center justify-center rounded-xl text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
          >
            <X className="size-5" />
          </button>
        </div>

        {/* Search Form */}
        <form onSubmit={handleSearch} className="mb-6">
          <div className="flex gap-3">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-5 text-slate-400" />
              <input
                type="text"
                value={searchParams.id}
                onChange={(e) => setSearchParams({ id: e.target.value })}
                placeholder="输入用户ID进行搜索"
                className="w-full rounded-xl border border-slate-200 bg-slate-50 py-3 pl-10 pr-4 text-slate-900 placeholder:text-slate-400 focus:border-[#E46342] focus:bg-white focus:outline-none focus:ring-2 focus:ring-[#E46342]/20 transition"
              />
            </div>
            <button
              type="submit"
              disabled={isSearching || !searchParams.id.trim()}
              className="px-6 py-3 rounded-xl bg-gradient-to-r from-[#E46342] to-[#E46342]/90 text-white font-medium shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed transition"
            >
              {isSearching ? '搜索中...' : '搜索'}
            </button>
          </div>
        </form>

        {/* 搜索结果 */}
        {searchedUser && (
          <div className="mb-6 rounded-2xl border border-slate-200 bg-gradient-to-br from-slate-50 to-white p-6">
            <div className="flex items-start gap-4">
              <div
                className="flex size-16 items-center justify-center rounded-2xl text-white text-2xl font-bold shadow-lg"
                style={{ backgroundColor: getAccentColor(searchedUser.id) }}
              >
                {searchedUser.nickname?.[0]?.toUpperCase() ||
                  searchedUser.username?.[0]?.toUpperCase() ||
                  'U'}
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-slate-900 mb-1">
                  {searchedUser.nickname || searchedUser.username}
                </h3>
                <p className="text-sm text-slate-500 mb-2">ID: {searchedUser.id}</p>
                {searchedUser.signature && (
                  <p className="text-sm text-slate-600 italic">"{searchedUser.signature}"</p>
                )}
              </div>
            </div>
          </div>
        )}

        {/* 验证消息输入 */}
        {searchedUser && (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="mb-2 block text-sm font-medium text-slate-700">
                验证消息（可选）
              </label>
              <textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="请输入验证消息..."
                rows={3}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-3 text-slate-900 placeholder:text-slate-400 focus:border-[#E46342] focus:bg-white focus:outline-none focus:ring-2 focus:ring-[#E46342]/20 transition resize-none"
              />
            </div>

            <div className="flex gap-3">
              <button
                type="button"
                onClick={onClose}
                className="flex-1 px-6 py-3 rounded-xl border border-slate-200 bg-white text-slate-700 font-medium hover:bg-slate-50 transition"
              >
                取消
              </button>
              <button
                type="submit"
                className="flex-1 px-6 py-3 rounded-xl bg-gradient-to-r from-[#E46342] to-[#E46342]/90 text-white font-medium shadow-lg hover:shadow-xl transition flex items-center justify-center gap-2"
              >
                <Send className="size-4" />
                发送请求
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  )
}
