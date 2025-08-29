package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/auth"
	"github.com/oktaharis/uji-teknis-godigi/internal/config"
	"github.com/oktaharis/uji-teknis-godigi/internal/handlers"
	"github.com/oktaharis/uji-teknis-godigi/internal/middleware"
	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()

	// Add middlewares
	r.Use(middleware.JSONRecovery())
	r.NoRoute(middleware.NotFoundHandler())

	// health
	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"}, "Health check")
	})

	ah := handlers.NewAuthHandler(cfg, db)
	uh := handlers.NewUserHandler()
	lh := handlers.NewLeadHandler(db)
	ph := handlers.NewProjectHandler(db)
	uah := handlers.NewUserAdminHandler(db)

	// Public
	pub := r.Group("/auth")
	{
		pub.POST("/register", ah.Register)
		pub.POST("/login", ah.Login)
		pub.POST("/forgot-password", ah.ForgotPassword)
		pub.POST("/reset-password", ah.ResetPassword)
	}

	// Protected
	api := r.Group("/")
	api.Use(auth.AuthRequired(cfg, db))
	{
		api.POST("/auth/logout", ah.Logout)
		api.GET("/me", uh.Me)

		api.POST("/leads", lh.Create)
		api.GET("/leads", lh.List)
		api.GET("/leads/summary", lh.Summary)
		api.GET("/leads/:id", lh.Get)
		api.PUT("/leads/:id", lh.Update)
		api.DELETE("/leads/:id", lh.Delete)
	}
	
	api.POST("/projects", ph.Create)
    api.GET("/projects", ph.List)
    api.GET("/projects/:id", ph.Get)
    api.PUT("/projects/:id", ph.Update)
    api.DELETE("/projects/:id", ph.Delete)

	admin := r.Group("/admin")
	admin.Use(auth.AdminOnly())
	{
		admin.POST("/users", uah.Create)
		admin.GET("/users", uah.List)
		admin.GET("/users/:id", uah.Get)
		admin.PUT("/users/:id", uah.Update)
		admin.DELETE("/users/:id", uah.Delete)
	}

	return r
}
