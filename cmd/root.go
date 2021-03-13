package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/google/go-github/v33/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dkpr",
		Short: "dkpr finds dekai pull request.",
		Long:  `dkpr finds dekai pull request on GitHub. dekai pull request means pull request has too large diff.`,
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 1 {
				fmt.Println("Requires repository name like 'inabajunmr/dkpr'.")
				os.Exit(1)
			}

			if len(strings.Split(args[0], "/")) != 2 {
				fmt.Println("Requires repository name like 'inabajunmr/dkpr'.")
				os.Exit(1)
			}

			org := strings.Split(args[0], "/")[0]
			repName := strings.Split(args[0], "/")[1]

			if len(org) == 0 || len(repName) == 0 {
				fmt.Println("Requires repository name like 'inabajunmr/dkpr'.")
				os.Exit(1)
			}

			ctx := context.Background()

			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)

			client := github.NewClient(tc)
			prs := getAllPrs(client, org, repName)

			sort.Slice(prs, func(i, j int) bool { return prs[i].GetAdditions() > prs[j].GetAdditions() })
			for _, pr := range prs {
				fmt.Println(pr.GetAdditions())
				fmt.Println(pr.GetURL())
			}

		},
	}
	token string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Your GitHub token.")
	rootCmd.MarkPersistentFlagRequired("token")
}

func getAllPrs(client *github.Client, owner string, repository string) []*github.PullRequest {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}

	ctx := context.Background()

	var allPrsWithoutDetails []*github.PullRequest
	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repository, opt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		allPrsWithoutDetails = append(allPrsWithoutDetails, prs...)
		opt.Page = resp.NextPage

		if resp.NextPage == 0 {
			break
		}
	}

	bar := pb.StartNew(len(allPrsWithoutDetails))

	var allPrs []*github.PullRequest
	for _, pr := range allPrsWithoutDetails {
		prd, _, err := client.PullRequests.Get(ctx, owner, repository, pr.GetNumber())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		allPrs = append(allPrs, prd)
		bar.Increment()
	}

	bar.Finish()
	return allPrs
}

// Execute is just root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
