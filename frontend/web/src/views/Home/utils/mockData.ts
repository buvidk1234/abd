import type {
  BlacklistItem,
  ContactItem,
  ConversationItem,
  FriendRequestItem,
  SavedGroupItem,
} from '../types'

const BASE_CONVERSATIONS: ConversationItem[] = [
  {
    id: 'wk-001',
    name: '产品讨论组',
    title: '需求澄清',
    avatar: 'PD',
    accent: '#f97316',
    lastMessage: '方案已经更新,晚上前确认下哈',
    time: '09:36',
    unread: 4,
    pinned: true,
    reminders: ['下午 2:00 评审'],
    messages: [
      {
        id: 'm-1',
        author: 'Alan',
        content: '新版需求排期已经发到了文档里,帮忙过一眼。',
        timestamp: '09:28',
        direction: 'out',
      },
      {
        id: 'm-2',
        author: 'Lee',
        content: '收到,流程图也同步了么?',
        timestamp: '09:32',
        direction: 'in',
      },
      {
        id: 'm-3',
        author: 'Alan',
        content: '是的,附件里有。我们等下开个 10 分钟的小会对齐。',
        timestamp: '09:36',
        direction: 'out',
      },
    ],
  },
  {
    id: 'wk-002',
    name: '设计同学',
    avatar: 'UX',
    accent: '#0ea5e9',
    lastMessage: '正在标注最新图稿...',
    time: '09:10',
    unread: 2,
    typing: true,
    description: '新配色尝试中',
    messages: [
      {
        id: 'm-4',
        author: 'Iris',
        content: 'Logo 线稿改完了,晚点发你看下。',
        timestamp: '08:58',
        direction: 'in',
      },
      {
        id: 'm-5',
        author: 'Alan',
        content: '好的,标注的时候帮忙把状态颜色也带上。',
        timestamp: '09:02',
        direction: 'out',
      },
    ],
  },
  {
    id: 'wk-003',
    name: '研发群',
    avatar: 'RD',
    accent: '#22c55e',
    lastMessage: 'CI 通过,准备用灰度包',
    time: '昨天',
    unread: 0,
    online: true,
    messages: [
      {
        id: 'm-6',
        author: 'DevOps',
        content: 'CI 已经全绿,网关规则跟主干一致。',
        timestamp: '昨天 21:10',
        direction: 'in',
      },
      {
        id: 'm-7',
        author: 'Alan',
        content: 'OK,灰度到 10% 先看回放。',
        timestamp: '昨天 21:12',
        direction: 'out',
      },
    ],
  },
  {
    id: 'wk-004',
    name: '文件助手',
    avatar: 'FH',
    accent: '#a855f7',
    lastMessage: '周报模板已保存',
    time: '周一',
    unread: 0,
    muted: true,
    messages: [
      {
        id: 'm-8',
        author: '助手',
        content: '周报模板和素材都在这里,有需要再叫我~',
        timestamp: '周一 10:21',
        direction: 'in',
      },
    ],
  },
  {
    id: 'wk-005',
    name: '市场同学',
    avatar: 'MK',
    accent: '#e11d48',
    lastMessage: '明天要给渠道版海报定稿',
    time: '周六',
    unread: 1,
    draft: '渠道报价单对齐中…',
    messages: [
      {
        id: 'm-9',
        author: 'Cathy',
        content: '明天 16:00 前要把渠道版海报定稿,别忘了。',
        timestamp: '周六 12:10',
        direction: 'in',
      },
    ],
  },
  {
    id: 'wk-006',
    name: '唐僧叨叨团队',
    title: '内部通知',
    avatar: 'TS',
    accent: '#e46342',
    lastMessage: '本周例会延迟 30 分钟',
    time: '周五',
    unread: 0,
    messages: [
      {
        id: 'm-10',
        author: '系统',
        content: '本周例会延迟 30 分钟,请提前准备周报。',
        timestamp: '周五 09:00',
        direction: 'in',
      },
    ],
  },
]

const BASE_CONTACTS: ContactItem[] = [
  {
    id: 'c-001',
    name: '林浅',
    title: '产品经理',
    avatar: 'LQ',
    accent: '#22c55e',
    department: '中台产品',
    email: 'linqian@tsdd.com',
    phone: '138-0000-1001',
    tags: ['Owner', '需求确认'],
    status: 'online',
    note: '负责需求澄清与排期跟进',
  },
  {
    id: 'c-002',
    name: '程远',
    title: '前端负责人',
    avatar: 'CY',
    accent: '#0ea5e9',
    department: 'FE 团队',
    email: 'chengyuan@tsdd.com',
    phone: '138-0000-1002',
    tags: ['Vue', 'React'],
    status: 'busy',
    note: '关注性能与体验统一',
  },
  {
    id: 'c-003',
    name: '顾冉',
    title: 'UI/UX 设计',
    avatar: 'GR',
    accent: '#e11d48',
    department: '设计体验',
    email: 'guran@tsdd.com',
    phone: '138-0000-1003',
    tags: ['交互', '组件规范'],
    location: '上海',
  },
  {
    id: 'c-004',
    name: '王希',
    title: '后端工程师',
    avatar: 'WX',
    accent: '#f97316',
    department: 'IM 服务',
    email: 'wangxi@tsdd.com',
    phone: '138-0000-1004',
    tags: ['Go', 'IM'],
    status: 'online',
  },
  {
    id: 'c-005',
    name: '张语',
    title: '测试工程师',
    avatar: 'ZY',
    accent: '#a855f7',
    department: '质量保障',
    email: 'zhangyu@tsdd.com',
    phone: '138-0000-1005',
    tags: ['自动化', '回归'],
  },
  {
    id: 'c-006',
    name: '陈曦',
    title: '市场运营',
    avatar: 'CX',
    accent: '#10b981',
    department: '市场',
    email: 'chenxi@tsdd.com',
    phone: '138-0000-1006',
    tags: ['活动', '渠道'],
  },
  {
    id: 'c-007',
    name: '刘琪',
    title: '客户成功',
    avatar: 'LQ',
    accent: '#6366f1',
    department: 'CS 团队',
    email: 'liuq@tsdd.com',
    phone: '138-0000-1007',
    tags: ['交付', '培训'],
    status: 'online',
  },
  {
    id: 'c-008',
    name: '邵宁',
    title: '运维工程师',
    avatar: 'SN',
    accent: '#eab308',
    department: '基础设施',
    email: 'shaoning@tsdd.com',
    phone: '138-0000-1008',
    tags: ['监控', '告警'],
    location: '成都',
  },
  {
    id: 'c-009',
    name: '何雨',
    title: '数据分析',
    avatar: 'HY',
    accent: '#14b8a6',
    department: '数据中台',
    email: 'heyu@tsdd.com',
    phone: '138-0000-1009',
    tags: ['看板', '洞察'],
  },
  {
    id: 'c-010',
    name: '苏瑶',
    title: '商务经理',
    avatar: 'SY',
    accent: '#ef4444',
    department: '商务拓展',
    email: 'suyao@tsdd.com',
    phone: '138-0000-1010',
    tags: ['合同', '报价'],
    status: 'offline',
  },
]

const BASE_NEW_FRIENDS: FriendRequestItem[] = [
  { id: 'f-001', from: '李安', note: '一起跟进需求评审', time: '09:12', status: 'pending' },
  { id: 'f-002', from: 'Tracy', note: '加个好友便于发素材', time: '昨天', status: 'pending' },
  { id: 'f-003', from: '刘鑫', note: '同一个项目组', time: '周二', status: 'accepted' },
]

const BASE_SAVED_GROUPS: SavedGroupItem[] = [
  { id: 'g-001', name: 'Web 性能讨论', members: 18, update: '昨天 21:00', accent: '#0ea5e9' },
  { id: 'g-002', name: '体验规范协同', members: 12, update: '周二 10:22', accent: '#22c55e' },
  { id: 'g-003', name: '渠道物料共享', members: 27, update: '周一 16:08', accent: '#e46342' },
]

const BASE_BLACKLIST: BlacklistItem[] = [
  { id: 'b-001', name: '陌生号码 178****9901', reason: '频繁骚扰', time: '08:12' },
  { id: 'b-002', name: '无备注用户', reason: '重复广告', time: '上周五' },
]

export const createConversations = () =>
  BASE_CONVERSATIONS.map((conversation) => ({
    ...conversation,
    messages: conversation.messages.map((message) => ({ ...message })),
    reminders: conversation.reminders ? [...conversation.reminders] : undefined,
  }))

export const createContacts = () =>
  BASE_CONTACTS.map((contact) => ({
    ...contact,
    tags: contact.tags ? [...contact.tags] : undefined,
  }))

export const createNewFriends = () => BASE_NEW_FRIENDS.map((item) => ({ ...item }))

export const createSavedGroups = () => BASE_SAVED_GROUPS.map((item) => ({ ...item }))

export const createBlacklist = () => BASE_BLACKLIST.map((item) => ({ ...item }))
