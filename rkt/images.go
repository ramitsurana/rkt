// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//+build linux

package main

import (
	"fmt"

	"github.com/coreos/rkt/Godeps/_workspace/src/github.com/appc/spec/schema/types"
	"github.com/coreos/rkt/cas"
	"github.com/coreos/rkt/common/apps"
	"github.com/coreos/rkt/pkg/keystore"
)

// findImages uses findImage to attain a list of image hashes using discovery if necessary
func findImages(al *apps.Apps, ds *cas.Store, ks *keystore.Keystore) error {
	return al.Walk(func(app *apps.App) error {
		h, err := findImage(app.Image, app.Asc, ds, ks, true)
		if err != nil {
			return err
		}
		app.ImageID = *h
		return nil
	})
}

// findImage will recognize a ACI hash and use that, import a local file, use
// discovery or download an ACI directly.
func findImage(img string, asc string, ds *cas.Store, ks *keystore.Keystore, discover bool) (*types.Hash, error) {
	// check if it is a valid hash, if so let it pass through
	h, err := types.NewHash(img)
	if err == nil {
		fullKey, err := ds.ResolveKey(img)
		if err != nil {
			return nil, fmt.Errorf("could not resolve key: %v", err)
		}
		h, err = types.NewHash(fullKey)
		if err != nil {
			// should never happen
			panic(err)
		}
		return h, nil
	}

	// try fetching the image, potentially remotely
	key, err := fetchImage(img, asc, ds, ks, discover)
	if err != nil {
		return nil, err
	}
	h, err = types.NewHash(key)
	if err != nil {
		// should never happen
		panic(err)
	}

	return h, nil
}

// TODO(vc): move image fetching out of fetch.go and into here?
