package handler

import (
	"bytes"
	aux "deployer/auxiliary"
	response_util "deployer/utils"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

var auth_key = aux.GetFromConfig("constants.auth_key")

func verifyAuth(w http.ResponseWriter, givenAuthKey string) bool {
	if givenAuthKey == auth_key {
		return true
	} else {
		response_util.SendUnauthorizedResponse(w, "Wrong auth key provided")
		return false
	}
}

func StartADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	res, err := DoDockerContainerAction(containerName, "start")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func StopADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	res, err := DoDockerContainerAction(containerName, "stop")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func RemoveADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	res, err := DoDockerContainerAction(containerName, "rm")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func PauseADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	res, err := DoDockerContainerAction(containerName, "pause")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func UnPauseADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	res, err := DoDockerContainerAction(containerName, "unpause")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func InspectADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	res, err := DoDockerContainerAction(containerName, "inspect")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func GetADockerLogs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	res, err := DoDockerContainerAction(containerName, "logs")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func GetDockerPS(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	cmd := exec.Command("docker", "ps", "-a")
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err0.String())
		return
	}
	response_util.SendOkResponse(w, out.String())
}

func GetDockerImages(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	cmd := exec.Command("docker", "images")
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err0.String())
		return
	}
	response_util.SendOkResponse(w, out.String())

}
func RemoveAnDockerImage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	imageName := r.FormValue("image")
	res, err := DoDockerImageAction(imageName, "rm")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func PullAnDockerImage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	imageName := r.FormValue("image")
	res, err := DoDockerImageAction(imageName, "pull")
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, err.Error())
		return
	}
	response_util.SendOkResponse(w, res)
}

func RunDockerContainer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	container_name := r.FormValue("name")
	container_image := r.FormValue("image")
	container_exists := checkDockerExistence(container_name)

	if !container_exists {
		fmt.Println("Container does not exist we have to create a new one")
		res, err := createContainer(container_image, container_name, r.FormValue("network"),
			r.FormValue("port_ex"), r.FormValue("port_in"), r.FormValue("volume_ex"), r.FormValue("volume_ex"), r.FormValue("v_map"))
		if err != nil {
			response_util.SendInternalServerErrorResponse(w, err.Error())
		} else {
			response_util.SendOkResponse(w, res)
		}
	} else {
		fmt.Println("Container already exists\n We have to kill first and create a new one")
		final_res := ""
		res1, err := DoDockerContainerAction(container_name, "stop")
		if err != nil {
			response_util.SendInternalServerErrorResponse(w, err.Error())
			return
		} else {
			final_res += res1
			res2, err := DoDockerContainerAction(container_name, "rm")
			if err != nil {
				response_util.SendInternalServerErrorResponse(w, err.Error())
				return
			} else {
				final_res += res2
				res3, err := DoDockerImageAction(container_image, "rm")
				if err != nil {
					response_util.SendInternalServerErrorResponse(w, err.Error())
					return
				} else {
					final_res += res3
					res4, err := DoDockerImageAction(container_image, "pull")
					if err != nil {
						response_util.SendInternalServerErrorResponse(w, err.Error())
						return
					} else {
						final_res += res4
					}
				}
			}
		}

		res5, err := createContainer(container_image, container_name, r.FormValue("network"),
			r.FormValue("port_ex"), r.FormValue("port_in"), r.FormValue("volume_ex"), r.FormValue("volume_ex"), r.FormValue("v_map"))

		if err != nil {
			response_util.SendInternalServerErrorResponse(w, err.Error())
		} else {
			final_res += res5
			response_util.SendOkResponse(w, final_res)
		}

	}

}

func checkDockerExistence(name string) bool {
	cmd := exec.Command("docker", "container", "logs", name)
	var out bytes.Buffer
	var err1 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err1
	err := cmd.Run()

	if err != nil {
		return false
	} else {
		if strings.Contains(out.String(), "No such container") {
			return false
		} else {
			return true
		}
	}
}

func DoDockerContainerAction(container string, action string) (string, error) {

	var out bytes.Buffer
	var err0 bytes.Buffer
	var cmd *exec.Cmd

	switch action {
	case "start":
		cmd = exec.Command("docker", "start", container)
	case "stop":
		cmd = exec.Command("docker", "stop", container)
	case "rm":
		cmd = exec.Command("docker", "rm", "-f", container)
	case "pause":
		cmd = exec.Command("docker", "pause", container)
	case "unpause":
		cmd = exec.Command("docker", "unpause", container)
	case "inspect":
		cmd = exec.Command("docker", "inspect", container)
	case "logs":
		cmd = exec.Command("docker", "logs", container)
	default:
		return "", errors.New("unknown action")
	}
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()
	if err != nil {
		return "", err
	} else {
		return out.String(), nil
	}
}

func DoDockerImageAction(image string, action string) (string, error) {

	var out bytes.Buffer
	var err0 bytes.Buffer
	var cmd *exec.Cmd

	switch action {
	case "rm":
		cmd = exec.Command("docker", "rmi", "-f", image)
	case "pull":
		cmd = exec.Command("docker", "pull", image)
	default:
		return "", errors.New("unknown action")
	}
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()
	if err != nil {
		return "", err
	} else {
		return out.String(), nil
	}
}

func createContainer(image string, name string, network string, port_ex string, port_in string, volume_ex string, volume_in string, v_map string) (string, error) {

	volume_mapping := volume_ex + ":" + volume_in
	port_mapping := port_ex + ":" + port_in
	var cmd *exec.Cmd
	if network == "host" {
		if v_map == "yes" {
			cmd = exec.Command("docker", "run", "--network", network, "--name", name, "-v", volume_mapping, "-d", image)
		} else {
			cmd = exec.Command("docker", "run", "--network", network, "--name", name, "-d", image)
		}

	} else {
		if v_map == "yes" {
			cmd = exec.Command("docker", "run", "-p", port_mapping, "--name", name, "-v", volume_mapping, "-d", image)
		} else {
			cmd = exec.Command("docker", "run", "-p", port_mapping, "--name", name, "-d", image)
		}

	}
	var out bytes.Buffer
	var err1 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err1
	err := cmd.Run()

	if err != nil {
		return "", err
	}
	return "Docker " + out.String() + " Successfully started!", nil
}
