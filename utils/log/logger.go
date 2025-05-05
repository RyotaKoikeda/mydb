package log

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var logId string

type logIdKey struct{}

func init() {
	isLogTypeDevelopment, err := strconv.ParseBool(os.Getenv("IS_LOG_TYPE_DEV"))
	if err != nil {
		isLogTypeDevelopment = false
	}

	callerSkip := zap.AddCallerSkip(1)
	if isLogTypeDevelopment {
		logger, _ = zap.NewDevelopment(callerSkip)
	} else {
		logger, _ = zap.NewProduction(callerSkip)
	}
}

func GetLogIdFromContext(ctx context.Context) (string, bool) {
	logId, ok := ctx.Value(logIdKey{}).(string)
	return logId, ok
}

func SetLogIdAndContext(r *http.Request, logId string) *http.Request {
	ctx := context.WithValue(r.Context(), logIdKey{}, logId)
	return r.WithContext(ctx)
}

func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Float(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func Object(objectKey string, object interface{}) zap.Field {
	jsonByte, _ := json.Marshal(object)
	return zap.String(objectKey, string(jsonByte))
}

func Info(msg string, fields ...zap.Field) {
	fields = append([]zapcore.Field{String("log-id", logId)}, fields...)
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	fields = append([]zapcore.Field{String("log-id", logId)}, fields...)
	logger.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	fields = append([]zapcore.Field{String("log-id", logId)}, fields...)
	logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	fields = append([]zapcore.Field{String("log-id", logId)}, fields...)
	logger.Error(msg, fields...)
}

func InfoV2(ctx context.Context, msg string, fields ...zap.Field) {
	logIdCtx, ok := GetLogIdFromContext(ctx)
	if !ok {
		logIdCtx = ""
	}
	fields = append([]zapcore.Field{String("log-id", logIdCtx)}, fields...)
	logger.Info(msg, fields...)
}

func WarnV2(ctx context.Context, msg string, fields ...zap.Field) {
	logIdCtx, ok := GetLogIdFromContext(ctx)
	if !ok {
		logIdCtx = ""
	}
	fields = append([]zapcore.Field{String("log-id", logIdCtx)}, fields...)
	logger.Warn(msg, fields...)
}

func DebugV2(ctx context.Context, msg string, fields ...zap.Field) {
	logIdCtx, ok := GetLogIdFromContext(ctx)
	if !ok {
		logIdCtx = ""
	}
	fields = append([]zapcore.Field{String("log-id", logIdCtx)}, fields...)
	logger.Debug(msg, fields...)
}

func ErrorV2(ctx context.Context, msg string, fields ...zap.Field) {
	logIdCtx, ok := GetLogIdFromContext(ctx)
	if !ok {
		logIdCtx = ""
	}
	fields = append([]zapcore.Field{String("log-id", logIdCtx)}, fields...)
	logger.Error(msg, fields...)
}

func GetLogId() string {
	return logId
}

func SetLogId(id string) {
	logId = id
}
