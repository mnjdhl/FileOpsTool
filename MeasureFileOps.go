/*
 * Author: Manoj Dahal
 */
package main

import (
	"fmt"
	"os"
	"sync"
)

func ProcessAllOpts(Args []string, startIndx int, pwg *sync.WaitGroup) {
	var arg string

	if pwg != nil {
		defer pwg.Done()
	}

	if len(Args) > 1 {
		arg = Args[startIndx]
	} else {
		arg = "default"
	}

	switch arg {

	case "h":
		ShowHelp()
	case "Help":
		ShowHelp()
	case "help":
		ShowHelp()
	case "ls":
		//TBD
		fmt.Println("TBD")
	case "rn":
		//TBD
		fmt.Println("TBD")
	case "crn": //Create the dir and then Rename
		//TBD
		fmt.Println("TBD")
	case "op":

		//TBD
		fmt.Println("TBD")

	case "io":

		MeasureFileIO(Args, startIndx)

	case "fio": //File I/O Workload
		//TBD
		fmt.Println("TBD")
	case "tnc":
		//TBD
		fmt.Println("TBD")

	case "tns":
		//TBD
		fmt.Println("TBD")

	default:
		ShowHelp()
	}
}

func ShowHelp() {
	fmt.Println("***************************FILE OPERATION HELP****************************************")
	fmt.Println("1. List dir/files: ls <dir path>")
	fmt.Println("2. Rename dir/file continously in a loop: rn <dir path> <loop count>")
	fmt.Println("3. Create & Rename dir/file continously in a loop: crn <dir path> <loop count>")
	fmt.Println("4. Open a dir/file: op <dir path>")
	fmt.Println("5. Write, Read & Remove file continously in a loop: io <dir path/filename pattern> <line count to read/write> <loop count> <test duration in seconds>")
	fmt.Println("6. File I/O Workload: fio <file Path> <line Count> <file count> <dir Path> <dir count>")
	fmt.Println("7. Test Net Client: tnc <IPv4 Addr> <tcp port>")
	fmt.Println("8. Test Net Server: tns <tcp port>") // <comand-line input string one of 1-5>")
	fmt.Println("9. Help or help or h: This Help Screen")
	fmt.Println("********************************************************************************")
}

func init() {
	//TBD
	fmt.Println("Starting File OPS Measurement Tool")
}

func main() {
	ProcessAllOpts(os.Args, 1, nil)
}
