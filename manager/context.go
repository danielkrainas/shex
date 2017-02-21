package manager

import (
	"context"

	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/shex/configuration"
	"github.com/danielkrainas/shex/profile"
)

/* execution context */
type ExecutionContext struct {
	context.Context
	Channels ChannelMap
	HomePath string
	Profiles map[string]*Profile
	Config   *configuration.Config
}

func (ctx *ExecutionContext) Value(key interface{}) interface{} {
	switch key {
	case "channels":
		return ctx.Channels
	case "channel":
		return ctx.Channel()
	case "profiles":
		return ctx.Profiles
	case "profile":
		return ctx.Profile()
	}

	return ctx.Context.Value(key)
}

func (ctx *ExecutionContext) Profile() *profile.Profile {
	return ctx.Profiles[ctx.Config.ActiveProfile]
}

func (ctx *ExecutionContext) Channel() *channel.Channel {
	return ctx.Channels[ctx.Config.ActiveRemote]
}

func Context(parent context.Context, homePath string) (*ExecutionContext, error) {
	config, err := configuration.Resolve(homePath)
	if err != nil {
		return nil, err
	}

	profiles, err := LoadAllProfiles(config.ProfilesPath)
	if err != nil {
		return nil, err
	}

	channels, err := LoadAllChannels(config.ChannelsPath)
	if err != nil {
		return nil, err
	}

	if config.IncludeDefaultChannel {
		channels[defaultChannel.Alias] = defaultChannel
	}

	ctx := &ExecutionContext{
		Context:  parent,
		HomePath: homePath,
		Profiles: profiles,
		Config:   config,
		Channels: channels,
	}

	return ctx, nil
}
