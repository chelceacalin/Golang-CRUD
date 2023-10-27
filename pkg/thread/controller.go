package thread

import (
	"First/model"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

func GetThreads(c *gin.Context, db *sql.DB) {

	rows, err := db.Query(`SELECT t.id, t.title, m.id, m.message, m.thread_id
        FROM thread t
        LEFT JOIN message m ON t.id = m.thread_id`)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error executing query")
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Rows error")
		}
	}(rows)

	threadsMap := make(map[int]*model.Thread)

	for rows.Next() {
		var thread model.Thread
		var message model.Message
		err := rows.Scan(&thread.Id, &thread.Title, &message.Id, &message.Message, &message.Thread_id)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error scanning rows")
			return
		}

		existingThread, ok := threadsMap[thread.Id]
		if ok {
			existingThread.Messages = append(existingThread.Messages, message)
		} else {
			thread.Messages = []model.Message{message}
			threadsMap[thread.Id] = &thread
		}
	}

	var threads []model.Thread
	for _, thread := range threadsMap {
		threads = append(threads, *thread)
	}

	c.HTML(http.StatusOK, "threads/thread.tmpl", gin.H{
		"threads": threads,
	})

}

func GetThreadById(c *gin.Context, db *sql.DB) {
	threadID := c.Param("id")
	id, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid thread ID")
		return
	}

	rows, err := db.Query(`SELECT t.id, t.title, m.id, m.message, m.thread_id
		FROM thread t
		LEFT JOIN message m ON t.id = m.thread_id
		WHERE t.id = $1`, id)

	if err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving thread and messages")
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var thread model.Thread
	for rows.Next() {
		var message model.Message
		err := rows.Scan(&thread.Id, &thread.Title, &message.Id, &message.Message, &message.Thread_id)
		if err != nil {
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
func AddNewThread(c *gin.Context, db *sql.DB) {
	threadId := rand.Intn(2147483647)
	title, ok := c.GetPostForm("title")
	if !ok {
		fmt.Println("EROARE ")
	}

	message, ok := c.GetPostForm("message")
	if !ok {
		fmt.Println("Nu am receptionat mesaju ")
	}
	messageId := rand.Intn(2147483647)

	// Convert messageId and threadId to integers
	messageIdStr := strconv.Itoa(messageId)
	threadIdStr := strconv.Itoa(threadId)

	_, err := db.Exec("INSERT INTO message VALUES ($1, $2, $3)", messageIdStr, message, threadIdStr)
	if err != nil {
		fmt.Println("Error inserting message for thread", err)
		c.String(http.StatusInternalServerError, "Error inserting message for thread", err)
		return
	}

	_, err = db.Exec("INSERT INTO Thread VALUES ($1, $2)", threadIdStr, title)
	if err != nil {
		fmt.Println("Error inserting thread", err)
		c.String(http.StatusInternalServerError, "Error inserting thread", err)
		return
	}
	c.Redirect(http.StatusFound, "/threads")
}

func DeleteThreadById(c *gin.Context, db *sql.DB) {
	threadId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid thread ID")
		return
	}

	_, err = db.Exec("DELETE FROM message WHERE thread_id = $1", threadId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting messages", err)
		return
	}

	_, err = db.Exec("DELETE FROM thread WHERE id = $1", threadId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting thread", err)
		return
	}
	c.Redirect(http.StatusFound, "/threads")
}

func GetAddThreadForm(c *gin.Context, db *sql.DB) {
	c.HTML(http.StatusOK, "addThread/thread.tmpl", gin.H{})
}

func EditThreadById(c *gin.Context, db *sql.DB) {
	threadID := c.Param("id")
	id, err := strconv.Atoi(threadID)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid thread ID", err)
		return
	}

	rows, err := db.Query(`SELECT t.id, t.title, m.id, m.message, m.thread_id
		FROM thread t
		LEFT JOIN message m ON t.id = m.thread_id
		WHERE t.id = $1`, id)

	if err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving thread and messages", err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var thread model.Thread
	for rows.Next() {
		var message model.Message
		err := rows.Scan(&thread.Id, &thread.Title, &message.Id, &message.Message, &message.Thread_id)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error scanning rows", err)
			return
		}
		thread.Messages = append(thread.Messages, message)
	}

	c.HTML(http.StatusOK, "editById/thread.tmpl", gin.H{
		"id":       thread.Id,
		"title":    thread.Title,
		"messages": thread.Messages,
	})
}

func UpdateThread(c *gin.Context, db *sql.DB) {
	threadId, ok := c.GetPostForm("id")
	title, ok := c.GetPostForm("title")

	if !ok {
		fmt.Println("EROARE ")
	}

	message, ok := c.GetPostForm("message")
	if !ok {
		fmt.Println("EROARE Mesage ")
	}
	messageId, ok := c.GetPostForm("messageId")

	_, err := db.Exec("update message set message=$1 where id=$2", message, messageId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting messages", err)
		return
	}

	_, err = db.Exec("update thread set title=$1 where id=$2", title, threadId)
	if err != nil {
		fmt.Println("Error updating thread", err)
		return
	}
	c.Redirect(http.StatusFound, string("/threads/"+threadId))

}
