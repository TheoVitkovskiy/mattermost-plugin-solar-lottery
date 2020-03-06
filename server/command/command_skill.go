// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/md"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
)

func (c *Command) skill(parameters []string) (md.MD, error) {
	subcommands := map[string]func([]string) (md.MD, error){
		commandNew:    c.newSkill,
		commandDelete: c.deleteSkill,
		commandList:   c.listSkills,
	}

	return c.handleCommand(subcommands, parameters)
}

func (c *Command) newSkill(parameters []string) (md.MD, error) {
	fs := c.assureFS()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	if len(fs.Args()) != 1 {
		return c.flagUsage(), errors.New("must specify skill")
	}
	skill := types.ID(fs.Arg(0))

	err = c.SL.AddKnownSkill(skill)
	if err != nil {
		return "", err
	}

	if c.outputJson {
		return md.JSONBlock(skill), nil
	}
	return md.Markdownf("Added **%s** to known skills.", skill), nil
}

func (c *Command) deleteSkill(parameters []string) (md.MD, error) {
	fs := c.assureFS()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	if len(fs.Args()) != 1 {
		return c.flagUsage(), errors.New("must specify skill")
	}
	skill := types.ID(fs.Arg(0))

	err = c.SL.DeleteKnownSkill(skill)
	if err != nil {
		return "", err
	}
	if c.outputJson {
		return md.JSONBlock(skill), nil
	}
	return md.Markdownf("Deleted **%s** from known skills. User profiles are not changed.", skill), nil
}

func (c *Command) listSkills(parameters []string) (md.MD, error) {
	fs := c.assureFS()
	err := fs.Parse(parameters)
	if err != nil {
		return c.flagUsage(), err
	}
	skills, err := c.SL.ListKnownSkills()
	if err != nil {
		return "", err
	}
	if c.outputJson {
		return md.JSONBlock(skills), nil
	}
	return "Known skills: " + md.JSONBlock(skills), nil
}
