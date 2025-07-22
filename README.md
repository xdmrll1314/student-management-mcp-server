# 学生管理 MCP Server

这是一个符合 Model Context Protocol (MCP) 标准的服务器，提供班级学生管理功能。

## 功能

- **get_student_list**: 获取班级所有学生的列表
  - 可选参数 `section`: 指定班级名称来过滤学生
- **get_student_info**: 根据学生ID获取单个学生的详细信息
  - 必需参数 `student_id`: 学生的唯一标识ID

## 使用方法

### 编译和运行

```bash
go build -o student-mcp-server hello.go
./student-mcp-server
```

### 在编程工具中使用

这个MCP Server可以被支持MCP协议的编程工具（如IDE、编辑器等）调用。服务器通过标准输入/输出进行JSON-RPC 2.0通信。

### 示例请求

#### 初始化
```json
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}
```

#### 获取工具列表
```json
{"jsonrpc":"2.0","id":2,"method":"tools/list"}
```

#### 调用工具 - 获取所有学生
```json
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_student_list","arguments":{}}}
```

#### 调用工具 - 获取特定班级学生
```json
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"get_student_list","arguments":{"section":"A班"}}}
```

#### 调用工具 - 获取单个学生信息
```json
{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_student_info","arguments":{"student_id":"1"}}}
```

## 数据结构

### 学生信息
- ID: 学生唯一标识
- Name: 学生姓名
- Age: 学生年龄
- Grade: 年级
- Section: 班级

## 协议支持

- MCP Protocol Version: 2024-11-05
- JSON-RPC 2.0
- 标准输入/输出通信