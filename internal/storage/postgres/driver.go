package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.yandex/hasql"
	"golang.yandex/hasql/checkers"
	"log"
	"social-network/cmd/social-network/configuration"
	"strconv"
	"strings"
)

const driverName = "pgx"

type PGStorage struct {
	logger  *zap.Logger
	cluster *hasql.Cluster
}

func NewPGStorage(logger *zap.Logger, cfg configuration.StorageConfig) (*PGStorage, error) {
	cluster, err := newPGCluster(cfg)
	if err != nil {
		return nil, fmt.Errorf("new pg cluster: %w", err)
	}

	return &PGStorage{
		logger:  logger,
		cluster: cluster,
	}, nil
}

func (s *PGStorage) Close() error {
	return s.cluster.Close()
}

func (s *PGStorage) Master() *sqlx.DB {
	db := s.cluster.Primary().DB()
	return sqlx.NewDb(db, driverName)
}

func (s *PGStorage) Slave() *sqlx.DB {
	db := s.cluster.StandbyPreferred().DB()
	return sqlx.NewDb(db, driverName)
}

func newPGCluster(cfg configuration.StorageConfig) (*hasql.Cluster, error) {
	nodes := make([]hasql.Node, 0, len(cfg.Hosts))
	for _, host := range cfg.Hosts {
		connString := constructConnectionString(host, cfg)

		parsedConnConfig, err := pgx.ParseConfig(connString)
		if err != nil {
			return nil, fmt.Errorf("parse connection config: %w", err)
		}
		db := sqlx.NewDb(stdlib.OpenDB(*parsedConnConfig), driverName)
		nodes = append(nodes, hasql.NewNode(host, db.DB))
	}

	cluster, err := hasql.NewCluster(nodes, checkers.PostgreSQL)
	if err != nil {
		return nil, fmt.Errorf("new cluster: %w", err)
	}

	ctx, cancelFunction := context.WithTimeout(context.Background(), cfg.InitializationTimeout)
	defer cancelFunction()
	_, err = cluster.WaitForPrimary(ctx)
	if err != nil {
		if closeErr := cluster.Close(); closeErr != nil {
			return nil, fmt.Errorf("cluster close error: %w", closeErr)
		}
		return nil, fmt.Errorf("wait for primary timeout exceed: %w", err)
	}
	return cluster, nil
}

func constructConnectionString(host string, cfg configuration.StorageConfig) string {
	connectionMap := map[string]string{
		"host":     host,
		"port":     strconv.Itoa(cfg.Port),
		"database": cfg.Database,
		"user":     cfg.Username,
		"password": cfg.Password,
	}
	if cfg.SSLMode != "" {
		connectionMap["sslmode"] = cfg.SSLMode
	}

	connectionSlice := make([]string, len(connectionMap))
	for k, v := range connectionMap {
		connectionSlice = append(connectionSlice, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(connectionSlice, " ")
}

// todo: delete
func example(connStr string) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database because %s", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot ping database because %s", err)
	}

	log.Println("Successfully connected to database and pinged it")
}
