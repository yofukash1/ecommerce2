package entity

type CreateUserRepo struct {
	//ID           string `bson:"_id,omitempty" json:"user_id,omitempty"`
	//PaymentID    string `bson:"payment_id" json:"payment_id"`
	FirstName    string `bson:"first_name" json:"first_name"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"password_hash" json:"-"`
}

type CreatePayment struct {
	UserID string `json:"user_id"`
}

type CreateCardRepo struct {
	PaymentID      string `json:"payment_id"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
}

type CheckBalanceRepo struct {
	CardNumber string `json:"card_number"`
}

type TransferRepo struct {
	SrcCardNumber string
	DstCardNumber string
	Amount        uint64
}

type AddMoneyRepo struct {
	CardNumber string
	Amount     uint64
}
