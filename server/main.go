package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	// database
	"database/sql"

	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	bin, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, err
	}
	return sql.Open("postgres", fmt.Sprintf("postgres://postgres:%s@db:5432/example?sslmode=disable", string(bin)))
}

func prepare() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS test"); err != nil {
		return err
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS test (id SERIAL, content VARCHAR)"); err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		if _, err := db.Exec("INSERT INTO test (content) VALUES ($1);", fmt.Sprintf("Blog post #%d", i)); err != nil {
			return err
		}
	}
	return nil
}

func getall(c *gin.Context) {
	db, err := connect()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT content FROM test")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	var contents []string
	for rows.Next() {
		var content string
		_ = rows.Scan(&content)
		contents = append(contents, content)
	}

	c.IndentedJSON(http.StatusOK, contents)
}

//go:embed dist
var frontend embed.FS

func main() {
	path := "./config.toml"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	err := parseConf(path)
	if err != nil {
		log.Fatal(err)
	}

	conf := getConf()
	router := gin.Default()

	// middleware for static files (frontend)
	router.Use(static.Serve("/", static.EmbedFolder(frontend, "dist")))

	router.GET("/api/ready", checkServer)

	// test db
	log.Print("Prepare db...")
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	router.GET("/api/db/", getall)

	// Read index.html into memory
	index, err := frontend.ReadFile("dist/index.html")
	if err != nil {
		log.Fatal(err)
	}

	// If no route match is probably a vue route. So we return the index.html
	// and the vue-router takes from here
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.String(), "/api") {
			c.JSON(http.StatusNotFound, "")
			return
		}
		c.Data(http.StatusOK, "text/html", index)
	})

	router.Run(conf.Server.Bind)
}
