# Config Parser
Difficulty: Medium
Estimated Interview Time: 35 min
Prerequisites: File I/O, parsing, regex, maps

## Problem Statement
Implement parsers for .env and .ini configuration file formats. Handle comments, quotes, multi-line values, and variable substitution.

## Requirements
- EnvParser: KEY=VALUE format, comments (#), quoted values, optional variable substitution
- IniParser: [section], key=value, comments (; and #), multi-line values (backslash)
- Line continuation support
- Whitespace trimming around keys and values

## Implementation Notes
- Uses regex for line pattern matching
- Multi-line values joined with newline
- Variable substitution supports ${VAR} syntax in .env
- Thread-safe (immutable result maps)

## Test Strategy
- Parse standard files, comments, quotes
- Multi-line values, empty values
- Error cases (invalid keys, null input)
- Round-trip verification

## Edge Cases
- Keys without sections (.ini)
- Values containing = signs
- Single vs double quotes
- Escaped characters

## Failure Cases
- null path/content → NullPointerException
- Invalid key format → IllegalArgumentException
- Malformed variable reference → preserved as-is

## Complexity (Time + Space)
- Parse: O(n) time, O(n) space where n = file size

## Progression Path
Implement simple key-value parsing first, then add sections, comments, multi-line support.

## Common Interview Follow-ups
- Nested section support (e.g., [section/subsection])
- Include/import directives
- Environment variable overrides

## Possible Production Improvements
- Schema validation for expected keys
- Type coercion (string → int, bool)
- Write-back support with formatting preservation
