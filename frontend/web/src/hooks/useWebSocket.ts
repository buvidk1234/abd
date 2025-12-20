import { sleep } from '@/views/Home/utils/async'
import { useCallback, useEffect, useRef } from 'react'
import { useImmer } from 'use-immer'

export type WSMessage<T = unknown> = {
  req_identifier: number
  data: T
}

type Handler<T = unknown> = (msg: WSMessage<T>) => void

export function useWebSocketCore(params: { url: string }) {
  const { url } = params

  const wsRef = useRef<WebSocket | null>(null)
  const [readyState, setReadyState] = useImmer<number>(WebSocket.CLOSED)

  // Map<reqId, Set<handler>>
  const listenersRef = useRef<Map<number, Set<Handler>>>(new Map())

  const dispatch = useCallback((msg: WSMessage) => {
    const set = listenersRef.current.get(msg.req_identifier)
    if (!set || set.size === 0) return // 防止 handler 里 unsubscribe 导致迭代异常
    ;[...set].forEach((fn) => {
      try {
        fn(msg)
      } catch (e) {
        console.error('WS handler error:', e)
      }
    })
  }, [])

  useEffect(() => {
    const ws = new WebSocket(url)
    wsRef.current = ws

    const syncState = () => setReadyState(ws.readyState)

    ws.onopen = () => syncState()
    ws.onclose = () => syncState()
    ws.onerror = () => syncState()

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        if (typeof msg?.req_identifier !== 'number') return
        dispatch(msg)
      } catch (e) {
        console.error('WS parse error:', e, event.data)
      }
    }

    syncState()

    return () => {
      ws.onopen = ws.onclose = ws.onerror = ws.onmessage = null
      ws.close()
      wsRef.current = null
      setReadyState(WebSocket.CLOSED)
    }
  }, [url, dispatch])

  const subscribe = useCallback((reqId: number, handler: Handler) => {
    let set = listenersRef.current.get(reqId)
    if (!set) {
      set = new Set()
      listenersRef.current.set(reqId, set)
    }
    set.add(handler)

    // 返回 unsubscribe，业务层 useEffect cleanup 直接用
    return () => {
      const cur = listenersRef.current.get(reqId)
      if (!cur) return
      cur.delete(handler)
      if (cur.size === 0) listenersRef.current.delete(reqId)
    }
  }, [])

  const send = useCallback(async (data: unknown) => {
    const payload = JSON.stringify(data)

    let lastState: number | undefined

    for (let attempt = 1; attempt <= 3; attempt++) {
      const ws = wsRef.current
      lastState = ws?.readyState

      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(payload)
        return
      }

      if (attempt < 3) {
        await sleep(1000)
      }
    }

    throw new Error(`WebSocket is not OPEN after 3 attempts. readyState=${lastState ?? 'null'}`)
  }, [])

  return {
    ws: wsRef.current,
    readyState,
    subscribe,
    send,
  }
}
