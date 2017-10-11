/*
MIT License

Copyright (c) 2017 Vishal Pandya

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package vlog 

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//Log Environment/Infrastructure model
type logEnv struct {
	logLevel       string
	logPath        string
	logFile        string
	logName        string
	logFileCounter int
}

//Globals
var (
	logConfig      logEnv
	logWritten *os.File
	loggerFile     *log.Logger
	logMutex       sync.RWMutex = sync.RWMutex{}
)

/*
Check existing files and find the latest log
to append based on file sizes. It's assumed that users
have not tinkered/edited these files in anyway.
They are not supposed to.
*/

func getLog(lf string) string {
	f := ""
	for {
		f = lf + strconv.Itoa(logConfig.logFileCounter) + ".log"
		if fInfo, err := os.Stat(f); err == nil {
			if fInfo.Size() >= math.MaxInt32 {
				logConfig.logFileCounter += 1
			} else {
				//filesize is less than the limit
				return f
			}
		} else {
			//file doesn't exist,just return the name needed
			return f
		}
	}
}

/*
Create a new log file. Called from Set/actualLog function
*/
func createLog(lpath string, level string, lname string) {
	l := ""
	if level != "DEBUG" {
		l = lpath + lname + "_"
	} else {
		l = lpath + lname + "_debug_"
	}
	//Get last updated
	l = getLog(l)
	//Create or append to existing file
	logWritten, err := os.OpenFile(l, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Unable to open/create log file:", err)
		os.Exit(-1)
	}
	//All is well, so increment our file counter
	logConfig.logPath = lpath
	logConfig.logFileCounter += 1
	logConfig.logFile = logWritten.Name()
	loggerFile = log.New(logWritten, "", 0)
}

/*
This function is called from the init function in the main file
after we read in our environment variables.
It accepts the path where the log files will be created and the loglevel.
*/
func Set(lpath string, level string, lname string) {
	//We don't want conflicting calls to this section
	logMutex.Lock()
	defer logMutex.Unlock()
	logConfig.logLevel = level
	logConfig.logFileCounter = 0
	logConfig.logName = lname
	createLog(lpath, level, lname)
}

/*
Returns a line of information of the caller.
Available only for DEBUG level.
*/
func getCaller() string {
	//skip 3 stack frames
	pc, f, lnum, ok := runtime.Caller(3)
	if ok {
		//Get details
		if d := runtime.FuncForPC(pc); d != nil {
			dir, file := filepath.Split(f)
			//grab the required directory name
			if dir != "" {
				dir = dir[:len(dir)-1]
				_, pkg := filepath.Split(dir)
				return pkg + "/" + file + ":" + d.Name() + ":" + strconv.Itoa(lnum)
			}
		}
	}
	return ""
}

/*
This is the actual function called by both Info and Debug functions.
The idea we use here is to wrap the log file once it reaches a limit
of currently around 2GB. The log files are named sequentially.
*/
func actualLog(level string, a ...interface{}) {

	now := time.Now()
	//Lock the code with a Mutex.
	logMutex.Lock()
	defer logMutex.Unlock()
	//Check the size of the file
	if logFileInfo, err := os.Stat(logConfig.logFile); err == nil {
		if logFileInfo.Size() >= math.MaxInt32 {
			logWritten.Close()
			createLog(logConfig.logPath, level, logConfig.logName)
		}
	}
	line := ""
	pmsg := ""
	rmsg := fmt.Sprintln(a...)
	if level == "DEBUG" {
		pmsg = getCaller()
		line = fmt.Sprintf("%s|%s|%s|%s", now.Format(time.RFC1123), level, pmsg, rmsg)
	} else {
		line = fmt.Sprintf("%s|%s|%s", now.Format(time.RFC1123), level, rmsg)
	}
	
	if loggerFile != nil {
		loggerFile.Print(line)
	}
	//Also write to Stdout
	fmt.Print(line)
}

//
func Debug(a ...interface{}) {
	if logConfig.logLevel == "DEBUG" {
		actualLog("DEBUG", a...)
	}

}
func Info(a ...interface{}) {
	if logConfig.logLevel != "DEBUG" {
		actualLog("INFO", a...)
	}

}

