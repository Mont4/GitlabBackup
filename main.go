package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

const groupApiUrl = "https://gitlab.com/api/v4/groups?per_page=999"
const projectApiUrl = "https://gitlab.com/api/v4/groups/%d/projects?per_page=999"

type group struct {
	Id      int
	Name    string
	Web_url string
	Path    string
}
type project struct {
	Id                  int
	Name                string
	Web_url             string
	Ssh_url_to_repo     string
	Http_url_to_repo    string
	Name_with_namespace string
	Path_with_namespace string
}

func main() {
	var token string

	flag.StringVar(&token, "token", "", "Gitlab api token from \"https://gitlab.com/profile/personal_access_tokens\"")
	flag.Parse()

	if token == "" {
		fmt.Println("Set token from https://gitlab.com/profile/personal_access_tokens and call again")
		os.Exit(1)
	}

	groups := getGroups(token)

	projects := getAllProjects(groups, token)

	cloneAll(projects)
}

func getGroups(token string) []group {
	body := request(groupApiUrl, token)

	var groups []group
	json.Unmarshal(body, &groups)

	return groups
}

func getAllProjects(groups []group, token string) []project {
	progressive := NewProgressive(len(groups))
	progressive.Start()

	var projects []project
	for _, group := range groups {
		os.MkdirAll(getPath(group.Path), 0777)

		newProjects := getProject(group.Id, token)

		projects = append(projects, newProjects...)

		progressive.Advance()
	}
	progressive.Done()

	fmt.Printf("\n\n List of Projets \n\n")
	for i, p := range projects {
		fmt.Printf("%3d) %s \n", i+1, p.Name_with_namespace)
	}

	return projects
}

func getProject(groupId int, token string) []project {
	url := fmt.Sprintf(projectApiUrl, groupId)
	body := request(url, token)

	var newProjects []project
	json.Unmarshal(body, &newProjects)

	return newProjects
}

func request(url string, token string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body
}

func cloneAll(projects []project) {
	progressive := NewProgressive(len(projects))
	progressive.Start()

	for _, project := range projects {
		if _, err := os.Stat(getPath(project.Path_with_namespace)); os.IsNotExist(err) {
			clone(project)
		}

		progressive.Advance()
	}
	progressive.Done()

	fmt.Println("\n\n Finished.")
}

func clone(p project) {
	cmd := exec.Command("git", "clone", p.Ssh_url_to_repo, getPath(p.Path_with_namespace))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func getPath(path string) string {
	return "./GitlabBackup/" + path
}
