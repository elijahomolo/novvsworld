package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"net/http"
	//"net/smtp"
	"os"
	"reflect"
	"text/template"
	_ "github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

// define a user model
type User struct {
	Id    int
	Email string
	SocialNetwork string
	Handle string
}

type Artist struct {
	Id int
	FirstName string
	LastName string
	Email string
	TwitterUserName string
	InstagramUserName string
	artLink string
}
// load .env file
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

// connect to the database and return it as an object
func dbConn() (db *sql.DB) {
	// pass the db credentials into variables
	host := goDotEnvVariable("DBHOST")
	port := goDotEnvVariable("DBPORT")
	dbUser := goDotEnvVariable("DBUSER")
	dbPass := goDotEnvVariable("DBPASS")
	dbname := goDotEnvVariable("DBNAME")
	sslmode := goDotEnvVariable("SSLMODE")
	caCert := "sslrootcert = " + goDotEnvVariable("CA_CERT")
	// create a connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s cacert=%s",
		 host, port, dbUser, dbPass, dbname, sslmode, caCert)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Index", nil)
}

//func contestSignUp(w http.ResponseWriter, r *http.Request) {
//	tmpl.ExecuteTemplate(w, "ContestSignUp", nil)
//}
//


func sendThanks(users_fn string, user_adr string) {
	from := mail.NewEmail("Novv", "novv@novvs.world") // Change to your verified sender
	subject := "Welcome to Novvs World"
	to := mail.NewEmail(users_fn, user_adr) // Change to your recipient
	plainTextContent := "Hey" + users_fn + "! Thanks for signing up! I'll be reaching out with an update on your status shortly! Let's see what you got and good luck!!"
	htmlContent := "<strong> Hey " + users_fn + "!! Thanks for signing up! I'll be reaching out with an update on your status shortly! Let's see what you got and good luck!!" + "\n Novv </strong> + "

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(goDotEnvVariable("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Headers)
	}

}

func thanks(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Thanks", nil)
}

var tmpl = template.Must(template.ParseGlob("assets/templates/*"))


func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	//If it's a post request, assign a variable to the value returned in each field of the New page.
	if r.Method == "POST" {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		twitterUsername := r.FormValue("twitter_username")
		instagramUsername := r.FormValue("instagram_username")
		artLink := r.FormValue("art_link")
		createdOn := time.Now().UTC()

		log.Println("INSERT: firstName: " + firstName + " | lastName: " + lastName + " | email : " + email + " | created on: " + createdOn.String() + " | twitter_username: " + twitterUsername + " | instagramUsername: " + instagramUsername )


		//prepare a query to insert the data into the database
		insForm, err := db.Prepare(`INSERT INTO public.art_contest(first_name, last_name, email, instagram_username, twitter_username, art_link, created_on) VALUES ($1,$2, $3, $4, $5, $6, $7)`)
		//check for  and handle any errors
		CheckError(err)
		//execute the query using the form data
		_, err = insForm.Exec(firstName, lastName, email, instagramUsername, twitterUsername, artLink, createdOn)
		CheckError(err)
		//print out added data in terminal
		log.Println("INSERT: email: " + email + " | social network: " + instagramUsername + " | social handle : " + twitterUsername + " | created on: " + createdOn.String() + " | createdOn is type: " + reflect.TypeOf(createdOn).String())
		sendThanks(firstName, email)
	}
	defer db.Close()

	//redirect to the index page
	http.Redirect(w, r, "/thanks", 301)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}


func main() {
	fileServer := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	//	http.Handle("/assets/static", http.StripPrefix("/assets/static", http.FileServer(http.Dir("assets/static"))))
	http.HandleFunc("/", index)
	//http.HandleFunc("/contest-sign-up", contestSignUp)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/thanks", thanks)
	fmt.Println("server starting on port 3000...")
	http.ListenAndServe(":3000", nil)


}

