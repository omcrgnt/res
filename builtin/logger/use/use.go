package use

import (
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/builtin/logger"
)

func init() {
	// Регистрируем дефолтный конфиг логгера в системном реестре
	res.Register(logger.DefaultConfig())
}
