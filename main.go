package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/sirinibin/go-mysql-rest/config"
	"gitlab.com/sirinibin/go-mysql-rest/controller"
)

func main() {
	fmt.Println("A GoLang / Myql Microservice [OAuth2,Redis & JWT used for token management]!")

	config.InitMysql()
	config.InitRedis()
	//defer config.DB.Close()

	router := mux.NewRouter()
	//API Info
	router.HandleFunc("/", controller.APIInfo).Methods("GET")

	// Register a new user account
	router.HandleFunc("/v1/register", controller.Register).Methods("POST")

	// OAuth2 Authentication
	router.HandleFunc("/v1/authorize", controller.Authorize).Methods("POST")
	router.HandleFunc("/v1/accesstoken", controller.Accesstoken).Methods("POST")

	// Refresh access token
	router.HandleFunc("/v1/refresh", controller.RefreshAccesstoken).Methods("POST")

	//Me
	router.HandleFunc("/v1/me", controller.Me).Methods("GET")
	// Logout
	router.HandleFunc("/v1/logout", controller.LogOut).Methods("DELETE")

	//Employees
	router.HandleFunc("/v1/employees", controller.CreateEmployee).Methods("POST")
	router.HandleFunc("/v1/employees", controller.UpdateEmployee).Methods("PUT")
	router.HandleFunc("/v1/employees/{id}", controller.DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/v1/employees/{id}", controller.ViewEmployee).Methods("GET")
	router.HandleFunc("/v1/employees", controller.ListEmployee).Methods("GET")

	go func() {
		log.Fatal(http.ListenAndServeTLS(":2001", "localhost.cert.pem", "localhost.key.pem", router))
	}()

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			log.Printf("Serving @ https://" + ip.String() + ":2001 /\n")
			log.Printf("Serving @ http://" + ip.String() + ":2000 /\n")
		}
	}
	log.Fatal(http.ListenAndServe(":2000", router))

}
