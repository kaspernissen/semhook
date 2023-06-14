package app

type Configuration struct {
	Port string
	Repo Repositories
	Auth Auth
}

type Repositories struct {
	RepoRoot    string
	GithubToken string
}

type Auth struct {
	Issuer   string
	Audience string
}


