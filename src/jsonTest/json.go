package main

import (
	//	"fmt"
	"githubJson"
	//	"html/template"
	"log"
	"os"
	"text/template"
	"time"
)

const tmpl = `{{.TotalCount}} issues:
				{{range .Items}}----------------------------------------
				Number: {{.Number}}
				User:   {{.User.Login}}
				Title:  {{.Title | printf "%.64s"}}
				Age:    {{.CreateAt | daysAgo}} days
				{{end}}`

func main() {
	//	result, err := githubJson.SearchIssues(os.Args[1:])
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("%d issues:\n", result.TotalCount)
	//	for _, item := range result.Items {
	//		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	//	}
	report, err := template.New("report").Funcs(template.FuncMap{"daysAgo": daysAgo}).Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	result, err := githubJson.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

//import (
//	"encoding/json"
//	"fmt"
//)

//type Movie struct {
//	Title  string
//	Year   int  `json:"released"`
//	Color  bool `json:"color,omitempty"` //omitempty选项，表示当Go语言结构体成员为空或零值时不生成JSON对象（这里false为零值）
//	Actors []string
//}

//func main() {
//	var movies = []Movie{
//		{Title: "Casablanca", Year: 1942, Color: false,
//			Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
//		{Title: "Cool Hand Luke", Year: 1967, Color: true,
//			Actors: []string{"Paul Newman"}},
//		{Title: "Bullitt", Year: 1968, Color: true,
//			Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
//	}

//	//	data, err := json.Marshal(movies)
//	data, err := json.MarshalIndent(movies, "", "   ")

//	if err != nil {
//		fmt.Println("JSON marshaling failed: %s", err)
//	}
//	fmt.Printf("%s\n", data)

//	var titles []struct{ Title string }

//	if err := json.Unmarshal(data, &titles); err != nil {
//		fmt.Println("JSON Unmarshaling failed: %s", err)
//	}
//	fmt.Println(titles)
//}
