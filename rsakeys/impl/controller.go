package impl

import (
	"github.com/Clay294/forum/config"
	"github.com/Clay294/forum/ioc"
	"github.com/Clay294/forum/rsakeys"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type serviceController struct {
	gdbKeys *gorm.DB
}

func (sc *serviceController) Init() error {
	gdbs, err := config.GlobalConf().CreateConnByORM()

	if err != nil {
		return err
	}

	sc.gdbKeys = gdbs[config.GlobalConf().MySQLKeysBase.Database].Debug()

	return nil
}

func (sc *serviceController) Name() string {
	return rsakeys.UnitName
}

func init() {
	sc := new(serviceController)

	err := ioc.Controllers().Registry(sc)

	if err != nil {
		log.Panic().Caller().Msgf("registering service controlelr of unit %s to ioc controllers failed: %s", sc.Name(), err)
	}
}
