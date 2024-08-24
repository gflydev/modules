package repository

// Repositories struct for collect all app repositories.
type Repositories struct {
	// Information Module
	*UserRepository
}

// Pool a repository pool to store all
var Pool = &Repositories{
	&UserRepository{},
}
