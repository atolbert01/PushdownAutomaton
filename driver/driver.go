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

	/********************************* Create new replica group ***********************************/

	fmt.Println()
	fmt.Println("Create new replica group pda with gid 0, member ids 0, 2, 4, 5: ")
	
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
	fmt.Println("Create new replica group pda with gid 2, member ids 0, 1, 3: ")
	
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
	if _, err = io.Copy(fw, strings.NewReader("0 1 3")); err != nil {
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
	fmt.Println("Return all currently defined replica ids:")

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
	fmt.Println("Return replica group members for gid, 0:")

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


	/*********************************** Get connect address **************************************/


	fmt.Println()
	fmt.Println("Return a random connection address for gid, 0:")

	// Set the http method, url, and request body
	req, err = http.NewRequest("GET", "http://localhost:8080/replica_pdas/0/connect", nil)
	
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

	/********************************* Create new pda with id 0/ **********************************/

	// fmt.Println()
	
	// fmt.Println("Create new pda with id 0 - ")
	
	// jsonText, err := ioutil.ReadFile("helloPda.json")

	// // Set the http method, url, and request body
	// req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0", bytes.NewBuffer(jsonText))
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// var pda PdaProcessor

	// // Read in the response body.
	// body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err := json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /********************************* Create new pda with id 1/ **********************************/

	// fmt.Println()
	
	// fmt.Println("Create new pda with id 1 - ")
	
	// jsonText, err = ioutil.ReadFile("goodbyePda.json")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/1", bytes.NewBuffer(jsonText))
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err = json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /********************************* Create new pda with id 3/ **********************************/

	// fmt.Println()
	
	// fmt.Println("Create new pda with id 3 - ")
	
	// jsonText, err = ioutil.ReadFile("alohaPda.json")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/3", bytes.NewBuffer(jsonText))
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err = json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /********************************* Check to see if pda exists *********************************/
	
	// fmt.Println()
	
	// fmt.Println("Get list of pdas - ")
	
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// req.Header.Set("Cache-Control","no-cache")

	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))

	// /********************************* Reset pda with id 3/ **********************************/

	// fmt.Println()
	
	// fmt.Println("Reset pda with id 3 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/3/reset", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// //req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err = json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /***************************** Present token 0, position 0 to PDA 0 ****************************/

	// fmt.Println()
	
	// fmt.Println("Present token 0 with position 0 to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/0", strings.NewReader("0"))
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err = json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /**************************** Present token 1, position 2 to PDA 0 *****************************/

	// fmt.Println()
	
	// fmt.Println("Present token 1 with position 2 to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/2", strings.NewReader("1"))
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err = json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /**************************** Present token 1, position 3 to PDA 0 *****************************/

	// fmt.Println()
	
	// fmt.Println("Present token 1 with position 3 to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/3", strings.NewReader("1"))
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err = json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /********************************** Send len to PDA 0 *********************************/

	// fmt.Println()
	
	// fmt.Println("Send len to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/stack/len", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))
	// /********************************** Send peek(2) to PDA 0 *********************************/

	// fmt.Println()
	
	// fmt.Println("Send peek(2) to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/stack/top/2", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))

	// /********************************** Send state to PDA 0 *********************************/

	// fmt.Println()
	
	// fmt.Println("Send state to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/state", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))

	// /********************************** Send /tokens to PDA 0 *********************************/

	// fmt.Println()

	// fmt.Println("Send /tokens to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/tokens", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))


	// /*********************************** Get snapshot of PDA 0 ************************************/

	// fmt.Println()
	
	// fmt.Println("Get snapshot of PDA 0 with top 2 elements - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/snapshot/2", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// var snap Snap

	// if err = json.Unmarshal(body, &snap); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(snap)

	// /**************************** Present token 0, position 1 to PDA 0 *****************************/

	// fmt.Println()
	
	// fmt.Println("Present token 0 with position 1 to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/tokens/1", strings.NewReader("0"))
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// if err = json.Unmarshal(body, &pda); err != nil {
	// 	panic(err)
	// }

	// fmt.Println(pda)

	// /********************************** Send is_accepted to PDA 0 *********************************/

	// fmt.Println()
	
	// fmt.Println("Send is_accepted to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/is_accepted", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))

	// /***************************** Send EOS to PDA 0 after position 3 *****************************/

	// fmt.Println()
	
	// fmt.Println("Send EOS to PDA 0 after position 3 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/pdas/0/eos/3", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))

	// /********************************** Send is_accepted to PDA 0 *********************************/

	// fmt.Println()
	
	// fmt.Println("Send is_accepted to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("GET", "http://localhost:8080/pdas/0/is_accepted", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// // Set the request header Content-Type for json
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))

	// /********************************** Send delete() to PDA 0 *********************************/

	// fmt.Println()
	
	// fmt.Println("Send delete() to PDA 0 - ")

	// // Set the http method, url, and request body
	// req, err = http.NewRequest("DELETE", "http://localhost:8080/pdas/0/delete", nil)
	
	// if err != nil {
	// 	panic(err)
	// }

	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// // Read in the response body.
	// body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(string(body))
}