package model

// GitStatus represents the Git status of a directory
type GitStatus struct {
	Branch         string
	IsDirty        bool
	AheadBy        int
	BehindBy       int
	ModifiedFiles  []string
	UntrackedFiles []string
	LastCommit     string
	CommitMessage  string
}


