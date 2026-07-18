# URL Parser
Difficulty: Medium
Estimated Interview Time: 35 min
Prerequisites: string manipulation, runes, net/url (for comparison only)

## Problem Statement
Implement a URL parser from scratch without using net/url.

## Requirements
- Parse into scheme, userinfo, host, port, path, query, fragment
- Handle IPv6 addresses
- Handle percent-encoded characters
- Reconstruct URL string
- Validate URL

## Implementation Notes
- Manual parsing with string operations
- IPv6 bracket handling
- Default ports for http (80) and https (443)
- Percent decode utility

## Test Strategy
- Table-driven tests for standard URLs
- Edge cases: IPv6, relative URLs, encoded chars
- Fuzz testing with random strings

## Edge Cases
- IPv6 address with and without port
- Empty URL
- Relative URLs (no scheme)
- Query and fragment without path

## Failure Cases
- Invalid scheme
- Unterminated IPv6 bracket
- Empty string

## Complexity (Time + Space)
- Parse: O(n) time, O(n) space
- String: O(n) time, O(n) space

## Progression Path
- Add percent encoding
- Support opaque URIs (mailto:, tel:)
- Add path normalization

## Common Interview Follow-ups
- URL normalization and canonicalization
- IDN/Unicode handling

## Possible Production Improvements
- Use net/url for production
- Add URL template matching
- Support custom scheme handlers
