package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgutz/dat.v2/dat"
)

type Post struct {
	ID        int64        `db:"id"`
	Title     string       `db:"title"`
	Body      string       `db:"body"`
	UserID    int64        `db:"user_id"`
	State     string       `db:"state"`
	UpdatedAt dat.NullTime `db:"updated_at"`
	CreatedAt time.Time    `db:"created_at"`
}

//dd:span
func main() {
	var post Post
	err := DB.Select("id, title").From("posts").Where("id = $1", 1).QueryStruct(&post)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", post)
}
