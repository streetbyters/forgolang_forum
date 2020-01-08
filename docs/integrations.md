# Forgolang.com integrations
Integrated third party software/applications.

## Github
Integrated for authentication and developer information.\
Information received from users via github ([see](https://github.com/akdilsiz/forgolang_forum/blob/master/thirdparty/github/github.go)):
```go
// UserInformation github basic user information structure
type UserInformation struct {
	Bio         string `db:"bio" json:"bio"`
	PublicRepos int `db:"public_repos" json:"public_repos"`
	PublicGists int `db:"public_gists" json:"public_gists"`
	Followers   int `db:"followers" json:"followers"`
	Following   int `db:"following" json:"following"`
}
```