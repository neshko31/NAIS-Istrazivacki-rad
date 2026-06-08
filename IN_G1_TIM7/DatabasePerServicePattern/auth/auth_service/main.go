package main

import (
	"net/http"
	"log"
	"auth-service/config"
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/jwt"
	"auth-service/saga"
	"auth-service/choreography"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Postgres database
	db, err := database.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.AlignUsersSequence(); err != nil {
		log.Fatal("Failed to align users sequence:", err)
	}
	defer db.DB.Close()

	// RabbitMQ
	rmqURL := cfg.GetRabbitMQURL()

	// Orkestrirana SAGA
	orchestrationRmq, err := saga.NewRabbitMQConnection(rmqURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ for orchestration:", err)
	}
	defer orchestrationRmq.Close()

	orchestrator := saga.NewRegistrationOrchestrator(db.DB, orchestrationRmq)

	// Koreografisana SAGA
	choreographyRmq, err := choreography.NewRabbitMQConnection(rmqURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ for choreography:", err)
	}
	defer choreographyRmq.Close()

	choreographyHandler := choreography.NewChoreographyHandler(db.DB, choreographyRmq)

	// Gin router
	if cfg.Enviroment == "production"{
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	})

	// HTTP Handlers
	authHandler := handlers.NewAuthHandler(db, []byte(cfg.JWT.Secret), orchestrator)
	choreographyAuthHandler := handlers.NewChoreographyAuthHandler(db, choreographyHandler)

	// Routes
	public := r.Group("/api")
	{
		// Orkestrirana SAGA
		public.POST("/register", authHandler.Register)
		public.GET("/register/status/:sagaId", authHandler.RegisterStatus)

		// Koreografisana SAGA
		public.POST("/choreography/register", choreographyAuthHandler.RegisterChoreography)
		public.GET("/choreography/register/status/:sagaId", choreographyAuthHandler.RegisterChoreographyStatus)

		// Auth endpoint
		public.POST("/login", authHandler.Login)
	}

	protected := r.Group("/api")
	protected.Use(jwt.AuthMiddleware([]byte(cfg.JWT.Secret)))
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/context", getUserContext)
	}

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Started server on %s", serverAddr)
	log.Println("Orkestrisana SAGA:    POST /api/register")
	log.Println("Koreografisana SAGA:  POST /api/choreography/register")

	srv := &http.Server{
		Addr:	serverAddr,
		Handler:	r,
	}
	go registerWithEureka()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
		log.Fatal("Server failed to start:", err)
	}
}

func getUserContext(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(200, gin.H{
		"user_id": userID,
		"username": username,
	})
}
