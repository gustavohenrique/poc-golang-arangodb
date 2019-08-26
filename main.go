package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/vst"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Teacher struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	City   string  `json:"city,omitempty"`
	Audios []Audio `json:"audios,omitempty"`
}

type Student struct {
	Teacher
	Playlists []Playlist `json:"playlists"`
}

type Subject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Audio struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Media     string     `json:"media"`
	Teachers  []Teacher  `json:"teachers,omitempty"`
	Playlists []Playlist `json:"playlists,omitempty"`
}

type Playlist struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Audios   []Audio   `json:"audios,omitempty"`
	Students []Student `json:"students,omitempty"`
	Subjects []Subject `json:"subjects,omitempty"`
}

func connect() (driver.Client, error) {
	host := "vst://127.0.0.1:8529"
	if os.Getenv("ARANGODB_HOST") != "" {
		host = fmt.Sprintf("vst://%s", os.Getenv("ARANGODB_HOST"))
	}
	user := "root"
	pass := "root"

	conn, err := vst.NewConnection(vst.ConnectionConfig{
		Endpoints: []string{host},
	})
	if err != nil {
		log.Println("Failed to connect in Arango.", err)
	}
	cli, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(user, pass),
	})
	if err != nil {
		log.Println("Failed to authenticate in Arango.", err)
	}

	return cli, err
}

func getHttpServer() *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.BodyLimit("5M"))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	e.HideBanner = true
	e.Debug = true
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	return e
}

func getSkipAndLimit(c echo.Context) (int, int) {
	page := c.QueryParams().Get("page")
	pageNumber, _ := strconv.Atoi(page)
	if pageNumber > 0 {
		pageNumber = pageNumber - 1
	}
	perPage := c.QueryParams().Get("per_page")
	perPageNumber, _ := strconv.Atoi(perPage)
	if perPageNumber < 1 {
		perPageNumber = 5
	}
	skip := pageNumber * perPageNumber
	return skip, perPageNumber
}

func main() {
	ctx := driver.WithQueryFullCount(context.Background())
	client, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	db, err := client.Database(ctx, "elearning")
	if err != nil {
		log.Fatal("Failed to get database:", err)
	}

	server := getHttpServer()

	// Get students with playlists, audios and teachers
	server.GET("/students", func(c echo.Context) error {
		skip, limit := getSkipAndLimit(c)
		query := `FOR student IN students LIMIT @skip, @limit
            LET playlists = (
                FOR p IN OUTBOUND student listen
                    LET audios = (
                        FOR a IN OUTBOUND p plays
                            let teachers = (FOR teacher IN INBOUND a records RETURN { id: teacher._key, name: teacher.name })
                        RETURN { id: a._key, title: a.title, media: a.media, teachers: teachers }
                    )
                    LET subjects = (FOR s IN OUTBOUND p tagged RETURN { id: s._key, name: s.name })
                RETURN { id: p._key, name: p.name, audios: audios, subjects: subjects }
            )
        RETURN { id: student._key, name: student.name, city: student.city, playlists: playlists }`
		params := map[string]interface{}{
			"skip":  skip,
			"limit": limit,
		}
		cursor, err := db.Query(ctx, query, params)
		if err != nil {
			log.Println("Failed to perform query", query, err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		defer cursor.Close()
		var students []Student
		var student Student
		for {
			_, err := cursor.ReadDocument(ctx, &student)
			if driver.IsNoMoreDocuments(err) {
				break
			}
			students = append(students, student)
		}
		stats := cursor.Statistics()
		total := stats.FullCount()
		res := map[string]interface{}{
			"data":  students,
			"total": total,
		}
		return c.JSON(http.StatusOK, res)
	})

	// Get teachers with his audios, playlists and students
	server.GET("/teachers", func(c echo.Context) error {
		skip, limit := getSkipAndLimit(c)
		query := `FOR t IN teachers LIMIT @skip, @limit
            LET audios = (
                FOR a IN OUTBOUND t records
                    LET playlists = (
                        FOR p IN INBOUND a plays
                            LET students = (FOR s IN INBOUND p listen RETURN { id: s._key, name: s.name, city: s.city })
                            LET subjects = (FOR s IN OUTBOUND p tagged RETURN { id: s._key, name: s.name })
                        RETURN { id: p._key, name: p.name, students: students, subjects: subjects }
                    )
                RETURN { id: a._key, title: a.title, media: a.media, playlists: playlists }
            )
        RETURN { id: t._key, name: t.name, audios: audios }`
		params := map[string]interface{}{
			"skip":  skip,
			"limit": limit,
		}
		cursor, err := db.Query(ctx, query, params)
		if err != nil {
			log.Println("Failed to perform query", query, err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		defer cursor.Close()
		var teachers []Teacher
		var teacher Teacher
		for {
			_, err := cursor.ReadDocument(ctx, &teacher)
			if driver.IsNoMoreDocuments(err) {
				break
			}
			teachers = append(teachers, teacher)
		}
		stats := cursor.Statistics()
		total := stats.FullCount()
		res := map[string]interface{}{
			"data":  teachers,
			"total": total,
		}
		return c.JSON(http.StatusOK, res)
	})

	// Get playlists by subject
	server.GET("/playlists/:subject", func(c echo.Context) error {
		skip, limit := getSkipAndLimit(c)
		query := `FOR s IN subjects FILTER LIKE (s.name, @subject, true)
            FOR p IN INBOUND s tagged LIMIT @skip, @limit
                LET audios = (FOR a IN OUTBOUND p plays RETURN { id: a._key, title: a.title, media: a.media })
                LET subjects = (FOR su IN OUTBOUND p tagged RETURN { id: su._key, name: su.name })
            RETURN { id: p._key, name: p.name, audios: audios, subjects: subjects }`
		params := map[string]interface{}{
			"skip":    skip,
			"limit":   limit,
			"subject": c.Param("subject"),
		}
		cursor, err := db.Query(ctx, query, params)
		if err != nil {
			log.Println("Failed to perform query", query, err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		defer cursor.Close()
		var playlists []Playlist
		var playlist Playlist
		for cursor.HasMore() {
			_, err := cursor.ReadDocument(ctx, &playlist)
			if err == nil {
				playlists = append(playlists, playlist)
			}
		}
		stats := cursor.Statistics()
		total := stats.FullCount()
		res := map[string]interface{}{
			"data":  playlists,
			"total": total,
		}
		return c.JSON(http.StatusOK, res)
	})

	log.Fatal(server.Start(":5001"))
}
