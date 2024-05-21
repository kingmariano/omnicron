# groq.com golang library

### Description

github.com/jpoz/groq is community maintained Go library for the [groq.com](https://console.groq.com) API.


### Installation

To install groq, run the following command in your terminal:

```shell
go get -u github.com/jpoz/groq
```

### Features

* Synchronous Chat completions
* Streaming Chat completions
* Zero dependencies


### Examples

```go
client := groq.NewClient() // will load API key from GROQ_API_KEY environment variable
client := groq.NewClient(WithAPIKey("YOUR_API_KEY"))

response, err := client.CreateChatCompletion(groq.CompletionCreateParams{
    Model: "llama3-8b-8192",
    Messages: []groq.Message{
        {
            Role:    "user",
            Content: "What is the meaning of life?",
        },
    },
})
if err != nil {
    panic(err)
}

println(response.Choices[0].Message.Content)
```

#### Streaming chat completions


```go
client := groq.NewClient(WithAPIKey("YOUR_API_KEY"))

chatCompletion, err := client.CreateChatCompletion(groq.CompletionCreateParams{
    Model: "llama3-70b-8192",
    Messages: []groq.Message{
        {
            Role:    "user",
            Content: "What is the meaning of life?",
        },
    },
    Stream:      true,
})
if err != nil {
    panic(err)
}

for delta := range chatCompletion.Stream {
    fmt.Print(delta.Choices[0].Delta.Content)
}

fmt.Print("\n")
```
