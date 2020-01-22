// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/sirupsen/logrus"
)

const (
	uploadURL   = "https://minecraft.curseforge.com/api/projects/%d/upload-file"
	forgesvcURL = "https://addons-ecs.forgesvc.net/api/v2/addon/%d"
)

// Settings for the plugin.
type Settings struct {
	APIKey    string
	Project   int
	File      string
	Title     string
	Release   string
	Note      string
	Type      string
	Games     []int
	Relations string
	Manifest  string
	Metadata  string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if p.settings.APIKey == "" {
		return errors.ExitMessage("no apikey provided")
	}

	if p.settings.Project == 0 {
		return errors.ExitMessage("no project provided")
	}

	if p.settings.File == "" {
		return errors.ExitMessage("no file provided")
	}

	if _, err := os.Stat(p.settings.File); os.IsNotExist(err) {
		return errors.ExitMessage("file doesn't exist")
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	metadata, err := p.metadata()

	if err != nil {
		return err
	}

	file, err := os.Open(p.settings.File)

	if err != nil {
		return errors.WithFields("failed to read and process file", logrus.Fields{
			"err": err,
		})
	}

	defer file.Close()

	content := bytes.NewBufferString("")
	multipart := multipart.NewWriter(content)

	fileWriter, err := multipart.CreateFormFile(
		"file",
		path.Base(p.settings.File),
	)

	if err != nil {
		return errors.WithFields("failed to create form file", logrus.Fields{
			"err": err,
		})
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return errors.WithFields("failed to copy form file", logrus.Fields{
			"err": err,
		})
	}

	fieldWriter, err := multipart.CreateFormField("metadata")

	if err != nil {
		return errors.WithFields("failed to create form field", logrus.Fields{
			"err": err,
		})
	}

	if _, err = fieldWriter.Write(metadata); err != nil {
		return errors.WithFields("failed to copy form field", logrus.Fields{
			"err": err,
		})
	}

	multipart.Close()

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			uploadURL,
			p.settings.Project,
		),
		content,
	)

	if err != nil {
		return errors.WithFields("failed to prepare request", logrus.Fields{
			"err": err,
		})
	}

	req.Header.Set("X-Api-Token", p.settings.APIKey)
	req.Header.Set("Content-Type", multipart.FormDataContentType())

	resp, err := p.network.Client.Do(req)

	if err != nil {
		return errors.WithFields("failed to submit request", logrus.Fields{
			"err": err,
		})
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return errors.WithFields("failed to read request", logrus.Fields{
			"err": err,
		})
	}

	result := Response{}

	if err := json.Unmarshal(body, &result); err != nil {
		return errors.WithFields("failed to parse response", logrus.Fields{
			"err":  err,
			"body": string(body),
		})
	}

	if resp.StatusCode != 200 {
		return errors.WithFields("failed to upload release", logrus.Fields{
			"code":    result.Code,
			"message": result.Message,
		})
	}

	logrus.WithFields(logrus.Fields{
		"id": result.ID,
	}).Info("successfully uploaded release")

	return nil
}

func (p *Plugin) metadata() ([]byte, error) {
	if p.settings.Metadata != "" {
		rawdata, err := readStringOrFile(p.settings.Metadata)

		if err != nil {
			return []byte(""), errors.WithFields("failed to read metadata from file", logrus.Fields{
				"err": err,
			})
		}

		return []byte(rawdata), nil
	}

	relations := Relations{
		Projects: make([]Project, 0),
	}

	if p.settings.Relations != "" {
		rawdata, err := readStringOrFile(p.settings.Relations)

		if err != nil {
			return []byte(""), errors.WithFields("failed to read relations from file", logrus.Fields{
				"err": err,
			})
		}

		if err := json.Unmarshal([]byte(rawdata), &relations.Projects); err != nil {
			return []byte(""), errors.WithFields("failed to parse relations", logrus.Fields{
				"err": err,
			})
		}
	}

	if p.settings.Manifest != "" {
		rawdata, err := readStringOrFile(p.settings.Manifest)

		if err != nil {
			return []byte(""), errors.WithFields("failed to read manifest from file", logrus.Fields{
				"err": err,
			})
		}

		manifest := Manifest{}

		if err := json.Unmarshal([]byte(rawdata), &manifest); err != nil {
			return []byte(""), errors.WithFields("failed to parse manifest", logrus.Fields{
				"err": err,
			})
		}

		for _, file := range manifest.Files {
			resp, err := p.network.Client.Get(
				fmt.Sprintf(
					forgesvcURL,
					file.ProjectID,
				),
			)

			if err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"project": file.ProjectID,
					"file":    file.FileID,
				}).Info("failed to gather dependency slug")

				continue
			}

			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"project": file.ProjectID,
					"file":    file.FileID,
				}).Info("failed to read dependency body")

				continue
			}

			forgesvc := ForgeSvc{}

			if err := json.Unmarshal(body, &forgesvc); err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"project": file.ProjectID,
					"file":    file.FileID,
				}).Info("failed to parse dependency body")

				continue
			}

			dependency := "requiredDependency"

			if !file.Required {
				dependency = "optionalDependency"
			}

			relations.Projects = append(relations.Projects, Project{
				Slug: forgesvc.Slug,
				Type: dependency,
			})

		}
	}

	release, err := readStringOrFile(p.settings.Release)

	if err != nil {
		return []byte(""), errors.WithFields("failed to read release from file", logrus.Fields{
			"err": err,
		})
	}

	note, err := readStringOrFile(p.settings.Note)

	if err != nil {
		return []byte(""), errors.WithFields("failed to read note from file", logrus.Fields{
			"err": err,
		})
	}

	metadata := Metadata{
		Title:     p.settings.Title,
		Release:   release,
		Note:      note,
		Type:      p.settings.Type,
		Games:     p.settings.Games,
		Relations: relations,
	}

	rawdata, err := json.Marshal(metadata)

	if err != nil {
		return []byte(""), errors.WithFields("failed to encode metadata", logrus.Fields{
			"err": err,
		})
	}

	return rawdata, nil
}
