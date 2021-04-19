# kittyimg

`kittyimg` is a Go library that allows to display images in terminal emulators implementing [kitty's *terminal graphics protocol*](https://sw.kovidgoyal.net/kitty/graphics-protocol.html).

[![Travis-CI](https://api.travis-ci.org/dolmen-go/kittyimg.svg?branch=master)](https://travis-ci.org/dolmen-go/kittyimg)
[![Codecov](https://img.shields.io/codecov/c/github/dolmen-go/kittyimg/master.svg)](https://codecov.io/gh/dolmen-go/kittyimg/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dolmen-go/kittyimg)](https://goreportcard.com/report/github.com/dolmen-go/kittyimg)

## Status

A [basic API](https://pkg.go.dev/github.com/dolmen-go/kittyimg/) (`Fprint`, `Fprintln`) allows to display an image (loaded with stdlib's [image](https://golang.org/pkg/image/) package) at the cursor position.

## License

Copyright 2021 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
