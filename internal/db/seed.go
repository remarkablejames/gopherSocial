package db

import (
	"context"
	"database/sql"
	"fmt"
	"gopherSocial/internal/store"
	"log"
	"math/rand"
)

var usernames = []string{
	"JamesSmith", "EmmaJohnson", "LiamBrown", "OliviaJones", "NoahGarcia",
	"AvaMartinez", "WilliamRodriguez", "SophiaHernandez", "BenjaminLopez",
	"IsabellaWilson", "LucasAnderson", "MiaThomas", "HenryTaylor", "AmeliaMoore",
	"AlexanderJackson", "CharlotteWhite", "MichaelHarris", "HarperMartin",
	"EthanThompson", "EvelynMartinez", "DanielClark", "AbigailLewis",
	"MatthewWalker", "EmilyYoung", "DavidHall", "ElizabethAllen", "JosephKing",
	"SofiaScott", "SamuelWright", "VictoriaTorres", "JacksonHill", "GraceGreen",
	"SebastianAdams", "ChloeBaker", "AnthonyNelson", "EllaCarter", "DylanMitchell",
	"AubreyPerez", "LeoRoberts", "ScarlettTurner", "JackPhillips", "ZoeyCampbell",
	"AidenParker", "LillianEvans", "OwenEdwards", "HannahCollins", "GabrielStewart",
	"AriaFlores", "CarterMorris", "PenelopeNguyen",
}

var titles = []string{
	"Getting Started with Go: A Beginner's Guide",
	"Top 10 Features of the Go Programming Language",
	"Understanding Goroutines and Concurrency in Go",
	"Building Web Applications with Go and Gin",
	"Error Handling in Go: Best Practices",
	"Go Modules Explained: Dependency Management Simplified",
	"Mastering Interfaces in Go: A Practical Guide",
	"Building REST APIs with Go: A Step-by-Step Tutorial",
	"Exploring the Power of Go's Standard Library",
	"Unit Testing in Go: A Complete Guide",
	"Understanding Slices and Arrays in Go",
	"Channels in Go: Communicating Between Goroutines",
	"How to Use the `context` Package in Go Effectively",
	"Writing High-Performance Code with Go",
	"Go vs Other Languages: Why Choose Go?",
	"Creating CLI Applications with Cobra in Go",
	"Memory Management and Garbage Collection in Go",
	"Working with JSON in Go: Tips and Tricks",
	"Structs and Methods in Go: The Basics",
	"Concurrency Patterns in Go: Pipelines and Fan-Out",
	"Building Real-Time Applications in Go",
	"Using SQL Databases in Go with GORM",
	"Best Practices for Logging in Go Applications",
	"Deploying Go Applications to the Cloud",
	"How to Use the `sync` Package in Go",
	"Creating Custom Middleware with Go",
	"Go vs Rust: Comparing Two Modern Languages",
	"Introduction to Generics in Go",
	"Profiling and Debugging Go Applications",
	"Advanced Goroutines: Managing Pooling and Throttling",
	"Exploring Go's Net/HTTP Package for Web Development",
	"Dependency Injection in Go: A Practical Approach",
	"Implementing Microservices Architecture in Go",
	"How to Optimize Go Code for Better Performance",
	"Using Reflection in Go: What You Need to Know",
	"Understanding Go's `defer`, `panic`, and `recover`",
	"Event-Driven Programming with Go",
	"Building GraphQL APIs with Go",
	"How to Write Secure Applications in Go",
	"Implementing Design Patterns in Go",
	"Using Go for Machine Learning and Data Science",
	"File Handling and Processing in Go",
	"Parsing Command-Line Flags with Go's `flag` Package",
	"Creating WebSocket Applications in Go",
	"Exploring Go's Time and Date Handling",
	"Understanding and Using Maps in Go",
	"Developing Cross-Platform Applications with Go",
	"How to Work with Third-Party Packages in Go",
	"Setting Up a Go Development Environment",
	"Tips for Writing Clean and Idiomatic Go Code",
}

var tags = []string{
	"Go", "Golang", "Programming", "Web Development", "Concurrency",
	"REST API", "Microservices", "Testing", "Error Handling", "CLI Tools",
	"Performance", "Generics", "Dependency Injection", "JSON", "Networking",
	"Cloud Deployment", "Security", "Debugging", "Best Practices", "Tutorial",
}

var comments = []string{
	"Great article! This really helped me understand Go better.",
	"I’ve been struggling with concurrency in Go—this post clarified so much. Thanks!",
	"Can you write more about building APIs with Go? This was super helpful.",
	"Nice breakdown of generics! I was waiting for someone to explain it like this.",
	"Awesome content! Keep up the good work.",
	"This is exactly what I needed to get started with Go modules.",
	"Thanks for the detailed examples. It made the concepts easy to follow.",
	"Is there a GitHub repo where I can find the code examples?",
	"Could you dive deeper into performance optimization in Go? This is a good start.",
	"This post on goroutines is pure gold. I feel much more confident now.",
	"How does this compare to other languages like Python or Rust?",
	"Great overview! I’d love to see more posts on using Go in production.",
	"Your explanation of the `defer` keyword was spot on. Thanks!",
	"This saved me hours of debugging. Excellent write-up!",
	"Do you have a recommendation for the best Go frameworks for web development?",
	"Really insightful! The section on error handling was especially useful.",
	"Any tips for using Go with Docker? A follow-up post would be great.",
	"I didn’t know Go could handle WebSockets this well. Great content!",
	"Looking forward to more tutorials like this. Subscribed!",
	"Thanks for sharing your expertise. This was very informative.",
}

func Seed(store store.Storage, db *sql.DB) {

	ctx := context.Background()

	users := generateUsers(100)

	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Printf("error creating user: %v", err)
		}
	}

	tx.Commit()

	posts := generatePosts(100, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("error creating post: %v", err)
		}
	}

	comments := generateComments(100, posts, users)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("error creating comment: %v", err)
		}
	}

	log.Println("Seeding completed ☑️")

}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)

	for i := 0; i < n; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)

	for i := 0; i < n; i++ {
		posts[i] = &store.Post{
			Title:   titles[rand.Intn(len(titles))],
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
			UserID: users[i%len(users)].ID,
		}
	}
	return posts
}

func generateComments(n int, posts []*store.Post, users []*store.User) []*store.Comment {
	commentList := make([]*store.Comment, n)

	for i := 0; i < n; i++ {
		commentList[i] = &store.Comment{
			Content: comments[rand.Intn(len(comments))],
			UserID:  users[i%len(users)].ID,
			PostID:  posts[i%len(posts)].ID,
		}
	}
	return commentList
}
