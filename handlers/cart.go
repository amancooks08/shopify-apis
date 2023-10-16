package handlers

import (
	"shopify-apis/domain"
	"sync"
)

var (
	// In-memory cart to store cart items by mobile number
	cart      = make(map[string][]domain.CartItem)
	cartMutex sync.Mutex // Mutex for concurrent access to the cart
)

func AddNewUserToCart(mobileNumber string) {
	// Lock the cart for concurrent access
	cartMutex.Lock()
	defer cartMutex.Unlock()

	// Add the user to the cart
	cart[mobileNumber] = []domain.CartItem{}
}
