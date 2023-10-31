package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Payload struct {
	Foo string `json:"foo"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s input\n", os.Args[0])
		os.Exit(1)
	}
	input := os.Args[1]

	bearerToken := os.Getenv("GONUSLY_TOKEN")
	if bearerToken == "" {
		fmt.Println("GONUSLY_TOKEN not set in the environment")
		os.Exit(1)
	}

	/*
	   data := Payload{
	       Foo: input,
	   }
	   payloadBytes, err := json.Marshal(data)
	   if err != nil {
	       fmt.Println(err.Error())
	       os.Exit(1)
	   }

	   body := bytes.NewReader(payloadBytes)
	*/

	encodedEmail := url.QueryEscape(input)

	req, err := http.NewRequest(
		"GET",
		"https://bonus.ly/api/v1/users/autocomplete?search="+encodedEmail,
		nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Fatal(err)
		}
		prettyJson, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(prettyJson))
	} else {
		fmt.Println("API call was not successful.")
	}
}
