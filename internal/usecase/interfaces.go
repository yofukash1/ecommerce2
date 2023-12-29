package usecase

import (
	"context"

	"github.com/yofukashi/e-commerce/internal/entity"
)

type (
	EcommerceRepoI interface {
		CreareUser(ctx context.Context, u *entity.CreateUserRepo) (userID string, err error)
		CreatePayment(ctx context.Context, p *entity.CreatePayment) (paymentID string, err error)
		CreateCard(ctx context.Context, c *entity.CreateCardRepo) (cardNumber string, err error)
		GetTransactions(ctx context.Context, b *entity.CheckBalanceRepo) (transactions []entity.Transaction, err error)
		Transfer(ctx context.Context, t *entity.TransferRepo) (success bool, err error)
		AddMoney(ctx context.Context, t *entity.AddMoneyRepo) (success bool, err error)
		IfCardExists(ctx context.Context, cardNumber string) (ex bool, err error)
	}

	EcommerceUseCaseI interface {
		CreateUser(ctx context.Context, firstName string, email string, passwordHash string) (userID string, err error)
		CreatePayment(ctx context.Context, userID string) (paymentID string, err error)
		CreateCard(ctx context.Context, paymentID string, cardNumber string, expirationDate string) (cardNum string, err error)
		CheckBalance(ctx context.Context, cardNumber string) (amount uint64, err error)
		Transfer(ctx context.Context, srcCardNumber string, dstCardNumber string, amount uint64) (success bool, err error)
		AddMoney(ctx context.Context, cardNumber string, amount uint64) (success bool, err error)
	}
)
