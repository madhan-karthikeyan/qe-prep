# REST API Tester
Difficulty: Hard
Estimated Interview Time: 45 min
Prerequisites: HTTP, REST, JSON, server implementation

## Problem Statement
Implement a lightweight REST API test server and testing utility. The server supports CRUD operations with JSON, and the tester provides request/assertion utilities.

## Requirements
- TestServer: HTTP server with GET/POST/DELETE, JSON responses, in-memory store
- ApiTester: sendRequest, assertStatus, assertJson, retry with exponential backoff
- TestResult record with status code, body, timing, pass/fail
- Proper error handling and status codes (200, 201, 400, 404, 405)

## Implementation Notes
- Uses com.sun.net.httpserver.HttpServer for the test server
- Uses java.net.HttpClient for the API tester
- Retry with exponential backoff + jitter
- JSON parsing uses simple string operations (no JSON library dependency)

## Test Strategy
- Unit: validation, assertion logic
- Integration: start server, run API operations
- Error cases: non-existent resources, wrong methods
- Retry behavior on failures

## Edge Cases
- Server not running when tester connects
- Non-existent endpoints (404)
- Invalid HTTP methods (405)
- Empty request/response bodies

## Failure Cases
- null/blank baseUrl → IllegalArgumentException
- null timeout → IllegalArgumentException
- assertion failures → AssertionError

## Complexity (Time + Space)
- Server: O(1) per request handling, O(items) space
- Tester: O(1) per request, O(1) per assertion

## Progression Path
Build the server first with basic GET/POST, then add DELETE, then build the testing utility with assertions.

## Common Interview Follow-ups
- Support for PUT and PATCH methods
- Authentication/authorization headers
- Request/response schema validation

## Possible Production Improvements
- Use a proper JSON library (Jackson, Gson)
- OpenAPI/Swagger schema validation
- Parallel test execution with test fixtures
