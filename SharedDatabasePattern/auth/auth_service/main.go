package main

import (
	"net/http"
	"log"
	"auth-service/config"
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/jwt"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.DB.Close()

	if cfg.Enviroment == "production"{
		gin.SetMode(gin.ReleaseMode)
	}

	//Router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	//CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		c.Next()
	})

	authHandler := handlers.NewAuthHandler(db, []byte(cfg.JWT.Secret))

	//Routes
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
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

	srv := &http.Server{
		Addr:	serverAddr,
		Handler:	r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
		log.Fatal("Server failed to start:", err)
	}
}

func viewUsers(){

}

func getUserContext(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(200, gin.H{
		"user_id": userID,
		"username": username,
	})
}
