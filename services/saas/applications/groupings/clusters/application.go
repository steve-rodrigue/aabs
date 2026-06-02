package clusters

import (
	"github.com/google/uuid"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

const rebuildBatchSize = 1000
const candidateAmount = 250

type application struct {
	repository            domain_clusters.Repository
	detector              domain_clusters.Detector
	clusterableRepository clusterables.Repository
	candidateRepository   clusterables.CandidateRepository
}

func createApplication(
	repository domain_clusters.Repository,
	detector domain_clusters.Detector,
	clusterableRepository clusterables.Repository,
	candidateRepository clusterables.CandidateRepository,
) Application {
	return &application{
		repository:            repository,
		detector:              detector,
		clusterableRepository: clusterableRepository,
		candidateRepository:   candidateRepository,
	}
}

func (app *application) BuildForTarget(
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	return app.detector.Detect(target, members)
}

func (app *application) FindByID(
	id uuid.UUID,
) (domain_clusters.Cluster, error) {
	return app.repository.FindByID(id)
}

func (app *application) FindByTarget(
	target clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	return app.repository.FindByTarget(target.Identifier())
}

func (app *application) FindByMember(
	member clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	return app.repository.FindByMember(member.Identifier())
}

func (app *application) RebuildAll() error {
	if err := app.RebuildPostClusters(); err != nil {
		return err
	}

	if err := app.RebuildUserClusters(); err != nil {
		return err
	}

	if err := app.RebuildCommunityClusters(); err != nil {
		return err
	}

	if err := app.RebuildPlatformClusters(); err != nil {
		return err
	}

	if err := app.RebuildCampaignClusters(); err != nil {
		return err
	}

	if err := app.RebuildTopicClusters(); err != nil {
		return err
	}

	if err := app.RebuildNarrativeClusters(); err != nil {
		return err
	}

	return nil
}

func (app *application) RebuildPostClusters() error {
	return app.rebuildClustersByKind(clusterables.PostKind)
}

func (app *application) RebuildUserClusters() error {
	return app.rebuildClustersByKind(clusterables.UserKind)
}

func (app *application) RebuildCommunityClusters() error {
	return app.rebuildClustersByKind(clusterables.CommunityKind)
}

func (app *application) RebuildPlatformClusters() error {
	return app.rebuildClustersByKind(clusterables.PlatformKind)
}

func (app *application) RebuildCampaignClusters() error {
	return app.rebuildClustersByKind(clusterables.CampaignKind)
}

func (app *application) RebuildTopicClusters() error {
	return app.rebuildClustersByKind(clusterables.TopicKind)
}

func (app *application) RebuildNarrativeClusters() error {
	return app.rebuildClustersByKind(clusterables.NarrativeKind)
}

func (app *application) rebuildClustersByKind(
	kind clusterables.Kind,
) error {
	cursor := uuid.Nil

	for {
		targets, err := app.clusterableRepository.FindByKindAfter(
			kind,
			cursor,
			rebuildBatchSize,
		)
		if err != nil {
			return err
		}

		if len(targets) == 0 {
			return nil
		}

		for _, target := range targets {
			candidates, err := app.candidateRepository.FindCandidates(
				target,
				kind,
				candidateAmount,
			)
			if err != nil {
				return err
			}

			clusters, err := app.detector.Detect(
				target,
				candidates,
			)
			if err != nil {
				return err
			}

			for _, cluster := range clusters {
				if err := app.repository.Save(cluster); err != nil {
					return err
				}
			}
		}

		cursor = targets[len(targets)-1].Identifier()
	}
}
