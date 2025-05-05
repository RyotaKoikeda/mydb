module mydb

go 1.23.3

replace (
	utils/log => ../utils/log
	utils/rdb/mysql => ../utils/rdb/mysql
)

require (
	github.com/go-sql-driver/mysql v1.9.2
	utils/log v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)
