package database

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

var validName = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// DB, "table1", "column1", {
// "id": "user-123762384",
// "username": "user",
// "email": "example@email.com",
// "password": "h5yf6874g84hf48"}
func Create(db *sqlx.DB, table string, rowData map[string]any) error {
	if !validName.MatchString(table) {
		return fmt.Errorf("Invalid table name: %s", table)
	}
	cols, placeholders := parseRowDataKeys(rowData)

	for column := range rowData {
		if !validName.MatchString(column) {
			return fmt.Errorf("Invalid column name in updates: %s", table)
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
	_, err := db.NamedExec(query, rowData)

	if err != nil {
		return err
	}
	return nil
}

// DB, "table1", {"column1", "column2"}, {"1", "2"}
func Read(db *sqlx.DB, table string, columns []string, values []any) ([]map[string]any, error) {
	if !validName.MatchString(table) {
		return nil, fmt.Errorf("Invalid table name: %s", table)
	}

	for _, column := range columns {
		if !validName.MatchString(column) {
			return nil, fmt.Errorf("Invalid column name: %s", column)
		}
	}

	if len(columns) != len(values) {
		return nil, fmt.Errorf("Mismatched number of columns and values")
	}

	whereClauses := make([]string, len(columns))
	for i, col := range columns {
		whereClauses[i] = fmt.Sprintf("%s = $%d", col, i+1)
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s",
		table,
		strings.Join(whereClauses, " AND "),
	)

	rows, err := db.Queryx(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]any
	for rows.Next() {
		rowData := make(map[string]any)
		if err := rows.MapScan(rowData); err != nil {
			return nil, err
		}
		result = append(result, rowData)
	}
	return result, nil
}

// DB, "table1", {
// "username": "newname",
// "email": "new@email.com",
// },
// {
// "id": "user-213764738"
// }
func Update(db *sqlx.DB, table string, updates map[string]any, where map[string]any) (int64, error) {
	if !validName.MatchString(table) {
		return 0, fmt.Errorf("Invalid table name: %s", table)
	}

	if len(updates) == 0 {
		return 0, fmt.Errorf("Updates is empty")
	}

	if len(where) == 0 {
		return 0, fmt.Errorf("Where is empty")
	}

	for column := range updates {
		if !validName.MatchString(column) {
			return 0, fmt.Errorf("Invalid column name in updates: %s", table)
		}
	}

	for column := range where {
		if !validName.MatchString(column) {
			return 0, fmt.Errorf("Invalid column name in where: %s", table)
		}
	}

	updateCols, updatePlaceholders := parseRowDataKeys(updates)
	whereCols, wherePlaceholders := parseRowDataKeys(where)

	setParts := make([]string, len(updateCols))
	for i, col := range updateCols {
		setParts[i] = fmt.Sprintf("%s = %s", col, updatePlaceholders[i])
	}
	whereParts := make([]string, len(whereCols))
	for i, col := range whereCols {
		whereParts[i] = fmt.Sprintf("%s = %s", col, wherePlaceholders[i])
	}
	params := map[string]any{}
	for k, v := range updates {
		params[k] = v
	}
	for k, v := range where {
		params[k] = v
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(setParts, ", "),
		strings.Join(whereParts, " AND "),
	)
	result, err := db.NamedExec(query, params)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// DB, "table1", {
// "id": "user-213764738"
// }
func Delete(db *sqlx.DB, table string, where map[string]any) (int64, error) {
	if !validName.MatchString(table) {
		return 0, fmt.Errorf("Invalid table name: %s", table)
	}

	if len(where) == 0 {
		return 0, fmt.Errorf("Where is empty")
	}

	for col := range where {
		if !validName.MatchString(col) {
			return 0, fmt.Errorf("Invalid column name: %s", col)
		}
	}

	cols, placeholders := parseRowDataKeys(where)

	query := fmt.Sprintf(
		"DELETE FROM % s WHERE %s = %s",
		table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, " AND "),
	)

	result, err := db.NamedExec(query, where)
	if err != nil {
		return 0, err
	}
	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsDeleted, nil
}
func parseRowDataKeys(rowData map[string]any) ([]string, []string) {
	cols := make([]string, 0, len(rowData))
	for col := range rowData {
		cols = append(cols, col)
	}
	sort.Strings(cols)
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = ":" + cols[i]
	}
	return cols, placeholders
}
