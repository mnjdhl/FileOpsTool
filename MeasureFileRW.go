/*
 * Author: Manoj Dahal
 */
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

func MeasureReadWriteOps(fname string, totLines int, wg *sync.WaitGroup, rch chan uint64, wch chan uint64, sw *sync.WaitGroup, sr *sync.WaitGroup, sd *sync.WaitGroup) {

	defer wg.Done()
	wbytes := uint64(0)
	rbytes := uint64(0)

	dataBytes := []byte("Namaste everyone!\nNow I am writing 1000s of bytes to a file on this file system.\nAnd later I am going to read them. So that I can do the benchmarking for Write & Read bytes per seconds...\n")
	readBytes := make([]byte, 107)
	flist := []string{}
	fplist := []*os.File{}
	fcnt := 0
	for fcnt < 100 {
		newFile := fname + strconv.Itoa(fcnt)
		//fmt.Println("Creating the file ", newFile)
		f, err := os.Create(newFile)
		if err != nil {
			fmt.Println("Error in Writing to the file ", newFile, "  Error:", err)
			return
		}
		flist = append(flist, newFile)
		fplist = append(fplist, f)

		fcnt++
	}

	sw.Wait()
	fcnt = 0
	offset := 0
	for _, fp := range fplist {

		//fmt.Println("Writing ", totLines, " lines to the file ", flist[i])

		for i := 1; i <= totLines; i++ {
			nwsz, _ := fp.WriteAt(dataBytes, int64(offset))
			offset = offset + nwsz
		}

		fp.Sync()

		//Sum total number of bytes written
		wbytes = wbytes + uint64(offset)
	}

	wch <- wbytes
	offset = 0
	sr.Wait()
	for _, fp := range fplist {
		//fmt.Println("Reading ", totLines, " lines from the file ", flist[i])
		for i := 1; i <= totLines; i++ {
			nsz, err2 := fp.ReadAt(readBytes, int64(offset))
			offset = offset + nsz
			if err2 == io.EOF {
				break
			}
			if err2 != nil {
				fmt.Println("Error in reading from the file ", flist[i], "  Error:", err2)
				break
			}

			if nsz == 0 {
				break
			}
		}
		fp.Close()

		//Sum total number of bytes read
		rbytes = rbytes + uint64(offset)
	}
	rch <- rbytes
	sd.Wait()
	for _, fl := range flist {
		//fmt.Println("Deleting the file ", fl)
		os.Remove(fl)
	}
}

func MeasureFileIO(Args []string, startIndx int) {

	fileName := Args[startIndx+1]

	lineCnt, _ := strconv.Atoi(Args[startIndx+2])
	loopCnt, _ := strconv.Atoi(Args[startIndx+3])
	dur, _ := strconv.Atoi(Args[startIndx+4]) // in seconds
	duration := time.Duration(dur) * time.Second
	rchan := make(chan uint64)
	wchan := make(chan uint64)
	fmt.Println("Duration ", duration, " minutes")

	wg := new(sync.WaitGroup)
	syncW := new(sync.WaitGroup)
	syncR := new(sync.WaitGroup)
	syncD := new(sync.WaitGroup)

	elapsed := time.Duration(0)
	elapsedw := time.Duration(0)
	elapsedr := time.Duration(0)
	tr := uint64(0)
	tw := uint64(0)
	for elapsed < duration {
		syncW.Add(1)
		syncR.Add(1)
		syncD.Add(1)
		fmt.Println("Creating the files... ")
		for i := 1; i <= loopCnt; i++ {
			newFileName := fileName + "_" + strconv.Itoa(i) + ".txt"
			wg.Add(1)
			go MeasureReadWriteOps(newFileName, lineCnt*i, wg, rchan, wchan, syncW, syncR, syncD)
		}

		startTm := time.Now()
		fmt.Println("Writing to the files... ")
		syncW.Done()

		for i := 1; i <= loopCnt; i++ {
			tw = tw + <-wchan
		}

		endTm := time.Now()
		diff := endTm.Sub(startTm)
		elapsedw = elapsedw + diff
		elapsed = elapsed + elapsedw
		fmt.Println("elapsed = ", elapsed.Seconds(), "  duration=", duration)
		startTm = time.Now()
		fmt.Println("Reading from the files... ")
		syncR.Done()
		for i := 1; i <= loopCnt; i++ {
			tr = tr + <-rchan
		}
		endTm = time.Now()
		diff = endTm.Sub(startTm)
		elapsedr = elapsedr + diff
		elapsed = elapsed + elapsedr
		fmt.Println("elapsed = ", elapsed.Seconds(), "  duration=", duration)
		fmt.Println("Deleting the files... ")
		syncD.Done()
		wg.Wait()
	}
	rth := float64(tr)
	wth := float64(tw)
	if elapsedw.Seconds() > 0 {
		wth = float64(tw) / elapsedw.Seconds()
	}
	if elapsedr.Seconds() > 0 {
		rth = float64(tr) / elapsedr.Seconds()
	}
	fmt.Println("***********************************************")
	fmt.Println("Elapsed for Write ", elapsedw, " seconds since the beginning of test!")
	fmt.Println("Elapsed for Read ", elapsedr, " seconds since the beginning of test!")
	fmt.Printf("\nTotal Bytes Written = %d Write Throughput = %.2f bytes/sec", tw, wth)
	fmt.Printf("\nTotal Bytes Read = %d Read Throughput = %.2f bytes/sec", tr, rth)
	fmt.Println("***********************************************")

}
