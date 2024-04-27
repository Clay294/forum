package impl

import (
	"github.com/Clay294/forum/apps/users"
	"github.com/Clay294/forum/config"
	"github.com/Clay294/forum/ioc"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

type serviceController struct {
	gdbForum *gorm.DB
}

func (sc *serviceController) Init() error {
	gdbs, err := config.GlobalConf().CreateConnByORM()
	if err != nil {
		return err
	}

	sc.gdbForum = gdbs[config.GlobalConf().MySQLForumBase.Database].Debug()

	return nil
}

func (sc *serviceController) Name() string {
	return users.UnitName
}

func init() {
	sc := new(serviceController)

	err := ioc.Controllers().Registry(sc)

	if err != nil {
		log.Panic().Caller().Msgf("registering service controlelr of unit %s to ioc controllers failed: %s", sc.Name(), err)
	}
}
