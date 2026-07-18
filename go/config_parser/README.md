# Config Parser
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: string manipulation, state machines

## Problem Statement
Implement parsers for .env and .ini configuration file formats.

## Requirements
- .env: KEY=VALUE, comments (#), quoted values, escapes, $VAR/${VAR} substitution
- .ini: [section] headers, key=value, comments (#, ;), blank lines

## Implementation Notes
- .env handles double-quoted strings with escape sequences
- Variable substitution uses values defined on earlier lines
- .ini keys outside sections are rejected

## Test Strategy
- Table-driven tests for standard cases
- Edge cases: comments, quotes, substitution, blank lines
- Error cases: unterminated strings, missing equals

## Edge Cases
- Empty values
- Values with internal spaces
- Escaped dollar signs (\$)
- Nested substitution
- Multiple sections with same name

## Failure Cases
- Unterminated quoted strings
- Keys without values
- Keys outside section in INI
- Malformed section headers

## Complexity (Time + Space)
- Both parsers: O(n) time, O(n) space where n = input length

## Progression Path
- Add nested section support (INI)
- Add typed values (int, bool, float)
- Add includes and conditionals

## Common Interview Follow-ups
- Config merging (multiple files)
- Environment variable override
- Type coercion

## Possible Production Improvements
- Validate against schema
- Support YAML/TOML alongside
- Hot-reload config changes
