// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	"os"

	"github.com/drone-plugins/drone-curseforge/plugin"
	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/urfave/cli/v2"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "drone-curseforge"
	app.Usage = "publish files to curseforge"
	app.Version = version
	app.Action = run
	app.Flags = append(settingsFlags(), urfave.Flags()...)

	if err := app.Run(os.Args); err != nil {
		errors.HandleExit(err)
	}
}

func run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)

	plugin := plugin.New(
		settingsFromContext(ctx),
		urfave.PipelineFromContext(ctx),
		urfave.NetworkFromContext(ctx),
	)

	if err := plugin.Validate(); err != nil {
		return err
	}

	if err := plugin.Execute(); err != nil {
		return err
	}

	return nil
}

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags() []cli.Flag {
	// Replace below with all the flags required for the plugin.
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "plugin.apikey",
			Usage: "api key to access curseforge",
			EnvVars: []string{
				"PLUGIN_API_KEY",
				"CURSEFORGE_API_KEY",
			},
		},
		&cli.IntFlag{
			Name:  "plugin.project",
			Usage: "project id used on curseforge",
			EnvVars: []string{
				"PLUGIN_PROJECT",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.file",
			Usage: "path to file for release",
			EnvVars: []string{
				"PLUGIN_FILE",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.title",
			Usage: "title of the release",
			EnvVars: []string{
				"PLUGIN_TITLE",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.release",
			Value: "release",
			Usage: "type of the release",
			EnvVars: []string{
				"PLUGIN_RELEASE",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.note",
			Usage: "changelog of the release",
			EnvVars: []string{
				"PLUGIN_NOTE",
				"PLUGIN_CHANGELOG",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.type",
			Value: "markdown",
			Usage: "type of changelog",
			EnvVars: []string{
				"PLUGIN_TYPE",
			},
		},
		&cli.IntSliceFlag{
			Name:  "plugin.games",
			Usage: "list of game ids or names",
			EnvVars: []string{
				"PLUGIN_GAMES",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.relations",
			Usage: "relations to other projects",
			EnvVars: []string{
				"PLUGIN_RELATIONS",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.manifest",
			Usage: "path to manifest file to parse",
			EnvVars: []string{
				"PLUGIN_MANIFEST",
			},
		},
		&cli.StringFlag{
			Name:  "plugin.metadata",
			Usage: "overwrite metadata payload",
			EnvVars: []string{
				"PLUGIN_METADATA",
			},
		},
	}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) plugin.Settings {
	return plugin.Settings{
		APIKey:    ctx.String("plugin.apikey"),
		Project:   ctx.Int("plugin.project"),
		File:      ctx.String("plugin.file"),
		Title:     ctx.String("plugin.title"),
		Release:   ctx.String("plugin.release"),
		Note:      ctx.String("plugin.note"),
		Type:      ctx.String("plugin.type"),
		Games:     ctx.IntSlice("plugin.games"),
		Relations: ctx.String("plugin.relations"),
		Manifest:  ctx.String("plugin.manifest"),
		Metadata:  ctx.String("plugin.metadata"),
	}
}
