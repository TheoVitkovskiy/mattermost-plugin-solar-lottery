// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
)

func (c *Command) showUser(parameters []string) (string, error) {
	fs := newFS()
	_ = fJSON(fs)
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(fs), err
	}

	mattermostUserIDs, err := c.resolveUsernames(fs.Args())
	if err != nil {
		return "", err
	}
	users, err := c.SL.LoadUsers(mattermostUserIDs)
	if err != nil {
		return "", err
	}

	if users.Len() == 1 {
		return md.JSONBlock(users.AsArray()[0]), nil
	}
	return md.JSONBlock(users), nil
}
