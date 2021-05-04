package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

var Urls = map[string]string{
	"/":         "Index",
	"/signed":   "Valid Signatures",
	"/unsigned": "Unsigned",
	"/unknown":  "Unknown Issuer",
	"/invalid":  "Invalid Signatures",
}

func RenderTmpl(c *gin.Context, entries *[]TLogCommit) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"entries": entries,
		"urls":    Urls,
		"active":  c.FullPath(),
	})
}

func RenderIndex(c *gin.Context) {
	var entry []TLogCommit
	db.Limit(200).Find(&entry)
	RenderTmpl(c, &entry)
}

func RenderInvalid(c *gin.Context) {
	var entry []TLogCommit
	db.Limit(200).Where("\"unknown\" = ? and \"signature\" = ? and \"valid\" = ?", false, true, false).Find(&entry)
	RenderTmpl(c, &entry)
}

func RenderSigned(c *gin.Context) {
	var entry []TLogCommit
	db.Limit(200).Where("\"valid\" = ?", true).Find(&entry)
	RenderTmpl(c, &entry)
}

func RenderUnsigned(c *gin.Context) {
	var entry []TLogCommit
	db.Limit(200).Where("\"signature\" = ?", false).Find(&entry)
	RenderTmpl(c, &entry)
}

func RenderUnknown(c *gin.Context) {
	var entry []TLogCommit
	db.Where("\"unknown\" = ?", true).Find(&entry)
	RenderTmpl(c, &entry)
}

func GetRevision(c *gin.Context) {
	var rev Revision
	commit := c.Param("commit")
	db.Preload(clause.Associations).Where("\"commit\" = ?", commit).Find(&rev)
	c.IndentedJSON(200, &rev)
}

func webpage(ctx context.Context) {
	router := gin.Default()
	router.Static("/css", "web/css")
	router.LoadHTMLGlob("web/templates/*")
	router.GET("/", RenderIndex)
	router.GET("/signed", RenderSigned)
	router.GET("/unsigned", RenderUnsigned)
	router.GET("/unknown", RenderUnknown)
	router.GET("/invalid", RenderInvalid)
	router.GET("/api/revision/:commit", GetRevision)

	srv := &http.Server{Addr: "localhost:8088", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
		log.Println("Server exiting")
	}
}
