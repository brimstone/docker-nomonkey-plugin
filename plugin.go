package main

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/docker/go-plugins-helpers/authorization"
)

func newPlugin() (*nomonkey, error) {
	/*
			This must be setup to use the docker client later
			client, err := dockerclient.NewEnvClient()
		   	if err != nil {
		   		return nil, err
		   	}
		   	return &nomonkey{client: client}, nil
	*/
	return &nomonkey{}, nil
}

var (
	startRegExp = regexp.MustCompile(`/containers/create$`)
)

type nomonkey struct {
	// May not be needed: client *dockerclient.Client
}
type request struct {
	HostConfig struct {
		Binds   []string `json:"Binds"`
		CapAdd  []string `json:"CapAdd"`
		Devices []struct {
			PathOnHost        string `json:"PathOnHost"`
			PathInContainer   string `json:"PathInContainer"`
			CgroupPermissions string `json:"CgroupPermissions"`
		} `json:"Devices"`
		Privileged bool `json:"Privileged"`
	} `json:"HostConfig"`
	VolumesFrom []string
}

func (p *nomonkey) AuthZReq(req authorization.Request) authorization.Response {
	// If this isn't a post, allow it
	if req.RequestMethod != "POST" {
		return authorization.Response{Allow: true}
	}
	// If this isn't to our url, allow it
	if !startRegExp.MatchString(req.RequestURI) {
		return authorization.Response{Allow: true}
	}
	/*
		If you want to go as far as inspecting the image the container wants to
		use, you'll have to expose it in the request struct.
		// setup context for usage later
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		image, _, err := p.client.ImageInspectWithRaw(ctx, container.Image, false)
	*/

	if req.RequestBody == nil {
		// Not sure you'd get here, but ok
		return authorization.Response{Allow: true}
	}

	// Marshell our json request into a struct
	container := &request{}
	if err := json.NewDecoder(bytes.NewReader(req.RequestBody)).Decode(container); err != nil {
		return authorization.Response{Err: err.Error()}
	}

	log.Println(string(req.RequestBody))
	log.Printf("Container config is: %#v\n", container)

	// Actually test the volume binds
	if len(container.HostConfig.Binds) > 0 {
		for _, b := range container.HostConfig.Binds {
			paths := strings.Split(b, ":")
			if !strings.HasPrefix(paths[0], "/home/") {
				return authorization.Response{Msg: paths[0] + " is not a whitelisted volume."}
			}
		}
	}

	// Test for privileged flag
	if container.HostConfig.Privileged {
		return authorization.Response{Msg: "privileged mode is not allowed"}
	}

	// Test for devices
	if len(container.HostConfig.Devices) > 0 {
		return authorization.Response{Msg: "devices are not allowed"}
	}

	// Test for capadd
	if len(container.HostConfig.CapAdd) > 0 {
		return authorization.Response{Msg: "adding capabilities is not allowed"}
	}

	// Since we got this far, everything is allowed
	return authorization.Response{Allow: true}

}

func (p *nomonkey) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}
