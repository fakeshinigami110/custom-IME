// db.go
package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const dbFileName = ".fcitx5-projects.json"

func GetDBPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get home directory: %v", err)
    }
    dbPath :=filepath.Join(homeDir,dbFileName)
    if _,err := os.Stat(dbPath) ; os.IsNotExist(err) {
        fmt.Printf("creating new db file on %s\n" , dbPath)
        emptyDB := Projects{}
        data, err := json.Marshal(emptyDB)
        if err != nil {
            return "" , fmt.Errorf("something went worng %v" , err)
        }
        err = os.WriteFile(dbPath ,data , 0600)
        if err != nil {
            return "" , err
        }
    }else if err != nil {
        return "", fmt.Errorf("failed to check db file: %v", err)
    }
    return  dbPath ,nil 


}



func LoadDB() (Projects, error) {
    path, err := GetDBPath()
    if err != nil {
        return nil, fmt.Errorf("failed to get DB path: %v", err)
    }
    
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read DB file: %v", err)
    }
    
    if len(data) == 0 {
        return make(Projects), nil
    }
    
    var projects Projects
    if err := json.Unmarshal(data, &projects); err != nil {
        return nil, fmt.Errorf("failed to parse DB JSON: %v", err)
    }
    
    return projects, nil
}

func (db *Projects) Save() error {
    data, err := json.MarshalIndent(db, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal projects: %v", err)
    }
    
    path, err := GetDBPath()
    if err != nil {
        return fmt.Errorf("failed to get DB path: %v", err)
    }
    
    err = os.WriteFile(path, data, 0600)
    if err != nil {
        return fmt.Errorf("failed to write DB file: %v", err)
    }
    
    return nil
}

func (db *Projects) AddProject(project Config) error {
 if _ , exists :=  db.GetProject(project.ProjectName) ; exists {
        return fmt.Errorf("project with name '%s' already exists", project.ProjectName)
    }
    if project.ProjectName == "" || project.IMEName == "" {
        return fmt.Errorf("opration needs project name and ime name , please try again using create option ")
    }
    project.CreatedAt = time.Now()
    project.UpdatedAt = time.Now()

    (*db)[project.ProjectName] = project
    return db.Save()
}


func (db *Projects) UpdateProject(project Config) error {
    if _, exists := db.GetProject(project.ProjectName)  ; !exists {
        return fmt.Errorf("project '%s' does not exist - use AddProject to create new projects", project.ProjectName)
    }
    
    project.UpdatedAt = time.Now()
    (*db)[project.ProjectName] = project
    return db.Save()
}


func (db *Projects) DeleteProject(projectName string) error {
    if _ , exists := db.GetProject(projectName) ; !exists {
        return fmt.Errorf("projects '%s' does not exist" , projectName)
    }
    delete(*db , projectName)
    return db.Save()
}



func (db *Projects) ListProjects() []string {
    var projectLists[]string
    for projectName ,_:= range *db {
        projectLists = append(projectLists , projectName)
    }
    return projectLists
}

func (db *Projects) GetProject(projectName string) (*Config,bool) {
    if conf, exists := (*db)[projectName]; exists {
        return &conf, true
    }
    return nil, false
}