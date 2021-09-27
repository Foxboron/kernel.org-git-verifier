package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"gorm.io/gorm/clause"
)

func GetLinuxRepo(earliest, latest string) {
	c := fmt.Sprintf("-C ./linux log --pretty=format:%%H %s..%s", earliest, latest)
	cmd := exec.Command("git", strings.Split(c, " ")...)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	commitFromLinux := strings.Split(string(stdout), "\n")
	exists := 0
	validsig := 0
	pushes := 0
	checkPush := map[string]bool{}
	for _, commit := range commitFromLinux {
		count := int64(0)
		var rev Revision
		err := db.Table("revisions").Where("\"commit\" = ?", commit).Count(&count).Error
		db.Preload(clause.Associations).Where("\"commit\" = ?", commit).Find(&rev)
		if _, ok := checkPush[rev.TLogCommit.Commit]; !ok {
			checkPush[rev.TLogCommit.Commit] = true
			pushes++
		}
		if rev.TLogCommit.Valid {
			validsig++
		}
		if err != nil {
			log.Fatal(err)
		}
		if count >= 1 {
			exists++
		} else {
			fmt.Println(commit)
		}
	}
	fmt.Printf("Out of %d commits, %d was found on the tlog between %s and %s\n", len(commitFromLinux), exists, earliest, latest)
	fmt.Printf("%d commits where done over %d pushes, %d push signatures where valid\n", len(commitFromLinux), pushes, validsig)
}
