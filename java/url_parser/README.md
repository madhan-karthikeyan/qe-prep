# URL Parser
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: RFC 3986, regex, parsing techniques

## Problem Statement
Implement an RFC 3986 compliant URL parser from scratch without using java.net.URI. Parse URLs into components and support reconstruction.

## Requirements
- Parse scheme, host, port, path, query, fragment
- Handle IPv6 addresses (bracketed notation)
- Default ports for known schemes (http=80, https=443, ftp=21)
- Query string parsing into key-value pairs
- URL reconstruction (canonical form)

## Implementation Notes
- Uses regex for initial parsing, then processes components
- Separate handling for IPv6 vs regular hosts
- Default port is omitted in reconstruction if matches scheme default
- Query parameters decoded with percent-decoding

## Test Strategy
- Standard URLs, edge cases, IPv6
- Cross-reference with java.net.URI
- Reconstruction round-trip tests
- Malformed input rejection

## Edge Cases
- IPv6 loopback [::1]
- Empty host or path
- Special characters and percent-encoding
- User info with passwords

## Failure Cases
- null input → NullPointerException
- Malformed URL (no scheme) → IllegalArgumentException
- Invalid port number → IllegalArgumentException

## Complexity (Time + Space)
- Parse: O(n) time, O(n) space where n = URL length
- Reconstruction: O(n) time, O(n) space

## Progression Path
Start with simple http:// URLs, add query/fragment parsing, then IPv6 support.

## Common Interview Follow-ups
- URL normalization and canonicalization
- Relative URL resolution
- Internationalized domain names (IDN)

## Possible Production Improvements
- Strict RFC 3986 compliance
- Percent-encoding validation
- Punycode IDN support
