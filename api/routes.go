package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "project/docs"
	"project/internal/handler"
	"project/internal/repository"
	"project/internal/service"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	repo := repository.NewHierarchyRepository(db)
	svc := service.NewHierarchyService(repo)
	h := handler.NewHierarchyHandler(svc)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.POST("/hierarchy/upload", h.UploadExcel)
		api.POST("/hierarchy", h.CreateSingleNode)
		api.GET("/hierarchy", h.GetAll)
		api.GET("/hierarchy/download", h.DownloadExcel)
		api.GET("/hierarchy/full", h.GetFullHierarchy)
		api.GET("/hierarchy/:id", h.GetHierarchyByNodeID)
		api.GET("/hierarchy/:id/subtree", h.GetSubTree)
	}
}
