package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var ErrInvalidID = errors.New("invalid ID")

const (
	NotFound            = "404 not found"
	InternalServerError = "500 internal server error"
	BadRequest          = "400 bad request"
)

type Request struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Category string   `json:"category" binding:"required"`
	Tags     []string `json:"tags" binding:"required"`
}

type Response struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreatePost(storage *pgxpool.Pool, logger *logrus.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		var request Request
		var response Response

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, BadRequest)
			logger.Error(err)
			return
		}

		stmt := `
			INSERT INTO posts (title, content, category, tags, created_at, updated_at)
			VALUES ($1, $2, $3, $4, timestamp 'now()', timestamp 'now()') RETURNING *;
		`

		row := storage.QueryRow(context.Background(), stmt,
			request.Title, request.Content, request.Category, request.Tags)
		if err := row.Scan(&response.Id, &response.Title, &response.Content, &response.Category,
			&response.Tags, &response.CreatedAt, &response.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, InternalServerError)
			logger.Error(err)
			return
		}

		c.JSON(http.StatusCreated, response)
	}
}

func UpdatePost(storage *pgxpool.Pool, logger *logrus.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		if err := validateID(c.Param("id")); err != nil {
			c.JSON(http.StatusNotFound, NotFound)
			logger.Error(err)
			return
		}

		var request Request
		var response Response

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, BadRequest)
			logger.Error(err)
			return
		}

		stmt := `
			UPDATE posts
			SET title = $1, content = $2, category = $3, tags = $4, updated_at = timestamp 'now()'
			WHERE id = $5 RETURNING *;
		`

		row := storage.QueryRow(context.Background(), stmt,
			request.Title, request.Content, request.Category, request.Tags, c.Param("id"))

		if err := row.Scan(&response.Id, &response.Title, &response.Content, &response.Category,
			&response.Tags, &response.CreatedAt, &response.UpdatedAt); err != nil {
			switch err {
			case pgx.ErrNoRows:
				c.JSON(http.StatusNotFound, NotFound)
				logger.Error(err)
				return
			default:
				c.JSON(http.StatusInternalServerError, InternalServerError)
				logger.Error(err)
				return
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// TODO: fix not existing ID's
func DeletePost(storage *pgxpool.Pool, logger *logrus.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		if err := validateID(c.Param("id")); err != nil {
			c.JSON(http.StatusNotFound, NotFound)
			logger.Error(err)
			return
		}

		if _, err := storage.Exec(context.Background(), "DELETE FROM posts WHERE id = $1", c.Param("id")); err != nil {
			switch err {
			case pgx.ErrNoRows:
				c.JSON(http.StatusNotFound, NotFound)
				logger.Error(err)
				return
			default:
				c.JSON(http.StatusInternalServerError, InternalServerError)
				logger.Error(err)
				return
			}
		}

		c.Status(http.StatusNoContent)
	}
}

func GetPost(storage *pgxpool.Pool, logger *logrus.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		if err := validateID(c.Param("id")); err != nil {
			c.JSON(http.StatusNotFound, NotFound)
			logger.Error(err)
			return
		}

		var response Response

		row := storage.QueryRow(context.Background(), "SELECT * FROM posts WHERE id = $1", c.Param("id"))

		if err := row.Scan(&response.Id, &response.Title, &response.Content, &response.Category,
			&response.Tags, &response.CreatedAt, &response.UpdatedAt); err != nil {
			switch err {
			case pgx.ErrNoRows:
				c.JSON(http.StatusNotFound, NotFound)
				logger.Error(err)
				return
			default:
				c.JSON(http.StatusInternalServerError, InternalServerError)
				logger.Error(err)
				return
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetAllPosts(storage *pgxpool.Pool, logger *logrus.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		var response []Response

		rows, err := storage.Query(context.Background(), "SELECT * FROM posts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, InternalServerError)
			logger.Error(err)
			return
		}
		for rows.Next() {
			var res Response
			if err := rows.Scan(&res.Id, &res.Title, &res.Content, &res.Category, &res.Tags, &res.CreatedAt, &res.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, InternalServerError)
				logger.Error(err)
				return
			}
			response = append(response, res)
		}

		filtered := responseFilter(response, c.Query("term"))
		if len(filtered) == 0 {
			c.JSON(http.StatusNotFound, NotFound)
			return
		}

		c.JSON(http.StatusOK, filtered)
	}
}
func validateID(id string) error {
	idint, err := strconv.Atoi(id)
	if err != nil {
		return ErrInvalidID
	}
	if idint < 1 {
		return ErrInvalidID
	}
	return nil
}

func responseFilter(posts []Response, search string) (filtered []Response) {
	for _, post := range posts {
		if strings.Contains(post.Title, search) ||
			strings.Contains(post.Content, search) ||
			strings.Contains(post.Category, search) {
			filtered = append(filtered, post)
		} else {
			for _, tag := range post.Tags {
				if strings.Contains(tag, search) {
					filtered = append(filtered, post)
				}
			}
		}
	}
	return
}
