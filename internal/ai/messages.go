package ai

// StreamChunkMsg is sent when a chunk of text is received from the stream
type StreamChunkMsg struct {
	Text string
}

// StreamDoneMsg is sent when streaming is complete
type StreamDoneMsg struct{}

// StreamErrorMsg is sent when an error occurs during streaming
type StreamErrorMsg struct {
	Err error
}

// StreamStartedMsg is sent when a stream has been successfully initiated.
// It contains the channel from which to read stream chunks.
type StreamStartedMsg struct {
	Stream <-chan StreamChunk
}

