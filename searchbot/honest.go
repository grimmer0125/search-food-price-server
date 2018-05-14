//Package searchbot implements query/crawler/bot part
package searchbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Product struct {
	ID    float64 `json:"id"`
	Title string  `json:"title"`
	// description     string
	// imageUrl        string
	PreviewImageUrl string `json:"previewImageUrl"`
	// size            string
	// status          string // "status_available",
	// currency        string
	Price string `json:"price"`
}

// type queryResult struct {
// 	Products   []map[string]interface{} `json:"products"`
// 	Meta       map[string]int           `json:"meta"`
// 	Categories []map[string]interface{} `json:"categories"`
// }

type queryResult struct {
	Products   []Product                `json:"products"` // using []Product does works, but vs code's debugger can not show it, a bug
	Meta       map[string]int           `json:"meta"`
	Categories []map[string]interface{} `json:"categories"`
}

// type meta struct {
// 	current_page int
// 	total_pages  int
// 	total_count  int
// }

// QueryProduct for carrefour etc
func QueryProduct(store, productName string) (rProduct Product) {

	if store == "carrefour" {

		var productList []Product
		totalPages := 0
		page := 1
		encodedName := url.QueryEscape(productName)

		for ; totalPages == 0 || page <= totalPages; page++ {
			queryURL := "https://www.honestbee.tw/api/api/stores/3932?q=" + encodedName + "&sort=relevance&page=" + strconv.Itoa(page)

			fmt.Println("start to query:%s; page:%d", queryURL, page)
			req, err := http.NewRequest("GET", queryURL, nil)
			if err != nil {
				return
			}
			req.Header.Set("Accept", "application/vnd.honestbee+json;version=2")
			req.Header.Set("Accept-Language", "zh-TW")

			c := http.Client{}
			res, err := c.Do(req)
			if err != nil {
				return
			}
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return
			}

			var resp queryResult //map[string]interface{}

			if err := json.Unmarshal(b, &resp); err != nil {
				fmt.Println("can not parse resp")
				fmt.Println(err)
				return
			}

			fmt.Println("ok")

			if _, ok := resp.Meta["total_pages"]; ok == false {
				fmt.Println("no totalpages info.")
				return
			}

			fmt.Println(resp.Meta["current_page"])
			fmt.Println(resp.Meta["total_pages"])
			fmt.Println(resp.Meta["total_count"])
			totalPages = resp.Meta["total_pages"]

			productList = append(productList, resp.Products...)

			fmt.Println(len(resp.Products))

			fmt.Println("query ok")

		}

		fmt.Println("finish all pages search")
		fmt.Println(len(productList))

		var lowerestProduct Product
		for i, p := range productList {
			if i == 0 {
				lowerestProduct = p
			} else {
				price1, err1 := strconv.ParseFloat(lowerestProduct.Price, 64)
				price2, err2 := strconv.ParseFloat(p.Price, 64)

				if err1 == nil && err2 == nil && price2 < price1 {
					lowerestProduct = p
				}
			}
		}

		fmt.Println("lowest product:")
		fmt.Println(lowerestProduct)
		rProduct = lowerestProduct
	}

	return
}
