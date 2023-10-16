package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
