package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/anirudhchy/gobank/types"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	// GetAccountByID(int) (*types.Account, error)
	GetAccountByNumber(int) (*types.Account, error)
	UpdateAccountBalanceByNumber(balance, number int64) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
				id serial primary key,
				first_name varchar(50),
				last_name varchar(50),
				number BIGINT,
				encrypted_password varchar(100),
				balance BIGINT,
				created_at timestamp
				)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *types.Account) error {
	query := `insert into account 
				(first_name, last_name, number, encrypted_password, balance, created_at)
				values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt)

	if err != nil {
		return err
	}

	// fmt.Printf("%+v\n", resp)

	return nil
}
func (s *PostgresStore) UpdateAccount(*types.Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(number int) error {
	_, err := s.db.Query("delete from account where number = $1", number)

	return err
}

func (s *PostgresStore) GetAccountByNumber(number int) (*types.Account, error) {

	rows, err := s.db.Query("select * from account where number = $1", number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", number)

}

func (s *PostgresStore) UpdateAccountBalanceByNumber(balance, number int64) error {

	_, err := s.db.Query("update account set balance = $1 where number = $2", balance, number)

	return err
}

// func (s *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
// 	rows, err := s.db.Query("select * from account where id = $1", id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		return scanIntoAccount(rows)
// 	}

// 	return nil, fmt.Errorf("account %d not found", id)

// }

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
	rows, err := s.db.Query("select * from account")

	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}

	for rows.Next() {

		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil

}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
