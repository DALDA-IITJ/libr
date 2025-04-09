package api

import "github.com/gorilla/mux"

// SetUpRoutes initializes and returns a new Gorilla Mux router with the following routes:
// - POST /db/saveMsg: Saves a message with its associated certificates to the database.
// - GET /db/getMsg: Retrieves messages from the database based on a given timestamp.
// - GET /db/isalive: Checks if the DB Node is active and running.
// - GET /db/timesign: Provides a signed timestamp for a given client timestamp.
//
// The following routes are defined:
//
// 1. POST /db/saveMsg
//    - Saves a message to the database after validating its associated certificates.
//    - Refer to the "Save Message" section in NODEAPIS.md for details.
//
// 2. GET /db/getMsg
//    - Retrieves messages from the database based on a provided timestamp.
//    - Refer to the "Retrieve Messages by Timestamp" section in NODEAPIS.md for details.
//
// 3. GET /db/isalive
//    - Checks if the DB Node is active and running.
//    - Refer to the "Check Node Health" section in NODEAPIS.md for details.
//
// Returns:
// - *mux.Router: A configured router with the defined API routes.

func SetUpRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/db/savemsg", SaveMsgHandler).Methods("POST")
	router.HandleFunc("/db/getmsg", GetMsgHandler).Methods("GET")
	router.HandleFunc("/db/isalive", IsAliveHandler).Methods("GET")
	// router.HandleFunc("/db/timesign", TimeSignHandler).Methods("GET")

	return router
}
