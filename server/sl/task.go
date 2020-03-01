// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package sl

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

type TaskStatus string

const (
	TaskStatusPending  = TaskStatus("pending")
	TaskStatusStarted  = TaskStatus("started")
	TaskStatusFinished = TaskStatus("finished")
)

type Task struct {
	PluginVersion string
	TaskID        types.ID
	Status        TaskStatus
	Created       types.Time
	Summary       string

	Scheduled         *types.Interval `json:",omitempty"`
	Min               Needs           `json:",omitempty"`
	Max               Needs           `json:",omitempty"`
	Actual            *types.Interval `json:",omitempty"`
	Grace             time.Duration   `json:",omitempty"`
	MattermostUserIDs *types.IDIndex  `json:",omitempty"`

	users Users
}

func NewTask() *Task {
	return &Task{
		Status:            TaskStatusPending,
		Created:           types.NewTime(),
		MattermostUserIDs: types.NewIDIndex(),
		users:             NewUsers(),
	}
}

func (t Task) MarkdownBullets(rotation *Rotation) string {
	out := fmt.Sprintf("- %s\n", t.Markdown())
	out += fmt.Sprintf("  - Status: **%s**\n", t.Status)
	out += fmt.Sprintf("  - Users: **%v**\n", t.MattermostUserIDs.Len())
	// for _, user := range rotation.TaskUsers(&t) {
	// 	out += fmt.Sprintf("    - %s\n", user.MarkdownWithSkills())
	// }
	return out
}

func (t Task) Markdown() string {
	return fmt.Sprintf("%s", t.TaskID)
}

func (t *Task) NewUnavailable() []*Unavailable {
	interval := t.Actual
	if interval.IsEmpty() {
		interval = t.Scheduled
	}
	if interval.IsEmpty() {
		now := types.NewTime()
		interval = &types.Interval{
			Start:  now,
			Finish: now,
		}
	}
	uu := []*Unavailable{
		{
			Reason:   ReasonTask,
			Interval: *interval,
			TaskID:   t.TaskID,
		},
	}

	if t.Grace > 0 {
		uu = append(uu, &Unavailable{
			Reason: ReasonGrace,
			Interval: types.Interval{
				Start:  interval.Finish,
				Finish: types.NewTime(interval.Finish.Add(t.Grace)),
			},
			TaskID: t.TaskID,
		})
	}

	return uu
}
