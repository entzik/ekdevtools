package gitlab

// RepositoryDescriptor a structure that contains information about a git repository, as returned by the Gitlab search API
type RepositoryDescriptor struct {
	// Name the name of the repository
	Name string `json:"name"`
	// Description the repository description
	Description string `json:"description"`
	// HttpUrlToRepo the repository URL, used for cloning
	HttpUrlToRepo string `json:"http_url_to_repo"`
	// RepositoryNamespace the namespace (group)to which the repository belongs
	Namespace RepositoryNamespace `json:"namespace"`
}

// RepositoryNamespace a repository namespace descriptor - can be a group.
type RepositoryNamespace struct {
	// Name the namespace name
	Name string `json:"name"`
	// Path gitlab path to the namespace
	Path string `json:"path"`
	// Kindthe namespace type (ex: group_
	Kind string `json:"kind"`
}
