package thread

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func GetThreadRoutes() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	db, err := sql.Open("pgx", "user=workshop_go password=pass host=localhost port=5433 database=workshop_go sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router.GET("/", func(c *gin.Context) {
		GetThreads(c, db)
	})

	router.GET("/threads", func(c *gin.Context) {
		GetThreads(c, db)
	})

	router.GET("/threads/:id", func(c *gin.Context) {
		GetThreadById(c, db)
	})

	router.POST("/threads/:id", func(c *gin.Context) {
		DeleteThreadById(c, db)
	})

	router.GET("/addThread/threads", func(c *gin.Context) {
		AddThread(c, db)
	})

	router.POST("/addThread/threads", func(c *gin.Context) {
		AddnewThread(c, db)
	})
	router.GET("/edit/:id", func(c *gin.Context) {
		EditThreadById(c, db)
	})

	router.POST("/update/:id", func(c *gin.Context) {
		UpdateThread(c, db)
	})

	return router
}
