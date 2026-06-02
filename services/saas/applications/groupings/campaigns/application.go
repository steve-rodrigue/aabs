package campaigns

import (
	"github.com/google/uuid"

	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

type application struct {
	repository       domain_campaigns.Repository
	posts            app_posts.Application
	participations   app_participations.Application
	classifier       domain_campaigns.Classifier
	rebuildBatchSize int
}

func createApplication(
	repository domain_campaigns.Repository,
	posts app_posts.Application,
	participations app_participations.Application,
	classifier domain_campaigns.Classifier,
	rebuildBatchSize int,
) Application {
	return &application{
		repository:       repository,
		posts:            posts,
		participations:   participations,
		classifier:       classifier,
		rebuildBatchSize: rebuildBatchSize,
	}
}

func (app *application) FindByID(
	id uuid.UUID,
) (domain_campaigns.Campaign, error) {
	return app.repository.FindByID(id)
}

func (app *application) Find(
	index int,
	amount int,
) ([]domain_campaigns.Campaign, error) {
	return app.repository.Find(index, amount)
}

func (app *application) FindAfter(
	cursor uuid.UUID,
	amount int,
) ([]domain_campaigns.Campaign, error) {
	return app.repository.FindAfter(cursor, amount)
}

func (app *application) FindCampaignsByUser(
	user users.User,
) ([]domain_campaigns.Campaign, error) {
	return app.findCampaignsByParticipant(user)
}

func (app *application) FindCampaignsByCommunity(
	community communities.Community,
) ([]domain_campaigns.Campaign, error) {
	return app.findCampaignsByParticipant(community)
}

func (app *application) FindCampaignsByPlatform(
	platform platforms.Platform,
) ([]domain_campaigns.Campaign, error) {
	return app.findCampaignsByParticipant(platform)
}

func (app *application) findCampaignsByParticipant(
	participant participatables.Participatable,
) ([]domain_campaigns.Campaign, error) {
	participations, err := app.participations.FindByParticipant(participant)
	if err != nil {
		return nil, err
	}

	out := make([]domain_campaigns.Campaign, 0, len(participations))

	for _, participation := range participations {
		if participation.Target().ParticipationKind() != participatables.CampaignKind {
			continue
		}

		campaign, err := app.repository.FindByID(
			participation.Target().Identifier(),
		)
		if err != nil {
			return nil, err
		}

		out = append(out, campaign)
	}

	return out, nil
}

func (app *application) Count() (int64, error) {
	return app.repository.Count()
}

func (app *application) RebuildCampaigns() error {
	cursor := uuid.Nil

	for {
		posts, err := app.posts.FindAfter(cursor, app.rebuildBatchSize)
		if err != nil {
			return err
		}

		if len(posts) == 0 {
			return nil
		}

		for _, post := range posts {
			campaign, _, err := app.classifier.Classify(post)
			if err != nil {
				return err
			}

			if err := app.repository.Save(campaign); err != nil {
				return err
			}
		}

		cursor = posts[len(posts)-1].Identifier()
	}
}
