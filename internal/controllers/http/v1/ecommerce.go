package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yofukashi/e-commerce/internal/usecase"
	"github.com/yofukashi/e-commerce/pkg/logging"
)

type ecommerceRoutes struct {
	e usecase.EcommerceUseCaseI
	l *logging.Logger
}

func newEcommerceRoutes(handler *gin.RouterGroup, e usecase.EcommerceUseCaseI, l *logging.Logger) {
	r := &ecommerceRoutes{e, l}

	h := handler.Group("/usr")
	{
		// h.GET("/history", r.history)
		h.POST("/create", r.createUser)
		h.POST("/balance", r.checkBalance)
		h.POST("/transfer", r.transfer)
	}
	h1 := handler.Group("/master")
	{
		h1.POST("/payment", r.createPayment)
		h1.POST("/card", r.createCard)
		h1.POST("/add", r.addMoney)
	}
}

// ---------------------------------------------
type createUserReq struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type createUserResp struct {
	UserID string `json:"user_id"`
}

func (er *ecommerceRoutes) createUser(c *gin.Context) {
	var request createUserReq
	if err := c.ShouldBindJSON(&request); err != nil {
		er.l.Errorf("error while binding json: %v", err)
		errorResponse(c, http.StatusBadRequest, "invalid request data")
		return
	}
	// TODO convert password to password hash
	userID, err := er.e.CreateUser(c.Request.Context(), request.FirstName, request.Email, request.Password)

	if err != nil {
		er.l.Errorf("error while creating user: %v", err)
		errorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, createUserResp{UserID: userID})

}

// ---------------------------------------------

type createPaymentReq struct {
	UserID string `json:"user_id"`
}

type createPaymentResp struct {
	PaymentID string `json:"payment_id"`
}

func (er *ecommerceRoutes) createPayment(c *gin.Context) {
	var request createPaymentReq
	if err := c.ShouldBindJSON(&request); err != nil {
		er.l.Errorf("error while binding json: %v", err)
		errorResponse(c, http.StatusBadRequest, "invalid request data")
		return
	}
	paymentID, err := er.e.CreatePayment(c.Request.Context(), request.UserID)

	if err != nil {
		er.l.Errorf("error while creating user: %v", err)
		errorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, createPaymentResp{PaymentID: paymentID})

}

// ---------------------------------------------

type CreateCardReq struct {
	PaymentID      string `json:"payment_id"`
	CardNumber     string `json:"card_number"`
	ExpirationDate string `json:"expiration_date"`
}

type CreateCardResp struct {
	CardNumber string `json:"card_number"`
}

func (er *ecommerceRoutes) createCard(c *gin.Context) {
	var request CreateCardReq
	if err := c.ShouldBindJSON(&request); err != nil {
		er.l.Errorf("error while binding json: %v", err)
		errorResponse(c, http.StatusBadRequest, "invalid request data")
		return
	}
	cardNumber, err := er.e.CreateCard(c.Request.Context(), request.PaymentID, request.CardNumber, request.ExpirationDate)
	if err != nil {
		er.l.Errorf("error while creating card: %v", err)
		errorResponse(c, http.StatusInternalServerError, "internal error")
		return
	}

	c.JSON(http.StatusOK, CreateCardResp{CardNumber: cardNumber})

}

// ---------------------------------------------

type CheckBalanceReq struct {
	CardNumber string `json:"card_number"`
}
type CheckBalanceResp struct {
	Amount uint64 `json:"amount"`
}

func (er *ecommerceRoutes) checkBalance(c *gin.Context) {
	var request CheckBalanceReq
	if err := c.ShouldBindJSON(&request); err != nil {
		er.l.Errorf("error while binding json: %v", err)
		errorResponse(c, http.StatusBadRequest, "invalid request data")
		return
	}

	amount, err := er.e.CheckBalance(c.Request.Context(), request.CardNumber)
	if err != nil {
		er.l.Errorf("error while checking balance: %v", err)
		errorResponse(c, http.StatusInternalServerError, "internal error, please try later")
		return
	}

	c.JSON(http.StatusOK, CheckBalanceResp{Amount: amount})

}

// ---------------------------------------------

type TransferReq struct {
	SrcCardNumber string `json:"src_card_number"`
	DstCardNumber string `json:"dst_card_number"`
	Amount        uint64 `json:"amount"`
}

type TranferResp struct {
	Success bool `json:"success"`
}

func (er *ecommerceRoutes) transfer(c *gin.Context) {
	var request TransferReq
	if err := c.ShouldBindJSON(&request); err != nil {
		er.l.Errorf("error while binding json: %v", err)
		errorResponse(c, http.StatusBadRequest, "invalid request data")
		return
	}

	success, err := er.e.Transfer(c.Request.Context(), request.SrcCardNumber, request.DstCardNumber, request.Amount)

	if err != nil {
		er.l.Errorf("error while transfering: %v", err)
		errorResponse(c, http.StatusInternalServerError, "internal error, please try later")
		return
	}
	c.JSON(http.StatusOK, TranferResp{Success: success})
}

// ---------------------------------------------

type AddMoneyReq struct {
	CardNumber string `json:"card_number"`
	Amount     uint64 `json:"amount"`
}

type AddMoneyResp struct {
	Success bool `json:"success"`
}

func (er *ecommerceRoutes) addMoney(c *gin.Context) {
	var request AddMoneyReq
	if err := c.ShouldBindJSON(&request); err != nil {
		er.l.Errorf("error while binding json: %v", err)
		errorResponse(c, http.StatusBadRequest, "invalid request data")
		return
	}

	success, err := er.e.AddMoney(c.Request.Context(), request.CardNumber, request.Amount)

	if err != nil {
		er.l.Errorf("error while transfering: %v", err)
		errorResponse(c, http.StatusInternalServerError, "internal error, please try later")
		return
	}
	c.JSON(http.StatusOK, AddMoneyResp{Success: success})

}
