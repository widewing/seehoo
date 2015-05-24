package util

import (
	"io"
	"os/exec"
	"bufio"
	"time"
	log "github.com/cihub/seelog"
)

var defaultLogStdout func(string) = func(text string) {
	log.Debug("stdout: " + text)
}
var defaultLogStderr func(string) = func(text string) {
	log.Error("stderr: " + text)
}

func ExecuteDefaultLogger(name string, args ...string) error {
	return Execute(name,args...)(defaultLogStdout,defaultLogStderr)
}

func Execute(name string,args ...string) func(func(string),func(string)) error {
	p := exec.Command(name,args...)
	return func (logStdout func(string),logStderr func(string)) error{
		stdoutPipe,_ := p.StdoutPipe()
		stderrPipe,_ := p.StderrPipe()
		if err:=p.Start();err!=nil{
			return err
		}
		log := func(pipe io.Reader,logger func(string)){
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				logger(scanner.Text())
			}
		}
		event := make(chan int)
		stdoutOk := false
		stderrOk := false
		if logStdout!=nil{
			go func(){
				log(stdoutPipe,logStdout)
				stdoutOk = true
				event<-1
			}()
		}
		if logStderr!=nil{
			go func(){
				log(stderrPipe,logStderr)
				stderrOk = true
				event<-2
			}()
		}
		for !(stdoutOk&&stderrOk){
			select {
				case <-event:
			}
		}
		return p.Wait()
	}
}

func Umount(busybox string,mount string) error {
	var err error = nil
	for i:=0;i<3;i++ {
		err = ExecuteDefaultLogger(busybox,"umount",mount)
		if err==nil { break }
		log.Error("umount error: "+err.Error())
		if i==2 { break }
		log.Error("Retrying in 500ms")
		time.Sleep(time.Duration(500)*time.Millisecond)
	}
	return err
}