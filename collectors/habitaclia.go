package collectors

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"alquiler-scrapping/database"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

func CollectHabitacliaEntries(db *database.Database, maxPrice int) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.habitaclia.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:126.0) Gecko/20100101 Firefox/126.0"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("accept-encoding", "gzip, deflate, br, zstd")
		r.Headers.Set("accept-language", "en-US,en;q=0.5")
		r.Headers.Set("cache-control", "no-cache")
		r.Headers.Set("sec-fetch-dest", "document")
		r.Headers.Set("sec-fetch-mode", "navigate")
		r.Headers.Set("sec-fetch-site", "same-origin")
		r.Headers.Set("sec-fetch-user", "?1")
		r.Headers.Set("sec-gpc", "1")
		r.Headers.Set("upgrade-insecure-requests", "1")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("DNT", "1")
	})

	c.OnHTML("article[id]", func(e *colly.HTMLElement) {

		price, err := strconv.Atoi(strings.TrimSuffix(e.ChildText("span[itemprop=price]"), " €"))
		if err != nil {
			price = -1
		}

		if price > maxPrice {
			return
		}

		// Elimina el filtro de la url. Empieza por "?". Ejemplo:
		// https://www.habitaclia.com/alquiler-piso-apartamento_sin_amueblar_en_les_corts_calle_del_comte_de_guell_48-barcelona-i27246000000071.htm?pmax=700&ady=1&f=&geo=p&from=list&lo=60
		url := strings.Split(e.Attr("data-href"), "?")[0]

		db.Insert(database.Entry{
			Title: e.ChildText("a[itemprop=name]"),
			Price: price,
			Url:   url,
		})
	})

	c.OnError(func(r *colly.Response, e error) {
		if os.Getenv("DEBUG") == "true" {
			log.Println("error:", e, r.Request.URL, string(r.Body))
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if os.Getenv("DEBUG") == "true" {
			log.Println("error:", r.Request.URL, string(r.Body))
		}
	})

	// Habitaclia muestra anuncios por encima del limite que especificamos, luego lo filtramos en el callback OnHTML
	c.Visit(fmt.Sprintf("https://www.habitaclia.com/alquiler-barcelona.htm?pmax=%d", maxPrice+150))
}
