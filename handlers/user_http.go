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
	"shopify-apis/utils"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request
	var user domain.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the User has entered valid fields send the request to Shopify to create a Customer
	if ok, err := utils.ValidateUser(user); ok && err == nil {

		// Create an URL for the Shopify store Create User endpoint
		createUserURL := constants.STORE_BASE_URL + constants.CREATE_USER_ENDPOINT

		//Create requestBody Shopify
		requestBody := map[string]interface{}{
			"customer": map[string]interface{}{
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"phone":      user.MobileNumber,
				"email":      user.Email,
			},
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", createUserURL, bytes.NewBuffer(jsonData))
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
			http.Error(w, err.Error(), resp.StatusCode)
			return
		}

		if resp.StatusCode != http.StatusCreated {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Output the response content and status
			w.WriteHeader(resp.StatusCode)
			w.Header().Set("Content-Type", "application/json") // Adjust the content type as needed
			w.Write(body)
			return
		}

		// Parse the JSON response
		var shopifyResonse map[string]domain.User
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&shopifyResonse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add the user to the cart
		AddNewUserToCart(user.MobileNumber[3:])

		// Return the JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shopifyResonse)
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetUser(phoneNumber string) string {
	// Create an URL for the Shopify store Create User endpoint
	getUserURL := constants.STORE_BASE_URL + constants.GET_USER_ENDPOINT + phoneNumber

	// Get the User Details from the Shopify store
	req, err := http.NewRequest("GET", getUserURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	// Set the admin token in the request headers
	req.Header.Set("X-Shopify-Access-Token", configs.AdminToken())

	// Send the request to Shopify
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", resp.Status)
		return ""
	}

	// Parse the JSON response
	var shopifyResonse map[string][]domain.User
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&shopifyResonse); err != nil {
		fmt.Println("Error decoding response:", err)
		return ""
	}

	return shopifyResonse["customers"][0].Email
}

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
			"email": userEmail,
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
}
