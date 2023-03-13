package webbuilder

import (
	"fmt"
	"html/template"
	"os"

	"github.com/brotherlogic/damprecordsthepast/core"
	pb "github.com/brotherlogic/damprecordsthepast/proto"
)

type MatchPage struct {
	MatchTitle string
	Matches    []*Match
}

type IndexPage struct {
	Matches []*Match
}

type Match struct {
	Username   string
	Percentage int32
	Matchname  string
	Shortname  string
}

func BuildMatchPage(users []*pb.User, matcher *pb.Matcher) error {
	MatchPage := &MatchPage{MatchTitle: matcher.GetName(), Matches: make([]*Match, 0)}
	for _, user := range users {
		match, _ := core.ComputeMatch(user, matcher)
		MatchPage.Matches = append(MatchPage.Matches, &Match{Username: user.GetName(), Percentage: int32(match), Shortname: matcher.GetSimpleName()})
	}

	template, err := template.ParseFiles("templates/complete.tmpl")
	if err != nil {
		return err
	}
	file, err := os.Create(fmt.Sprintf("public/%v.html", matcher.GetSimpleName()))
	if err != nil {
		return err
	}
	defer file.Close()
	return template.Execute(file, MatchPage)
}

type UserMatchPage struct {
	Username   string
	MatchTitle string
	Percentage float64
	Missing    []*pb.Release
}

func BuildUserMatchPage(user *pb.User, matcher *pb.Matcher) error {
	userMatchPage := &UserMatchPage{MatchTitle: matcher.GetName(), Missing: make([]*pb.Release, 0)}
	match, missing := core.ComputeMatch(user, matcher)

	userMatchPage.Percentage = match
	userMatchPage.Missing = missing

	template, err := template.ParseFiles("templates/user.tmpl")
	if err != nil {
		return err
	}
	file, err := os.Create(fmt.Sprintf("public/%v-%v.html", user.GetName(), matcher.GetSimpleName()))
	if err != nil {
		return err
	}
	defer file.Close()
	return template.Execute(file, userMatchPage)
}

func BuildIndexPage(users []*pb.User, matchers []*pb.Matcher) error {
	IndexPage := &IndexPage{Matches: make([]*Match, 0)}
	for _, match := range matchers {
		bestUser := ""
		bestPerc := float64(0)
		for _, user := range users {
			matchv, _ := core.ComputeMatch(user, match)
			if matchv > bestPerc {
				bestUser = user.Name
				bestPerc = matchv
			}
		}

		IndexPage.Matches = append(IndexPage.Matches, &Match{
			Matchname:  match.GetName(),
			Username:   bestUser,
			Percentage: int32(bestPerc),
			Shortname:  match.GetSimpleName(),
		})
	}

	template, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		return err
	}
	file, err := os.Create("public/index.html")
	if err != nil {
		return err
	}
	defer file.Close()
	return template.Execute(file, IndexPage)
}
