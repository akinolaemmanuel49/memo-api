package repository

// Repositories encapsulates all available repositories for easy reuse.
type Repositories struct {
	Social SocialRepository
	Users  UserRepository
	File   FileRepository
	Memo   MemoRepository
}
