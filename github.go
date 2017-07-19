package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func dumpToGithub() {
	for {
		time.Sleep(5 * time.Minute)
		r, err := getAll()
		if err != nil {
			fmt.Println("!!!", err)
			continue
		}
		fmt.Printf("Sending %d results to github\n", len(r))
		dat, _ := json.MarshalIndent(r, "", "  ")
		//fmt.Println(string(dat))

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GH_TOK")},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		gOpts := &github.RepositoryContentGetOptions{
			Ref: "master",
		}
		_, fInfo, _, err := client.Repositories.GetContents(context.Background(), "captncraig", "wildchef", "docs", gOpts)
		if err != nil {
			fmt.Println("!!!", err)
			continue
		}
		var sha string
		for _, f := range fInfo {
			if f.GetPath() == "docs/results.json" {
				sha = f.GetSHA()
				break
			}
		}
		msg := fmt.Sprintf("automated upload %d results", len(r))
		opts := &github.RepositoryContentFileOptions{
			Content: dat,
			SHA:     &sha,
			Message: &msg,
		}
		_, _, err = client.Repositories.UpdateFile(context.Background(), "captncraig", "wildchef", "docs/results.json", opts)
		fmt.Println(err)
	}
}
