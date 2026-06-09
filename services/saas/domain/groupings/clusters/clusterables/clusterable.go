package clusterables

import "github.com/google/uuid"

type clusterable struct {
	identifier  uuid.UUID
	clusterKind Kind
}

func (clusterable *clusterable) Identifier() uuid.UUID {
	return clusterable.identifier
}

func (clusterable *clusterable) ClusterKind() Kind {
	return clusterable.clusterKind
}
