package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Result1 map[string]string

func myrecover(w http.ResponseWriter) {
	if r := recover(); r != nil {
		fmt.Fprintln(w, "\"Wrong input.\"")
	}
}

func product(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// In case of wrong type of input we should recover.
		defer myrecover(w)
		decoder := json.NewDecoder(r.Body)
		pdata := make(map[string]string)
		err := decoder.Decode(&pdata)
		if err != nil {
			panic(err)
		}
		user := pdata["user"]
		password := pdata["password"]
		name := pdata["name"]
		desc := pdata["description"]
		if authenticate_redis(user, password) {
			fmt.Println(user, password, name, desc)
			id, _ := insert_product(name, desc)
			res := Result1{"id": id, "name": name, "description": desc}
			res_json, _ := json.Marshal(res)
			fmt.Fprintln(w, string(res_json))

		} else {
			fmt.Fprintln(w, "\"Authentication failure.\"")
		}

	}
}

func component(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// In case of wrong type of input we should recover.
		defer myrecover(w)
		decoder := json.NewDecoder(r.Body)
		pdata := make(map[string]interface{})
		err := decoder.Decode(&pdata)
		if err != nil {
			panic(err)
		}
		user := pdata["user"].(string)
		password := pdata["password"].(string)
		name := pdata["name"].(string)
		desc := pdata["description"].(string)
		product_id := int(pdata["product_id"].(float64))
		owner := int(pdata["owner_id"].(float64))
		if authenticate_redis(user, password) {
			fmt.Println(user, password, name, desc, product_id, owner)
			id, _ := insert_component(name, desc, product_id, owner)
			res := Result1{"id": id, "name": name, "description": desc}
			res_json, _ := json.Marshal(res)
			fmt.Fprintln(w, string(res_json))

		} else {
			fmt.Fprintln(w, "\"Authentication failure.\"")
		}

	}
}

func components(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		pdata := make(map[string]string)
		err := decoder.Decode(&pdata)
		if err != nil {
			panic(err)
		}
		// name := pdata["name"].(string)
		product_id := pdata["product_id"]
		if product_id != "" {
			m := get_components_by_id(product_id)
			res_json, _ := json.Marshal(m)
			fmt.Fprintln(w, string(res_json))
		}

	}
}

func bug(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// In case of wrong type of input we should recover.
		defer myrecover(w)
		decoder := json.NewDecoder(r.Body)
		pdata := make(map[string]interface{})
		err := decoder.Decode(&pdata)
		if err != nil {
			panic(err)
		}

		user := pdata["user"].(string)
		password := pdata["password"].(string)
		if authenticate_redis(user, password) {
			user_id := get_user_id(user)
			pdata["reporter"] = user_id
			id, err := new_bug(pdata)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Fprintln(w, id)
		}
	}
}

func updatebug(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// In case of wrong type of input we should recover.
		defer myrecover(w)
		decoder := json.NewDecoder(r.Body)
		pdata := make(map[string]interface{})
		err := decoder.Decode(&pdata)
		if err != nil {
			panic(err)
		}

		user := pdata["user"].(string)
		password := pdata["password"].(string)
		if authenticate_redis(user, password) {
			update_bug(pdata)
		}
	}
}

func comment(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		pdata := make(map[string]interface{})
		err := decoder.Decode(&pdata)
		if err != nil {
			panic(err)
		}
		user := pdata["user"].(string)
		password := pdata["password"].(string)
		desc := pdata["desc"].(string)
		bug_id := int(pdata["bug_id"].(float64))
		if authenticate_redis(user, password) {
			user_id := get_user_id(user)
			id, err := new_comment(user_id, bug_id, desc)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Fprintln(w, id)
		}
	}
}

func bug_cc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// In case of wrong type of input we should recover.
		defer myrecover(w)
		decoder := json.NewDecoder(r.Body)
		pdata := make(map[string]interface{})
		err := decoder.Decode(&pdata)
		if err != nil {
			panic(err)
		}
		user := pdata["user"].(string)
		password := pdata["password"].(string)
		if authenticate_redis(user, password) {
			bug_id := int(pdata["bug_id"].(float64))
			emails := pdata["emails"]
			action := pdata["action"].(string)
			if action == "add" {
				add_bug_cc(bug_id, emails)
			} else if action == "remove" {
				remove_bug_cc(bug_id, emails)
			} else {
				fmt.Fprintln(w, "\"No vaild action provided.\"")
			}
		} else {
			fmt.Fprintln(w, "\"Authentication failure.\"")
		}
	}
}

func main() {
	load_config("config/bugspad.ini")
	load_users()
	http.HandleFunc("/component/", component)
	http.HandleFunc("/components/", components)
	http.HandleFunc("/product/", product)
	http.HandleFunc("/bug/", bug)
	http.HandleFunc("/bug/cc/", bug_cc)
	http.HandleFunc("/updatebug/", updatebug)
	http.HandleFunc("/comment/", comment)
	http.ListenAndServe(":9998", nil)
}
