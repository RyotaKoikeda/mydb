package mysql

import (
	"database/sql"
	"time"

	"utils/log"

	_ "github.com/go-sql-driver/mysql"
)

var logId = ""

type SqlHandler interface {
	Execute(string, ...interface{}) (Result, error)
	Query(string, ...interface{}) (Rows, error)
	Close() error
	CheckInUse() (int, int)
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

type SqlHandlerImpl struct {
	Conn     *sql.DB
	ConnRead *sql.DB
}

func (handler *SqlHandlerImpl) Execute(statement string, args ...interface{}) (Result, error) {
	res := SqlResult{}
	result, err := handler.Conn.Exec(statement, args...)
	if err != nil {
		log.Error(
			err.Error(),
			log.String("statement", statement),
			log.Object("args", args),
		)
		return nil, err
	}
	res.Result = result

	// 追加: ログIDが一致していない場合、最新のログIDをセット
	if logId != log.GetLogId() {
		logId = log.GetLogId()
	}
	return res, nil
}

func (handler *SqlHandlerImpl) Query(statement string, args ...interface{}) (Rows, error) {
	// 追加: ログIDが一致している場合、ライターエンドポイントを使用
	//      していない場合は、リーダーエンドポイントを使用
	conn := handler.ConnRead
	if logId == log.GetLogId() {
		conn = handler.Conn
	}
	rows, err := conn.Query(statement, args...)
	if err != nil {
		log.Error(
			err.Error(),
			log.String("statement", statement),
			log.Object("args", args),
		)
		return nil, err
	}
	return SqlRow{rows}, nil
}

func (handler *SqlHandlerImpl) CheckInUse() (int, int) {
	connMax := handler.Conn.Stats().MaxOpenConnections
	connInUse := handler.Conn.Stats().InUse
	connReadMax := handler.ConnRead.Stats().MaxOpenConnections
	connReadInUse := handler.ConnRead.Stats().InUse
	now := time.Now()

	// 04:00~04:10(JST)の間のリクエストのみチェックを行う
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 19, 0, 0, 0, time.UTC)
	endTime := startTime.Add(10 * time.Minute)

	if now.After(startTime) && now.Before(endTime) {
		if connMax > 1 {
			if connInUse > connMax/2 {
				log.Warn("number of db connections in use has exceeded the threshold: ConnWriter")
			}
		}
		if connReadMax > 1 {
			if connReadInUse > connReadMax/2 {
				log.Warn("number of db connections in use has exceeded the threshold: ConnReader")
			}
		}
	}

	return connInUse, connReadInUse
}

func (handler *SqlHandlerImpl) Close() error {
	err := handler.Conn.Close()
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

type SqlResult struct {
	Result sql.Result
}

func (r SqlResult) LastInsertId() (int64, error) {
	lastInsertId, err := r.Result.LastInsertId()
	if err != nil {
		log.Error(err.Error())
	}
	return lastInsertId, err
}

func (r SqlResult) RowsAffected() (int64, error) {
	rowsAffected, err := r.Result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
	}
	return rowsAffected, err
}

type SqlRow struct {
	Rows *sql.Rows
}

func (r SqlRow) Scan(dest ...interface{}) error {
	err := r.Rows.Scan(dest...)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func (r SqlRow) Next() bool {
	return r.Rows.Next()
}

func (r SqlRow) Close() error {
	err := r.Rows.Close()
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

type SqlHandlerParamsGetter interface {
	GetMysqlUser() string
	GetMysqlHost() string
	GetMysqlReadHost() string
	GetMysqlPort() string
	GetMysqlPassword() string
	GetMysqlDB() string
	GetDbConnMaxIdleTime() time.Duration
	GetDbConnMaxLifetime() time.Duration
	GetDbMaxIdleConns() int
	GetDbMaxOpenConns() int
}

func NewSqlHandler(params SqlHandlerParamsGetter) *SqlHandlerImpl {

	// ２段階にしている理由は、DBが存在しない場合に作成するため
	// (主にapitestでデータベースが選択されていないことを阻止するため)

	// TODO: apitest用DB接続処理
	// // データベース選択せずに接続
	// conn, err := sql.Open("mysql", getConnectionString(
	// 	params.GetMysqlUser(),
	// 	params.GetMysqlPassword(),
	// 	params.GetMysqlHost(),
	// 	params.GetMysqlPort(),
	// 	"",
	// ))
	// if err != nil {
	// 	panic(err)
	// }
	// // envのDB名が存在しなければ作成
	// _, err = conn.Exec("CREATE DATABASE IF NOT EXISTS " + params.GetMysqlDB())
	// if err != nil {
	// 	panic(err)
	// }
	// // 一旦閉じる
	// conn.Close()

	// 作成したDBに接続
	connV2, err := sql.Open("mysql", getConnectionString(
		params.GetMysqlUser(),
		params.GetMysqlPassword(),
		params.GetMysqlHost(),
		params.GetMysqlPort(),
		params.GetMysqlDB(),
	))

	if err != nil {
		panic(err)
	}
	connV2.SetConnMaxIdleTime(params.GetDbConnMaxIdleTime())
	connV2.SetConnMaxLifetime(params.GetDbConnMaxLifetime())
	connV2.SetMaxIdleConns(params.GetDbMaxIdleConns())
	connV2.SetMaxOpenConns(params.GetDbMaxOpenConns())

	// 読み取り用クラスターに接続
	connRead, err := sql.Open("mysql", getConnectionString(
		params.GetMysqlUser(),
		params.GetMysqlPassword(),
		params.GetMysqlReadHost(),
		params.GetMysqlPort(),
		params.GetMysqlDB(),
	))

	if err != nil {
		panic(err)
	}
	connRead.SetConnMaxIdleTime(params.GetDbConnMaxIdleTime())
	connRead.SetConnMaxLifetime(params.GetDbConnMaxLifetime())
	connRead.SetMaxIdleConns(params.GetDbMaxIdleConns())
	connRead.SetMaxOpenConns(params.GetDbMaxOpenConns())

	return &SqlHandlerImpl{connV2, connRead}
}

func getConnectionString(user string, pass string, host string, port string, db string) string {
	return user + ":" + pass + "@(" + host + ":" + port + ")/" + db
}
