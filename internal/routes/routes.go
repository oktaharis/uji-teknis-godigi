package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/auth"
	"github.com/oktaharis/uji-teknis-godigi/internal/config"
	"github.com/oktaharis/uji-teknis-godigi/internal/handlers"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger())
    // r.Use(middleware.RecoveryJSON(), middleware.NotFoundJSON()) // kalau kamu pakai

    // Public (tanpa auth)
    ah  := handlers.NewAuthHandler(cfg, db)
    uh  := handlers.NewUserHandler()
    lh  := handlers.NewLeadHandler(db)
    ph  := handlers.NewProjectHandler(db)
    uah := handlers.NewUserAdminHandler(db)

    pub := r.Group("/auth")
    {
        pub.POST("/register", ah.Register)
        pub.POST("/login", ah.Login)
        pub.POST("/forgot-password", ah.ForgotPassword)
        pub.POST("/reset-password", ah.ResetPassword)
    }

    // Protected (WAJIB AuthRequired agar `user` ada di context)
    api := r.Group("/")
    api.Use(auth.AuthRequired(cfg, db))
    {
        api.POST("/auth/logout", ah.Logout)
        api.GET("/me", uh.Me)

        // Leads
        api.POST("/leads", lh.Create)
        api.GET("/leads", lh.List)
        api.GET("/leads/summary", lh.Summary)
        api.GET("/leads/:id", lh.Get)
        api.PUT("/leads/:id", lh.Update)
        api.DELETE("/leads/:id", lh.Delete)

        // Projects
        api.POST("/projects", ph.Create)
        api.GET("/projects", ph.List)
        api.GET("/projects/:id", ph.Get)
        api.PUT("/projects/:id", ph.Update)
        api.DELETE("/projects/:id", ph.Delete)

        // ADMIN — HARUS di dalam `api` supaya AuthRequired jalan lebih dulu
        admin := api.Group("/admin")
        admin.Use(auth.AdminOnly())
        {
            admin.POST("/users", uah.Create)
            admin.GET("/users", uah.List)
            admin.GET("/users/:id", uah.Get)
            admin.PUT("/users/:id", uah.Update)
            admin.DELETE("/users/:id", uah.Delete)
        }

        // (opsional) endpoint debug
        api.GET("/debug/whoami", func(c *gin.Context) {
            u := c.MustGet("user").(models.User)
            response.OK(c, gin.H{"id": u.ID, "email": u.Email, "role": u.Role}, "whoami")
        })
    }

    // r.NoRoute(func(c *gin.Context){ response.NotFound(c, "route not found") })
    return r
}
