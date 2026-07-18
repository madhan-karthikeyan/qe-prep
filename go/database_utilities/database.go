package database_utilities

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type table struct {
	columns []string
	rows    []map[string]interface{}
}

// Rows represents a query result set.
type Rows struct {
	columns []string
	rows    []map[string]interface{}
	pos     int
}

// Next advances the cursor to the next row. Returns false when exhausted.
func (r *Rows) Next() bool {
	if r.pos >= len(r.rows) {
		return false
	}
	r.pos++
	return true
}

// Scan copies the current row's columns into dest. Each dest must be a pointer
// to a compatible type (string, int64, float64, bool, or []byte).
func (r *Rows) Scan(dest ...interface{}) error {
	if r.pos == 0 || r.pos > len(r.rows) {
		return fmt.Errorf("no row available")
	}
	row := r.rows[r.pos-1]
	for i, d := range dest {
		if i >= len(r.columns) {
			break
		}
		col := r.columns[i]
		val, ok := row[col]
		if !ok {
			continue
		}
		if err := scanAssign(d, val); err != nil {
			return fmt.Errorf("column %d (%s): %w", i, col, err)
		}
	}
	return nil
}

func scanAssign(dest interface{}, val interface{}) error {
	switch d := dest.(type) {
	case *string:
		*d = fmt.Sprintf("%v", val)
	case *int64:
		switch v := val.(type) {
		case int64:
			*d = v
		case float64:
			*d = int64(v)
		case string:
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			*d = n
		default:
			return fmt.Errorf("cannot convert %T to int64", val)
		}
	case *float64:
		switch v := val.(type) {
		case float64:
			*d = v
		case int64:
			*d = float64(v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			*d = f
		default:
			return fmt.Errorf("cannot convert %T to float64", val)
		}
	case *bool:
		v, ok := val.(bool)
		if !ok {
			return fmt.Errorf("cannot convert %T to bool", val)
		}
		*d = v
	case *[]byte:
		*d = []byte(fmt.Sprintf("%v", val))
	default:
		return fmt.Errorf("unsupported scan type %T", dest)
	}
	return nil
}

// Columns returns the column names of the result set.
func (r *Rows) Columns() []string {
	return r.columns
}

// Tx represents an in-memory database transaction with snapshot isolation.
type Tx struct {
	db       *Database
	snapshot map[string]*table
	active   bool
}

// Database is a simple in-memory map-based database supporting basic SQL
// operations. Thread-safe via sync.RWMutex.
type Database struct {
	mu     sync.RWMutex
	tables map[string]*table
}

// Open creates a new in-memory database.
func Open() *Database {
	return &Database{
		tables: make(map[string]*table),
	}
}

// Close clears all tables and releases resources.
func (db *Database) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.tables = nil
	return nil
}

// TableColumns returns the column names for the given table.
func (db *Database) TableColumns(tableName string) ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	t, ok := db.tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %q not found", tableName)
	}
	cols := make([]string, len(t.columns))
	copy(cols, t.columns)
	return cols, nil
}

// Execute runs a SQL statement. Supported statements:
//
//	CREATE TABLE name (col1 type, col2 type, ...)
//	INSERT INTO name VALUES (val1, val2, ...)
//	UPDATE name SET col=val [WHERE cond [AND cond ...]]
//	DELETE FROM name [WHERE cond [AND cond ...]]
//	DROP TABLE name
//
// Use "?" placeholders in values and conditions; actual values are passed via
// args.
func (db *Database) Execute(query string, args ...interface{}) error {
	tokens, err := tokenizeSQL(query)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return fmt.Errorf("empty query")
	}
	cmd := strings.ToUpper(tokens[0])
	switch cmd {
	case "CREATE":
		return db.execCreate(tokens[1:])
	case "INSERT":
		return db.execInsert(tokens[1:], args)
	case "UPDATE":
		return db.execUpdate(tokens[1:], args)
	case "DELETE":
		return db.execDelete(tokens[1:], args)
	case "DROP":
		return db.execDrop(tokens[1:])
	default:
		return fmt.Errorf("unknown statement type %q", cmd)
	}
}

// Query runs a SELECT statement and returns a result set. Supports WHERE with
// AND conditions and "?" placeholders.
func (db *Database) Query(query string, args ...interface{}) (*Rows, error) {
	tokens, err := tokenizeSQL(query)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 || strings.ToUpper(tokens[0]) != "SELECT" {
		return nil, fmt.Errorf("query must start with SELECT")
	}
	return db.execSelect(tokens[1:], args)
}

func (db *Database) execCreate(tokens []string) error {
	if len(tokens) < 4 || strings.ToUpper(tokens[0]) != "TABLE" ||
		tokens[2] != "(" || tokens[len(tokens)-1] != ")" {
		return fmt.Errorf("invalid CREATE TABLE syntax")
	}
	name := tokens[1]
	colTokens := tokens[3 : len(tokens)-1]
	var columns []string
	for i := 0; i < len(colTokens); i++ {
		if colTokens[i] == "," {
			continue
		}
		columns = append(columns, colTokens[i])
		if i+1 < len(colTokens) && colTokens[i+1] != "," {
			i++
		}
	}
	if len(columns) == 0 {
		return fmt.Errorf("no columns defined for table %q", name)
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.tables[name]; ok {
		return fmt.Errorf("table %q already exists", name)
	}
	db.tables[name] = &table{
		columns: columns,
		rows:    make([]map[string]interface{}, 0),
	}
	return nil
}

func (db *Database) execInsert(tokens []string, args []interface{}) error {
	if len(tokens) < 4 || strings.ToUpper(tokens[0]) != "INTO" ||
		strings.ToUpper(tokens[2]) != "VALUES" || tokens[3] != "(" {
		return fmt.Errorf("invalid INSERT syntax")
	}
	name := tokens[1]
	// Collect values up to closing paren
	var valTokens []string
	i := 4
	for i < len(tokens) && tokens[i] != ")" {
		valTokens = append(valTokens, tokens[i])
		i++
	}
	if i >= len(tokens) {
		return fmt.Errorf("missing closing paren in INSERT")
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	t, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("table %q not found", name)
	}
	row := make(map[string]interface{})
	argIdx := 0
	colIdx := 0
	for _, vt := range valTokens {
		if vt == "," {
			continue
		}
		if colIdx >= len(t.columns) {
			return fmt.Errorf("too many values for table %q", name)
		}
		val := resolveValue(vt, &argIdx, args)
		row[t.columns[colIdx]] = val
		colIdx++
	}
	t.rows = append(t.rows, row)
	return nil
}

func (db *Database) execSelect(tokens []string, args []interface{}) (*Rows, error) {
	fromIdx := -1
	for i, tok := range tokens {
		if strings.ToUpper(tok) == "FROM" {
			fromIdx = i
			break
		}
	}
	if fromIdx < 0 || fromIdx+1 >= len(tokens) {
		return nil, fmt.Errorf("invalid SELECT syntax, expected FROM")
	}
	name := tokens[fromIdx+1]
	whereIdx := -1
	for i, tok := range tokens[fromIdx+2:] {
		if strings.ToUpper(tok) == "WHERE" {
			whereIdx = fromIdx + 2 + i
			break
		}
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	t, ok := db.tables[name]
	if !ok {
		return nil, fmt.Errorf("table %q not found", name)
	}
	var conditions []condition
	if whereIdx >= 0 {
		argIdx := 0
		var err error
		conditions, err = parseConditions(tokens[whereIdx+1:], &argIdx, args)
		if err != nil {
			return nil, err
		}
	}
	var results []map[string]interface{}
	for _, row := range t.rows {
		if matchConditions(row, conditions) {
			results = append(results, row)
		}
	}
	return &Rows{columns: t.columns, rows: results}, nil
}

func (db *Database) execUpdate(tokens []string, args []interface{}) error {
	if len(tokens) < 3 || strings.ToUpper(tokens[1]) != "SET" {
		return fmt.Errorf("invalid UPDATE syntax")
	}
	name := tokens[0]
	setTokens := tokens[2:]
	whereIdx := -1
	for i, tok := range setTokens {
		if strings.ToUpper(tok) == "WHERE" {
			whereIdx = i
			break
		}
	}
	var setPart []string
	if whereIdx >= 0 {
		setPart = setTokens[:whereIdx]
	} else {
		setPart = setTokens
	}
	updates := make(map[string]interface{})
	argIdx := 0
	for i := 0; i < len(setPart); i++ {
		if setPart[i] == "," {
			continue
		}
		if i+2 >= len(setPart) {
			return fmt.Errorf("invalid SET clause")
		}
		col := setPart[i]
		if setPart[i+1] != "=" {
			return fmt.Errorf("expected '=' in SET clause")
		}
		updates[col] = resolveValue(setPart[i+2], &argIdx, args)
		i += 2
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	t, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("table %q not found", name)
	}
	var conditions []condition
	if whereIdx >= 0 {
		var err error
		conditions, err = parseConditions(setTokens[whereIdx+1:], &argIdx, args)
		if err != nil {
			return err
		}
	}
	for _, row := range t.rows {
		if matchConditions(row, conditions) {
			for k, v := range updates {
				row[k] = v
			}
		}
	}
	return nil
}

func (db *Database) execDelete(tokens []string, args []interface{}) error {
	fromIdx := -1
	for i, tok := range tokens {
		if strings.ToUpper(tok) == "FROM" {
			fromIdx = i
			break
		}
	}
	if fromIdx < 0 || fromIdx+1 >= len(tokens) {
		return fmt.Errorf("invalid DELETE syntax")
	}
	name := tokens[fromIdx+1]
	whereIdx := -1
	for i, tok := range tokens[fromIdx+2:] {
		if strings.ToUpper(tok) == "WHERE" {
			whereIdx = fromIdx + 2 + i
			break
		}
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	t, ok := db.tables[name]
	if !ok {
		return fmt.Errorf("table %q not found", name)
	}
	var conditions []condition
	if whereIdx >= 0 {
		argIdx := 0
		var err error
		conditions, err = parseConditions(tokens[whereIdx+1:], &argIdx, args)
		if err != nil {
			return err
		}
	}
	var kept []map[string]interface{}
	for _, row := range t.rows {
		if !matchConditions(row, conditions) {
			kept = append(kept, row)
		}
	}
	t.rows = kept
	return nil
}

func (db *Database) execDrop(tokens []string) error {
	if len(tokens) < 2 || strings.ToUpper(tokens[0]) != "TABLE" {
		return fmt.Errorf("invalid DROP syntax")
	}
	name := tokens[1]
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.tables[name]; !ok {
		return fmt.Errorf("table %q not found", name)
	}
	delete(db.tables, name)
	return nil
}

type condition struct {
	column string
	op     string
	value  interface{}
}

func parseConditions(tokens []string, argIdx *int, args []interface{}) ([]condition, error) {
	var conds []condition
	for i := 0; i < len(tokens); i++ {
		if strings.ToUpper(tokens[i]) == "AND" {
			continue
		}
		if i+2 >= len(tokens) {
			return nil, fmt.Errorf("incomplete condition at %q", tokens[i])
		}
		col := tokens[i]
		op := tokens[i+1]
		val := resolveValue(tokens[i+2], argIdx, args)
		conds = append(conds, condition{column: col, op: op, value: val})
		i += 2
	}
	return conds, nil
}

func resolveValue(tok string, argIdx *int, args []interface{}) interface{} {
	if tok == "?" {
		if *argIdx < len(args) {
			v := args[*argIdx]
			*argIdx++
			return v
		}
		return nil
	}
	if tok == "NULL" {
		return nil
	}
	if len(tok) >= 2 && (tok[0] == '\'' || tok[0] == '"') {
		return tok[1 : len(tok)-1]
	}
	if n, err := strconv.ParseInt(tok, 10, 64); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(tok, 64); err == nil {
		return f
	}
	if strings.ToUpper(tok) == "TRUE" {
		return true
	}
	if strings.ToUpper(tok) == "FALSE" {
		return false
	}
	return tok
}

func matchConditions(row map[string]interface{}, conds []condition) bool {
	if len(conds) == 0 {
		return true
	}
	for _, c := range conds {
		rv, rok := row[c.column]
		if !rok {
			return false
		}
		if !compareValues(rv, c.value, c.op) {
			return false
		}
	}
	return true
}

func compareValues(a, b interface{}, op string) bool {
	switch op {
	case "=":
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	case "!=", "<>":
		return fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b)
	case ">":
		af, aok := toFloat(a)
		bf, bok := toFloat(b)
		if aok && bok {
			return af > bf
		}
		return fmt.Sprintf("%v", a) > fmt.Sprintf("%v", b)
	case "<":
		af, aok := toFloat(a)
		bf, bok := toFloat(b)
		if aok && bok {
			return af < bf
		}
		return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
	case ">=":
		af, aok := toFloat(a)
		bf, bok := toFloat(b)
		if aok && bok {
			return af >= bf
		}
		return fmt.Sprintf("%v", a) >= fmt.Sprintf("%v", b)
	case "<=":
		af, aok := toFloat(a)
		bf, bok := toFloat(b)
		if aok && bok {
			return af <= bf
		}
		return fmt.Sprintf("%v", a) <= fmt.Sprintf("%v", b)
	}
	return false
}

func toFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	}
	return 0, false
}

// Begin starts a new transaction with snapshot isolation. All subsequent
// operations on the returned Tx see a consistent snapshot of the database.
func (db *Database) Begin() (*Tx, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	snapshot := make(map[string]*table, len(db.tables))
	for name, t := range db.tables {
		tCopy := &table{
			columns: make([]string, len(t.columns)),
			rows:    make([]map[string]interface{}, len(t.rows)),
		}
		copy(tCopy.columns, t.columns)
		for i, row := range t.rows {
			tCopy.rows[i] = make(map[string]interface{}, len(row))
			for k, v := range row {
				tCopy.rows[i][k] = v
			}
		}
		snapshot[name] = tCopy
	}
	return &Tx{db: db, snapshot: snapshot, active: true}, nil
}

// Execute runs a SQL statement within the transaction snapshot.
func (tx *Tx) Execute(query string, args ...interface{}) error {
	tokens, err := tokenizeSQL(query)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return fmt.Errorf("empty query")
	}
	cmd := strings.ToUpper(tokens[0])
	// Reuse Database's methods by temporarily swapping tables
	origTables := tx.db.tables
	tx.db.mu.Lock()
	tx.db.tables = tx.snapshot
	tx.db.mu.Unlock()

	defer func() {
		tx.db.mu.Lock()
		tx.db.tables = origTables
		tx.db.mu.Unlock()
	}()

	switch cmd {
	case "CREATE":
		return tx.db.execCreate(tokens[1:])
	case "INSERT":
		return tx.db.execInsert(tokens[1:], args)
	case "UPDATE":
		return tx.db.execUpdate(tokens[1:], args)
	case "DELETE":
		return tx.db.execDelete(tokens[1:], args)
	case "DROP":
		return tx.db.execDrop(tokens[1:])
	default:
		return fmt.Errorf("unknown statement type %q", cmd)
	}
}

// Query runs a SELECT within the transaction snapshot.
func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error) {
	tokens, err := tokenizeSQL(query)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 || strings.ToUpper(tokens[0]) != "SELECT" {
		return nil, fmt.Errorf("query must start with SELECT")
	}
	origTables := tx.db.tables
	tx.db.mu.Lock()
	tx.db.tables = tx.snapshot
	tx.db.mu.Unlock()

	defer func() {
		tx.db.mu.Lock()
		tx.db.tables = origTables
		tx.db.mu.Unlock()
	}()

	return tx.db.execSelect(tokens[1:], args)
}

// Commit writes the transaction snapshot to the database.
func (tx *Tx) Commit() error {
	if !tx.active {
		return fmt.Errorf("transaction is not active")
	}
	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()
	tx.db.tables = tx.snapshot
	tx.active = false
	return nil
}

// Rollback discards the transaction snapshot.
func (tx *Tx) Rollback() error {
	if !tx.active {
		return fmt.Errorf("transaction is not active")
	}
	tx.active = false
	return nil
}

func tokenizeSQL(s string) ([]string, error) {
	var tokens []string
	var buf strings.Builder
	flush := func() {
		if buf.Len() > 0 {
			tokens = append(tokens, buf.String())
			buf.Reset()
		}
	}
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			flush()
			i++
		case c == '(' || c == ')' || c == ',' || c == '*' || c == '?':
			flush()
			tokens = append(tokens, string(c))
			i++
		case c == '=' || c == '!' || c == '>' || c == '<':
			flush()
			if i+1 < len(s) && s[i+1] == '=' {
				tokens = append(tokens, string(c)+"=")
				i += 2
			} else if c == '<' && i+1 < len(s) && s[i+1] == '>' {
				tokens = append(tokens, "<>")
				i += 2
			} else {
				tokens = append(tokens, string(c))
				i++
			}
		case c == '\'' || c == '"':
			flush()
			j := i + 1
			for j < len(s) && s[j] != c {
				if s[j] == '\\' {
					j += 2
				} else {
					j++
				}
			}
			if j >= len(s) {
				return nil, fmt.Errorf("unterminated string literal")
			}
			tokens = append(tokens, s[i:j+1])
			i = j + 1
		default:
			buf.WriteByte(c)
			i++
		}
	}
	flush()
	return tokens, nil
}
