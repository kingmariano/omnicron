package groq

// ChoiceDeltaFunctionCall struct for details of function calls
type ChoiceDeltaFunctionCall struct {
	Arguments *string `json:"arguments,omitempty"`
	Name      *string `json:"name,omitempty"`
}

// ChoiceDeltaToolCallFunction struct for tool call function details
type ChoiceDeltaToolCallFunction struct {
	Arguments *string `json:"arguments,omitempty"`
	Name      *string `json:"name,omitempty"`
}

// ChoiceDeltaToolCall struct for tool call details within deltas
type ChoiceDeltaToolCall struct {
	Index    int                         `json:"index"`
	ID       *string                     `json:"id,omitempty"`
	Function ChoiceDeltaToolCallFunction `json:"function,omitempty"`
	Type     *string                     `json:"type,omitempty"`
}

// ChoiceDelta struct for deltas within choices
type ChoiceDelta struct {
	Content      string                   `json:"content"`
	Role         string                   `json:"role"`
	FunctionCall *ChoiceDeltaFunctionCall `json:"functionCall,omitempty"`
	ToolCalls    []ChoiceDeltaToolCall    `json:"toolCalls,omitempty"`
}

// Choice struct for handling choices in completion chunks
type ChoiceChunk struct {
	Delta        ChoiceDelta    `json:"delta"`
	FinishReason string         `json:"finishReason"`
	Index        int            `json:"index"`
	Logprobs     ChoiceLogprobs `json:"logprobs"`
}

// XGroq struct for additional external data
type XGroq struct {
	Usage Usage `json:"usage"`
}

// ChatCompletion struct to represent the overall completion
type ChatChunkCompletion struct {
	ID                *string       `json:"id"`
	Choices           []ChoiceChunk `json:"choices"`
	Created           *int          `json:"created"`
	Model             *string       `json:"model"`
	Object            *string       `json:"object"`
	SystemFingerprint *string       `json:"systemFingerprint"`
	XGroq             *XGroq        `json:"xGroq,omitempty"`

	stream chan *ChoiceChunk
}
