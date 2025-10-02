// Copyright (c) 2024-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

const (
	resetDataCommand    = "resetdata"
	listSurveysCommand  = "listsurveys"
	listSessionsCommand = "listsessions"
)

func (p *Plugin) registerDebugCommands() error {
	err := p.API.RegisterCommand(&model.Command{
		Trigger:      resetDataCommand,
		AutoComplete: true,
	})

	if err != nil {
		p.API.LogError("registerDebugCommands: failed to register reset data command", "error", err.Error())
		return errors.Wrap(err, "registerDebugCommands: failed to register reset data command")
	}

	err = p.API.RegisterCommand(&model.Command{
		Trigger:      listSurveysCommand,
		AutoComplete: true,
	})
	if err != nil {
		p.API.LogError("registerDebugCommands: failed to register list surveys command", "error", err.Error())
		return errors.Wrap(err, "registerDebugCommands: failed to register list surveys command")
	}

	err = p.API.RegisterCommand(&model.Command{
		Trigger:      listSessionsCommand,
		AutoComplete: true,
	})
	if err != nil {
		p.API.LogError("registerDebugCommands: failed to register list sessions command", "error", err.Error())
		return errors.Wrap(err, "registerDebugCommands: failed to register list sessions command")
	}

	return nil
}

func (p *Plugin) ExecuteCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	if len(split) == 0 {
		return nil, nil
	}
	command := split[0]

	switch command {
	case "/" + resetDataCommand:
		return p.executeResetDataCommand(ctx, args)
	case "/" + listSurveysCommand:
		return p.executeListSurveysCommand(ctx, args)
	case "/" + listSessionsCommand:
		return p.executeListSessionsCommand(ctx, args)
	}
	return nil, nil
}

func (p *Plugin) executeResetDataCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	user, appErr := p.API.GetUser(args.UserId)
	if appErr != nil {
		p.API.LogError("executeResetDataCommand: failed to get user by id", "userID", args.UserId, "error", appErr.Error())

		return &model.CommandResponse{
			Text: "There was an error executing the command",
		}, nil
	}

	if !user.IsSystemAdmin() {
		return &model.CommandResponse{}, nil
	}

	p.API.LogWarn("Processing request to reset all user survey data. Requested by user ID: " + args.UserId)

	var message string

	err := p.app.ResetData()
	if err != nil {
		message = err.Error()
	} else {
		message = "Successfully reset survey data"
	}

	return &model.CommandResponse{
		Text: message,
	}, nil
}

func (p *Plugin) executeListSurveysCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	surveys, err := p.app.GetSurveys()
	if err != nil {
		return &model.CommandResponse{Text: "Failed to list surveys: " + err.Error()}, nil
	}
	if len(surveys) == 0 {
		return &model.CommandResponse{Text: "No surveys found."}, nil
	}
	var ids []string
	for _, s := range surveys {
		ids = append(ids, s.ID)
	}
	return &model.CommandResponse{Text: "Survey IDs: " + strings.Join(ids, ", ")}, nil
}

func (p *Plugin) executeListSessionsCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	sessions, err := p.store.GetSessions()
	if err != nil {
		return &model.CommandResponse{Text: "Failed to list surveys: " + err.Error()}, nil
	}
	if len(sessions) == 0 {
		return &model.CommandResponse{Text: "No sessions found."}, nil
	}
	var ids []string
	for _, s := range sessions {
		ids = append(ids, s.ID)
	}
	return &model.CommandResponse{Text: "Session IDs: " + strings.Join(ids, ", ")}, nil
}
