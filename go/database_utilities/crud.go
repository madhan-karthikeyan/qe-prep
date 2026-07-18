package database_utilities

import "fmt"

// Create inserts a new row into the given table with the provided data. The
// data keys must match the table's column names.
func Create(db *Database, table string, data map[string]interface{}) error {
	cols, err := db.TableColumns(table)
	if err != nil {
		return err
	}
	// Build placeholders
	placeholders := make([]string, len(cols))
	args := make([]interface{}, len(cols))
	for i, c := range cols {
		placeholders[i] = "?"
		if v, ok := data[c]; ok {
			args[i] = v
		} else {
			args[i] = nil
		}
	}
	query := fmt.Sprintf("INSERT INTO %s VALUES (%s)", table, join(placeholders, ", "))
	return db.Execute(query, args...)
}

// Read retrieves rows from the given table matching the conditions.
func Read(db *Database, table string, conditions map[string]interface{}) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	var args []interface{}
	if len(conditions) > 0 {
		query += " WHERE "
		var conds []string
		for k, v := range conditions {
			conds = append(conds, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
		query += join(conds, " AND ")
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for rows.Next() {
		dests := make([]interface{}, len(rows.Columns()))
		for i := range dests {
			dests[i] = new(string)
		}
		if err := rows.Scan(dests...); err != nil {
			return nil, err
		}
		row := make(map[string]interface{})
		for i, col := range rows.Columns() {
			row[col] = *(dests[i].(*string))
		}
		result = append(result, row)
	}
	return result, nil
}

// Update modifies rows in the given table matching the conditions.
func Update(db *Database, table string, updates map[string]interface{}, conditions map[string]interface{}) (int, error) {
	query := fmt.Sprintf("UPDATE %s SET", table)
	var setClauses []string
	var args []interface{}
	for k, v := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s=?", k))
		args = append(args, v)
	}
	query += " " + join(setClauses, ", ")
	if len(conditions) > 0 {
		query += " WHERE "
		var conds []string
		for k, v := range conditions {
			conds = append(conds, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
		query += join(conds, " AND ")
	}
	if err := db.Execute(query, args...); err != nil {
		return 0, err
	}
	return 0, nil
}

// Delete removes rows from the given table matching the conditions.
func Delete(db *Database, table string, conditions map[string]interface{}) (int, error) {
	query := fmt.Sprintf("DELETE FROM %s", table)
	var args []interface{}
	if len(conditions) > 0 {
		query += " WHERE "
		var conds []string
		for k, v := range conditions {
			conds = append(conds, fmt.Sprintf("%s = ?", k))
			args = append(args, v)
		}
		query += join(conds, " AND ")
	}
	if err := db.Execute(query, args...); err != nil {
		return 0, err
	}
	return 0, nil
}

func join(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	result := elems[0]
	for _, e := range elems[1:] {
		result += sep + e
	}
	return result
}
