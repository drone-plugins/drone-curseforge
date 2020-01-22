// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"io/ioutil"
	"os"
)

func readStringOrFile(input string) (string, error) {
	if len(input) > 255 {
		return input, nil
	}

	if _, err := os.Stat(input); err != nil && os.IsNotExist(err) {
		return input, nil
	} else if err != nil {
		return "", err
	}

	result, err := ioutil.ReadFile(input)

	if err != nil {
		return "", err
	}

	return string(result), nil
}
