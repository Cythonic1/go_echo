package server

import (
	"database/sql"
	"net/http"

    _ "github.com/mattn/go-sqlite3" 
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var file string = "./db/projects.db";

type db_t struct{
    db *sql.DB;
}

const create string = `
  CREATE TABLE IF NOT EXISTS projects (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  project_name VARCHAR(255),
  project_desc VARCHAR(500),
  project_url VARCHAR(255)
);`

type Project struct {
	ID          int    `json:"id"`
	ProjectName string `json:"project_name"`
	ProjectDesc string `json:"project_desc"`
	ProjectURL  string `json:"project_url"`
}

func DBInit() (*db_t, error){
    db, err := sql.Open("sqlite3", file);
    if err != nil {
        return nil, err;
    }
    if _, err := db.Exec(create); err!= nil{
        return nil, err;
    }
    return &db_t{
        db: db,
    },nil;

}

func get_prjects(c echo.Context) error{
    connection ,err := DBInit();
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"Error": err.Error()});
    }

    data, err := connection.db.Query("select * from projects");
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"errro": err.Error()});
    }

    var projects []Project;
    for data.Next() {
        var project Project
        if err := data.Scan(&project.ID, &project.ProjectName, &project.ProjectDesc, &project.ProjectURL); err != nil {
            return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
        }
        projects = append(projects, project)
    }
    return c.JSON(http.StatusOK, echo.Map{"data": projects});
}

func post_projects(c echo.Context) error{
    connection ,err := DBInit();
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"Error": err.Error()});
    }

    var project Project;
    if err := c.Bind(&project); err != nil {

        return c.JSON(http.StatusBadRequest, echo.Map{"Error": err.Error()});
    }

    query := `insert into projects (project_name, project_desc, project_url) VALUES (?,?,?)`
    result, err := connection.db.Exec(query, project.ProjectName, project.ProjectDesc, project.ProjectURL);
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"Error": err.Error()});
    }
    id , err := result.LastInsertId();
    
    return c.JSON(http.StatusCreated, echo.Map{"added": id});
}


func ServerInit(){
    e := echo.New();
    e.Use(middleware.Logger());

    e.GET("/", get_prjects);
    e.POST("/add", post_projects);

    e.Logger.Fatal(e.Start(":1234"));


}
