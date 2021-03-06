package dbr

// DeleteStmt builds `DELETE ...`
type DeleteStmt struct {
	raw

	Table string

	WhereCond []Builder
	JoinTable []Builder
}

// Build builds `DELETE ...` in dialect
func (b *DeleteStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	buf.WriteString("DELETE ")
	if len(b.JoinTable) > 0 {
		buf.WriteString(d.QuoteIdent(b.Table) + " ")
	}
	buf.WriteString("FROM ")
	buf.WriteString(d.QuoteIdent(b.Table))

	if len(b.JoinTable) > 0 {
		for _, join := range b.JoinTable {
			err := join.Build(d, buf)
			if err != nil {
				return err
			}
		}
	}

	if len(b.WhereCond) > 0 {
		buf.WriteString(" WHERE ")
		err := And(b.WhereCond...).Build(d, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteFrom creates a DeleteStmt
func DeleteFrom(table string) *DeleteStmt {
	return &DeleteStmt{
		Table: table,
	}
}

// DeleteBySql creates a DeleteStmt from raw query
func DeleteBySql(query string, value ...interface{}) *DeleteStmt {
	return &DeleteStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
	}
}

// Where adds a where condition
func (b *DeleteStmt) Where(query interface{}, value ...interface{}) *DeleteStmt {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Builder:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}

// Join joins table on condition
func (b *DeleteStmt) Join(table, on interface{}) *DeleteStmt {
	b.JoinTable = append(b.JoinTable, join(inner, table, on))
	return b
}

// Join joins table on condition
func (b *DeleteStmt) LeftJoin(table, on interface{}) *DeleteStmt {
	b.JoinTable = append(b.JoinTable, join(left, table, on))
	return b
}
