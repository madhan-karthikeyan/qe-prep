# URL Parser
Difficulty: Medium
Estimated Interview Time: 25 min
Prerequisites: string parsing, RFC 3986

## Problem Statement
Parse URLs according to RFC 3986, extracting scheme, userinfo, host, port, path, query, and fragment.

## Requirements
- Parse all URL components
- Handle default ports for common schemes
- Support IPv6 addresses (bracketed)
- Handle encoded characters
- Reconstruct URL from parsed components
- Validate URLs

## Implementation Notes
- Manual recursive-descent-style parser (no regex for the full grammar)
- IPv6 detection via leading bracket
- Default port mapping for well-known schemes

## Test Strategy (Unit/Fuzz)
- Unit: standard URLs, edge cases, malformed inputs, round-trip reconstruction
- Fuzz: 10000 random strings ensuring no crash

## Edge Cases
- Empty URL, missing scheme, missing host
- IPv6 with/without port
- Empty query/fragment components
- Path-only URLs (no scheme)
- Port out of range, non-numeric port

## Failure Cases
- Unclosed IPv6 bracket → ValueError
- Invalid port string → ValueError
- Empty URL → ValueError

## Complexity
- O(n) parse time, O(n) space (n = URL length)

## Progression Path
- Basic: split on ://
- Intermediate: full RFC 3986 grammar
- Advanced: percent-encoding validation
- Production: URL normalization, punycode

## Common Interview Follow-ups
- How would you handle relative URL resolution?
- How would you normalize URLs?
- How would you validate percent-encoding?

## Possible Production Improvements
- Add URL normalization (lowercase scheme/host, remove default ports)
- Add relative URL resolution per RFC 3986
- Support punycode IDN
- Performance optimization with zero-copy parsing
