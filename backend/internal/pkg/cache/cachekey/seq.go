package cachekey

const (
	ConvSeq    = "ConvSeq_SEQ:"
	ConvMinSeq = "ConvSeq_MIN_SEQ:"

	SeqUserMaxSeq  = "SEQ_USER_MAX:"
	SeqUserMinSeq  = "SEQ_USER_MIN:"
	SeqUserReadSeq = "SEQ_USER_READ:"
)

func GetSeqConvKey(conversationID string) string {
	return ConvSeq + conversationID
}

func GetSeqConvMinSeqKey(conversationID string) string {
	return ConvMinSeq + conversationID
}

func GetSeqUserMaxSeqKey(conversationID string, userID string) string {
	return SeqUserMaxSeq + conversationID + ":" + userID
}

func GetSeqUserMinSeqKey(conversationID string, userID string) string {
	return SeqUserMinSeq + conversationID + ":" + userID
}

func GetSeqUserReadSeqKey(conversationID string, userID string) string {
	return SeqUserReadSeq + conversationID + ":" + userID
}
