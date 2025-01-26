package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

const (
	defaultMaxIdleConn     = 10
	defaultMaxOpenConn     = 100
	defaultConnMaxLifetime = time.Hour
)

type options struct {
	readDSN         string
	writeDSN        string
	maxIdleConn     int
	maxOpenConn     int
	connMaxLifetime time.Duration
	dsn             DSN
	logLevel        logger.LogLevel
}

// DSN represents the db connection params.
type DSN struct {
	Username string
	Password string
	Host     string
	Name     string
	Port     uint32
	SslMode  string
	Timezone string
}

// Option represents configurable db option.
type Option func(o *options)

// DB represents the database instance for the whole app.
type DB struct {
	db      *gorm.DB
	dsn     DSN
	options *options
}

// WithReadConfig sets the read db.
func WithReadConfig(dsn DSN) Option {
	return func(o *options) {
		o.readDSN = buildDSN(dsn)
	}
}

// WithWriteConfig sets the write db.
func WithWriteConfig(dsn DSN) Option {
	return func(o *options) {
		o.writeDSN = buildDSN(dsn)
	}
}

// WithReadDSN sets the read db.
func WithReadDSN(dsn string) Option {
	return func(o *options) {
		o.readDSN = dsn
	}
}

// WithWriteDSN sets the write db.
func WithWriteDSN(dsn string) Option {
	return func(o *options) {
		o.writeDSN = dsn
	}
}

// WithMaxIdleConn sets the max idle connections.
func WithMaxIdleConn(value int) Option {
	return func(o *options) {
		o.maxIdleConn = value
	}
}

// WithMaxOpenConn sets the max life .
func WithMaxOpenConn(value int) Option {
	return func(o *options) {
		o.maxOpenConn = value
	}
}

// WithConnMaxLifetime sets the max life for connection.
func WithConnMaxLifetime(t time.Duration) Option {
	return func(o *options) {
		o.connMaxLifetime = t
	}
}

// WithGormLogLevel sets the gorm log level.
func WithGormLogLevel(level logger.LogLevel) Option {
	return func(o *options) {
		o.logLevel = level
	}
}

// New returns a new instance of database object with connection.
func New(ctx context.Context, opts ...Option) (*DB, error) {
	rand.NewSource(time.Now().UnixNano())

	opt := options{
		maxIdleConn:     defaultMaxIdleConn,
		maxOpenConn:     defaultMaxOpenConn,
		connMaxLifetime: defaultConnMaxLifetime,
		logLevel:        logger.Error,
	}

	for _, o := range opts {
		o(&opt)
	}

	conn, err := sql.Open("postgres", opt.writeDSN)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{
		Logger: logger.Default.LogMode(opt.logLevel),
	})
	if err != nil {
		return nil, err
	}

	dsn, err := parseDSN(opt.writeDSN)
	if err != nil {
		return nil, err
	}

	opt.dsn = dsn

	if err := configure(db, opt); err != nil {
		return nil, err
	}

	return &DB{
		db:      db.WithContext(ctx),
		options: &opt,
	}, nil
}

// ReadDSN returns the read db dsn.
func (db *DB) ReadDSN() string {
	return db.options.readDSN
}

// WriteDSN returns the read db dsn.
func (db *DB) WriteDSN() string {
	return db.options.writeDSN
}

// Close closes the database connection.
func (db *DB) Close() error {
	if db.db == nil {
		return nil
	}

	pqDB, err := db.db.DB()
	if err != nil {
		return err
	}

	return pqDB.Close()
}

// DB returns the gorm db instance.
func (db *DB) DB() *gorm.DB {
	return db.db
}

func parseDSN(strDSN string) (DSN, error) {
	config, err := pgx.ParseConfig(strDSN)
	if err != nil {
		return DSN{}, err
	}

	return DSN{
		Host:     config.Host,
		Username: config.User,
		Name:     config.Database,
		Port:     uint32(config.Port),
	}, nil
}

func configure(db *gorm.DB, opts options) error {
	// use db replica for reading.
	err := db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(opts.readDSN)},
	}))
	if err != nil {
		return err
	}

	// Get generic database object sql.DB to use its functions.
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.maxIdleConn)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.maxOpenConn)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.connMaxLifetime)

	return nil
}

func buildDSN(dsn DSN) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dsn.Host,
		dsn.Username,
		dsn.Password,
		dsn.Name,
		dsn.Port,
		dsn.SslMode,
		dsn.Timezone,
	)
}
