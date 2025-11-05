package handlers

import (
    "math"
    "strconv"

    response "bismillah/src/utils/response_helper"
    "github.com/gin-gonic/gin"
)

// GetAllUsers handler
func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
    // Ambil query params (default: page=1, limit=10)
    page := ctx.DefaultQuery("page", "1")
    limit := ctx.DefaultQuery("limit", "10")
    search := ctx.DefaultQuery("search", "")

    // Konversi ke int
    pageInt, err := strconv.Atoi(page)
    if err != nil || pageInt < 1 {
        response.BadRequest(ctx, "Invalid page number", nil)
        return
    }

    limitInt, err := strconv.Atoi(limit)
    if err != nil || limitInt < 1 {
        response.BadRequest(ctx, "Invalid limit value", nil)
        return
    }

    // Ambil data user dari service
    users, total, err := h.service.GetAllUsers(pageInt, limitInt, search)
    if err != nil {
        response.InternalServerError(ctx, "Failed to fetch users", err.Error())
        return
    }

    // Hitung pagination info
    pagination := map[string]interface{}{
        "page":       pageInt,
        "limit":      limitInt,
        "total":      total,
        "totalPages": int(math.Ceil(float64(total) / float64(limitInt))),
    }

    // Response sukses dengan pagination
    response.Paginated(ctx, users, pagination, "Users retrieved successfully")
}
