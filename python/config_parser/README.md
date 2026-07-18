# Config Parser
Difficulty: Medium
Estimated Interview Time: 30 min
Prerequisites: string parsing, file I/O

## Problem Statement
Implement parsers for .env and .ini configuration file formats.

## Requirements
- Env: KEY=VALUE format, comments (#), empty lines, quoted values, escaped chars, variable substitution ($VAR or ${VAR})
- Ini: [sections], key=value, comments (; and #), multi-line support, type coercion

## Implementation Notes
- Env parser uses regex for variable substitution, manual quoting logic
- Ini parser uses fallback to DEFAULT section for key lookup
- Type coercion for int, float, bool, str

## Test Strategy (Unit)
- Both parsers: standard cases, comments, quotes, substitution, errors, file parsing

## Edge Cases
- Empty lines and comments in both formats
- Quoted values with embedded special characters
- Variable substitution with braces vs without
- Fallback to OS environment variables
- Circular variable references (not handled, would resolve to empty string)
- DEFAULT section inheritance in INI

## Failure Cases
- Type coercion failure → ValueError
- Missing section access → KeyError
- Malformed file (partial parsing accepted)

## Complexity
- O(n) parse time, O(k) space (k = number of keys)

## Progression Path
- Basic: flat key-value parsing
- Intermediate: sections, types
- Advanced: variable substitution, escape handling
- Production: validation schema, nested sections

## Common Interview Follow-ups
- How would you handle circular references in variable substitution?
- How would you add type validation schemas?
- How would you handle multi-line values in INI?

## Possible Production Improvements
- Schema validation for config values
- Watch for file changes and reload
- Hierarchical config (env → ini → default)
- Profile-based config loading
