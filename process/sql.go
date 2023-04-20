package process

import (
	log "malawi-callback/logger"
	"malawi-callback/models"
)

func (c *Controller) getMbtPerId(mbtId string) (*models.MbtEntity, error) {
	log.Debugf(*c.requestId, "trying to get MBT from DB", mbtId)
	mbt, err := (*c.repository).GetMbt(mbtId)

	if err != nil {
		return mbt, err
	}

	log.Debugf(*c.requestId, "payment still pending", mbt)
	return mbt, nil
}
