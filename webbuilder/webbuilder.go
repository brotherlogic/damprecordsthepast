package webbuilder

import (
	"html/template"
	"os"

	pb "github.com/brotherlogic/damprecordsthepast/proto"
)

type MatchPage struct {
	MatchTitle string
	Matches    []*Match
}

type Match struct {
	Username   string
	Percentage int32
}

func BuildMatchPage(users []*pb.User, matcher *pb.Matcher) error {
	MatchPage := &MatchPage{MatchTitle: matcher.GetName(), Matches: make([]*Match, 0)}
	for _, user := range users {
		MatchPage.Matches = append(MatchPage.Matches, &Match{Username: user.GetName(), Percentage: 56})
	}

	template, err := template.ParseFiles("templates/complete.tmpl")
	if err != nil {
		return err
	}
	file, err := os.Create("public/complete.html")
	if err != nil {
		return err
	}
	defer file.Close()
	return template.Execute(file, MatchPage)
}
