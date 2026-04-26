package gvfs


// preparation for moving to goroutine
type Job struct {
    FilePath     string
    SideRepoPath string
    URLKey       string
}
