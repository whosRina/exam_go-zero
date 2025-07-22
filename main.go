package main

import (
	"log"
	"os/exec"
)

func main() {
	services := []struct {
		Name string
		Cmd  *exec.Cmd
	}{
		{
			Name: "user",
			Cmd:  exec.Command("go", "run", "users.go"),
		},
		{
			Name: "exam",
			Cmd:  exec.Command("go", "run", "exam.go"),
		},
		{
			Name: "class",
			Cmd:  exec.Command("go", "run", "class.go"),
		},
		{
			Name: "questionBank",
			Cmd:  exec.Command("go", "run", "question.go"),
		},
	}

	// 设置每个服务的工作目录
	services[0].Cmd.Dir = "./users"
	services[1].Cmd.Dir = "./exam"
	services[2].Cmd.Dir = "./classes"
	services[3].Cmd.Dir = "./questionBank"

	for _, service := range services {
		s := service // 避免闭包引用错误
		go func() {
			log.Printf("启动服务: %s\n", s.Name)
			s.Cmd.Stdout = log.Writer()
			s.Cmd.Stderr = log.Writer()
			err := s.Cmd.Run()
			if err != nil {
				log.Printf("服务 %s 启动失败: %v\n", s.Name, err)
			}
		}()
	}

	select {}
}
