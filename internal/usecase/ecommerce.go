package usecase

import (
	"context"
	"errors"

	"github.com/yofukashi/e-commerce/internal/entity"
	"github.com/yofukashi/e-commerce/internal/usecase/middleware"
	"github.com/yofukashi/e-commerce/pkg/logging"
)

type EcommerceUseCase struct {
	l    *logging.Logger
	repo EcommerceRepoI
}

func New(l *logging.Logger, repo EcommerceRepoI) *EcommerceUseCase {
	return &EcommerceUseCase{l: l, repo: repo}
}

//-----------------------------------

//-----------------------------------

func (e *EcommerceUseCase) CreateUser(ctx context.Context, firstName string, email string, passwordHash string) (userID string, err error) {
	if f := middleware.EmptyFields(firstName, email, passwordHash); f {
		e.l.Errorf("invalid data")
		return "", errors.New("invalid data")
	}
	u := &entity.CreateUserRepo{FirstName: firstName, Email: email, PasswordHash: passwordHash}
	userID, err = e.repo.CreareUser(ctx, u)
	if err != nil {
		e.l.Errorf("error while creating user: %v", err)
		return "", err
	}
	e.l.Tracef("user %s has been created", userID)
	return userID, nil
}

//-----------------------------------

func (e *EcommerceUseCase) CreatePayment(ctx context.Context, userID string) (paymentID string, err error) {
	if f := middleware.EmptyFields(userID); f {
		e.l.Errorf("invalid data")
		return "", errors.New("invalid data")
	}

	p := &entity.CreatePayment{UserID: userID}
	paymentID, err = e.repo.CreatePayment(ctx, p)
	if err != nil {
		e.l.Errorf("error while creating payment: %v", err)
		return "", err
	}
	e.l.Tracef("payment %s for user %s has been created", paymentID, userID)
	return paymentID, nil

}

//-----------------------------------

// TODO check date middleware
// card format middleware
func (e *EcommerceUseCase) CreateCard(ctx context.Context, paymentID string, cardNumber string, expirationDate string) (cardNum string, err error) {
	// checks
	if f := middleware.EmptyFields(paymentID, cardNumber, expirationDate); f {
		e.l.Errorf("invalid data")
		return "", errors.New("invalid data")
	}

	if err := middleware.CheckExpirationDate(expirationDate); err != nil {
		return "", err
	}

	ex, err := e.repo.IfCardExists(ctx, cardNumber)
	if err != nil {
		e.l.Errorf("error while checking: %v", err)
		return "", err
	}
	if ex {
		e.l.Errorf("error while adding card. Card %s already exists", cardNumber)
		return "", errors.New("card already exists")
	}

	c := &entity.CreateCardRepo{PaymentID: paymentID, CardNumber: cardNumber, ExpirationDate: expirationDate}
	cardNum, err = e.repo.CreateCard(ctx, c)
	if err != nil {
		e.l.Errorf("error while creating card: %v", err)
		return "", err
	}
	e.l.Tracef("card %s for payment accont %s has been created", cardNumber, paymentID)
	return cardNum, nil
}

//-----------------------------------

func (e *EcommerceUseCase) CheckBalance(ctx context.Context, cardNumber string) (amount uint64, err error) {
	if f := middleware.EmptyFields(cardNumber); f {
		e.l.Errorf("invalid data")
		return 0, errors.New("invalid data")
	}
	b := &entity.CheckBalanceRepo{CardNumber: cardNumber}
	transactions, err := e.repo.GetTransactions(ctx, b)
	if err != nil {
		e.l.Errorf("error while getting transactions: %v", err)
		return 0, err
	}

	var s int64
	for _, t := range transactions {
		if t.TType == "charge" {
			s -= int64(t.Amount)
		}
		if t.TType == "recharge" {
			s += int64(t.Amount)
		}
	}
	var res uint64
	if s >= 0 {
		res = uint64(s)
	} else {
		e.l.Errorf("negative balance of card: %v", cardNumber)
		return 0, errors.New("internal error")

	}

	return res, nil
}

//-----------------------------------

// TODO error with balance
func (e *EcommerceUseCase) Transfer(ctx context.Context, srcCardNumber string, dstCardNumber string, amount uint64) (success bool, err error) {
	if f := middleware.EmptyFields(srcCardNumber, dstCardNumber); f {
		e.l.Errorf("invalid data")
		return false, errors.New("invalid data")
	}
	balance, err := e.CheckBalance(ctx, srcCardNumber)
	if err != nil {
		e.l.Errorf("error while checking balance: %v", err)
		return false, err
	}
	if (int64(balance) - int64(amount)) < 0 {
		return false, errors.New("not enough money")
	}

	t := &entity.TransferRepo{SrcCardNumber: srcCardNumber, DstCardNumber: dstCardNumber, Amount: amount}
	success, err = e.repo.Transfer(ctx, t)
	if err != nil {
		e.l.Errorf("error while transfering money: %v", err)
		return false, err
	}
	return success, nil
}

//-----------------------------------

func (e *EcommerceUseCase) AddMoney(ctx context.Context, cardNumber string, amount uint64) (success bool, err error) {
	if f := middleware.EmptyFields(cardNumber); f {
		e.l.Errorf("invalid data")
		return false, errors.New("invalid data")
	}

	a := &entity.AddMoneyRepo{CardNumber: cardNumber, Amount: amount}
	success, err = e.repo.AddMoney(ctx, a)
	if err != nil {
		e.l.Errorf("error while adding money: %v", err)
		return false, err
	}
	return success, nil
}
