package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	atomicdocs "github.com/yourusername/atomicdocs/fiber"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Category string  `json:"category"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 999.99, Stock: 10, Category: "electronics"},
	{ID: 2, Name: "Mouse", Price: 25.50, Stock: 50, Category: "electronics"},
}

var nextUserID = 3
var nextProductID = 3

func main() {
	app := fiber.New()
	
	
	// Auth routes
	app.Post("/auth/login", login)
	app.Post("/auth/register", register)
	
	// User routes
	app.Get("/users", getUsers)
	app.Get("/users/:id", getUser)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)
	
	// Product routes
	app.Get("/products", getProducts)
	app.Get("/products/:id", getProduct)
	app.Post("/products", createProduct)
	app.Put("/products/:id", updateProduct)
	app.Delete("/products/:id", deleteProduct)
	
	// Stats
	app.Get("/stats/summary", getStats)
	
	port := 3000
	app.Use(atomicdocs.New(port))
	
	// Register routes AFTER all routes are defined
	atomicdocs.Register(app, port)
	
	fmt.Printf("Server: http://localhost:%d\n", port)
	app.Listen(fmt.Sprintf(":%d", port))
}

func login(c *fiber.Ctx) error {
	var body map[string]string
	c.BodyParser(&body)
	
	if body["username"] == "admin" && body["password"] == "secret" {
		return c.JSON(fiber.Map{
			"token": "fiber-jwt-token-xyz",
			"user":  fiber.Map{"username": "admin", "role": "admin"},
		})
	}
	return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
}

func register(c *fiber.Ctx) error {
	var body map[string]interface{}
	c.BodyParser(&body)
	
	user := User{
		ID:    nextUserID,
		Name:  body["name"].(string),
		Email: body["email"].(string),
		Age:   int(body["age"].(float64)),
	}
	nextUserID++
	users = append(users, user)
	
	return c.Status(201).JSON(fiber.Map{"message": "User registered", "userId": user.ID})
}

func getUsers(c *fiber.Ctx) error {
	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	for _, u := range users {
		if u.ID == id {
			return c.JSON(u)
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}

func createUser(c *fiber.Ctx) error {
	var user User
	c.BodyParser(&user)
	user.ID = nextUserID
	nextUserID++
	users = append(users, user)
	return c.Status(201).JSON(user)
}

func updateUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	for i, u := range users {
		if u.ID == id {
			c.BodyParser(&users[i])
			users[i].ID = id
			return c.JSON(users[i])
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}

func deleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			return c.JSON(fiber.Map{"message": "User deleted"})
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}

func getProducts(c *fiber.Ctx) error {
	return c.JSON(products)
}

func getProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	for _, p := range products {
		if p.ID == id {
			return c.JSON(p)
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
}

func createProduct(c *fiber.Ctx) error {
	var product Product
	c.BodyParser(&product)
	product.ID = nextProductID
	nextProductID++
	products = append(products, product)
	return c.Status(201).JSON(product)
}

func updateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	for i, p := range products {
		if p.ID == id {
			c.BodyParser(&products[i])
			products[i].ID = id
			return c.JSON(products[i])
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
}

func deleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return c.JSON(fiber.Map{"message": "Product deleted"})
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
}

func getStats(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"totalUsers":    len(users),
		"totalProducts": len(products),
	})
}
