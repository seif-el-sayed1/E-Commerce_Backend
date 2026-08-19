package middlewares

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var fieldInParens = regexp.MustCompile(`\(([^)]+)\)`)

func sendError(c *gin.Context, err *utils.ApiError) {
	c.JSON(err.StatusCode, gin.H{
		"success": err.Success,
		"status":  err.Status,
		"message": err.Message,
	})
}

// JWT Errors
func handleInvalidJwtSignature() *utils.ApiError {
	return utils.NewApiError(
		"Invalid token, Please login again ...",
		http.StatusUnauthorized,
	)
}

func handleJwtExpired() *utils.ApiError {
	return utils.NewApiError(
		"Expired token, Please login again ...",
		http.StatusUnauthorized,
	)
}

//PostgreSQL Errors

// extractField tries ColumnName first, then falls back to whatever's inside
func extractField(pgErr *pgconn.PgError) string {
	if pgErr.ColumnName != "" {
		return pgErr.ColumnName
	}

	if match := fieldInParens.FindStringSubmatch(pgErr.Detail); len(match) > 1 {
		return match[1]
	}

	return "field"
}

func handleDuplicatedFieldsDB(pgErr *pgconn.PgError) *utils.ApiError {
	field := utils.CapitalizeFirstLetter(extractField(pgErr))

	return utils.NewApiError(
		field+" is already used",
		http.StatusBadRequest,
	)
}

func handleForeignKeyConstraint(pgErr *pgconn.PgError) *utils.ApiError {
	if strings.Contains(pgErr.Detail, "is not present in table") {
		field := extractField(pgErr)

		return utils.NewApiError(
			"Invalid relation: "+field,
			http.StatusBadRequest,
		)
	}

	return utils.NewApiError(
		"Cannot delete this record because it has related data",
		http.StatusBadRequest,
	)
}

func handleValueTooLong(pgErr *pgconn.PgError) *utils.ApiError {
	field := utils.CapitalizeFirstLetter(extractField(pgErr))

	return utils.NewApiError(
		field+" value is too long",
		http.StatusBadRequest,
	)
}

func handleNullConstraint(pgErr *pgconn.PgError) *utils.ApiError {
	field := utils.CapitalizeFirstLetter(extractField(pgErr))

	return utils.NewApiError(
		field+" is required",
		http.StatusBadRequest,
	)
}

// GORM Errors
func handleRecordNotFound() *utils.ApiError {
	return utils.NewApiError(
		"Record not found",
		http.StatusNotFound,
	)
}

func GlobalError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var apiErr *utils.ApiError
		var pgErr *pgconn.PgError

		switch {

		// JWT
		case errors.Is(err, jwt.ErrTokenExpired):
			apiErr = handleJwtExpired()

		case errors.Is(err, jwt.ErrTokenSignatureInvalid),
			errors.Is(err, jwt.ErrTokenMalformed):
			apiErr = handleInvalidJwtSignature()

		// GORM
		case errors.Is(err, gorm.ErrRecordNotFound):
			apiErr = handleRecordNotFound()

		// PostgreSQL
		case errors.As(err, &pgErr):
			switch pgErr.Code {

			// Unique constraint violation
			case "23505":
				apiErr = handleDuplicatedFieldsDB(pgErr)

			// Foreign key constraint violation
			case "23503":
				apiErr = handleForeignKeyConstraint(pgErr)

			// String too long
			case "22001":
				apiErr = handleValueTooLong(pgErr)

			// NOT NULL violation
			case "23502":
				apiErr = handleNullConstraint(pgErr)

			// Unknown PostgreSQL error
			default:
				apiErr = utils.NewApiError(
					"Something went wrong",
					http.StatusInternalServerError,
				)

				apiErr.IsOperational = false
			}

		case errors.As(err, &apiErr):

		// Unknown Error
		default:
			apiErr = utils.NewApiError(
				"Something went wrong",
				http.StatusInternalServerError,
			)

			apiErr.IsOperational = false
		}

		sendError(c, apiErr)
	}
}
