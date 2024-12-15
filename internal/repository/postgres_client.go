package repository

import (
	"context"
	"email-verification-service/model"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ClientPostgres struct {
	pool *pgxpool.Pool
}

func NewClientPostgres(pool *pgxpool.Pool) *ClientPostgres {
	return &ClientPostgres{
		pool: pool,
	}
}

func (c *ClientPostgres) Register(ctx context.Context, client model.ClientInput) error {
	query := fmt.Sprintf(`INSERT INTO %s(
               first_name, last_name, username, hashed_password, email, code)
               VALUES($1, $2, $3, $4, $5, $6)`, clientTable)
	_, err := c.pool.Exec(ctx, query, client.FirstName, client.LastName, client.Username, client.Password, client.Email, client.VerificationCode)
	return err
}

func (c *ClientPostgres) Verify(ctx context.Context, client model.ClientVerification) error {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	output, err := c.IsVerified(ctx, client)
	if err != nil {
		return err
	}
	if output.IsVerified {
		return fmt.Errorf("client not found or already verified")
	}
	if time.Now().UTC().After(*output.Exp) {
		return fmt.Errorf("verification code expired")
	}

	query := fmt.Sprintf(`UPDATE %s SET is_verified = true, exp = NULL, code = ''
          WHERE email = $1 AND code = $2`, clientTable)
	res, err := tx.Exec(ctx, query, client.Email, client.Code)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("client not found or already verified")
	}
	err = tx.Commit(ctx)
	return err
}

func (c *ClientPostgres) IsVerified(ctx context.Context, client model.ClientVerification) (model.VerificationCode, error) {
	var output model.VerificationCode
	query := fmt.Sprintf(`SELECT code, is_verified, exp FROM %s WHERE email = $1`, clientTable)
	err := c.pool.QueryRow(ctx, query, client.Email).Scan(&output.Code, &output.IsVerified, &output.Exp)
	return output, err
}

func (c *ClientPostgres) Refresh(ctx context.Context, client model.ClientVerification) (model.VerificationCode, error) {
	query := fmt.Sprintf(`UPDATE %s SET exp = CURRENT_TIMESTAMP + '1 minute'::interval, code = $1
          WHERE email = $2`, clientTable)
	res, err := c.pool.Exec(ctx, query, client.VerificationCode.Code, client.Email)
	if err != nil {
		return model.VerificationCode{}, err
	}
	if res.RowsAffected() == 0 {
		return model.VerificationCode{}, fmt.Errorf("client not found")
	}
	return client.VerificationCode, nil
}
