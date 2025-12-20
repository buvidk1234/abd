import type { Friend } from '@/modules'
import type { ConversationItem } from '../types'

export function filterConversations(
  conversations: ConversationItem[],
  keyword: string
): ConversationItem[] {
  if (!keyword.trim()) return conversations
  const lowerKeyword = keyword.trim().toLowerCase()
  return conversations.filter(
    (item) =>
      item.name.toLowerCase().includes(lowerKeyword) ||
      item.lastMessage.toLowerCase().includes(lowerKeyword)
  )
}

export function filterFriends(friends: Friend[], keyword: string): Friend[] {
  if (!keyword.trim()) return friends
  const lowerKeyword = keyword.trim().toLowerCase()
  return friends.filter((item) => {
    const text =
      `${item.friendUser.username}${item.remark}${item.friendUser.nickname}${item.friendUser.signature ?? ''}`.toLowerCase()
    return text.includes(lowerKeyword)
  })
}

export function groupFriendsByInitial(friends: Friend[]): [string, Friend[]][] {
  const groups: Record<string, Friend[]> = {}

  friends.forEach((item) => {
    const key = item.friendUser.username.slice(0, 1).toUpperCase()
    if (!groups[key]) {
      groups[key] = []
    }
    groups[key].push(item)
  })

  return Object.entries(groups)
    .map(
      ([key, list]) =>
        [key, list.sort((a, b) => a.friendUser.username.localeCompare(b.friendUser.username))] as [
          string,
          Friend[],
        ]
    )
    .sort(([a], [b]) => a.localeCompare(b))
}
