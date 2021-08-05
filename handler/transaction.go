package handler

import (
	"moyu/helper"
	"moyu/transaction"
	"moyu/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// parameter di uri
// tangkap parameter dan mapping input struct
// panggil service, input struct sebagai parameter
// service, dengan campaign id bisa panggil repo
// repo mencari data transaction suatu campaign

type transactionHanler struct {
	service transaction.Service
}

func NewTransactionHandler(transactionService transaction.Service) *transactionHanler {
	return &transactionHanler{transactionService}
}

func (h *transactionHanler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transaction", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transaction", "error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	transactionsFormatter := transaction.FormatTransactions(transactions)

	response := helper.APIResponse("Successfuly to get campaign's transaction", "success", http.StatusOK, transactionsFormatter)
	c.JSON(http.StatusOK, response)
}