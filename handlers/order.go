package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"shopify-apis/configs"
	"shopify-apis/constants"
	"shopify-apis/domain"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract the mobile number from the URL path
	mobileNumber := r.URL.Path[len("/order/create/"):]

	// Lock the cart for concurrent access
	cartMutex.Lock()
	defer cartMutex.Unlock()

	// Check if a cart exists for the mobile number
	lineItems := cart[mobileNumber]

	if len(lineItems) == 0 {
		http.Error(w, "Cart is empty", http.StatusNotFound)
		return
	}
	// Convert lineItems into JSON string

	userEmail := GetUser(mobileNumber)
	fmt.Println(userEmail)

	// Create an URL for the Shopify store Create User endpoint
	createOrderURL := constants.STORE_BASE_URL + constants.CREATE_ORDER_ENDPOINT + constants.UTM_PARAMS

	//Create requestBody Shopify
	requestBody := map[string]interface{}{
		"order": map[string]interface{}{
			"email":      userEmail,
			"line_items": lineItems,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", createOrderURL, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the admin token in the request headers
	req.Header.Set("X-Shopify-Access-Token", configs.AdminToken())
	req.Header.Set("Content-Type", "application/json")

	// Send the request to Shopify
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Request failed with status: %s", resp.Status)
		// Output the response content and status
		w.WriteHeader(resp.StatusCode)
		w.Header().Set("Content-Type", "application/json") // Adjust the content type as needed
		w.Write(body)
		return
	}

	// Parse the JSON response
	var shopifyResonse map[string]domain.Order
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&shopifyResonse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var createOrderResponse domain.OrderResponse
	createOrderResponse.OrderID = fmt.Sprintf("%d", shopifyResonse["order"].ID)

	// Return the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createOrderResponse)

	// Empty the cart
	cart[mobileNumber] = []domain.CartItem{}
}

func ViewOrdersByUTM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract the utm_source from the url parameter "source"
	utmSource := r.URL.Query().Get("source")

	// Create an URL for the Shopify store orders endpoint
	// 
	// Blocked: I was unable to find the correct endpoint to get orders by utm_source.
	// I found a python snippet in that it told that the endpoint is /orders.json?source=<utm_source>
	// but it did not work for me. Since I think that for the utm tags to be linked to the order it 
	// needs to be paid. I was unable to find a way to create a paid order in Shopify.

	// Get the orders from the Shopify store
	// TODO


}
