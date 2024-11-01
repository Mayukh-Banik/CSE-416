// controllers/transactionController.go
package controllers

import (
	"encoding/json"
	"go-server/services"
	"net/http"
	"go-server/models"
)

// CreateTransaction handles creating a new transaction
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service to create a new transaction
	createdTransaction, err := services.NewTransaction(transaction.Sender, transaction.Receiver, transaction.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the created transaction
	json.NewEncoder(w).Encode(createdTransaction)
}
