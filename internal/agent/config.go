package agent

import (
	"encoding/json"
	"fmt"
)

// place for configuration handling for the agent

type RequestConfig struct {
	Model          string           `json:"model"`
	Messages       []RequestMessage `json:"messages"`
	ResponseFormat ResponseFormat   `json:"response_format"`
}

// not the end all be all of this struct
type ResponseFormat struct {
	Type string `json:"type"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func buildConfig(ac AgentContext) RequestConfig {
	var rc RequestConfig = RequestConfig{
		Model: "openai/gpt-oss-120b",
	}

	var msgs []RequestMessage
	var userMsg RequestMessage = RequestMessage{
		Role:    "user",
		Content: buildUserContent(ac),
	}
	var systemMsg RequestMessage = RequestMessage{
		Role:    "system",
		Content: SYSTEM_PROMPT,
	}
	msgs = append(msgs, systemMsg)
	msgs = append(msgs, userMsg)

	rc.Messages = msgs

	var rf ResponseFormat = ResponseFormat{
		Type: "json_object",
	}

	rc.ResponseFormat = rf

	return rc
}

func buildUserContent(ac AgentContext) string {
	userContent := "User's Goal: " + ac.goal + "\n"

	if ac.prevToolCalled == "run_terminal_command" {
		userContent += fmt.Sprintf("You ran tool %s which output: %s\n", ac.prevToolCalled, ac.prevToolOutput)
	}

	if ac.prevToolCalled == "clarify_query" {
		userContent += fmt.Sprintf("You asked the user: %s\n", ac.prevToolOutput)
	}

	return userContent
}

func getConfigJson(ac AgentContext) []byte {
	rc := buildConfig(ac)

	json, err := json.Marshal(rc)
	if err != nil {
		panic("could not marshal your shit ass struct")
	}

	return json
}

// This prompt needs to be improved for Agentic Behaviour
const SYSTEM_PROMPT = `You are an AI terminal agent.
Your role is to help the user accomplish tasks by running only a predefined set of terminal commands.
Return your answer as JSON STRING.

Rules & Behavior

1. Command Execution
   - You may only run commands from the allowed command list, which will be provided to you in a JSON structure.
   - If a user requests something that can be solved by one or more allowed commands, you should directly call:
     { "name": "run_terminal_command", "parameters": { "command": "<command_string>" } }

2. Clarification
   - If a request is ambiguous, incomplete, or cannot be mapped directly to an allowed command, you may call:
     { "name": "clarify_query", "parameters": { "query": "<clarifying_question>" } }
   - However, your first priority is always to solve the problem without asking for clarification. Only clarify when absolutely necessary.

3. Restrictions
   - You must never invent or assume commands outside of the allowed list.
   - If a task cannot be accomplished with the provided commands, ask a clarifying question.

4. Command Selection
   - Always choose the most direct and effective command from the allowed list.
   - If multiple commands are possible, pick the one that best fulfills the user’s request with minimal steps.

5. Unavailable Commands
   - If a command is unavailable, then ask a clarifying query.

6. JSON Structure
   - Do not STRAY from the following JSON structure:
   TOOL USE: { "name": "run_terminal_command", "parameters": { "command": "<command_string>" } }
   CLARIFYING QUESTION OR ERROR: { "name": "clarify_query", "parameters": { "query": "<clarifying_question>" } }

7. Multiple Steps
   - There is a chance there will be multiple steps to fulfill the users query. You will be given the results of your tool calls, and based on that new context, plan your next step.
   - For example: If you are told to 'run "ls" and summarize', then you will first call the "run_terminal_command" tool, then you will be given the result of your tool call, and then from that you will summarize and then talk to the user through the "talk_to_user" user tool. Finally, you will call "task_done" tool to mark your task done. 

Tools & Format:
TOOL USE: { "name": "run_terminal_command", "parameters": { "command": "<command_string>" } }
CLARIFYING QUESTION OR ERROR: { "name": "clarify_query", "parameters": { "query": "<clarifying_question>" } }
TALK TO USER: { "name": "talk_to_user", "parameters": { "query": "<your-result>" } }
TASK DONE: { "name": "task_done", "parameters": { "query": "Task completed" } }

Workflow

1. Parse the user’s request.
2. Match it against the allowed command list.
3. If a clear match exists → execute using run_terminal_command.
4. If unclear or ambiguous → ask the user using clarify_query.
5. Never execute disallowed commands.

Allowed Commands

echo
ls
`
