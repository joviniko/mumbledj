/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/shuffle_test.go
 * Copyright (c) 2014, 2015 Matthieu Grieger (MIT License)
 */

package commands

import (
	"testing"

	"github.com/matthieugrieger/mumbledj/state"
	"github.com/matthieugrieger/mumbledj/testutils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ShuffleCommandTestSuite struct {
	Command ShuffleCommand
	State   *state.BotState
	suite.Suite
}

func (suite *ShuffleCommandTestSuite) SetupSuite() {
	viper.Set("aliases.shuffle", []string{"shuffle", "shuf", "sh"})
	viper.Set("descriptions.shuffle", "shuffle")
	viper.Set("permissions.shuffle", true)
}

func (suite *ShuffleCommandTestSuite) SetupTest() {
	suite.State = new(state.BotState)
	suite.State.Queue = state.NewAudioQueue()
}

func (suite *ShuffleCommandTestSuite) TestAliases() {
	suite.Equal([]string{"shuffle", "shuf", "sh"}, suite.Command.Aliases())
}

func (suite *ShuffleCommandTestSuite) TestDescription() {
	suite.Equal("shuffle", suite.Command.Description())
}

func (suite *ShuffleCommandTestSuite) TestIsAdmin() {
	suite.True(suite.Command.IsAdmin())
}

func (suite *ShuffleCommandTestSuite) TestExecuteWithEmptyQueue() {
	state, message, isPrivateMessage, err := suite.Command.Execute(suite.State, nil)

	suite.Nil(state, "No state should be returned since an error occurred.")
	suite.Equal("", message, "No message should be returned since an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned for attempting to shuffle an empty queue.")
}

func (suite *ShuffleCommandTestSuite) TestExecuteWithNotEnoughTracks() {
	suite.State.Queue.AddTracks(new(testutils.MockedAudioTrack))

	state, message, isPrivateMessage, err := suite.Command.Execute(suite.State, nil)
	suite.Equal(1, len(suite.State.Queue.Queue))
	suite.Nil(state, "No state should be returned since an error occurred.")
	suite.Equal("", message, "No message should be returned since an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned for attempting to shuffle a queue with only one track.")

	suite.State.Queue.AddTracks(new(testutils.MockedAudioTrack))

	state, message, isPrivateMessage, err = suite.Command.Execute(suite.State, nil)
	suite.Equal(2, len(suite.State.Queue.Queue))
	suite.Nil(state, "No state should be returned since an error occurred.")
	suite.Equal("", message, "No message should be returned since an error occurred.")
	suite.True(isPrivateMessage, "This should be a private message.")
	suite.NotNil(err, "An error should be returned for attempting to shuffle a queue with only two tracks.")
}

func (suite *ShuffleCommandTestSuite) TestExecuteWithValidQueue() {
	suite.State.Queue.AddTracks(new(testutils.MockedAudioTrack))
	suite.State.Queue.AddTracks(new(testutils.MockedAudioTrack))
	suite.State.Queue.AddTracks(new(testutils.MockedAudioTrack))

	state, message, isPrivateMessage, err := suite.Command.Execute(suite.State, nil)
	suite.Equal(3, len(suite.State.Queue.Queue))
	suite.NotNil(state, "An updated state should be returned since the execution was successful.")
	suite.NotEqual("", message, "A message should be returned since the execution was successful.")
	suite.False(isPrivateMessage, "This message should not be private.")
	suite.Nil(err, "No error should be returned.")
}

func TestShuffleCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ShuffleCommandTestSuite))
}
