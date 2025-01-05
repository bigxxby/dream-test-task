package shortener

import (
	"github.com/bigxxby/dream-test-task/internal/api/service/shortener"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IShortenerController interface {
	CreateShortLink(ctx *gin.Context)
	Redirect(ctx *gin.Context)
	GetLinks(ctx *gin.Context)
	GetLink(ctx *gin.Context)
	DeleteLink(ctx *gin.Context)
}

func NewShortenerController(shortenerService shortener.IShortenerService) IShortenerController {
	return &ShortenerController{ShortenerService: shortenerService}
}

type ShortenerController struct {
	ShortenerService shortener.IShortenerService
}

func (sc *ShortenerController) DeleteLink(ctx *gin.Context) {
	shortID := ctx.Param("shortID")
	if shortID == "" {
		ctx.JSON(400, gin.H{
			"error":   "ShortID is empty",
			"message": "Bad request",
			"success": false,
		})
		return
	}

	status, err := sc.ShortenerService.DeleteLink(shortID)
	if err != nil {
		switch status {
		case 404:
			ctx.JSON(404, gin.H{
				"error":   err.Error(),
				"message": "Not found",
				"success": false,
			})
			return
		case 500:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"message": "Link deleted",
		"success": true,
	})
}
func (sc *ShortenerController) GetLink(ctx *gin.Context) {
	shortID := ctx.Param("shortID")
	if shortID == "" {
		ctx.JSON(400, gin.H{
			"error":   "ShortID is empty",
			"message": "Bad request",
			"success": false,
		})
		return
	}

	link, status, err := sc.ShortenerService.GetLink(shortID)
	if err != nil {
		switch status {
		case 404:
			ctx.JSON(404, gin.H{
				"error":   err.Error(),
				"message": "Not found",
				"success": false,
			})
			return
		case 500:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"link":    link,
		"message": "Link found",
		"success": true,
	})
}

func (sc *ShortenerController) GetLinks(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(string)
	if userId == "" {
		ctx.JSON(401, gin.H{
			"error":   "Unauthorized",
			"message": "Unauthorized",
			"success": false,
		})
		return
	}
	userIDUUID, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(401, gin.H{
			"error":   "Unauthorized",
			"message": "Unauthorized",
			"success": false,
		})
		return
	}

	links, status, err := sc.ShortenerService.GetLinks(&userIDUUID)
	if err != nil {
		switch status {
		case 500:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"links":   links,
		"message": "Links found",
		"success": true,
	})
}
func (sc *ShortenerController) CreateShortLink(ctx *gin.Context) {
	type createShortLinkRequest struct {
		Url string `json:"url"`
	}
	userId := ctx.MustGet("user_id").(string)
	if userId == "" {
		ctx.JSON(401, gin.H{
			"error":   "Unauthorized",
			"message": "Unauthorized",
			"success": false,
		})
		return
	}
	userIDUUID, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(401, gin.H{
			"error":   "Unauthorized",
			"message": "Unauthorized",
			"success": false,
		})
		return
	}

	var req createShortLinkRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Bad request",
			"success": false,
		})
		return
	}
	if req.Url == "" {
		ctx.JSON(400, gin.H{
			"error":   "Link is empty",
			"message": "Bad request",
			"success": false,
		})
		return
	}

	link, status, err := sc.ShortenerService.CreateShortLink(req.Url, &userIDUUID)
	if err != nil {
		switch status {
		case 400:
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "Bad request",
				"success": false,
			})
			return
		case 500:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"short_link": link,
		"message":    "Short link created",
		"success":    true,
	})

}
func (sc *ShortenerController) Redirect(ctx *gin.Context) {
	shortID := ctx.Param("shortID")
	if shortID == "" {
		ctx.JSON(400, gin.H{
			"error":   "ShortID is empty",
			"message": "Bad request",
			"success": false,
		})
		return
	}

	link, status, err := sc.ShortenerService.Redirect(shortID)
	if err != nil {
		switch status {
		case 404:
			ctx.JSON(404, gin.H{
				"error":   err.Error(),
				"message": "Not found",
				"success": false,
			})
			return
		case 500:
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"message": "Internal server error",
				"success": false,
			})
			return
		}
	}

	ctx.Redirect(301, link)
}
