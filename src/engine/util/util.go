package util

import (
	"io"
	"os/exec"
	"bufio"
	"time"
	"log"
)

var defaultLogStdout func(string) = func(text string) {
	log.Print("stdout: " + text)
}
var defaultLogStderr func(string) = func(text string) {
	log.Print("stderr: " + text)
}

func ExecuteDefaultLogger(name string, args ...string) error {
	return Execute(name,args...)(defaultLogStdout,defaultLogStderr)
}

func Execute(name string,args ...string) func(func(string),func(string)) error {
	p := exec.Command(name,args...)
	return func (logStdout func(string),logStderr func(string)) error{
		stdoutPipe,_ := p.StdoutPipe()
		stderrPipe,_ := p.StderrPipe()
		p.Start()
		log := func(pipe io.Reader,logger func(string)){
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				logger(scanner.Text())
			}
		}
		if logStdout!=nil{
			go log(stdoutPipe,logStdout)
		}
		if logStderr!=nil{
			go log(stderrPipe,logStderr)
		}
		return p.Wait()
	}
}

func Umount(busybox string,mount string) error {
	var err error = nil
	for i:=0;i<3;i++ {
		err = ExecuteDefaultLogger(busybox,"umount",mount)
		if err==nil { break }
		if i==2 { break }
		log.Printf("umount error: %s, retrying in 500ms",err.Error())
		time.Sleep(500*time.Millisecond)
	}
	return err
}