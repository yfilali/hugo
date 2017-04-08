// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"path"

	"github.com/dchest/cssmin"
	"github.com/spf13/hugo/source"
)

func init() {
	RegisterHandler(new(cssHandler))
	RegisterHandler(new(defaultHandler))
}

type basicFileHandler Handle

func (h basicFileHandler) Read(f *source.File, s *Site) HandledResult {
	return HandledResult{file: f}
}

func (h basicFileHandler) PageConvert(*Page) HandledResult {
	return HandledResult{}
}

type defaultHandler struct{ basicFileHandler }

func (h defaultHandler) Extensions() []string { return []string{"*"} }
func (h defaultHandler) FileConvert(f *source.File, s *Site) HandledResult {
	err := s.publish(f.Path(), f.Contents)
	if err != nil {
		return HandledResult{err: err}
	}
	return HandledResult{file: f}
}

type cssHandler struct{ basicFileHandler }

func (h cssHandler) Extensions() []string { return []string{"css"} }
func (h cssHandler) FileConvert(f *source.File, s *Site) HandledResult {
	fBytes := f.Bytes()
	fPath := f.Path()
	x := cssmin.Minify(fBytes)

	hash := md5.New()
	hash.Write(fBytes)
	fileHash := hex.EncodeToString(hash.Sum([]byte{}))

	ext := path.Ext(fPath)
	hashedPath := fPath[0:len(fPath)-len(ext)] + "-" + fileHash + ext

	herr := s.publish(hashedPath, bytes.NewReader(x))
	if herr != nil {
		return HandledResult{err: herr}
	}

	err := s.publish(fPath, bytes.NewReader(x))
	if err != nil {
		return HandledResult{err: err}
	}
	return HandledResult{file: f}
}
