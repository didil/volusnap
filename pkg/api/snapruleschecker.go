package api

import (
	"fmt"
	"time"

	"github.com/didil/volusnap/pkg/models"
	"github.com/sirupsen/logrus"
)

type snapRulesChecker struct {
	snapRuleSvc snapRuleSvcer
	snapshotSvc snapshotSvcer
	accountSvc  accountSvcer
	shooter     snapshooter
	ticker      *time.Ticker
}

func newSnapRulesChecker(snapRuleSvc snapRuleSvcer, snapshotSvc snapshotSvcer, accountSvc accountSvcer, shooter snapshooter) *snapRulesChecker {
	return &snapRulesChecker{
		snapRuleSvc: snapRuleSvc,
		snapshotSvc: snapshotSvc,
		accountSvc:  accountSvc,
		shooter:     shooter,
	}
}

func (checker *snapRulesChecker) Start() {
	logrus.Infof("Starting snapRulesChecker ...")
	checker.ticker = time.NewTicker(5 * time.Minute)
	go func() {
		for range checker.ticker.C {
			logrus.Infof("Checking SnapRules ...")
			err := checker.checkAll()
			if err != nil {
				logrus.Errorf("checkall snaprules err: %v", err)
			}
		}
	}()
}

func (checker *snapRulesChecker) Stop() {
	logrus.Infof("Stopping snapRulesChecker ...")
	checker.ticker.Stop()
}

func (checker *snapRulesChecker) checkAll() error {
	snapRules, err := checker.snapRuleSvc.ListAll()
	if err != nil {
		return fmt.Errorf("list snaprules err: %v", err)
	}

	for _, snapRule := range snapRules {
		err := checker.check(snapRule)
		if err != nil {
			logrus.Errorf("check snaprule %v err: %v", snapRule.ID, err)
		}
	}

	return nil
}

func (checker *snapRulesChecker) check(snapRule *models.SnapRule) error {
	createdAfter := time.Now().Add(time.Duration(-1*snapRule.Frequency) * time.Hour)
	exists, err := checker.snapshotSvc.ExistsFor(snapRule.ID, createdAfter)
	if err != nil {
		return fmt.Errorf("snapshots exists query err: %v", err)
	}

	if exists {
		// snapshot exists
		return nil
	}

	account, err := checker.accountSvc.Get(snapRule.AccountID)
	if err != nil {
		return fmt.Errorf("get account err: %v", err)
	}

	providerSnapshotID, err := checker.shooter.Take(account, snapRule)
	if err != nil {
		return fmt.Errorf("take snapshot err: %v", err)
	}

	id, err := checker.snapshotSvc.Create(snapRule.ID, providerSnapshotID)
	if err != nil {
		return fmt.Errorf("create snapshot err: %v", err)
	}

	logrus.Infof("Created snapshot #%d for snaprule #%v", id, snapRule.ID)
	return nil
}
