package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

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

			// get pull requests from GitHub
			ctx := context.Background()

			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)

			client := github.NewClient(tc)
			prs := getAllPrs(client, org, repName)

			// ranking
			sort.Slice(prs, func(i, j int) bool { return prs[i].GetAdditions() > prs[j].GetAdditions() })
			addTop10 := append([]*github.PullRequest{}, prs...)
			if len(prs) >= numberOfRanking {
				addTop10 = append([]*github.PullRequest{}, prs[0:numberOfRanking]...)
			}
			sort.Slice(prs, func(i, j int) bool { return prs[i].GetDeletions() > prs[j].GetDeletions() })
			delTop10 := append([]*github.PullRequest{}, prs...)
			if len(prs) >= numberOfRanking {
				delTop10 = append([]*github.PullRequest{}, prs[0:numberOfRanking]...)
			}
			sort.Slice(prs, func(i, j int) bool {
				return prs[i].GetDeletions()+prs[i].GetAdditions() > prs[j].GetDeletions()+prs[j].GetAdditions()
			})
			addAndDelTop10 := append([]*github.PullRequest{}, prs...)
			if len(prs) >= numberOfRanking {
				addAndDelTop10 = append([]*github.PullRequest{}, prs[0:numberOfRanking]...)
			}

			// author ranking
			var byAuthorsAddMap map[string]int = map[string]int{}
			var byAuthorsDelMap map[string]int = map[string]int{}
			var byAuthorsAddAndDelMap map[string]int = map[string]int{}
			var byAuthorPrCount map[string]int = map[string]int{}
			for _, pr := range prs {
				surl := strings.Split(pr.GetUser().GetURL(), "/")
				name := surl[len(surl)-1]
				byAuthorsAddMap[name] = byAuthorsAddMap[name] + pr.GetAdditions()
				byAuthorsDelMap[name] = byAuthorsDelMap[name] + pr.GetDeletions()
				byAuthorsAddAndDelMap[name] = byAuthorsAddAndDelMap[name] + pr.GetAdditions() + pr.GetDeletions()
				byAuthorPrCount[name] = byAuthorPrCount[name] + 1
			}

			for name, count := range byAuthorPrCount {
				byAuthorsAddMap[name] = byAuthorsAddMap[name] / count
				byAuthorsDelMap[name] = byAuthorsDelMap[name] / count
				byAuthorsAddAndDelMap[name] = byAuthorsAddAndDelMap[name] / count
			}

			var byAuthorsAddTop10 []userAndCount
			for k, v := range byAuthorsAddMap {
				byAuthorsAddTop10 = append(byAuthorsAddTop10, userAndCount{k, v})
			}
			sort.Slice(byAuthorsAddTop10, func(i, j int) bool { return byAuthorsAddTop10[i].count > byAuthorsAddTop10[j].count })
			if len(byAuthorsAddTop10) >= numberOfRanking {
				byAuthorsAddTop10 = byAuthorsAddTop10[0:numberOfRanking]
			}

			var byAuthorsDelTop10 []userAndCount
			for k, v := range byAuthorsDelMap {
				byAuthorsDelTop10 = append(byAuthorsDelTop10, userAndCount{k, v})
			}
			sort.Slice(byAuthorsDelTop10, func(i, j int) bool { return byAuthorsDelTop10[i].count > byAuthorsDelTop10[j].count })
			if len(byAuthorsDelTop10) >= numberOfRanking {
				byAuthorsDelTop10 = byAuthorsDelTop10[0:numberOfRanking]
			}

			var byAuthorsAddAndDelTop10 []userAndCount
			for k, v := range byAuthorsAddAndDelMap {
				byAuthorsAddAndDelTop10 = append(byAuthorsAddAndDelTop10, userAndCount{k, v})
			}
			sort.Slice(byAuthorsAddAndDelTop10, func(i, j int) bool { return byAuthorsAddAndDelTop10[i].count > byAuthorsAddAndDelTop10[j].count })
			if len(byAuthorsAddAndDelTop10) >= numberOfRanking {
				byAuthorsAddAndDelTop10 = byAuthorsAddAndDelTop10[0:numberOfRanking]
			}

			printPrsRanking(addTop10, "👑 Additions Top"+strconv.Itoa(numberOfRanking))
			printPrsRanking(delTop10, "👑 Deletions Top"+strconv.Itoa(numberOfRanking))
			printPrsRanking(addAndDelTop10, "👑 Additions+Deletions Top"+strconv.Itoa(numberOfRanking))

			printAuthorsRanking(byAuthorsAddTop10, byAuthorsAddMap, byAuthorsDelMap, "👑 Additions average by User Top"+strconv.Itoa(numberOfRanking))
			printAuthorsRanking(byAuthorsDelTop10, byAuthorsAddMap, byAuthorsDelMap, "👑 AdditioDeletionsns average by User Top"+strconv.Itoa(numberOfRanking))
			printAuthorsRanking(byAuthorsAddAndDelTop10, byAuthorsAddMap, byAuthorsDelMap, "👑 Additions+Deletions average by User Top"+strconv.Itoa(numberOfRanking))
		},
	}
	token           string
	numberOfRanking int
)

func printPrsRanking(prs []*github.PullRequest, label string) {
	fmt.Println("=================================================")
	fmt.Println(label)
	for i, pr := range prs {
		if i == 0 {
			fmt.Println("=================================================")
			fmt.Printf("%v. %v\n", i+1, pr.GetTitle())
		} else {
			fmt.Println("-------------------------------------------------")
			fmt.Printf("%v. %v\n", i+1, pr.GetTitle())
		}
		fmt.Printf("Additions: %v Deletions: %v\n", pr.GetAdditions(), pr.GetDeletions())
		fmt.Printf("%v\n", pr.GetHTMLURL())
	}

	fmt.Println()
}

func printAuthorsRanking(ucs []userAndCount, addMap map[string]int, delMap map[string]int, label string) {
	fmt.Println("=================================================")
	fmt.Println(label)
	for i, uc := range ucs {
		if i == 0 {
			fmt.Println("=================================================")
			fmt.Printf("👑%v. %v\n", i+1, uc.name)
		} else {
			fmt.Println("-------------------------------------------------")
			fmt.Printf("%v. %v\n", i+1, uc.name)
		}
		fmt.Printf("Additions(Average): %v Deletions(Average): %v\n", addMap[uc.name], delMap[uc.name])
	}

	fmt.Println()
}

type userAndCount struct {
	name  string
	count int
}

func init() {
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Your GitHub token.")
	rootCmd.MarkPersistentFlagRequired("token")
	rootCmd.PersistentFlags().IntVar(&numberOfRanking, "numberOfRanking", 3, "Number of ranking.")
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

	ch := make(chan *github.PullRequest, len(allPrsWithoutDetails))
	var allPrs []*github.PullRequest
	semaphore := make(chan int, 8)
	var wg sync.WaitGroup

	for _, pr := range allPrsWithoutDetails {
		wg.Add(1)
		go func(p *github.PullRequest) {
			defer wg.Done()
			semaphore <- 1

			prd, _, err := client.PullRequests.Get(ctx, owner, repository, p.GetNumber())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			ch <- prd
			bar.Increment()
			<-semaphore
		}(pr)
	}

	for i := 0; i < len(allPrsWithoutDetails); i++ {
		allPrs = append(allPrs, <-ch)
	}
	close(ch)

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
