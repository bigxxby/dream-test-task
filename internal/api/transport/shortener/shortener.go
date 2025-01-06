package shortener

import (
	"github.com/bigxxby/dream-test-task/internal/api/service/shortener"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Ответы на успешные запросы
type SuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Ответы на ошибки
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Структура запроса для создания короткой ссылки
type CreateShortLinkRequest struct {
	Url string `json:"url" binding:"required"`
}

// Ответ для создания короткой ссылки
type CreateShortLinkResponse struct {
	ShortLink string `json:"short_link"`
	Message   string `json:"message"`
	Success   bool   `json:"success"`
}

// Ответ для получения ссылки
type GetLinkResponse struct {
	Link    string `json:"link"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Ответ для списка ссылок
type GetLinksResponse struct {
	Links   []string `json:"links"`
	Message string   `json:"message"`
	Success bool     `json:"success"`
}

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

// DeleteLink godoc
//	@Summary		Delete a shortened link
//	@Description	Deletes a shortened link by its shortID
//	@Tags			Shortener
//	@Param			shortID	path	string	true	"Shortened Link ID"
//	@Security		BearerAuth
//	@Success		200	{object}	SuccessResponse	"Link deleted successfully"
//	@Failure		400	{object}	ErrorResponse	"ShortID is empty or invalid"
//	@Failure		404	{object}	ErrorResponse	"Shortened link not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/shortener/{shortID} [delete]
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

// GetLink godoc
//	@Summary		Get original link stats from short URL
//	@Description	Retrieves the original URL statistics based on the provided shortened link ID.
//	@Tags			Shortener
//	@Param			shortID	path	string	true	"Shortened Link ID"
//	@Security		BearerAuth
//	@Success		200	{object}	GetLinkResponse	"Link stats retrieved successfully"
//	@Failure		400	{object}	ErrorResponse	"Invalid ShortID"
//	@Failure		404	{object}	ErrorResponse	"Link not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/shortener/stats/{shortID} [get]
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

// GetLinks godoc
//	@Summary		Get all shortened links for a user
//	@Description	Retrieves all the shortened links associated with the authenticated user.
//	@Tags			Shortener
//	@Security		BearerAuth
//	@Success		200	{object}	GetLinksResponse	"Links retrieved successfully"
//	@Failure		401	{object}	ErrorResponse		"Unauthorized"
//	@Failure		500	{object}	ErrorResponse		"Internal server error"
//	@Router			/shortener [get]
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

// CreateShortLink godoc
//	@Summary		Create a shortened link
//	@Description	Creates a new shortened link from the provided URL.
//	@Tags			Shortener
//	@Param			request	body	CreateShortLinkRequest	true	"Request body for creating short link"
//	@Security		BearerAuth
//	@Success		200	{object}	CreateShortLinkResponse	"Link created successfully"
//	@Failure		400	{object}	ErrorResponse			"Invalid URL or missing parameters"
//	@Failure		401	{object}	ErrorResponse			"Unauthorized"
//	@Failure		500	{object}	ErrorResponse			"Internal server error"
//	@Router			/shortener [post]
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

// Redirect godoc
//	@Summary		Redirect to the original URL
//	@Description	Redirects the user to the original URL from a shortened link ID.
//	@Tags			Shortener
//	@Param			shortID	path	string	true	"Shortened Link ID"
//	@Security		BearerAuth
//	@Success		302	{string}	"Redirected to the original URL"
//	@Failure		400	{object}	ErrorResponse	"ShortID is empty or invalid"
//	@Failure		404	{object}	ErrorResponse	"Link not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/shortener/{shortID} [get]
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
