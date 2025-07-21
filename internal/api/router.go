package api

import "github.com/gorilla/mux"

func (handler *SubscriptionHandler) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/subscriptions/total", handler.GetTotalCost).Methods("GET")
	router.HandleFunc("/subscriptions", handler.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions/{id:[0-9a-fA-F-]{36}}", handler.GetSubscription).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", handler.UpdateSubscription).Methods("PATCH")
	router.HandleFunc("/subscriptions/{id}", handler.DeleteSubscription).Methods("DELETE")
	router.HandleFunc("/subscriptions", handler.GetListSubscription).Methods("GET")
}
