package handler

import (
	"bytes"
	aux "deployer/auxiliary"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

var auth_key = aux.GetFromConfig("constants.auth_key")

func SendResponse(w http.ResponseWriter, res string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`["` + res + `"]`))
}

func verifyAuth(w http.ResponseWriter, givenAuthKey string) bool {
	if givenAuthKey == auth_key {
		return true
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Unauthorized"}`))
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
	cmd := exec.Command("docker", "start", containerName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)
}

func StopADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	cmd := exec.Command("docker", "stop", containerName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)

}

func RemoveADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	cmd := exec.Command("docker", "rm", "-f", containerName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)

}

func PauseADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	cmd := exec.Command("docker", "pause", containerName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)

}

func UnPauseADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	cmd := exec.Command("docker", "unpause", containerName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)

}

func InspectADocker(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	cmd := exec.Command("docker", "inspect", containerName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)
}

func GetADockerLogs(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	containerName := r.FormValue("name")
	cmd := exec.Command("docker", "logs", containerName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)
}
func GetDockerPS(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	cmd := exec.Command("docker", "ps", "-a")
	res := "Nothing!"
	var out bytes.Buffer
	var err1 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err1
	err := cmd.Run()

	if err != nil {
		res = err1.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)
}

func GetDockerImages(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	cmd := exec.Command("docker", "images")
	res := "Nothing!"
	var out bytes.Buffer
	var err1 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err1
	err := cmd.Run()

	if err != nil {
		res = err1.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)
}
func RemoveAnDockerImage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	imageName := r.FormValue("image")
	cmd := exec.Command("docker", "rmi", "-f", imageName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)
}

func PullAnDockerImage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	//params := mux.Vars(r)
	imageName := r.FormValue("image")
	cmd := exec.Command("docker", "pull", imageName)
	var out bytes.Buffer
	var err0 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err0
	err := cmd.Run()

	res := ""
	if err != nil {
		res = err0.String()
	} else {
		res = out.String()
	}
	SendResponse(w, res)
}

func RunDockerContainer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !verifyAuth(w, r.FormValue("auth_key")) {
		return
	}
	container_name := r.FormValue("name")
	container_image := r.FormValue("image")
	container_exists := checkDockerExistence(container_name)
	resp := ""

	if !container_exists {
		fmt.Println("Container does not exist we have to create a new one")
		resp = createContainer(container_image, container_name, r.FormValue("network"),
			r.FormValue("port_ex"), r.FormValue("port_in"), r.FormValue("volume_ex"), r.FormValue("volume_ex"), r.FormValue("v_map"))
	} else {
		fmt.Println("Container already exists\n We have to kill first and create a new one")
		//var err error
		s_cmd := exec.Command("docker", "stop", container_name)
		s_cmd.Run()

		r_cmd := exec.Command("docker", "rm", "-f", container_name)
		r_cmd.Run()

		rmi_cmd := exec.Command("docker", "rmi", "-f", container_image)
		rmi_cmd.Run()

		pull_cmd := exec.Command("docker", "pull", container_image)
		pull_cmd.Run()

		resp = createContainer(container_image, container_name, r.FormValue("network"),
			r.FormValue("port_ex"), r.FormValue("port_in"), r.FormValue("volume_ex"), r.FormValue("volume_ex"), r.FormValue("v_map"))

	}

	SendResponse(w, resp)
}

/*
	func checkDockerImageExistence(name string) bool {
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
*/
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

/**
 */
func createContainer(image string, name string, network string, port_ex string, port_in string, volume_ex string, volume_in string, v_map string) string {

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

	res := ""
	if err != nil {
		res = err1.String()
	} else {
		res = "Docker " + out.String() + " Successfully started!"
	}

	return res
}
