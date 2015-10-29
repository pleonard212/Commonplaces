package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

type resultObject struct {
	Author       string         `json:"author"`
	Title        string         `json:"title"`
	Date         int16          `json:"date"`
	LeftContext  string         `json:"leftContext"`
	RightContext string         `json:"rightContext"`
	MatchContext string         `json:"matchContext"`
	ContextLink  string         `json:"contextLink"`
	PassageID    string         `json:"passageID"`
	OtherTitles  map[string]int `json:"otherTitles,omitempty"`
}

type results struct {
	Commonplace resultObject   `json:"commonplace"`
	PassageList []resultObject `json:"passageList"`
	TitleList   []resultObject `json:"titleList"`
}

type metadataResultObject struct {
	Author       string `json:"author"`
	Title        string `json:"title"`
	Date         int16  `json:"date"`
	LeftContext  string `json:"leftContext"`
	RightContext string `json:"rightContext"`
	MatchContext string `json:"matchContext"`
	ContextLink  string `json:"contextLink"`
	PassageID    string `json:"passageID"`
}

type metadataResults struct {
	Count        int                    `json:"count"`
	MetadataList []metadataResultObject `json:"metadataList"`
}

type urlKeyValue struct {
	Key   string
	Value []string
}

type byDate []resultObject

func (slice byDate) Len() int {
	return len(slice)
}

func (slice byDate) Less(i, j int) bool {
	return slice[i].Date < slice[j].Date
}

func (slice byDate) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

var defaultConnConfig pgx.ConnConfig
var pool = createConnPool()

func createConnPool() *pgx.ConnPool {
	defaultConnConfig.Host = "localhost"
	defaultConnConfig.Database = "philologic"
	defaultConnConfig.User = "postgres"
	defaultConnConfig.Password = "***REMOVED***"
	config := pgx.ConnPoolConfig{ConnConfig: defaultConnConfig, MaxConnections: 10}
	pool, err := pgx.NewConnPool(config)
	if err != nil {
		fmt.Printf("Unable to create connection pool: %v", err)
	}
	return pool
}

func findCommonPlaces(c *gin.Context) {
	passageID := c.Query("passageident")
	rows, err := pool.Query("select sourceauthor, sourcetitle, sourcedate, sourceleftcontext, sourcematchcontext, sourcerightcontext, sourcecontextlink, targetauthor, targettitle, targetdate, targetleftcontext, targetmatchcontext, targetrightcontext, targetcontextlink from eccotcpxallgale where passageident=$1", passageID)
	if err != nil {
		c.JSON(200, results{})
	}

	defer rows.Close()

	filteredAuthors := make(map[string]resultObject, 0)
	filteredTitles := make(map[string]resultObject, 0)
	for rows.Next() {
		var author string
		var targetAuthor string
		var title string
		var targetTitle string
		var date int16
		var targetDate int16
		var leftContext string
		var targetLeftContext string
		var rightContext string
		var targetRightContext string
		var matchContext string
		var targetMatchContext string
		var contextLink string
		var targetContextLink string
		err := rows.Scan(&author, &title, &date, &leftContext, &matchContext, &rightContext, &contextLink, &targetAuthor, &targetTitle, &targetDate, &targetLeftContext, &targetMatchContext, &targetRightContext, &targetContextLink)
		if err != nil {
			fmt.Println(err)
		}
		otherTitles := make(map[string]int, 0)
		sourceObject := resultObject{author, title, date, leftContext, rightContext, matchContext, contextLink, passageID, otherTitles}
		if _, ok := filteredAuthors[author]; !ok {
			filteredAuthors[author] = sourceObject
		} else if _, ok := filteredAuthors[author]; ok {
			if filteredAuthors[author].Date > date {
				sourceObject.OtherTitles = filteredAuthors[author].OtherTitles
				filteredAuthors[author] = sourceObject
			} else if filteredAuthors[author].Date == sourceObject.Date && len(filteredAuthors[author].MatchContext) < len(sourceObject.MatchContext) {
				sourceObject.OtherTitles = filteredAuthors[author].OtherTitles
				filteredAuthors[author] = sourceObject
			}
			if filteredAuthors[author].Date != sourceObject.Date {
				filteredAuthors[author].OtherTitles[sourceObject.Title] = 1
			}
		}
		if _, ok := filteredTitles[sourceObject.Title]; !ok {
			filteredTitles[sourceObject.Title] = sourceObject
		} else if filteredTitles[sourceObject.Title].Date > sourceObject.Date {
			filteredTitles[sourceObject.Title] = sourceObject
		}
		// Process target results
		targetOtherTitles := make(map[string]int, 0)
		targetObject := resultObject{targetAuthor, targetTitle, targetDate, targetLeftContext, targetRightContext, targetMatchContext, targetContextLink, passageID, targetOtherTitles}
		if _, ok := filteredAuthors[targetAuthor]; !ok {
			filteredAuthors[targetAuthor] = targetObject
		} else if _, ok := filteredAuthors[targetAuthor]; ok {
			if filteredAuthors[targetAuthor].Date > date {
				targetObject.OtherTitles = filteredAuthors[targetAuthor].OtherTitles
				filteredAuthors[targetAuthor] = targetObject
			} else if filteredAuthors[targetAuthor].Date == targetObject.Date && len(filteredAuthors[targetAuthor].MatchContext) < len(targetObject.MatchContext) {
				targetObject.OtherTitles = filteredAuthors[targetAuthor].OtherTitles
				filteredAuthors[targetAuthor] = targetObject
			}
			if filteredAuthors[targetAuthor].Date != targetObject.Date && len(filteredAuthors[targetAuthor].OtherTitles) > 0 {
				filteredAuthors[targetAuthor].OtherTitles[targetObject.Title] = 1
			}
		}
		if _, ok := filteredTitles[targetTitle]; !ok {
			filteredTitles[targetObject.Title] = targetObject
		} else if filteredTitles[targetTitle].Date > targetObject.Date {
			filteredTitles[targetObject.Title] = targetObject
		}
	}
	var passageList []resultObject
	for _, value := range filteredAuthors {
		passageList = append(passageList, value)
	}
	sort.Sort(byDate(passageList))
	var titleList []resultObject
	for _, value := range filteredTitles {
		titleList = append(titleList, value)
	}
	sort.Sort(byDate(titleList))
	fullResults := results{passageList[0], passageList[1:], titleList}
	c.JSON(200, fullResults)
}

func metadataQuery(c *gin.Context) {
	queryStringMap, _ := url.ParseQuery(c.Request.URL.RawQuery)
	fmt.Println(queryStringMap)
	query := "select sourceauthor, sourcetitle, sourcedate, sourceleftcontext, sourcematchcontext, sourcerightcontext, sourcecontextlink, passageident from eccotcpxallgale where "
	count := 1
	var params []string
	var values []interface{}
	for param, v := range queryStringMap {
		for _, value := range v {
			params = append(params, param+" = $"+strconv.Itoa(count))
			values = append(values, value)
			count++
		}
	}
	query += strings.Join(params, " and ")
	fmt.Printf("query is:%s\n", query)
	fmt.Println(values)
	rows, err := pool.Query(query, values...)
	if err != nil {
		var emptyResults []metadataResultObject
		fmt.Println("query failed")
		c.JSON(200, metadataResults{0, emptyResults})
	}

	defer rows.Close()

	var results metadataResults
	for rows.Next() {
		var author string
		var title string
		var date int16
		var leftContext string
		var rightContext string
		var matchContext string
		var contextLink string
		var passageID string
		err := rows.Scan(&author, &title, &date, &leftContext, &matchContext, &rightContext, &contextLink, &passageID)
		if err != nil {
			var emptyResults []metadataResultObject
			fmt.Println("retrieving results of query failed")
			c.JSON(200, metadataResults{0, emptyResults})
		}
		sourceResults := metadataResultObject{author, title, date, leftContext, matchContext, rightContext, contextLink, passageID}
		results.MetadataList = append(results.MetadataList, sourceResults)
	}
	results.Count = len(results.MetadataList)
	c.JSON(200, results)
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func main() {
	router := gin.Default()

	// Static files
	router.LoadHTMLFiles("index.html")
	router.Static("public", "./public")
	router.Static("components", "./public/components")
	router.Static("css", "./public/css")
	// Routes
	router.GET("/DiggingIntoData/", index)
	router.GET("/DiggingIntoData/query", index)
	router.GET("/DiggingIntoData/commonplaces", findCommonPlaces)

	router.GET("/DiggingIntoData/metadata", metadataQuery)

	router.Run(":3000")
}