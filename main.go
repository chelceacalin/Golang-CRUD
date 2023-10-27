package main

import (
	"First/structs"
	"context"
	"database/sql"
	"fmt"
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"math/rand"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	db, err := sql.Open("pgx", "user=workshop_go password=pass host=localhost port=5433 database=workshop_go sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router.GET("/threads", func(c *gin.Context) {
		getThreads(c, db)
	})

	router.GET("/threads/:id", func(c *gin.Context) {
		getThreadById(c, db)
	})

	router.POST("/threads/:id", func(c *gin.Context) {
		deleteThreadById(c, db)
	})

	router.GET("/addThread/threads", func(c *gin.Context) {
		addThread(c, db)
	})

	router.POST("/addThread/threads", func(c *gin.Context) {
		addnewThread(c, db)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
func getThreads(c *gin.Context, db *sql.DB) {

	rows, err := db.Query(`SELECT t.id, t.title, m.id, m.message, m.thread_id
        FROM thread t
        LEFT JOIN message m ON t.id = m.thread_id`)
	if err != nil {
		fmt.Println("uncu", err)
		c.String(http.StatusInternalServerError, "Error executing query")
		return
	}
	defer rows.Close()

	threadsMap := make(map[int]*structs.Thread)

	for rows.Next() {
		var thread structs.Thread
		var message structs.Message
		err := rows.Scan(&thread.Id, &thread.Title, &message.Id, &message.Message, &message.Thread_id)
		if err != nil {
			fmt.Println("ciuciu", err)
			c.String(http.StatusInternalServerError, "Error scanning rows")
			return
		}

		existingThread, ok := threadsMap[thread.Id]
		if ok {
			existingThread.Messages = append(existingThread.Messages, message)
		} else {
			thread.Messages = []structs.Message{message}
			threadsMap[thread.Id] = &thread
		}
	}

	var threads []structs.Thread
	for _, thread := range threadsMap {
		threads = append(threads, *thread)
	}

	c.HTML(http.StatusOK, "threads/thread.tmpl", gin.H{
		"threads": threads,
	})

}

func getThreadById(c *gin.Context, db *sql.DB) {
	threadID := c.Param("id")
	id, err := strconv.Atoi(threadID)
	if err != nil {
		fmt.Println("tactu", err)
		c.String(http.StatusBadRequest, "Invalid thread ID")
		return
	}

	rows, err := db.Query(`SELECT t.id, t.title, m.id, m.message, m.thread_id
		FROM thread t
		LEFT JOIN message m ON t.id = m.thread_id
		WHERE t.id = $1`, id)

	if err != nil {
		fmt.Println("rahat", err)
		c.String(http.StatusInternalServerError, "Error retrieving thread and messages")
		return
	}
	defer rows.Close()

	var thread structs.Thread
	for rows.Next() {
		var message structs.Message
		err := rows.Scan(&thread.Id, &thread.Title, &message.Id, &message.Message, &message.Thread_id)
		if err != nil {
			fmt.Println("gucci", err)
			c.String(http.StatusInternalServerError, "Error scanning rows")
			return
		}
		thread.Messages = append(thread.Messages, message)
	}

	c.HTML(http.StatusOK, "threadsById/thread.tmpl", gin.H{
		"id":       thread.Id,
		"title":    thread.Title,
		"messages": thread.Messages,
	})
}

func deleteThreadById(c *gin.Context, db *sql.DB) {
	threadId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("Invalid thread ID")
		c.String(http.StatusBadRequest, "Invalid thread ID")
		return
	}

	_, err = db.Exec("DELETE FROM message WHERE thread_id = $1", threadId)
	if err != nil {
		fmt.Println("Error deleting messages")
		c.String(http.StatusInternalServerError, "Error deleting messages")
		return
	}

	_, err = db.Exec("DELETE FROM thread WHERE id = $1", threadId)
	if err != nil {
		fmt.Println("Error deleting thread")
		c.String(http.StatusInternalServerError, "Error deleting thread")
		return
	}

	getThreads(c, db)
}

func addThread(c *gin.Context, db *sql.DB) {
	c.HTML(http.StatusOK, "addThread/thread.tmpl", gin.H{})
}

func addnewThread(c *gin.Context, db *sql.DB) {
	id := rand.Intn(2147483647)
	title, ok := c.GetPostForm("title")
	if !ok {
		fmt.Println("EROARE ")
	}

	message, ok := c.GetPostForm("message")
	if !ok {
		fmt.Println("EROARE Mesage ")
	}
	messageId := strconv.Itoa(rand.Intn(2147483647))

	fmt.Println("message_id:", messageId)
	fmt.Println("message:", message)
	_, err := db.Exec("insert into message values($1,$2,$3)", messageId, message, id)
	if err != nil {
		fmt.Println("Error inserting message", err)
		c.String(http.StatusInternalServerError, "Error deleting messages")
		return
	}

	_, err = db.Exec("insert into Thread values($1,$2)", id, title)
	if err != nil {
		fmt.Println("Error inserting thread", err)
		c.String(http.StatusInternalServerError, "Error deleting messages")
		return
	}
	getThreads(c, db)
}
