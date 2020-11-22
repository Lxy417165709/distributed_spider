package cache

import (
	"go.uber.org/zap"
	"spider/common/logger"
	"time"
)

type spider struct {
	supplierLock string
}

var Spider = &spider{
	supplierLock: "supplierLock",
}

func (s *spider) GetSupplierLock() bool {
	ok, err := mainRedis.SetNX(s.supplierLock, 1, time.Second).Result()
	if err != nil {
		logger.Error("Fail to finish mainRedis.SetNX", zap.String("key", s.supplierLock), zap.Error(err))
		return false
	}
	return ok
}

func (s *spider) ReleaseSupplierLock() {
	logger.Info("ReleaseSupplierLock")
	if err := mainRedis.Del(s.supplierLock).Err(); err != nil {
		logger.Error("Fail to finish mainRedi", zap.String("key", s.supplierLock), zap.Error(err))
		return
	}
}
