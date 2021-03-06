// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package solarlottery

import (
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/bot"
	"github.com/pkg/errors"
)

func (sl *solarLottery) Guess(rotation *Rotation, startingShiftNumber int, numShifts int) ([]*Shift, error) {
	err := sl.Filter(
		withActingUserExpanded,
		withRotationExpanded(rotation),
	)
	if err != nil {
		return nil, err
	}
	logger := sl.Logger.Timed().With(bot.LogContext{
		"Location":       "sl.Guess",
		"ActingUsername": sl.actingUser.MattermostUsername(),
		"NumShifts":      numShifts,
		"ShiftNumber":    startingShiftNumber,
		"RotationID":     rotation.RotationID,
	})
	rotation = rotation.Clone(true)

	logger.Debugf("...running guess for\n%s", rotation.MarkdownBullets())
	var shifts []*Shift
	for shiftNumber := startingShiftNumber; shiftNumber < startingShiftNumber+numShifts; shiftNumber++ {
		var shift *Shift
		shift, _, err := sl.getShiftForGuess(rotation, shiftNumber)
		if err != nil {
			if err == store.ErrNotFound {
				shifts = append(shifts, nil)
				continue
			}
			return nil, err
		}

		if shift.Status == store.ShiftStatusOpen {
			autofiller := sl.Dependencies.Autofillers[rotation.Type]
			if autofiller == nil {
				return nil, errors.Errorf("unsupported rotation type %s", rotation.Type)
			}
			var added UserMap
			added, err = autofiller.FillShift(rotation, shiftNumber, shift, logger)
			if err != nil {
				return nil, err
			}

			_, err = sl.joinShift(rotation, shiftNumber, shift, added, false)
			if err != nil {
				return nil, err
			}
		}

		rotation.markShiftUsersEvents(shiftNumber, shift)
		rotation.markShiftUsersServed(shiftNumber, shift)
		shifts = append(shifts, shift)
	}

	logger.Debugf("Ran guess for %s", rotation.Markdown())
	return shifts, nil
}
