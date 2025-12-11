import type { ContactItem, ConversationItem } from '../types'

export function filterConversations(conversations: ConversationItem[], keyword: string): ConversationItem[] {
  if (!keyword.trim()) return conversations
  const lowerKeyword = keyword.trim().toLowerCase()
  return conversations.filter(
    (item) =>
      item.name.toLowerCase().includes(lowerKeyword) || item.lastMessage.toLowerCase().includes(lowerKeyword)
  )
}

export function filterContacts(contacts: ContactItem[], keyword: string): ContactItem[] {
  if (!keyword.trim()) return contacts
  const lowerKeyword = keyword.trim().toLowerCase()
  return contacts.filter((contact) => {
    const text = `${contact.name}${contact.title}${contact.department}${contact.tags?.join('') ?? ''}`.toLowerCase()
    return text.includes(lowerKeyword)
  })
}

export function groupContactsByInitial(contacts: ContactItem[]): [string, ContactItem[]][] {
  const groups: Record<string, ContactItem[]> = {}
  
  contacts.forEach((contact) => {
    const key = contact.name.slice(0, 1).toUpperCase()
    if (!groups[key]) {
      groups[key] = []
    }
    groups[key].push(contact)
  })
  
  return Object.entries(groups)
    .map(([key, list]) => [key, list.sort((a, b) => a.name.localeCompare(b.name))] as [string, ContactItem[]])
    .sort(([a], [b]) => a.localeCompare(b))
}
