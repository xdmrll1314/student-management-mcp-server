package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// MCP Protocol structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool definitions
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// Student data structure
type Student struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Grade   string `json:"grade"`
	Section string `json:"section"`
}

// Mock student data
var students = []Student{
	{ID: "1", Name: "张三", Age: 18, Grade: "高三", Section: "A班"},
	{ID: "2", Name: "李四", Age: 17, Grade: "高二", Section: "B班"},
	{ID: "3", Name: "王五", Age: 16, Grade: "高一", Section: "C班"},
	{ID: "4", Name: "赵六", Age: 17, Grade: "高二", Section: "A班"},
	{ID: "5", Name: "孙七", Age: 16, Grade: "高一", Section: "B班"},
}

// MCP Server implementation
type MCPServer struct{}

func (s *MCPServer) handleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func (s *MCPServer) handleInitialize(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "student-management-server",
				"version": "1.0.0",
			},
		},
	}
}

func (s *MCPServer) handleToolsList(req MCPRequest) MCPResponse {
	tools := []Tool{
		{
			Name:        "get_student_list",
			Description: "获取班级所有学生的列表",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"section": map[string]interface{}{
						"type":        "string",
						"description": "班级名称，可选参数。如果提供，则只返回该班级的学生",
					},
				},
			},
		},
		{
			Name:        "get_student_info",
			Description: "根据学生ID获取单个学生的详细信息",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"student_id": map[string]interface{}{
						"type":        "string",
						"description": "学生的唯一标识ID",
					},
				},
				"required": []string{"student_id"},
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func (s *MCPServer) handleToolsCall(req MCPRequest) MCPResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Tool name is required",
			},
		}
	}

	toolArgs, _ := params["arguments"].(map[string]interface{})

	switch toolName {
	case "get_student_list":
		return s.getStudentList(req.ID, toolArgs)
	case "get_student_info":
		return s.getStudentInfo(req.ID, toolArgs)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Tool not found",
			},
		}
	}
}

func (s *MCPServer) getStudentList(id interface{}, args map[string]interface{}) MCPResponse {
	section, _ := args["section"].(string)

	var filteredStudents []Student
	if section != "" {
		for _, student := range students {
			if student.Section == section {
				filteredStudents = append(filteredStudents, student)
			}
		}
	} else {
		filteredStudents = students
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("找到 %d 名学生:\n%s", len(filteredStudents), formatStudentList(filteredStudents)),
				},
			},
		},
	}
}

func (s *MCPServer) getStudentInfo(id interface{}, args map[string]interface{}) MCPResponse {
	studentID, ok := args["student_id"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "student_id is required",
			},
		}
	}

	for _, student := range students {
		if student.ID == studentID {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      id,
				Result: map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": formatStudentInfo(student),
						},
					},
				},
			}
		}
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    -32603,
			Message: "Student not found",
		},
	}
}

func formatStudentList(students []Student) string {
	result := ""
	for _, student := range students {
		result += fmt.Sprintf("ID: %s, 姓名: %s, 年龄: %d, 年级: %s, 班级: %s\n",
			student.ID, student.Name, student.Age, student.Grade, student.Section)
	}
	return result
}

func formatStudentInfo(student Student) string {
	return fmt.Sprintf("学生详细信息:\nID: %s\n姓名: %s\n年龄: %d\n年级: %s\n班级: %s",
		student.ID, student.Name, student.Age, student.Grade, student.Section)
}

func main() {
	server := &MCPServer{}
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req MCPRequest
		if err := decoder.Decode(&req); err != nil {
			log.Printf("Error decoding request: %v", err)
			break
		}

		resp := server.handleRequest(req)
		if err := encoder.Encode(resp); err != nil {
			log.Printf("Error encoding response: %v", err)
			break
		}
	}
}