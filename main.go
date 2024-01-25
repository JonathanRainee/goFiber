package main

import(
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	"net/http"
	"log"
	"os"
)

type Book struct{
	Author			string			`json:"author"`
	Title				string			`json:"title"`
	Publisher		string			`json:"publisher"`
}	

type Repository struct{
	DB *gorm.db
}

func (r * Repository) CreateBook(context *fiber.ctx) error{
	book := Book{}

	//ngeparse json tanpa json parser (pakai function bawaan fiber)
	err := context.BodyParser(&book)

	if err != nil{
		context.Status(http.StatisUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failes"})
			return err
	}

	err = r.DB.Create(&book).Error
	
	if err != nil{
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not create book"})
			return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"book has been added"
	})
	return nil

}

func(r *Repository) GetBooks(context *fiber.Ctx) error{
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil{
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"could not get books"})
			return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"books fetched successfully",
		"data":bookModels
	})
	return nil
}



func(r *Repository) SetupRotes(app *fiber.app){
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.get("/get_books/:id", r.GetBookByID)
	api.get("/Books", r.GetBooks)

}

func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User: os.Getenv("DB_USER"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE")
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load database")
	}

	r:= Repository(
		DB: db,
	)

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080");
}