package main

import (
	"GameApp/repository/mysql"
	"GameApp/service/userservice"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/users/login", userLoginHandler)
	http.HandleFunc("/users/profile", user_profile_handler)
	http.HandleFunc("/health", health_check_handler)

	log.Print("server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func health_check_handler(write http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(write, `{"message": "every thing is gooood."}`)
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//fmt.Fprintf(w, `{"error": "invalid method."}`)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			//w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			w.Write([]byte(fmt.Sprintf(`{"error": "#{err.Error()}"}`)))
		}

		var req userservice.RegisterRequest
		err = json.Unmarshal(data, &req)
		if err != nil {
			//w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			w.Write([]byte(fmt.Sprintf(`{"error": "#{err.Error()}"}`)))
		}

		mysqlRepo := mysql.New()
		userSVC := userservice.New(mysqlRepo)

		_, err = userSVC.Register(&req)
		//err = json.Unmarshal(data, &req)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}
		w.Write([]byte(`{"message": "user created."}`))
	} else {
		fmt.Fprintf(w, `{"error": "invalid method."}`)
	}

}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error": "invalid method."}`)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "#{err.Error()}"}`)))
		return
	}
	var loginReq userservice.LoginRequest
	err = json.Unmarshal(data, &loginReq)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "#{err.Error()}"}`)))
		return
	}

	mysqlRepo := mysql.New()
	userSVC := userservice.New(mysqlRepo)

	_, err = userSVC.Login(&loginReq)

	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	w.Write([]byte(`{"message": "user credential is ok."}`))
}

func user_profile_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, `{"error": "invalid method."}`)
	}
	preq := userservice.ProfileRequest{UserID: 0}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "#{err.Error()}"}`)))
		return
	}
	err = json.Unmarshal(data, &preq)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "#{err.Error()}"}`)))
		return
	}

	mysqlRepo := mysql.New()
	userSVC := userservice.New(mysqlRepo)

	resp, err := userSVC.Profile(preq)

	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	data, err = json.Marshal(resp)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error":  "%s"}`, err.Error())))
		return
	}

	w.Write(data)
}

// curl -X POST http://localhost:8080/user/register -h "Content-Type: application/json" -d '{"Name": "ALi", "password": "09123312313", "PhoneNumber": "09192751975"}'
