package util

import (
	"io"
	"os/exec"
	"bufio"
	log "github.com/cihub/seelog"
)

var LogStdout func(string) = func(text string) {
	log.Debug("stdout: " + text)
}
var LogStderr func(string) = func(text string) {
	log.Error("stderr: " + text)
}

func LogStream(stream io.Reader,logFunc func(string)){
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		logFunc(scanner.Text())
	}
}

func ExecuteLogger(name string, args ...string) error {
	return Execute(name,args...)(LogStdout,LogStderr)
}

func Execute(name string,args ...string) func(func(string),func(string)) error {
	p := exec.Command(name,args...)
	return func (logStdout func(string),logStderr func(string)) error{
		stdoutPipe,_ := p.StdoutPipe()
		stderrPipe,_ := p.StderrPipe()
		if err:=p.Start();err!=nil{
			return err
		}
		Wait(
			func(){LogStream(stdoutPipe,logStdout)},
			func(){LogStream(stderrPipe,logStderr)},
		)
		return p.Wait()
	}
}

func Umount(mount string) error {
	return ExecuteLogger("/bin/busybox","umount","-l",mount)
}

func WaitLast(funcs ...func()) {
	n := len(funcs)
	waitChan := make(chan int, 1)
	for i,f := range funcs {
		go func(){
			f()
			waitChan <- i
		}()
	}
	for {
		select {
			case <-waitChan: 
				n-=1
		}
		if n==0 { break }
	}
}

func WaitFirst(funcs ...func()) {
	waitChan := make(chan int, 1)
	for i,f := range funcs {
		go func(){
			f()
			waitChan <- i
		}()
	}
	<-waitChan
}

var Wait = WaitLast