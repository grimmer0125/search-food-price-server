package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grimmer0125/bee/searchbot"
	"github.com/grimmer0125/bee/util"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

// each http request will live in an individual goroutine
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	_, _ = slackToken, slackChannel

	//r := gin.Default()
	r := gin.New()
	r.Use(func(context *gin.Context) {
		// add header Access-Control-Allow-Origin
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Next()
	})
	r.POST("/rpc", func(c *gin.Context) {

		fmt.Println("get rpc call")

		// alternative: gin's bind
		b, _ := ioutil.ReadAll(c.Request.Body) // []unit8
		fmt.Printf("[request body] %s", string(b))

		// alternative: interface{}/struct
		var reqBody map[string]interface{}
		if err := json.Unmarshal(b, &reqBody); err != nil {
			fmt.Println("can not parse request body")
		}

		method := util.GetStringProperty(reqBody, "method")
		if method == "queryProduct" {
			if reqBody["params"] != nil {
				params := reqBody["params"].(map[string]interface{})

				store := util.GetStringProperty(params, "store")
				name := util.GetStringProperty(params, "productName")
				shownPrice := util.GetStringProperty(params, "shownPrice")
				sourceURL := util.GetStringProperty(params, "sourceURL")
				if store != "" && name != "" && shownPrice != "" && sourceURL != "" {

					product := searchbot.QueryProduct(store, name)
					fmt.Println(product)

					remoteURL := "https://www.honestbee.tw/zh-TW/groceries/stores/" + store + "/products/" + strconv.FormatFloat(product.ID, 'f', 0, 64)
					c.JSON(200, gin.H{
						"Price":           product.Price,
						"PreviewImageUrl": product.PreviewImageUrl,
						"Title":           product.Title,
						"RemoteURL":       remoteURL,
					})

					fmt.Printf("shown price:%s", shownPrice)
					if product.Price != "" {
						priceHonest, _ := strconv.ParseFloat(product.Price, 64)
						price, _ := strconv.ParseFloat(shownPrice, 64)

						if priceHonest > price {
							message := "honesetbee is not cheaper than " + sourceURL
							// send a slack mesasge
							notifyAdmin(slackToken, slackChannel, message)
						}
					}

					return
				}
			}
		}

		c.JSON(404, gin.H{"code": "SEARCH NOT_FOUND", "message": "SEARCH not found"})

	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

// use the method in https://stackoverflow.com/a/44883343/7354486 to get channelID
func notifyAdmin(token, channelID, message string) {

	if token == "" || channelID == "" {
		return
	}

	api := slack.New(token)

	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api.SetDebug(true)

	// Ref: https://github.com/nlopes/slack/blob/master/examples/messages/messages.go#L26:2
	params := slack.PostMessageParameters{}
	channelID, timestamp, err := api.PostMessage(channelID, message, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

}
