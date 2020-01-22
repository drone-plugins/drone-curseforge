// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

// ForgeSvc defines the forgesvc api type.
type ForgeSvc struct {
	Slug string `json:"slug"`
}

// File defines a record within manifest.
type File struct {
	ProjectID int  `json:"projectID"`
	FileID    int  `json:"fileID"`
	Required  bool `json:"required"`
}

// Manifest defines the manifest input.
type Manifest struct {
	Files []File `json:"files"`
}

// Project is par of the api payload.
type Project struct {
	Slug string `json:"slug"`
	Type string `json:"type"`
}

// Relations is par of the api payload.
type Relations struct {
	Projects []Project `json:"projects"`
}

// Metadata implements the api payload.
type Metadata struct {
	Title     string    `json:"displayName"`
	Release   string    `json:"releaseType"`
	Note      string    `json:"changelog"`
	Type      string    `json:"changelogType"`
	Games     []int     `json:"gameVersions"`
	Relations Relations `json:"relations"`
}

// Response implements the api response.
type Response struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
	ID      int    `json:"id"`
}
