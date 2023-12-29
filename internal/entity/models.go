package entity

// Collection Users
type User struct {
	ID           string `bson:"_id,omitempty" json:"user_id,omitempty"`
	//PaymentID    string `bson:"payment_id" json:"payment_id"`
	FirstName    string `bson:"first_name" json:"first_name"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"password_hash" json:"-"`
	//LastConnection time.Time `bson:"last_connection" json:"-"`
}

// Collection Cards
type Card struct {
	CardNumber     string `bson:"card_number" json:"card_number"`
	ExpirationDate string `bson:"expiration_date" json:"expiration_date"`
}

// Collection Payments
type Payment struct {
	ID          string   `bson:"_id,omitempty" json:"payment_id"`
	CardNumbers []string `bson:"cards" json:"cards"`
	//UserID      string `bson:"user_id" json:"-"`
}

// Collection Transactions
type Transaction struct {
	ID         string `bson:"_id,omitempty" json:"transaction_id"`
	CardNumber string `bson:"card_number" json:"card_numebr"`
	TType      string `bson:"ttype" json:"ttype"` // charge / recharge
	Amount     uint64 `bson:"amount" json:"amount"`
}
