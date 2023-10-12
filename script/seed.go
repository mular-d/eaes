package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

type Result struct {
	Name      string
	Image     string
	Chemistry int
	Biology   int
	Maths     int
	Civic     int
	English   int
	Aptitude  int
	Physics   int
}

var names = []string{"Abebe", "Kebede", "Ayele"}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/result?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100000; i++ {
		result := generateRandomResult()
		err := insertResult(db, result)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data inserted successfully!")
}

func generateRandomResult() Result {
	name := names[rand.Intn(len(names))]
	image := "https://eaes.com/image/image.jpg"
	chemistry := rand.Intn(100) + 1
	biology := rand.Intn(100) + 1
	maths := rand.Intn(100) + 1
	civic := rand.Intn(100) + 1
	english := rand.Intn(100) + 1
	aptitude := rand.Intn(100) + 1
	physics := rand.Intn(100) + 1

	return Result{
		Name:      name,
		Image:     image,
		Chemistry: chemistry,
		Biology:   biology,
		Maths:     maths,
		Civic:     civic,
		English:   english,
		Aptitude:  aptitude,
		Physics:   physics,
	}
}

func insertResult(db *sql.DB, result Result) error {
	_, err := db.Exec("INSERT INTO results (name, image, chemistry, biology, maths, civic, english, aptitude, physics) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		result.Name, result.Image, result.Chemistry, result.Biology, result.Maths, result.Civic, result.English, result.Aptitude, result.Physics)
	return err
}
