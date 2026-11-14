package main

import (
  "bytes"
  "encoding/json"
  "io"
  "log"
  "net/http"
  "os"
  "blaiddddrwg/go-harness/skills"
  "strings"
)

var MODEL_ENDPOINT = "http://localhost:11434/api/chat"
var SKILLS_DIR     = "skills"

type FunctionAction func(input string) ([]string, error)

type PARAMETERS struct {
  Type       string                       `json:"type"`
  Properties map[string]map[string]string `json:"properties"`
  Required   []string                     `json:"required"`
}

type FUNCTION struct {
  Name        string                      `json:"name"`
  Description string                      `json:"description"`
  Parameters  PARAMETERS                  `json:"parameters"`
  Action      FunctionAction              `json:"-"`
}

type TOOL struct {
  Type     string                         `json:"type"`
  Function FUNCTION                       `json:"function"`
}

type TOOL_CALL struct {
  ID       string                         `json:"id"`
  Function struct {
    Index      int                        `json:"index"`
    Name       string                     `json:"name"`
    Arguments  map[string]interface{}     `json:"arguments"`
  } `json:"function"`
}

type MESSAGE struct {
  Role    string                          `json:"role"`
  ToolName string                         `json:"tool_name,omitempty"`
  ID      string                          `json:"id,omitempty"`
  Content string                          `json:"content"`
  Thinking string                         `json:"thinking,omitempty"`
  ToolCalls []TOOL_CALL                   `json:"tool_calls,omitempty"`
}

type CHAT_REQUEST struct {
  Model    string                         `json:"model"`
  Messages []MESSAGE                      `json:"messages"`
  Tools    []TOOL                         `json:"tools,omitempty"`
  Stream   bool                           `json:"stream"`
}

func list_files_in_directory(path string) ([]string, error) {
  files, err := os.ReadDir(path)
  if err != nil {
    return nil, err
  }

  var fileNames []string
  for _, file := range files {
    fileNames = append(fileNames, file.Name())
  }

  return fileNames, nil
}

var musicReccommendationsCallCount int = 0
var MUSIC_RECOMMENDATIONS = map[string][]string{
  "rock": {"led zeppelin", "pink floyd", "fleetwood mac"},
  "classical": {"beethoven", "brahms", "vivaldi"},
  "alternative": {"cheb mami", "santana", "tiniarawen"},
}

func find_music(category string) ([]string, error) {
  musicReccommendationsCallCount++
  value, ok := MUSIC_RECOMMENDATIONS[category]
  if ok {
    return value, nil
  }

  return []string{"karnatic", "hindustani", "rai"}, nil
}

func execute(instructions string) ([]string, error) {
  log.Printf("Executing skill with instructions: %s", instructions)
  return []string{"one", "two", "three"}, nil
}

// Ideally, we should use some kind of vector search to find the best skill, but
// for now, we will just return ALL the skills
func find_skills(description string) ([]string, error) {
  skills, err := skill.FindSkills(SKILLS_DIR)
  if err != nil {
    log.Printf("Got error trying to find skills: %s", err)
    return nil, err
  }

  var skillNames []string
  for _, skill := range skills {
    skillNames = append(skillNames, skill.Name())
    // Then add the skill to the set of tools that will get returned.
    SKILLS = append(SKILLS, TOOL{
      Type: "function",
      Function: FUNCTION{
        Name:        skill.Name(),
        Description: skill.Description(),
        Parameters: PARAMETERS{
          Type: "object",
          Properties: map[string]map[string]string{
            "path": {
              "type":        "string",
              "description": "Instructions for executing the skill",
            },
          },
          Required: []string{"path"},
        },
        Action: execute,
      },
    })
  }

  log.Print("---------------------------------------------------------------")
  log.Printf("Returning skills: %s", strings.Join(skillNames, ","))
  log.Print("---------------------------------------------------------------")
  return skillNames, nil
}

// Placeholder for skills
var SKILLS = []TOOL{}

// Basic tools that are packaged with the basic harness.
var TOOLS = []TOOL{
  {
    Type: "function",
    Function: FUNCTION{
      Name:        "list_dir",
      Description: "List the files in a directory",
      Parameters: PARAMETERS{
        Type: "object",
        Properties: map[string]map[string]string{
          "path": {
            "type":        "string",
            "description": "The directory to list",
          },
        },
        Required: []string{"path"},
      },
      Action: list_files_in_directory,
    },
  },
  {
    Type: "function",
    Function: FUNCTION {
      Name:         "find_music",
      Description:  "Find music of a particular genre",
      Parameters: PARAMETERS{
        Type: "object",
        Properties: map[string]map[string]string{
          "path": {
            "category":     "string",
            "description":  "A general description of the vibe",
          },
        },
        Required: []string{"path"},
      },
      Action: find_music,
    },
  },
  {
    Type: "function",
    Function: FUNCTION {
      Name:           "more_tools",
      Description:    "Check for additional tools. Choose this tool if existing tools provided to you are inadequate to answer the users request",
      Parameters: PARAMETERS{
        Type: "skills",
        Properties: map[string]map[string]string{
          "path": {
            "category":       "string",
            "description":    "Skill finder from user defined skills",
          },
        },
        Required: []string{"path"},
      },
      Action: find_skills,
    },
  },
}

type CHAT_RESPONSE struct {
  Model string                        `json:"model"`
  CreatedAt string                    `json:"created_at"`
  Message MESSAGE                     `json:"message"`
  Done bool                           `json:"done"`
  DoneReason string                   `json:"done_reason"`
  TotalDuration float64               `json:"total_duration"`
  LoadDuration float64                `json:"load_duration"`
  PromptEvalCount int                 `json:"prompt_eval_count"`
  PromptEvalCachedCount int           `json:"prompt_eval_cached_count"`
  PromptEvalDuration float64          `json:"prompt_eval_duration"`
  EvalCount int                       `json:"eval_count"`
  EvalDuration float64                `json:"eval_duration"`
}

func find_tool_index_by_name(name string) int {
  for i, tool := range TOOLS {
    if tool.Function.Name == name {
      return i
    }
  }
  return -1
}

func model_call(client *http.Client, data CHAT_REQUEST) (*http.Response, error) {
  payload, err := json.Marshal(data)
  if err != nil {
    return nil, err
  }

  req, err := http.NewRequest("POST", MODEL_ENDPOINT, bytes.NewBuffer(payload))
  if err != nil {
    return nil, err
  }
  req.Header.Set("Content-Type", "application/json")

  println("Request payload:", string(payload))

  return client.Do(req)
}

func execute_tool_calls(tool_calls []TOOL_CALL) []MESSAGE {
  results := []MESSAGE{}
  for _, toolCall := range tool_calls {
    log.Printf("Tool call: %s - %s", toolCall.ID, toolCall.Function.Name)
    log.Printf("Arguments: %v", toolCall.Function.Arguments)
    toolIndex:= find_tool_index_by_name(toolCall.Function.Name)
    if toolIndex == -1 {
      toolMessage := MESSAGE{
        Role:    "tool",
        ToolName: toolCall.Function.Name,
        Content: "Tool call failed - tool not found",
      }
      results = append(results, toolMessage)
      continue
    }

    toolFunction := TOOLS[toolIndex].Function
    // The model's arguments are untyped JSON, so check that "path"
    // exists and is a string before using it.
    path, ok := toolCall.Function.Arguments["path"].(string)
    if !ok {
      toolMessage := MESSAGE{
        Role:    "tool",
        ToolName: toolCall.Function.Name,
        Content: "Tool call failed - missing or non-string \"path\" argument",
      }
      results = append(results, toolMessage)
      continue
    }

    tool_call_result, err := toolFunction.Action(path)
    if err != nil {
      log.Printf("Error executing tool call: %v", err)
      toolMessage := MESSAGE{
        Role:    "tool",
        ToolName: toolCall.Function.Name,
        Content: "Tool call failed - error executing tool function",
      }
      results = append(results, toolMessage)
      continue
    }

    toolMessage := MESSAGE{
      Role:    "tool",
      ToolName: toolCall.Function.Name,
      ID:       toolCall.ID,
      Content: strings.Join(tool_call_result[:], ","),
    }
    results = append(results, toolMessage)
  }

  return results
}

func main() {
  client := &http.Client{}
  // Initialize product call count
  musicReccommendationsCallCount = 0

  MessageHistory := []MESSAGE{
    {
      Role: "system",
      Content: "Try to answer the users question to the best of your abilities using the existing tools that are provided. Remember that there are other tools as well, that can be provided on demand. If none of the existing tools provided look appropriate to handle the users question, ask the harness for a list of additional tools to help you handle the user question. If you ask for additional tools with a brief summary of what you are looking for, the harness should send you back an addtional list of tools that you can use. You can search for additional tools by using the 'more_tools' tool - send a request back to that tool with a search term, and you should get back additional tools that you can use to answer the users question",
    },
    {
      Role: "user",
      Content: "Find new emails from Pinju",
    },
  }

  OriginalChatRequest := CHAT_REQUEST{
    Model: "qwen3.5:9b",
    Messages: MessageHistory[:],
    Tools:  TOOLS,
    Stream: false,
  }

  resp, err := model_call(client, OriginalChatRequest)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  // Read the response body
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  var _resp CHAT_RESPONSE
  err = json.Unmarshal(body, &_resp)
  if err != nil {
    log.Fatal(err)
  }

  // Append the model's response to the message history
  MessageHistory = append (MessageHistory, _resp.Message)

  toolCalls := _resp.Message.ToolCalls
  for len(toolCalls) > 0 {
    tool_call_results := execute_tool_calls(toolCalls)
    if len(tool_call_results) > 0 {
      MessageHistory = append(MessageHistory, tool_call_results...)
      // Send the tool call results back to the model
      resp, err = model_call(client, CHAT_REQUEST{
        Model: "qwen3.5:9b",
        Messages: MessageHistory,
        Tools:  append(TOOLS, SKILLS...),
        Stream: false,
      })
      if err != nil {
        log.Fatal(err)
      }

      // Read the response body
      body, err = io.ReadAll(resp.Body)
      if err != nil {
        log.Fatal(err)
      }
      resp.Body.Close()

      var toolResp CHAT_RESPONSE
      err = json.Unmarshal(body, &toolResp)
      if err != nil {
        log.Fatal(err)
      }

      // Append the model's response to the message history
      MessageHistory = append (MessageHistory, toolResp.Message)

      toolCalls = toolResp.Message.ToolCalls
    }
  }

  log.Printf("Response status: %s", resp.Status)
  log.Printf("Response body: %s", string(body))
  log.Printf("Music reccommendations call count: %d", musicReccommendationsCallCount)
}
