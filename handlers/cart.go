package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shopify-apis/constants"
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

func AddItemToCart(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)
	// Extract the mobile number from the URL path
	mobileNumber := r.URL.Path[len("/cart/"):]

	if r.Method != http.MethodPost {
		// Return an error for unsupported methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Parse the request body to get the variant to add
	var item domain.CartItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Lock the cart for concurrent access
	cartMutex.Lock()
	defer cartMutex.Unlock()

	// Check if a cart exists for the mobile number
	if _, ok := cart[mobileNumber]; !ok {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Add the item to the cart
	cart[mobileNumber] = append(cart[mobileNumber], item)

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Variant added to the cart"})
}

func RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	var message string
	// Extract the mobile number from the URL path
	mobileNumber := r.URL.Path[len("/cart/remove/"):]

	if r.Method != http.MethodPost {
		// Return an error for unsupported methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var item domain.CartItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Lock the cart for concurrent access
	cartMutex.Lock()
	defer cartMutex.Unlock()

	// Check if a cart exists for the mobile number
	if _, ok := cart[mobileNumber]; !ok {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	// Remove the item from the cart
	userCart := cart[mobileNumber]
	for i, cartItem := range userCart {
		if cartItem.VariantID == item.VariantID {
			cartItem.Quantity -= item.Quantity
		}
		if cartItem.Quantity == 0 {
			userCart = append(userCart[:i], userCart[i+1:]...)
		}
	}
	if len(userCart) == 0 {
		message = constants.CART_IS_EMPTY
	} else {
		message = constants.VARIANT_REMOVED
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func ViewCart(w http.ResponseWriter, r *http.Request) {
	// Extract the mobile number from the URL path
	mobileNumber := r.URL.Path[len("/cart/"):]

	if r.Method != http.MethodGet {
		// Return an error for unsupported methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Lock the cart for concurrent access
	cartMutex.Lock()
	defer cartMutex.Unlock()

	// Check if a cart exists for the mobile number
	if _, ok := cart[mobileNumber]; !ok {
		http.Error(w, "User does not exist", http.StatusNotFound)
		return
	}

	var viewCartResponse []domain.ViewCartItem
	for _, cartItem := range cart[mobileNumber] {
		variant := GetVariantFromID(cartItem.VariantID)
		cartResponse := domain.ViewCartItem{
			VariantID:  cartItem.VariantID,
			VariantTitle: variant.Title,
		}
		viewCartResponse = append(viewCartResponse, cartResponse)
	}

	// Respond with the cart items
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(viewCartResponse)
}