import { useParams, useNavigate } from 'react-router'
import { useImmer } from 'use-immer'

import { THEME_COLOR, SPECIAL_KEYS } from '../constants'
import { ContactsSidebar } from '../components/contact/ContactsSidebar'
import { FriendDetail } from '../components/contact/ContactDetail'
import { ContactSpecialPanel } from '../components/contact/ContactSpecialPanel'
import { NewFriendsPanel } from '../components/contact/NewFriendsPanel'
import { useMockData } from '../hooks/useMockData'
import { useResponsive } from '../hooks/useResponsive'
import { AddContactDialog } from '../components/dialogs/AddContactDialog'
import { applyToAddFriend, searchFriend } from '@/modules'
import type { FriendRequestItem } from '../types'
import { useFriendApplies } from '../hooks'
import { useFriends } from '../hooks/useFriends'

interface UIState {
  showSidebarOnMobile: boolean
  showAddDialog: boolean
}

interface LayoutContext {
  onOpenGlobalSearch: () => void
}

export function ContactPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isMobile = useResponsive(1024)
  const { applies, loading, handleApply, refresh } = useFriendApplies()
  const { friends } = useFriends()

  const { savedGroups, blacklist } = useMockData()
  const selectedFriend = friends.find((f) => f.friendUser.id === id) ?? null

  const newFriends: FriendRequestItem[] = applies.map((apply) => ({
    id: apply.id,
    from: apply.fromUser.nickname || apply.fromUser.username,
    note: apply.message,
    time: apply.createdAt,
    status:
      apply.handleResult === 0 ? 'pending' : apply.handleResult === 1 ? 'accepted' : 'rejected',
  }))

  const [uiState, setUIState] = useImmer<UIState>({
    showSidebarOnMobile: true,
    showAddDialog: false,
  })

  // Sync mobile state
  if (isMobile && !uiState.showSidebarOnMobile && !id) {
    setUIState((draft) => {
      draft.showSidebarOnMobile = true
    })
  }

  const handleSelectFriend = (friendId: string) => {
    navigate(`/contact/${friendId}`)
    if (isMobile) {
      setUIState((draft) => {
        draft.showSidebarOnMobile = false
      })
    }
  }

  // Render content based on selection
  const renderContent = () => {
    if (id === SPECIAL_KEYS.newFriends) {
      return (
        <NewFriendsPanel
          themeColor={THEME_COLOR}
          applies={applies}
          loading={loading}
          isMobile={isMobile}
          onBack={() =>
            setUIState((draft) => {
              draft.showSidebarOnMobile = true
            })
          }
          onAccept={(applyId) => handleRespondApply(applyId, 'accept')}
          onReject={(applyId) => handleRespondApply(applyId, 'reject')}
        />
      )
    }

    if (id === SPECIAL_KEYS.savedGroups) {
      return (
        <ContactSpecialPanel
          type="saved-groups"
          themeColor={THEME_COLOR}
          savedGroups={savedGroups}
          onBack={() =>
            setUIState((draft) => {
              draft.showSidebarOnMobile = true
            })
          }
        />
      )
    }

    if (id === SPECIAL_KEYS.blacklist) {
      return (
        <ContactSpecialPanel
          type="blacklist"
          themeColor={THEME_COLOR}
          blacklist={blacklist}
          onBack={() =>
            setUIState((draft) => {
              draft.showSidebarOnMobile = true
            })
          }
        />
      )
    }

    return (
      <FriendDetail
        themeColor={THEME_COLOR}
        friend={selectedFriend!}
        isMobile={isMobile}
        onBack={() =>
          setUIState((draft) => {
            draft.showSidebarOnMobile = true
          })
        }
      />
    )
  }

  const handleAddFriend = async (userId: string, message?: string) => {
    await applyToAddFriend({
      toUserID: userId,
      message: message,
    })
    // 添加成功后刷新申请列表
    await refresh()
  }

  const handleRespondApply = async (applyId: string, action: 'accept' | 'reject') => {
    await handleApply(applyId, action)
  }

  return (
    <div className="relative flex min-w-0 flex-1">
      <ContactsSidebar
        themeColor={THEME_COLOR}
        friends={friends}
        selectedId={id ?? null}
        onSelect={handleSelectFriend}
        onSelectSpecial={handleSelectFriend}
        specialCounts={{
          newFriends: newFriends.filter((i) => i.status === 'pending').length,
          savedGroups: savedGroups.length,
          blacklist: blacklist.length,
        }}
        showSidebar={uiState.showSidebarOnMobile}
        isMobile={isMobile}
        onOpenAddDialog={() =>
          setUIState((draft) => {
            draft.showAddDialog = true
          })
        }
      />

      {renderContent()}

      <AddContactDialog
        open={uiState.showAddDialog}
        onSearch={searchFriend}
        onClose={() =>
          setUIState((draft) => {
            draft.showAddDialog = false
          })
        }
        onSubmit={handleAddFriend}
      />
    </div>
  )
}
