// Search service - micro service for gruzoperevozki-rf.com
//
// Copyright 2025 Ivanov Vladimir Vyacheslavovich aka Ronin1024. All rights reserved.

package logger

import "time"

func TimeStampInFile(buf *[]byte) {
	t := time.Now()
	year, month, day := t.Date()
	Itoa(buf, year, 4)
	*buf = append(*buf, '/')
	Itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	Itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	hour, min, sec := t.Clock()
	Itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	Itoa(buf, min, 2)
	*buf = append(*buf, ':')
	Itoa(buf, sec, 2)
	*buf = append(*buf, '.')
	Itoa(buf, t.Nanosecond()/1e3, 6)
	*buf = append(*buf, ' ')
}

func TimeStampInFilename(buf *[]byte) {
	t := time.Now()
	year, month, day := t.Date()

	Itoa(buf, year, 4)
	*buf = append(*buf, '.')
	Itoa(buf, int(month), 2)
	*buf = append(*buf, '.')
	Itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	hour, min, sec := t.Clock()
	Itoa(buf, hour, 2)
	*buf = append(*buf, '_')
	Itoa(buf, min, 2)
	*buf = append(*buf, '_')
	Itoa(buf, sec, 2)
}

func GetStringTimeStampInFile() string {
	var nowTime []byte
	TimeStampInFile(&nowTime)
	return string(nowTime)
}

func GetStringTimeStampInFilename() string {
	var nowTime []byte
	TimeStampInFilename(&nowTime)
	return string(nowTime)
}
