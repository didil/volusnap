package api

import "database/sql"

type appController struct {
	authCtrl     *authController
	accountCtrl  *accountController
	volumeCtrl   *volumeController
	snapRuleCtrl *snapRuleController
	snapshotCtrl *snapshotController
}

func buildAppController(db *sql.DB) *appController {
	authSvc := newAuthService(db)
	authCtrl := newAuthController(authSvc)

	accountSvc := newAccountService(db)
	accountCtrl := newAccountController(accountSvc)

	volumeCtrl := newVolumeController(accountSvc)

	snapRuleSvc := newSnapRuleService(db)
	snapRuleCtrl := newSnapRuleController(snapRuleSvc, accountSvc)

	snapshotSvc := newSnapshotService(db)
	snapshotCtrl := newSnapshotController(snapshotSvc, accountSvc)

	return &appController{
		authCtrl:     authCtrl,
		accountCtrl:  accountCtrl,
		volumeCtrl:   volumeCtrl,
		snapRuleCtrl: snapRuleCtrl,
		snapshotCtrl: snapshotCtrl,
	}
}
