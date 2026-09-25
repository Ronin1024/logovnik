// Search service - micro service for gruzoperevozki-rf.com
//
// Copyright 2025 Ivanov Vladimir Vyacheslavovich aka Ronin1024. All rights reserved.

package logger

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func Itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
