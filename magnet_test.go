// Copyright 2013 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package magnet

import (
	"fmt"
	"testing"
)

var TestMagnets = []string{
	"magnet:?xt=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C",
	"magnet:?xt=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C&dn=Great+Speeches+-+Martin+Luther+King+Jr.+-+I+Have+A+Dream.mp3",
	"magnet:?kt=martin+luther+king+mp3",
	"magnet:?xt.1=urn:sha1:YNCKHTQCWBTRNJIV4WNAE52SJUQCZO5C&xt.2=urn:sha1:TXGCZQTH26NL6OUQAJJPFALHG2LTGBC7",
	"magnet:?mt=http://weblog.foo/all-my-favorites.rss",
	"magnet:?xt=urn:btih:9480ac31b43e6219f2109c7877e48aeb47dfc7ac&dn=Of+Montreal+-+Lousy+with+Sylvianbriar+%282013%29+%5BFLAC%5D&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Ftracker.publicbt.com%3A80&tr=udp%3A%2F%2Ftracker.istole.it%3A6969&tr=udp%3A%2F%2Ftracker.ccc.de%3A80&tr=udp%3A%2F%2Fopen.demonii.com%3A1337",
}

func TestMain(t *testing.T) {
	//for _, v := range TestMagnets {
	//	m, err := NewMagnet(v)
	//	if err != nil {
	//		t.Error("NewMagnet() failed.")
	//	}
	//	fmt.Printf("%+v\n", m)
	//}
	m, err := NewMagnet(TestMagnets[5])
	if err != nil {
		t.Errorf("NewMagnet() failed: %v\n", err)
	}
	fmt.Printf("%+v\n", m)
}
