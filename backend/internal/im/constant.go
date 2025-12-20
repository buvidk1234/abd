package im

import "time"

const (
	// Websocket URL parameters
	WsUserID                = "sendID"
	PlatformID              = "platformID"
	Token                   = "token"
	Compression             = "compression" //compression == "gzip" means use gzip compression
	GzipCompressionProtocol = "gzip"

	// Additional parameters in context
	ConnID           = "connID"
	OperationID      = "operationID"
	BackgroundStatus = "isBackground"
	SendResponse     = "isMsgResp"
)

const (
	// Websocket Protocol.
	// WSSendMsg           = 1001
	// WSPullSpecifiedConv = 1002
	// WSPullConvList      = 1003

	WSGetNewestSeq        = 1001
	WSPullMsgBySeqList    = 1002
	WSSendMsg             = 1003
	WSSendSignalMsg       = 1004
	WSPullMsg             = 1005
	WSGetConvMaxReadSeq   = 1006
	WsPullConvLastMessage = 1007
	WSPushMsg             = 2001
	WSKickOnlineMsg       = 2002
	WsLogoutMsg           = 2003
	WsSetBackgroundStatus = 2004
	WsSubUserOnlineStatus = 2005
	WSDataError           = 3001
	WSTest                = 4001
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10000 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30000 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 51200
)

const (
	// MessageText is for UTF-8 encoded text messages like JSON.
	MessageText = iota + 1
	// MessageBinary is for binary messages like protobufs.
	MessageBinary
	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

const (
	WebPlatformID = 1
)
