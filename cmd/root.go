package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/google/go-github/v33/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var rootCmd = &cobra.Command{
	Use:   "dkpr",
	Short: "dkpr finds dekai pull request.",
	Long:  `dkpr finds dekai pull request on GitHub. dekai pull request means pull request has too large diff.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "TODO"},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)
		prs := getAllPrs(client, "inabajumr", "pctts")

		sort.Slice(prs, func(i, j int) bool { return prs[i].GetAdditions() > prs[j].GetAdditions() })
		for _, pr := range prs {
			fmt.Println(pr.GetAdditions())
			fmt.Println(pr.GetURL())
		}

	},
}

func getAllPrs(client *github.Client, owner string, repository string) []*github.PullRequest {
	var allPrs []*github.PullRequest
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}

	ctx := context.Background()

	index := 0
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repository, opt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		opt.Page = resp.NextPage

		for _, pr := range prs {
			prd, _, err := client.PullRequests.Get(ctx, owner, repository, pr.GetNumber())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			allPrs = append(allPrs, prd)
			index++
			if index%10 == 0 {
				fmt.Println(index)
			}
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return allPrs
}

// Execute is just root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
