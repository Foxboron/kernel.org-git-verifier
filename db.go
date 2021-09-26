package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

// These are the commits from the transparency logs
type TLogCommit struct {
	Commit     string `gorm:"primaryKey"`
	CommitMsg  string
	CommitDate time.Time
	User       string
	SigIssuer  string
	Valid      bool       `gorm:"default:false"`
	Signature  bool       `gorm:"default:false"`
	Unknown    bool       `gorm:"default:false"`
	Revisions  []Revision `gorm:"foreignKey:TLogCommitID"`
}

type TLogCommits []TLogCommit

func AddCommit(c *TLogCommit) {
	db.Create(c)
}

func AddCommits(c TLogCommits) {
	db.Create(&c)
}

type Revision struct {
	Commit       string `gorm:"primaryKey"`
	Who          string
	Repository   string
	TLogCommitID string
	TLogCommit   TLogCommit
}

func AddRevision(c *TLogCommit) {
	db.Create(c)
}

// :D that name
func AddRevisions(c TLogCommits) {
	db.Create(&c)
}

type Scan struct {
	Time    time.Time
	Commits int
}

func AddScan(c *Scan) {
	db.Create(c)
}

func LastTimestamp() (time.Time, error) {
	var s Scan
	if err := db.Preload(clause.Associations).Last(&s).Error; err != nil {
		return time.Time{}, err
	}
	return s.Time, nil
}

func InitDB(name string) {
	var err error
	db, err = gorm.Open(sqlite.Open(name), &gorm.Config{})
	db.AutoMigrate(&TLogCommit{}, &Revision{}, &Scan{})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	InitDB("kernel.db")
}
