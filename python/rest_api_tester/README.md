# REST API Tester
Difficulty: Medium
Estimated Interview Time: 35 min
Prerequisites: HTTP, JSON, testing concepts

## Problem Statement
Build a minimal test HTTP server and an API testing utility that can send requests, assert responses, and generate test reports.

## Requirements
- Minimal test server with GET, POST, DELETE handlers
- JSON request/response
- Configurable status codes
- In-memory storage
- API testing utility with send_request, assert_status, assert_json
- Retry on failure
- Test report (pass/fail count, timing)

## Implementation Notes
- Uses http.server.HTTPServer with a custom handler
- API tester uses urllib.request for HTTP calls
- Test report tracks total, passed, failed, errors, and timing

## Test Strategy (Unit/Integration)
- Unit: test server + assertions, error cases
- Integration: full API workflow end-to-end, multiple data types

## Edge Cases
- Missing key in JSON body → 400
- Deleting nonexistent key → 404
- Getting nonexistent key → 404
- Invalid JSON → 400
- Nested JSON path assertions

## Failure Cases
- Assertion mismatch → AssertionError
- Request failure after retries → RuntimeError
- Empty base URL → ValueError

## Complexity
- Server: O(1) per request with O(k) storage
- Tester: O(r) per request (r = retries), O(1) per assertion

## Progression Path
- Basic: single-endpoint HTTP test
- Intermediate: full CRUD test server
- Advanced: test suite with reporting
- Production: parallel test execution, plugin architecture

## Common Interview Follow-ups
- How would you add authentication testing?
- How would you handle file uploads?
- How would you parallelize test execution?

## Possible Production Improvements
- Parallel test execution
- Authentication/authorization testing
- Response schema validation
- HTML report generation
- Support for custom assertions
