package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func containerRun(dir string) {
	cli, err := client.NewEnvClient()
	f := fmt.Sprintf(string(dir))
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}
	fmt.Println(f)
	volumes := fmt.Sprintf("/Users/alessio/code/oniondrop/src/%v:/site/www/", f)
	ctx := context.Background()
	imageName := "docker.io/einyx/tor"
	hostConfig := &container.HostConfig{
		Binds: []string{volumes},
	}
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, hostConfig, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.ListenAndServe(":8080", nil)
}
