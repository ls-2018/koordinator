//go:build !linux
// +build !linux

/*
Copyright 2022 The Koordinator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package over_pleg

import (
	"k8s.io/utils/inotify"
)

func NewWatcher() (Watcher, error) {
	return nil, errNotSupported
}

func TypeOf(event *inotify.Event) EventType {
	if event.Mask&IN_CREATE != 0 && event.Mask&IN_ISDIR != 0 {
		return DirCreated
	}
	if event.Mask&IN_DELETE != 0 && event.Mask&IN_ISDIR != 0 {
		return DirRemoved
	}
	return UnknownType
}