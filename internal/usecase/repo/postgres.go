package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yofukashi/e-commerce/internal/entity"
	"github.com/yofukashi/e-commerce/pkg/logging"
	"github.com/yofukashi/e-commerce/pkg/postgresql"
)

type EcommerceRepo struct {
	pg *postgresql.Postgres
	l  *logging.Logger
}

func New(pg *postgresql.Postgres, l *logging.Logger) *EcommerceRepo {
	return &EcommerceRepo{pg, l}
}

func (e *EcommerceRepo) CreareUser(ctx context.Context, u *entity.CreateUserRepo) (userID string, err error) {
	userID = uuid.NewString()
	sql, args, err := e.pg.Builder.
		Insert("users").
		Columns("ID, FirstName, Email, PasswordHash").
		Values(userID, u.FirstName, u.Email, u.PasswordHash).
		ToSql()
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return "", err
	}

	_, err = e.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		e.l.Errorf("Error while creating user: %v", err)
		return "", err
	}
	return userID, nil
}

// TODO paymentID check, delete if user cannot be matched
func (e *EcommerceRepo) CreatePayment(ctx context.Context, p *entity.CreatePayment) (paymentID string, err error) {
	// payment create
	paymentID = uuid.NewString()
	sql, args, err := e.pg.Builder.
		Insert("payments").
		Columns("ID").
		Values(paymentID).
		ToSql()
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return "", err
	}

	_, err = e.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		e.l.Errorf("Error while creating payment account: %v", err)
		return "", err
	}

	// assign to user
	sql = fmt.Sprintf("UPDATE users SET PaymentID = '%s' WHERE ID = '%s';", paymentID, p.UserID)
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return "", err
	}

	_, err = e.pg.Pool.Exec(ctx, sql)
	if err != nil {
		e.l.Errorf("Error while adding payment account to user: %v", err)
		return "", err
	}
	return paymentID, nil

}

// TODO delete card if user could not be assigned
func (e *EcommerceRepo) CreateCard(ctx context.Context, c *entity.CreateCardRepo) (cardNumber string, err error) {
	//creating in card table
	sql, args, err := e.pg.Builder.
		Insert("cards").
		Columns("CardNumber, ExpirationDate").
		Values(c.CardNumber, c.ExpirationDate).
		ToSql()
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return "", err
	}

	_, err = e.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		e.l.Errorf("Error while creating card: %v", err)
		return "", err
	}

	//updating in payments
	sql = fmt.Sprintf("UPDATE payments SET CardNumbers = ARRAY_APPEND(CardNumbers, '%s') WHERE ID = '%s';", c.CardNumber, c.PaymentID)
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return "", err
	}

	_, err = e.pg.Pool.Exec(ctx, sql)
	if err != nil {
		e.l.Errorf("Error while adding card to payment account: %v", err)
		return "", err
	}

	return c.CardNumber, nil

}

func (e *EcommerceRepo) GetTransactions(ctx context.Context, b *entity.CheckBalanceRepo) (transactions []entity.Transaction, err error) {
	sql := fmt.Sprintf("SELECT * FROM transactions WHERE CardNumber = '%s';", b.CardNumber)
	rows, err := e.pg.Pool.Query(ctx, sql)
	if err != nil {
		e.l.Errorf("Error while geting transactions: %v", err)
		return nil, err
	}

	defer rows.Close()
	//TODO get the rid of hardcode
	transactions = make([]entity.Transaction, 0, 100)
	for rows.Next() {
		t := entity.Transaction{}

		err = rows.Scan(&t.ID, &t.CardNumber, &t.TType, &t.Amount)
		if err != nil {
			e.l.Errorf("error while reading rows: %v", err)
			return nil, err
		}

		transactions = append(transactions, t)
	}
	return transactions, nil
}

// assuming that src has enough money
func (e *EcommerceRepo) Transfer(ctx context.Context, t *entity.TransferRepo) (success bool, err error) {
	// transfering from srcCard
	transactionID := uuid.NewString()
	sql, args, err := e.pg.Builder.
		Insert("transactions").
		Columns("ID, CardNumber, TType, Amount").
		Values(transactionID, t.SrcCardNumber, "charge", t.Amount).
		ToSql()
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return false, err
	}

	_, err = e.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		e.l.Errorf("Error while transfering from card %s: %v", t.SrcCardNumber, err)
		return false, err
	}

	// transfering to dstCard
	transactionID = uuid.NewString()
	sql, args, err = e.pg.Builder.
		Insert("transactions").
		Columns("ID, CardNumber, TType, Amount").
		Values(transactionID, t.DstCardNumber, "recharge", t.Amount).
		ToSql()
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return false, err
	}

	_, err = e.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		e.l.Errorf("Error while transfering to card %s: %v", t.DstCardNumber, err)
		return false, err
	}
	return true, nil

}

func (e *EcommerceRepo) AddMoney(ctx context.Context, t *entity.AddMoneyRepo) (success bool, err error) {
	// transfering from srcCard
	transactionID := uuid.NewString()
	sql, args, err := e.pg.Builder.
		Insert("transactions").
		Columns("ID, CardNumber, TType, Amount").
		Values(transactionID, t.CardNumber, "recharge", t.Amount).
		ToSql()
	if err != nil {
		e.l.Errorf("Error while converting to sql: %v", err)
		return false, err
	}

	_, err = e.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		e.l.Errorf("Error while transfering to card %s: %v", t.CardNumber, err)
		return false, err
	}
	return true, nil
}

func (e *EcommerceRepo) IfCardExists(ctx context.Context, cardNumber string) (ex bool, err error) {
	sql := fmt.Sprintf("SELECT * FROM payments WHERE '%s' = ANY(CardNumbers);", cardNumber)
	rows, err := e.pg.Pool.Query(ctx, sql)
	if err != nil {
		e.l.Errorf("Error while geting transactions: %v", err)
		return false, err
	}

	defer rows.Close()
	//TODO get the rid of hardcode
	if rows.Next() { // если есть минимум один ряд
		return true, nil
	}
	return false, nil
}
