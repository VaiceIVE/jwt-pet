package main_db_provider

import (
	"auth-testcase/library/loggerhelper"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JWTRepository struct {
	logger       *loggerhelper.CustomLogger
	dbConnection *pgxpool.Pool
}

type IJWTRepository interface {
	CreateAccess(guid string, hash string, exp int64, ip string) (int, error)
	GetDataByGuid(guid string) (JWTModel, error)
	UpdateAccess(guid string, exp int64, ip string) (error)
}

func NewJWTRepository(logger *loggerhelper.CustomLogger, dbConfig DatabaseConfiguration) IJWTRepository {
	conn, err := pgxpool.New(context.Background(), "postgres://" + dbConfig.Username + ":" + dbConfig.Password + "@" + dbConfig.Host + ":5432/" + dbConfig.DatabaseName)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &JWTRepository{
		logger: logger,
		dbConnection: conn,
	}
}

func (repo *JWTRepository) CreateAccess(guid string, hash string, exp int64, ip string) (int, error) {
	repo.logger.SugarWithTracing().Infof("LoanRepository.CreateJWT")


	id := 0
	candidate := repo.dbConnection.QueryRow(context.Background(), "SELECT \"id\" FROM \"jwt-tokens\" WHERE \"guid\"=$1", guid)

	candidate.Scan(&id)

	if id != 0{
		repo.logger.SugarWithTracing().Infof(strconv.Itoa(id))
		repo.dbConnection.QueryRow(context.Background(), "DELETE FROM \"jwt-tokens\" WHERE \"guid\"=$1", guid)
	}
	
	row := repo.dbConnection.QueryRow(context.Background(), "INSERT INTO \"jwt-tokens\"(\"hash\", \"guid\", \"expiry\", last_login_ip) VALUES ($1, $2, $3, $4) RETURNING \"id\"", hash, guid, exp, ip)

	row.Scan(&id)

	return id, nil
}



func (repo *JWTRepository) GetDataByGuid(guid string) (JWTModel, error) {
	
	var id int32
	var hash string
	var expiry int32
	var ip string
	candidate := repo.dbConnection.QueryRow(context.Background(), "SELECT * FROM \"jwt-tokens\" WHERE \"guid\"=$1", guid)

	candidate.Scan(&id, &hash, nil, &expiry, &ip)

	jwtModel := JWTModel{
		Id: id,
		Hash: hash,
		Expiry: expiry,
		Guid: guid,
		LastLoginIp: ip,
	}

	return jwtModel, nil
}

func (repo *JWTRepository) UpdateAccess(guid string, exp int64, ip string) (error) {

 	repo.dbConnection.QueryRow(context.Background(), "UPDATE \"jwt-tokens\" SET \"expiry\"=$2, last_login_ip=$3 WHERE \"guid\"=$1", guid, exp, ip)

	return nil

}	