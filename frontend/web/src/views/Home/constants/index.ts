export const THEME_COLOR = '#E46342'

export const SPECIAL_KEYS = {
  newFriends: 'special:new-friends',
  savedGroups: 'special:saved-groups',
  blacklist: 'special:blacklist',
} as const

export type SpecialKeyType = (typeof SPECIAL_KEYS)[keyof typeof SPECIAL_KEYS]
