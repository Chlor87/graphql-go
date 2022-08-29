package resolvers

import "github.com/Chlor87/graphql/domain"

type Root struct {
	*todoResolver
	*userResolver
}

func NewRoot(
	domain *domain.Domain,
) *Root {
	return &Root{
		&todoResolver{Domain: domain},
		&userResolver{Domain: domain},
	}
}
