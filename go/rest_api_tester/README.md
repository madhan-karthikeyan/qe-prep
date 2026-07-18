# REST API Tester
Difficulty: Medium
Estimated Interview Time: 40 min
Prerequisites: net/http, json, REST API basics

## Problem Statement
Implement a minimal HTTP test server and an API testing utility.

## Requirements
- Test server with GET/POST/DELETE and in-memory JSON storage
- API tester: SendRequest, AssertStatus, AssertJSON, Retry, RunTests
- Full integration test suite

## Implementation Notes
- Server uses Go 1.22+ pattern routing
- JSONPath supports dot notation for nested access
- Retry with exponential backoff and jitter

## Test Strategy
- Unit tests for server and tester independently
- Integration: start test server, run full API test suite

## Edge Cases
- Missing items (404)
- Bad JSON bodies
- Empty storage
- Concurrent requests

## Failure Cases
- Server not running
- Invalid JSON path
- Response timeout

## Complexity (Time + Space)
- Server: O(1) per operation
- Tester: O(retries) per request
- JSONPath: O(depth) per lookup

## Progression Path
- Add PUT/PATCH support
- Add response body schema validation
- Add test report generation

## Common Interview Follow-ups
- HTTP status code conventions
- RESTful API design
- Testing strategies for APIs

## Possible Production Improvements
- Use testify/httptest for production
- Add OpenAPI/Swagger validation
- Generate HTML test reports
