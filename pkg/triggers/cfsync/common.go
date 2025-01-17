package cfsync

import (
	"strings"
	"time"

	"github.com/Notifiarr/notifiarr/pkg/apps"
	"github.com/Notifiarr/notifiarr/pkg/triggers/common"
	"github.com/Notifiarr/notifiarr/pkg/website/clientinfo"
	"golift.io/cnfg"
)

/* CF Sync means Custom Format Sync. This is a premium feature that allows syncing
   TRaSH's custom Radarr formats and Sonarr Release Profiles.
	 The code in this file deals with sending data and getting updates at an interval.
*/

const (
	randomMilliseconds = 5000
)

// New configures the library.
func New(config *common.Config) *Action {
	return &Action{
		cmd: &cmd{
			Config: config,
		},
	}
}

// Action contains the exported methods for this package.
type Action struct {
	cmd *cmd
}

type cmd struct {
	*common.Config
}

// Create initializes the library.
func (a *Action) Create() {
	a.cmd.create()
}

func (c *cmd) create() {
	ci := clientinfo.Get()
	c.setupRadarr(ci)
	c.setupSonarr(ci)

	// Check each instance and enable only if needed.
	if ci != nil && ci.Actions.Sync.Interval.Duration > 0 {
		if len(ci.Actions.Sync.RadarrInstances) > 0 {
			c.Printf("==> Radarr TRaSH Sync: interval: %s, %s ",
				ci.Actions.Sync.Interval, strings.Join(ci.Actions.Sync.RadarrSync, ", "))
		}

		if len(ci.Actions.Sync.SonarrInstances) > 0 {
			c.Printf("==> Sonarr TRaSH Sync: interval: %s, %s ",
				ci.Actions.Sync.Interval, strings.Join(ci.Actions.Sync.SonarrSync, ", "))
		}
	}

	// These aggregate  triggers have no timers. Used to sync "all the things" at once.
	c.Add(&common.Action{
		Name: TrigCFSyncRadarr,
		Fn:   c.syncRadarr,
		C:    make(chan *common.ActionInput, 1),
	}, &common.Action{
		Name: TrigRPSyncSonarr,
		Fn:   c.syncSonarr,
		C:    make(chan *common.ActionInput, 1),
	})
}

type radarrApp struct {
	app *apps.RadarrConfig
	cmd *cmd
	idx int
}

func (c *cmd) setupRadarr(ci *clientinfo.ClientInfo) {
	if ci == nil {
		return
	}

	for idx, app := range c.Apps.Radarr {
		instance := idx + 1
		if !app.Enabled() || !ci.Actions.Sync.RadarrInstances.Has(instance) {
			continue
		}

		var dur cnfg.Duration

		if ci != nil && ci.Actions.Sync.Interval.Duration > 0 {
			randomTime := time.Duration(c.Config.Rand().Intn(randomMilliseconds)) * time.Millisecond
			dur = cnfg.Duration{Duration: ci.Actions.Sync.Interval.Duration + randomTime}
		}

		c.Add(&common.Action{
			Hide: true,
			D:    dur,
			Name: TrigCFSyncRadarrInt.WithInstance(instance),
			Fn:   (&radarrApp{app: app, cmd: c, idx: idx}).syncRadarr,
			C:    make(chan *common.ActionInput, 1),
		})
	}
}

type sonarrApp struct {
	app *apps.SonarrConfig
	cmd *cmd
	idx int
}

func (c *cmd) setupSonarr(ci *clientinfo.ClientInfo) {
	if ci == nil {
		return
	}

	for idx, app := range c.Apps.Sonarr {
		instance := idx + 1
		if !app.Enabled() || !ci.Actions.Sync.SonarrInstances.Has(instance) {
			continue
		}

		var dur cnfg.Duration

		if ci != nil && ci.Actions.Sync.Interval.Duration > 0 {
			randomTime := time.Duration(c.Config.Rand().Intn(randomMilliseconds)) * time.Millisecond
			dur = cnfg.Duration{Duration: ci.Actions.Sync.Interval.Duration + randomTime}
		}

		c.Add(&common.Action{
			Hide: true,
			D:    dur,
			Name: TrigRPSyncSonarrInt.WithInstance(instance),
			Fn:   (&sonarrApp{app: app, cmd: c, idx: idx}).syncSonarr,
			C:    make(chan *common.ActionInput, 1),
		})
	}
}
