package exception

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func ResponseForumException(gctx *gin.Context, err error) {
	defer gctx.Abort()

	var feErr *ForumException
	if errors.As(err, &feErr) {
		gctx.JSON(feErr.HttpCode, feErr)
		return
	}

	otherErr := NewNonForumException(err)
	gctx.JSON(otherErr.HttpCode, otherErr)

}
