package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQL struct {
	DB *gorm.DB
}

// New はGORMを使用してMySQLデータベースに接続するための新しいMySQLインスタンスを生成します。
func New() (*MySQL, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	// SQL非表示設定か確認する。
	if os.Getenv("IS_HIDDEN_SQL") == "true" {
		newLogger = logger.Default.LogMode(logger.Silent)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.ExpandEnv("${MYSQL_USER}"), os.ExpandEnv("${MYSQL_PASSWORD}"),
		os.ExpandEnv("${DB_HOST}"), os.ExpandEnv("${DB_PORT}"),
		os.ExpandEnv("${MYSQL_DATABASE}"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &MySQL{DB: db}, nil
}

// NewTest はGORMを使用してMySQLデータベースに接続するための新しいMySQLインスタンスを生成します。
func NewTest() (*MySQL, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	if os.Getenv("IS_HIDDEN_TEST_SQL") == "true" {
		newLogger = logger.Default.LogMode(logger.Silent)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.ExpandEnv("${MYSQL_USER}"), os.ExpandEnv("${MYSQL_PASSWORD}"),
		os.ExpandEnv("${DB_HOST}"), os.ExpandEnv("${DB_PORT}"),
		os.ExpandEnv("${MYSQL_TEST_DATABASE}"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &MySQL{DB: db}, nil
}

// CreateNewTestDB はランダムな名前でDBを作成し、そのインスタンスを返します。
// また、deferで削除を予約します。
func CreateNewTestDB() (*MySQL, func() error, error) {
	randomDbName, err := generateUniqueID()
	if err != nil {
		return nil, nil, err
	}

	dbName := fmt.Sprintf("%s_test", randomDbName)
	if err := createMySQLDatabase(dbName); err != nil {
		return nil, nil, fmt.Errorf("failed to create database: %w", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("DB_HOST"),
		dbName,
	)

	// 作成したDBに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 作成したDBを削除
	cleanUp := func() error {
		err := deleteMySQLDatabase(dbName)
		if err != nil {
			fmt.Printf("failed to create database: %v", err)
		}
		return err
	}

	return &MySQL{DB: db}, cleanUp, nil
}

// createMySQLDatabase DBを作成する。
func createMySQLDatabase(dbName string) (err error) {
	// rootユーザーでの接続情報を設定
	dsn := fmt.Sprintf("root:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_PASSWORD"), os.Getenv("DB_HOST"))

	// GormでrootユーザーとしてDBに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database as root: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get generic database object: %w", err)
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close db: %w", closeErr)
		}
	}()

	// テスト用のDBを作成
	if err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// テストで使うユーザーに権限を付与
	if err = db.Exec(fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%%'", dbName, os.Getenv("MYSQL_USER"))).Error; err != nil {
		return fmt.Errorf("failed to grant privileges: %w", err)
	}

	return nil
}

func deleteMySQLDatabase(dbName string) error {
	// rootユーザーでの接続情報を設定
	dsn := fmt.Sprintf("root:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_PASSWORD"), os.Getenv("DB_HOST"))

	// GormでrootユーザーとしてDBに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database as root: %w", err)
	}

	// データベース削除
	if err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)).Error; err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}

	return nil
}

// generateUniqueID はランダムな文字列を生成します。
func generateUniqueID() (string, error) {
	// カスタムアルファベットを定義 (特定の文字を除外)
	alphabet := "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz0123456789"

	// IDの長さを設定
	size := 21 // UUIDと同じ長さに設定

	// カスタムアルファベットとサイズを使用してIDを生成
	id, err := nanoid.GenerateString(alphabet, size)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Transactional は新しいトランザクションを開始しインスタンスを返す関数です。
// 戻り値の関数はcleanUPとして受け取り、"defer cleanUP()"を直下の行に記述してください。
func Transactional(db *gorm.DB) (*gorm.DB, func()) {
	tx := db.Begin()
	if tx.Error != nil {
		panic("トランザクションの開始に失敗しました。")
	}

	// エラーハンドリングとロールバックを行うクロージャを返す。
	return tx, func() {
		// panicによるエラーの場合
		if r := recover(); r != nil {
			fmt.Printf("panic recovered: %v\n", r) // ← 追加
			fmt.Printf("予期せぬエラーが発生したため、トランザクションを巻き戻しました。\n")
			tx.Rollback()
		}

		// tx.Errorが設定されている場合（明示的なエラー設定）
		if tx.Error != nil {
			fmt.Printf("エラーが発生したため、トランザクションを巻き戻しました:\n %v\n", tx.Error)
			tx.Rollback()
		}

		tx.Commit()
	}
}
