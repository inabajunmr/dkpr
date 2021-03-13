# dkpr

Show dekai Pull Request ranking from GitHub repository.

## ⚠Caution

This tool doesn't consider [GitHub API Rate limit](https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting).

When you use this tool to too much repository or repository with too much Pull Requests, GitHub will return 403 and further requests will restricted.

## Usage

### Install

```
go install inabajunmr/dkpr
```

### Show Ranking
Issue Personal access tokens from https://github.com/settings/tokens.

```
dkpr [Repository name] --token [Your access token]
```

If you want to show inabajunmr/dkpr ranking.

```
dkpr inabajunmr/dkpr --token xxxxx
```

Default ranking will be shown until top 3.
If you want to get top 5, you can do like following command.

```
dkpr inabajunmr/dkpr --token xxxxx --numberOfRanking 
```
### Sample

```
$ dkpr spring-projects/spring-security --token xxxxx
1586 / 1586 [----------------------------------------------------------------------------------------------------------------------------------------------] 100.00% 15 p/s
=================================================
👑 Additions Top3
=================================================
1. Fix typo 'authenticated' -> 'isAuthenticated' in code comment
Additions: 97451 Deletions: 29148
https://github.com/spring-projects/spring-security/pull/7579
-------------------------------------------------
2. one bug fix cherry up from master
Additions: 90048 Deletions: 41586
https://github.com/spring-projects/spring-security/pull/5365
-------------------------------------------------
3. Add spring-javaformat checkstyle and formatting 
Additions: 74808 Deletions: 92232
https://github.com/spring-projects/spring-security/pull/8946

=================================================
👑 Deletions Top3
=================================================
1. Add spring-javaformat checkstyle and formatting 
Additions: 74808 Deletions: 92232
https://github.com/spring-projects/spring-security/pull/8946
-------------------------------------------------
2. one bug fix cherry up from master
Additions: 90048 Deletions: 41586
https://github.com/spring-projects/spring-security/pull/5365
-------------------------------------------------
3. Fix typo 'authenticated' -> 'isAuthenticated' in code comment
Additions: 97451 Deletions: 29148
https://github.com/spring-projects/spring-security/pull/7579

=================================================
👑 Additions+Deletions Top3
=================================================
1. Add spring-javaformat checkstyle and formatting 
Additions: 74808 Deletions: 92232
https://github.com/spring-projects/spring-security/pull/8946
-------------------------------------------------
2. one bug fix cherry up from master
Additions: 90048 Deletions: 41586
https://github.com/spring-projects/spring-security/pull/5365
-------------------------------------------------
3. Fix typo 'authenticated' -> 'isAuthenticated' in code comment
Additions: 97451 Deletions: 29148
https://github.com/spring-projects/spring-security/pull/7579

=================================================
👑 Additions average by User Top3
=================================================
👑1. Choi-JinHwan
Additions(Average): 97451 Deletions(Average): 29148
-------------------------------------------------
2. philwebb
Additions(Average): 26014 Deletions(Average): 32166
-------------------------------------------------
3. json20080301
Additions(Average): 22529 Deletions(Average): 10397

=================================================
👑 AdditioDeletionsns average by User Top3
=================================================
👑1. philwebb
Additions(Average): 26014 Deletions(Average): 32166
-------------------------------------------------
2. Choi-JinHwan
Additions(Average): 97451 Deletions(Average): 29148
-------------------------------------------------
3. json20080301
Additions(Average): 22529 Deletions(Average): 10397

=================================================
👑 Additions+Deletions average by User Top3
=================================================
👑1. Choi-JinHwan
Additions(Average): 97451 Deletions(Average): 29148
-------------------------------------------------
2. philwebb
Additions(Average): 26014 Deletions(Average): 32166
-------------------------------------------------
3. json20080301
Additions(Average): 22529 Deletions(Average): 10397
```
