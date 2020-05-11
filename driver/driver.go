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
	"mime/multipart"
)

func main() {

	client := &http.Client{Timeout: time.Second * 10}

	// This will be updated throughout our interactions with the replica server
	sessionCookie := ""

	// Stores pda id with server response from connect requests
	connectId := ""

	/********************************* Create new replica group ***********************************/
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Create new replica group pda with gid 0, member ids 0, 2, 4, 5: ")
	fmt.Println("*******************************************************************************")

	jsonText, err := ioutil.ReadFile("helloPda.json")

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	
	// First field
	if fw, err = w.CreateFormField("pda_code"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader(string(jsonText))); err != nil {
		return
	}

	// Next field
	if fw, err = w.CreateFormField("members"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader("0 2 4 5")); err != nil {
		return
	}
	w.Close()


	// Set the http method, url, and request body
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/replica_pdas/0", &b)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	//req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status:", resp.Status)
	}


	/********************************* Create new replica group 2 *********************************/
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Create new replica group pda with gid 2, member ids 1, 3, 6: ")
	fmt.Println("*******************************************************************************")
	
	jsonText, err = ioutil.ReadFile("goodbyePda.json")

	//var b bytes.Buffer
	w = multipart.NewWriter(&b)
	//var fw io.Writer
	
	// First field
	if fw, err = w.CreateFormField("pda_code"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader(string(jsonText))); err != nil {
		return
	}

	// Next field
	if fw, err = w.CreateFormField("members"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader("1 3 6")); err != nil {
		return
	}
	w.Close()


	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/replica_pdas/2", &b)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	//req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status:", resp.Status)
	}

	/********************************* Get replica group ids ***********************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Return all currently defined replica ids:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/replica_pdas", nil)
	
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
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))



	/**************************** Get replica group member addresses ******************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Return replica group members for gid, 0:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/replica_pdas/0/members", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	 var addresses GroupMemberAddresses

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &addresses); err != nil {
		panic(err)
	}

	fmt.Println(addresses)


	/************************* Get list of all pdas, regardless of group **************************/
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Get list of pdas: ")
	fmt.Println("*******************************************************************************")

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

	/************************************ Delete group: gid 2 *************************************/
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Delete group 2")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("DELETE", "http://localhost:8080/replica_pdas/2/delete", nil)
	
	if err != nil {
		panic(err)
	}

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

	/********************************** Verify group 2 deleted ************************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Return all currently defined replica ids:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/replica_pdas", nil)
	
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


	/************************* Get list of all pdas, regardless of group **************************/
	fmt.Println()		
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Get list of pdas: ")
	fmt.Println("*******************************************************************************")

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


	/********************************* Create new pda with id 9 ***********************************/
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Create new pda with id 9. No group id: ")
	fmt.Println("*******************************************************************************")
	
	jsonText, err = ioutil.ReadFile("directionPda.json")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/9", bytes.NewBuffer(jsonText))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var pda PdaProcessor

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &pda); err != nil {
		panic(err)
	}

	fmt.Println(pda)


	/********************************** Join pda id 9 to group 0 **********************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Join pda with id 9 to group 0: ")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/9/join", 
		strings.NewReader("localhost:8080/replica_pdas/0"))
	
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

	/********************** Verify join. Get replica group member addresses ***********************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Return replica group members for gid, 0:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/replica_pdas/0/members", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var addresses2 GroupMemberAddresses

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &addresses2); err != nil {
		panic(err)
	}

	fmt.Println(addresses2)


	/******** Get PDA Code for pda 9. Note: it should have changed since it was created ***********/
	/******** REMINDER: This is the base spec for the pda, not the pda itself *********************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Get PDA code for PDA 9:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/9/code", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var pdaCode PdaProcessor

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &pdaCode); err != nil {
		panic(err)
	}

	fmt.Println(pdaCode)


	/********************************** Get c3 state for pda 9 ************************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Return C3 state info for PDA 9:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/9/c3state", nil)
	
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

	/********************************** Get c3 state for pda 2 ************************************/
	/* Note: This check is to ensure that the other pdas in the group were updated after the join */

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Return C3 state info for PDA 2:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/2/c3state", nil)
	
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

	// Let's update the session_cookie here that we will send next.
	sessionCookie = string(body)
	fmt.Println(sessionCookie)


	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/************************ Present token 0, position 0, pda @ connectId ************************/
	
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Present token 0, position 0, pda @ connect ID")
	fmt.Println("*******************************************************************************")

	w = multipart.NewWriter(&b)
	
	// First field
	if fw, err = w.CreateFormField("token_value"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader("0")); err != nil {
		return
	}

	// Next field
	if fw, err = w.CreateFormField("session_cookie"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader(sessionCookie)); err != nil {
		return
	}
	w.Close()


	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, 
		"http://localhost:8080/pdas/" + connectId + "/tokens/0", &b)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for FormData
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status:", resp.Status)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie = string(body)

	fmt.Println(sessionCookie)

	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/************************ Present token 1, position 2, pda @ connectId ************************/
	
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Present token 1, position 2, pda @ connect ID")
	fmt.Println("*******************************************************************************")

	w = multipart.NewWriter(&b)
	
	// First field
	if fw, err = w.CreateFormField("token_value"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader("1")); err != nil {
		return
	}

	// Next field
	if fw, err = w.CreateFormField("session_cookie"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader(sessionCookie)); err != nil {
		return
	}
	w.Close()


	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, 
		"http://localhost:8080/pdas/" + connectId + "/tokens/2", &b)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for FormData
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status:", resp.Status)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie = string(body)

	fmt.Println(sessionCookie)

	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)
	
	/********************************* Get snapshot of connect id *********************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Get snapshot of connect id with top 2 elements:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/" + connectId + "/snapshot/2", 
		strings.NewReader(sessionCookie))

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

	var snap Snap

	if err = json.Unmarshal(body, &snap); err != nil {
		panic(err)
	}

	sessionCookie = snap.Cookie
	fmt.Println(snap)

	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	
	/************************ Present token 1, position 3, pda @ connectId ************************/
	
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Present token 1, position 3, pda @ connect ID")
	fmt.Println("*******************************************************************************")

	w = multipart.NewWriter(&b)
	
	// First field
	if fw, err = w.CreateFormField("token_value"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader("1")); err != nil {
		return
	}

	// Next field
	if fw, err = w.CreateFormField("session_cookie"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader(sessionCookie)); err != nil {
		return
	}
	w.Close()


	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, 
		"http://localhost:8080/pdas/" + connectId + "/tokens/3", &b)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for FormData
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status:", resp.Status)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie = string(body)

	fmt.Println(sessionCookie)

	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/********************************* Send peek(2) to connect id *********************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Send peek(2) to pda @ connect id:")
	fmt.Println("*******************************************************************************")
	
	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/" + connectId + "/stack/top/2", 
		strings.NewReader(sessionCookie))

	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var peekResp PeekResponse

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &peekResp); err != nil {
		panic(err)
	}

	sessionCookie = peekResp.Cookie
	fmt.Println(peekResp.Tokens)
	fmt.Println(sessionCookie)

	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/*********************************** Send len to connect id ***********************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Send len to pda @ connect id:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/" + connectId + "/stack/len", 
		strings.NewReader(sessionCookie))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var lenResp PeekResponse

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &lenResp); err != nil {
		panic(err)
	}

	sessionCookie = lenResp.Cookie
	fmt.Println(lenResp.Tokens)
	fmt.Println(sessionCookie)


	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	
	/************************ Present token 1, position 3, pda @ connectId ************************/
	
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Present token 0, position 1, pda @ connect ID")
	fmt.Println("*******************************************************************************")

	w = multipart.NewWriter(&b)
	
	// First field
	if fw, err = w.CreateFormField("token_value"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader("0")); err != nil {
		return
	}

	// Next field
	if fw, err = w.CreateFormField("session_cookie"); err != nil {
		return
	}
	if _, err = io.Copy(fw, strings.NewReader(sessionCookie)); err != nil {
		return
	}
	w.Close()


	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, 
		"http://localhost:8080/pdas/" + connectId + "/tokens/1", &b)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for FormData
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status:", resp.Status)
	}

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie = string(body)

	fmt.Println(sessionCookie)

	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/********************************** Send state to connect id **********************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Send state to pda @ connect id:")
	fmt.Println("*******************************************************************************")
	
	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/" + connectId + "/state", 
		strings.NewReader(sessionCookie))

	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var stateResp PeekResponse

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &stateResp); err != nil {
		panic(err)
	}

	sessionCookie = stateResp.Cookie
	fmt.Println(stateResp.Tokens)
	fmt.Println(sessionCookie)



	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/********************************** Send tokens to connect id *********************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Send /tokens/ to pda @ connect id:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/" + connectId + "/tokens", 
		strings.NewReader(sessionCookie))
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	var tokensResp PeekResponse

	// Read in the response body.
	body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &tokensResp); err != nil {
		panic(err)
	}

	sessionCookie = tokensResp.Cookie
	fmt.Println(tokensResp.Tokens)
	fmt.Println(sessionCookie)


	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/************************** Send EOS to connect id after position 3 ***************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Send EOS to pda @ connect id after position 3:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/" + connectId + "/eos/3", 
		strings.NewReader(sessionCookie))

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

	sessionCookie = string(body)

	if resp.StatusCode != http.StatusCreated {
		panic("Error, input stream not accepted")
	} else {
		fmt.Println("Input stream is accepted. Language recognized.")
		fmt.Println(sessionCookie)
	}

	/*********************************** Get connect address **************************************/

	connectId = GetConnectId(client)

	/************************ Send is_accepted to pda @ connect id ****************************/

	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Send is_accepted to pda @ connect id:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/pdas/" + connectId + "/is_accepted", 
		strings.NewReader(sessionCookie))

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

	sessionCookie = string(body)

	if resp.StatusCode != http.StatusAccepted {
		panic("Error, PDA not in accepted status")
	} else {
		fmt.Println("Success. PDA is in accepted status")
		fmt.Println(sessionCookie)
	}
}

func GetConnectId(client *http.Client) (string){
	fmt.Println()
	fmt.Println()
	fmt.Println("*******************************************************************************")
	fmt.Println("Return a random connection address for gid, 0:")
	fmt.Println("*******************************************************************************")

	// Set the http method, url, and request body
	req, err := http.NewRequest("GET", "http://localhost:8080/replica_pdas/0/connect", nil)
	
	if err != nil {
		panic(err)
	}

	// Set the request header Content-Type for json
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Read in the response body.
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		panic(err)
	}
	
	fmt.Println(string(body))

	connectAddress := strings.Split(string(body), "/")
	return connectAddress[len(connectAddress) - 1]
}