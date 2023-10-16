package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shopify-apis/configs"
	"shopify-apis/constants"
	"shopify-apis/domain"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := domain.PingResponse{Message: "pong"}
	json.NewEncoder(w).Encode(response)
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Create an URL for the Shopify store products endpoint
	productsURL := constants.STORE_BASE_URL + constants.PRODUCTS_ENDPOINT

	// Get the products from the Shopify store
	req, err := http.NewRequest("GET", productsURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(configs.AdminToken())
	// Set the admin token in the request headers
	req.Header.Set("X-Shopify-Access-Token", configs.AdminToken())

	// Send the request to Shopify
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Request failed with status: %s", resp.Status), resp.StatusCode)
		return
	}

	// Parse the JSON response
	var shopifyResonse map[string][]domain.Product
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&shopifyResonse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shopifyResonse)
}

func GetVariantFromID(variantID string) domain.Variant {

	// Create an URL for the Shopify store products endpoint
	variantURL := constants.STORE_BASE_URL + constants.VARIANT_ENDPOINT + variantID + ".json"

	// Get the Variant from the Shopify store
	req, err := http.NewRequest("GET", variantURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return domain.Variant{}
	}

	// Set the admin token in the request headers
	req.Header.Set("X-Shopify-Access-Token", configs.AdminToken())

	// Send the request to Shopify
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return domain.Variant{}
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Println("Variant was not fetched with status:", resp.Status)
		return domain.Variant{}
	}

	fmt.Println(resp.StatusCode)

	var response domain.GetVariantResponse
	// Parse the JSON response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		log.Println("Error decoding JSON:", err)
		return domain.Variant{}
	}

	return response.Variant
}
