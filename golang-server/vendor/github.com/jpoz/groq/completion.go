package groq

// ChoiceLogprobsContentTopLogprob struct to represent detailed log probabilities
type ChoiceLogprobsContentTopLogprob struct {
	Token   *string  `json:"token,omitempty"`
	Bytes   []int    `json:"bytes,omitempty"`
	Logprob *float64 `json:"logprob,omitempty"`
}

// ChoiceLogprobsContent struct to handle content log probabilities
type ChoiceLogprobsContent struct {
	Token       *string                           `json:"token,omitempty"`
	Bytes       []int                             `json:"bytes,omitempty"`
	Logprob     *float64                          `json:"logprob,omitempty"`
	TopLogprobs []ChoiceLogprobsContentTopLogprob `json:"topLogprobs,omitempty"`
}

// ChoiceLogprobs struct to contain multiple log probability contents
type ChoiceLogprobs struct {
	Content []ChoiceLogprobsContent `json:"content,omitempty"`
}

// ChoiceMessageToolCallFunction struct to describe the function called by the tool
type ChoiceMessageToolCallFunction struct {
	Arguments *string `json:"arguments,omitempty"`
	Name      *string `json:"name,omitempty"`
}

// ChoiceMessageToolCall struct to represent a tool call in a message
type ChoiceMessageToolCall struct {
	ID       *string                       `json:"id,omitempty"`
	Function ChoiceMessageToolCallFunction `json:"function,omitempty"`
	Type     *string                       `json:"type,omitempty"`
}

// ChoiceMessage struct for messages within choices
type ChoiceMessage struct {
	Content   string                  `json:"content"`
	Role      string                  `json:"role"`
	ToolCalls []ChoiceMessageToolCall `json:"toolCalls,omitempty"`
}

// Choice struct to handle individual choices in a completion
type Choice struct {
	FinishReason string         `json:"finishReason"`
	Index        int            `json:"index"`
	Logprobs     ChoiceLogprobs `json:"logprobs"`
	Message      ChoiceMessage  `json:"message"`
}

// Usage struct to handle usage data
type Usage struct {
	CompletionTime   *float64 `json:"completionTime,omitempty"`
	CompletionTokens *int     `json:"completionTokens,omitempty"`
	PromptTime       *float64 `json:"promptTime,omitempty"`
	PromptTokens     *int     `json:"promptTokens,omitempty"`
	QueueTime        *float64 `json:"queueTime,omitempty"`
	TotalTime        *float64 `json:"totalTime,omitempty"`
	TotalTokens      *int     `json:"totalTokens,omitempty"`
}

// ChatCompletion struct to represent the overall completion
type ChatCompletion struct {
	ID                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Object            string   `json:"object"`
	SystemFingerprint string   `json:"systemFingerprint"`
	Usage             Usage    `json:"usage,omitempty"`

	Stream chan *ChatChunkCompletion `json:"-"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
