package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/saif404/go-postgres/models"
	"github.com/saif404/go-postgres/storage"
	"gorm.io/gorm"
)
type Book struct{
	Author		string 	`json:"author"`
	Title		string 	`json:"title"`
	Publisher	string 	`json:"publisher"`	
}

type Repo struct{
	DB *gorm.DB
}
func main(){
	err:=godotenv.Load(".env")
	if err!=nil{
		log.Fatal(err)
	}
	config:=&storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password : os.Getenv("DB_PASSWORD"),
		SSLMode: os.Getenv("DB_SSL"),
		DBName: os.Getenv("DB_NAME"),
	}
	db,err:=storage.NewConnection(config)
	if err!=nil{
		log.Fatal("could not load the DB!")
	}

	r:=&Repo{
		DB: db,
	}

	errr:= models.MigrateBooks(db)
	if errr!=nil{
		log.Fatal("could not migrate to db!")
	}
	app:= fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}

func (r *Repo)SetupRoutes(app *fiber.App){
	api:= app.Group("/api")
	api.Get("/book/:id",r.GetBook)
	api.Get("/books",r.GetBooks)
	api.Post("/book/create",r.CreateBook)
	api.Delete("/book/:id",r.DeleteBook)
	// api.Put("/book/:àid/update",r.UpdateBook)

}
func (r *Repo)CreateBook(ctx *fiber.Ctx)error{
	book :=&Book{}
	if err:=ctx.BodyParser(&book);err!=nil{
		ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message":"request failed",
		})
	}

	if err:=r.DB.Create(&book).Error;err!=nil{
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"could not create a book",
			},
		)
		return err
	}
	return ctx.Status(http.StatusOK).JSON(book)

}

func (r *Repo)GetBooks(ctx *fiber.Ctx)error{
	bookModels:= &[]models.Books{}

	if err:=r.DB.Find(bookModels).Error;err!=nil{
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"could not get all the books",
			},
		)
		return err
	}
	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":"books fetched successfully",
			"data":bookModels,
		},
	)
	return nil
}

func (r *Repo)GetBook(ctx *fiber.Ctx)error{
	Book:=&models.Books{}
	id:=ctx.Params("id")
	err:=r.DB.Where("id = ?",id).First(Book).Error
	if err!=nil{
		ctx.Status(http.StatusBadRequest).JSON("could not get this book")
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"data":Book,
	})
	return nil
}

// func (r *Repo)UpdateBook(ctx *fiber.Ctx)error{

// }

func (r *Repo)DeleteBook(ctx *fiber.Ctx)error{
	bookModel:= &models.Books{}
	id,err:= ctx.ParamsInt("id")
	if err!=nil{
		ctx.Status(http.StatusBadGateway).JSON("invalid id!")
	}
	errr:= r.DB.Delete(bookModel,id)
	if errr.Error!=nil{
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message":"could not delete this book",
		})
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"the book has been deleted!",
		
	})
	return nil
}