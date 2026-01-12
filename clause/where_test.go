package clause_test

import (
	"fmt"
	"testing"

	"gorm.io/gorm/clause"
)

func TestWhere(t *testing.T) {
	results := []struct {
		Clauses []clause.Interface
		Result  string
		Vars    []interface{}
	}{
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}, clause.Or(clause.Neq{Column: "name", Value: "jinzhu"})},
			}},
			"SELECT * FROM `users` WHERE `users`.`id` = ? AND `age` > ? OR `name` <> ?",
			[]interface{}{"1", 18, "jinzhu"},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Or(clause.Neq{Column: "name", Value: "jinzhu"}), clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}},
			}},
			"SELECT * FROM `users` WHERE `users`.`id` = ? OR `name` <> ? AND `age` > ?",
			[]interface{}{"1", "jinzhu", 18},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Or(clause.Neq{Column: "name", Value: "jinzhu"}), clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}},
			}},
			"SELECT * FROM `users` WHERE `users`.`id` = ? OR `name` <> ? AND `age` > ?",
			[]interface{}{"1", "jinzhu", 18},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Or(clause.Eq{Column: clause.PrimaryColumn, Value: "1"}), clause.Or(clause.Neq{Column: "name", Value: "jinzhu"})},
			}},
			"SELECT * FROM `users` WHERE `users`.`id` = ? OR `name` <> ?",
			[]interface{}{"1", "jinzhu"},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}, clause.Or(clause.Neq{Column: "name", Value: "jinzhu"})},
			}, clause.Where{
				Exprs: []clause.Expression{clause.Or(clause.Gt{Column: "score", Value: 100}, clause.Like{Column: "name", Value: "%linus%"})},
			}},
			"SELECT * FROM `users` WHERE `users`.`id` = ? AND `age` > ? OR `name` <> ? AND (`score` > ? OR `name` LIKE ?)",
			[]interface{}{"1", 18, "jinzhu", 100, "%linus%"},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Not(clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}), clause.Or(clause.Neq{Column: "name", Value: "jinzhu"})},
			}, clause.Where{
				Exprs: []clause.Expression{clause.Or(clause.Not(clause.Gt{Column: "score", Value: 100}), clause.Like{Column: "name", Value: "%linus%"})},
			}},
			"SELECT * FROM `users` WHERE (`users`.`id` <> ? AND `age` <= ?) OR `name` <> ? AND (`score` <= ? OR `name` LIKE ?)",
			[]interface{}{"1", 18, "jinzhu", 100, "%linus%"},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.And(clause.Eq{Column: "age", Value: 18}, clause.Or(clause.Neq{Column: "name", Value: "jinzhu"}))},
			}},
			"SELECT * FROM `users` WHERE `age` = ? OR `name` <> ?",
			[]interface{}{18, "jinzhu"},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Not(clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}), clause.And(clause.Expr{SQL: "`score` <= ?", Vars: []interface{}{100}, WithoutParentheses: false})},
			}},
			"SELECT * FROM `users` WHERE (`users`.`id` <> ? AND `age` <= ?) AND `score` <= ?",
			[]interface{}{"1", 18, 100},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Not(clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}), clause.Expr{SQL: "`score` <= ?", Vars: []interface{}{100}, WithoutParentheses: false}},
			}},
			"SELECT * FROM `users` WHERE (`users`.`id` <> ? AND `age` <= ?) AND `score` <= ?",
			[]interface{}{"1", 18, 100},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Not(clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}), clause.Or(clause.Expr{SQL: "`score` <= ?", Vars: []interface{}{100}, WithoutParentheses: false})},
			}},
			"SELECT * FROM `users` WHERE (`users`.`id` <> ? AND `age` <= ?) OR `score` <= ?",
			[]interface{}{"1", 18, 100},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.And(clause.Not(clause.Eq{Column: clause.PrimaryColumn, Value: "1"}),
						clause.And(clause.Expr{SQL: "`score` <= ?", Vars: []interface{}{100}, WithoutParentheses: false})),
				},
			}},
			"SELECT * FROM `users` WHERE `users`.`id` <> ? AND `score` <= ?",
			[]interface{}{"1", 100},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Not(clause.Eq{Column: clause.PrimaryColumn, Value: "1"},
					clause.And(clause.Expr{SQL: "`score` <= ?", Vars: []interface{}{100}, WithoutParentheses: false}))},
			}},
			"SELECT * FROM `users` WHERE (`users`.`id` <> ? AND NOT `score` <= ?)",
			[]interface{}{"1", 100},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{clause.Not(clause.Expr{SQL: "`score` <= ?", Vars: []interface{}{100}},
					clause.Expr{SQL: "`age` <= ?", Vars: []interface{}{60}})},
			}},
			"SELECT * FROM `users` WHERE NOT (`score` <= ? AND `age` <= ?)",
			[]interface{}{100, 60},
		},
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Not(clause.AndConditions{
						Exprs: []clause.Expression{
							clause.Eq{Column: clause.PrimaryColumn, Value: "1"},
							clause.Gt{Column: "age", Value: 18},
						}}, clause.OrConditions{
						Exprs: []clause.Expression{
							clause.Lt{Column: "score", Value: 100},
						},
					}),
				}}},
			"SELECT * FROM `users` WHERE NOT ((`users`.`id` = ? AND `age` > ?) OR `score` < ?)",
			[]interface{}{"1", 18, 100},
		},
		// Test single-line OR still works (space)
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1 OR b = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1 OR b = 2)",
			[]interface{}{"1"},
		},
		// Test single-line AND still works (space)
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Or(clause.Expr{SQL: "a = 1 AND b = 2"}),
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? OR (a = 1 AND b = 2)",
			[]interface{}{"1"},
		},
		// Test OR with newline after
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1 OR\nb = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1 OR\nb = 2)",
			[]interface{}{"1"},
		},
		// Test OR with newline before
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1\nOR b = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1\nOR b = 2)",
			[]interface{}{"1"},
		},
		// Test OR with tab after
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1 OR\tb = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1 OR\tb = 2)",
			[]interface{}{"1"},
		},
		// Test OR with tab before
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1\tOR b = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1\tOR b = 2)",
			[]interface{}{"1"},
		},
		// Test OR with carriage return
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1 OR\rb = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1 OR\rb = 2)",
			[]interface{}{"1"},
		},
		// Test OR with form feed
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1 OR\fb = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1 OR\fb = 2)",
			[]interface{}{"1"},
		},
		// Test OR with vertical tab
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1 OR\vb = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1 OR\vb = 2)",
			[]interface{}{"1"},
		},
		// Test AND with newline
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Or(clause.Expr{SQL: "a = 1 AND\nb = 2"}),
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? OR (a = 1 AND\nb = 2)",
			[]interface{}{"1"},
		},
		// Test multi-line with multiple whitespace characters
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "a = 1 OR\n\t  b = 2"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND (a = 1 OR\n\t  b = 2)",
			[]interface{}{"1"},
		},
		// Test that ORACLE,ORPORATION etc. are not matched (OR must be surrounded by whitespace)
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "name = 'ORACLE'"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND name = 'ORACLE'",
			[]interface{}{"1"},
		},
		// Test that ANDROID, COMMAND etc. are not matched (AND must be surrounded by whitespace)
		{
			[]clause.Interface{clause.Select{}, clause.From{}, clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{Column: "org_id", Value: "1"},
					clause.Expr{SQL: "name = 'ANDROID'"},
				}}},
			"SELECT * FROM `users` WHERE `org_id` = ? AND name = 'ANDROID'",
			[]interface{}{"1"},
		},
	}

	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			checkBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}
