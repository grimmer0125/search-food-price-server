package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grimmer0125/bee/searchbot"
	"github.com/grimmer0125/bee/util"
)

func main() {
	r := gin.Default()
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
				_, _ = store, name
				if store != "" && name != "" {
					product := searchbot.QueryProduct(store, name)
					fmt.Println(product)

					remoteURL := "https://www.honestbee.tw/zh-TW/groceries/stores/carrefour/products/" + strconv.FormatFloat(product.ID, 'f', 0, 64)
					c.JSON(200, gin.H{
						"Price":           product.Price,
						"PreviewImageUrl": product.PreviewImageUrl,
						"Title":           product.Title,
						"RemoteURL":       remoteURL,
					})

					return
				}
			}
		}

		c.JSON(404, gin.H{"code": "SEARCH NOT_FOUND", "message": "SEARCH not found"})

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
