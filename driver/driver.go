package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"bytes"
	"fmt"
	"net/http"
	"time"
	"strings"
)

func main() {

	client := &http.Client{Timeout: time.Second * 10}

	/********************************* Create new pda with id 0/ **********************************/

	fmt.Println("Create new pda with id 0 - ")
	jsonText, err := ioutil.ReadFile("helloPda.json")

	// Set the http method, url, and request body
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0", bytes.NewBuffer(jsonText))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var pda PdaProcessor

	// Read in the response body.
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)

	/********************************* Create new pda with id 1/ **********************************/

	fmt.Println("Create new pda with id 1 - ")
	jsonText, err = ioutil.ReadFile("goodbyePda.json")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/1", bytes.NewBuffer(jsonText))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)

	/********************************* Create new pda with id 3/ **********************************/

	fmt.Println("Create new pda with id 3 - ")
	jsonText, err = ioutil.ReadFile("alohaPda.json")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/3", bytes.NewBuffer(jsonText))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)

	/********************************* Check to see if pda exists *********************************/
	fmt.Println("Get list of pdas - ")
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Cache-Control","no-cache")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	/********************************* Reset pda with id 3/ **********************************/

	fmt.Println("Reset pda with id 3 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/3/reset", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	//req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)

	/***************************** Present token 0, position 0 to PDA 0 ****************************/

	fmt.Println("Present token 0 with position 0 to PDA 0 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/0", strings.NewReader("0"))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)

	/**************************** Present token 1, position 2 to PDA 0 *****************************/

	fmt.Println("Present token 1 with position 2 to PDA 0 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/2", strings.NewReader("1"))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)

	/**************************** Present token 1, position 3 to PDA 0 *****************************/

	fmt.Println("Present token 1 with position 3 to PDA 0 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/3", strings.NewReader("1"))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)
	
	/**************************** Present token 0, position 1 to PDA 0 *****************************/

	fmt.Println("Present token 0 with position 1 to PDA 0 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/1", strings.NewReader("0"))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)

	/********************************** Send is_accepted to PDA 0 *********************************/

	fmt.Println("Send is_accepted to PDA 0 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/is_accepted", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	/***************************** Send EOS to PDA 0 after position 3 *****************************/

	fmt.Println("Send EOS to PDA 0 after position 3 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/eos/3", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	/********************************** Send is_accepted to PDA 0 *********************************/

	fmt.Println("Send is_accepted to PDA 0 - ")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/is_accepted", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}