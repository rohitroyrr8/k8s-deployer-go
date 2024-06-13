package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// randomID generates a random hex string in the UUID format:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func randomID() (string, error) {
	var u [16]byte
	if _, err := rand.Read(u[:]); err != nil {
		return "", fmt.Errorf("genUUID: %v", err)
	}

	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf), nil
}

func PrepareCustomError(ctx *gin.Context, errorCode int, functionName string, displayMessage string, details string) {
	fName := "utils/PrepareCustomError"
	tracer := otel.Tracer("api")
	_, span := tracer.Start(ctx.Request.Context(), fName)
	defer span.End()

	errID, err := randomID()
	if err != nil {
		errID = "unknown(" + err.Error() + ")"
	}

	ctx.Header("X-ErrID", errID)
	fmt.Printf("api error id=%q,path=%q,function name= %s,code=%d,err=%q,detail=%q", errID, ctx.Request.URL, functionName, errorCode, displayMessage, details)

	span.RecordError(fmt.Errorf(fmt.Sprintf("api error id=%q,path=%q,function name= %s,code=%d,err=%q,detail=%q", errID, ctx.Request.URL, functionName, errorCode, displayMessage, details)))
	span.SetStatus(codes.Error, details)

	ctx.JSON(errorCode, gin.H{
		"code":    http.StatusBadRequest,
		"message": displayMessage,
		"error":   errID,
	})

	// TODO: need to abort the execution and return the error object back
	return
}

func PrepareCustomResponse(ctx *gin.Context, displayMessage string, body any) {
	fName := "utils/PrepareCustomResponse"
	tracer := otel.Tracer("api")
	_, span := tracer.Start(ctx.Request.Context(), fName)
	defer span.End()

	span.SetStatus(codes.Ok, "")
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": displayMessage,
		"body":    body,
	})
}
