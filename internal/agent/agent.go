package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type GroqChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		ToolCall     []ToolCall `json:"tool_calls"`
		Logprobs     any        `json:"logprobs"`
		FinishReason string     `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		QueueTime        float64 `json:"queue_time"`
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	UsageBreakdown interface{} `json:"usage_breakdown"`
	// SystemFingerprint string      `json:"system_fingerprint"`
	// XGroq             struct {
	// 	ID string `json:"id"`
	// } `json:"x_groq"`
	// ServiceTier string `json:"service_tier"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"` // keep as string; parse separately if needed
	} `json:"function"`
}

func Symphony() {
	fmt.Println("It has begun")
}

const MAX_STEPS = 3 // This is probably for when it fails, maybe, idk have to architect it better
const GROQ_URL = "https://api.groq.com/openai/v1/chat/completions"

type UserMessage struct {
	msg string
}

// This needs to be thought out more, but im tryna get the agentic loop going first
// probably will be other datatypes, especially for tools, then have a message builder that takes this struct,
// and creates a message ready for the LLM
type AgentContext struct {
	goal        string
	toolResults string
	tools       string
}

type TempResponse struct {
	Name       string     `json:"name"`
	Parameters Parameters `json:"parameters"`
}

type Parameters struct {
	Query   string `json:"query"`
	Command string `json:"command"`
}

func StartLoop() {
	for i := range MAX_STEPS {
		fmt.Println(i)
	}

	var ac AgentContext = AgentContext{
		goal:        "run git status",
		toolResults: "",
		tools:       "see_command_history",
	}

	// Alright so here is how it will go
	// Step 1, feed the tool result, and user goal into the LLM
	// Step 2, check if response is a tool call, or action to stop
	// Step 3, if stop, then stop, else add to tool result
	// Step 4, GO TO Step 1

	response := requestLLM(ac)

	var tr TempResponse
	err := json.Unmarshal([]byte(response), &tr)
	if err != nil {
		panic(err)
	}

	fmt.Println(tr.Name, tr.Parameters.Query, tr.Parameters.Command)
	// runTerminalCommand()
}

func requestLLM(ac AgentContext) string {
	jsonData := getConfigJson(ac)

	fmt.Println(string(jsonData) + "*\n\n")

	req, err := http.NewRequest("POST", GROQ_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln("New Request err:", err)
	}

	GROQ_KEY := os.Getenv("GROQ_API_KEY")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", GROQ_KEY))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("http DO:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Read ALL err:", err)
	}

	// fmt.Println("Resp: "+string(body))
	// var res map[string]interface{}

	// err = json.NewDecoder(resp.Body).Decode(&res)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)
	// fmt.Println(res["choices"])

	fmt.Println(string(body))

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Fatalln("JSON unmarshal error:", err)
	}

	choices := res["choices"].([]interface{})
	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})
	content := message["content"]
	fmt.Println("Full message:", content)

	return content.(string)
}
