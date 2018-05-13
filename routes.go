package main

func initializeRoutes() {

	// Handle the index route
	router.GET("/", showHomePage)
	// Handle the login route
	router.GET("/login", showLoginPage)
	// Handle the login route
	router.GET("/signup", showSignupPage)
	// Get all todos
	router.GET("/todos", getAllTodos)
	// Create a todo
	router.POST("/todos", createATodo)
	// Update a todo
	router.PUT("/todos/:todo_id", updateATodo)
	// Delete a todo
	router.DELETE("/todos/:todo_id", deleteATodo)
	// Get a todo by ID
	router.GET("/todos/:todo_id", getATodo)

	//Method specific to form submitalls
	router.POST("/todos/:todo_id", updateOrDeleteTodo)
}
