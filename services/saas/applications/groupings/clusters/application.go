package clusters

import (
	"context"

	"github.com/google/uuid"

	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type application struct {
	repository            domain_clusters.Repository
	detector              domain_clusters.Detector
	clusterableRepository clusterables.Repository
	candidateRepository   clusterables.CandidateRepository
	rebuildBatchSize      int
	candidateAmount       int
}

func createApplication(
	repository domain_clusters.Repository,
	detector domain_clusters.Detector,
	clusterableRepository clusterables.Repository,
	candidateRepository clusterables.CandidateRepository,
	rebuildBatchSize int,
	candidateAmount int,
) Application {
	return &application{
		repository:            repository,
		detector:              detector,
		clusterableRepository: clusterableRepository,
		candidateRepository:   candidateRepository,
		rebuildBatchSize:      rebuildBatchSize,
		candidateAmount:       candidateAmount,
	}
}

func (app *application) BuildForTarget(
	ctx context.Context,
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	return app.detector.Detect(ctx, target, members)
}

func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_clusters.Cluster, error) {
	return app.repository.FindByID(ctx, id)
}

func (app *application) FindByTarget(
	ctx context.Context,
	target clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	return app.repository.FindByTarget(ctx, target.Identifier())
}

func (app *application) FindByMember(
	ctx context.Context,
	member clusterables.Clusterable,
) ([]domain_clusters.Cluster, error) {
	return app.repository.FindByMember(ctx, member.Identifier())
}

func (app *application) RebuildAll(
	ctx context.Context,
) error {
	if err := app.RebuildPostClusters(ctx); err != nil {
		return err
	}

	if err := app.RebuildUserClusters(ctx); err != nil {
		return err
	}

	if err := app.RebuildCommunityClusters(ctx); err != nil {
		return err
	}

	if err := app.RebuildPlatformClusters(ctx); err != nil {
		return err
	}

	if err := app.RebuildCampaignClusters(ctx); err != nil {
		return err
	}

	if err := app.RebuildTopicClusters(ctx); err != nil {
		return err
	}

	if err := app.RebuildNarrativeClusters(ctx); err != nil {
		return err
	}

	return nil
}

func (app *application) RebuildPostClusters(
	ctx context.Context,
) error {
	return app.rebuildClustersByKind(ctx, clusterables.PostKind)
}

func (app *application) RebuildUserClusters(
	ctx context.Context,
) error {
	return app.rebuildClustersByKind(ctx, clusterables.UserKind)
}

func (app *application) RebuildCommunityClusters(
	ctx context.Context,
) error {
	return app.rebuildClustersByKind(ctx, clusterables.CommunityKind)
}

func (app *application) RebuildPlatformClusters(
	ctx context.Context,
) error {
	return app.rebuildClustersByKind(ctx, clusterables.PlatformKind)
}

func (app *application) RebuildCampaignClusters(
	ctx context.Context,
) error {
	return app.rebuildClustersByKind(ctx, clusterables.CampaignKind)
}

func (app *application) RebuildTopicClusters(
	ctx context.Context,
) error {
	return app.rebuildClustersByKind(ctx, clusterables.TopicKind)
}

func (app *application) RebuildNarrativeClusters(
	ctx context.Context,
) error {
	return app.rebuildClustersByKind(ctx, clusterables.NarrativeKind)
}

func (app *application) rebuildClustersByKind(
	ctx context.Context,
	kind clusterables.Kind,
) error {
	cursor := uuid.Nil

	for {
		targets, err := app.clusterableRepository.FindByKindAfter(
			ctx,
			kind,
			cursor,
			app.rebuildBatchSize,
		)
		if err != nil {
			return err
		}

		if len(targets) == 0 {
			return nil
		}

		for _, target := range targets {
			retCandidates, err := app.candidateRepository.FindCandidates(
				ctx,
				target,
				kind,
				app.candidateAmount,
			)
			if err != nil {
				return err
			}

			clusters, err := app.detector.Detect(
				ctx,
				target,
				retCandidates,
			)
			if err != nil {
				return err
			}

			for _, cluster := range clusters {
				if err := app.repository.Save(ctx, cluster); err != nil {
					return err
				}
			}
		}

		cursor = targets[len(targets)-1].Identifier()
	}
}
