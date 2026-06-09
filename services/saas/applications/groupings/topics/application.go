package topics

import (
	"github.com/google/uuid"

	app_participations "github.com/steve-rodrigue/aabs/services/saas/applications/groupings/participations"
	app_posts "github.com/steve-rodrigue/aabs/services/saas/applications/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
	domain_topics "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
)

type application struct {
	repository       domain_topics.Repository
	posts            app_posts.Application
	participations   app_participations.Application
	builder          domain_topics.Builder
	rebuildBatchSize int
}

func createApplication(
	repository domain_topics.Repository,
	posts app_posts.Application,
	participations app_participations.Application,
	builder domain_topics.Builder,
	rebuildBatchSize int,
) Application {
	return &application{
		repository:       repository,
		posts:            posts,
		participations:   participations,
		builder:          builder,
		rebuildBatchSize: rebuildBatchSize,
	}
}

func (app *application) FindByID(id uuid.UUID) (domain_topics.Topic, error) {
	return app.repository.FindByID(id)
}

func (app *application) Find(index int, amount int) ([]domain_topics.Topic, error) {
	return app.repository.Find(index, amount)
}

func (app *application) FindAfter(cursor uuid.UUID, amount int) ([]domain_topics.Topic, error) {
	return app.repository.FindAfter(cursor, amount)
}

func (app *application) FindTopicsByUser(user users.User) ([]domain_topics.Topic, error) {
	return app.findTopicsByParticipant(user)
}

func (app *application) FindTopicsByCommunity(
	community communities.Community,
) ([]domain_topics.Topic, error) {
	return app.findTopicsByParticipant(community)
}

func (app *application) findTopicsByParticipant(
	participant participatables.Participatable,
) ([]domain_topics.Topic, error) {
	participations, err := app.participations.FindByParticipant(participant)
	if err != nil {
		return nil, err
	}

	out := make([]domain_topics.Topic, 0, len(participations))

	for _, participation := range participations {
		if participation.Target().ParticipationKind() != participatables.TopicKind {
			continue
		}

		topic, err := app.repository.FindByID(
			participation.Target().Identifier(),
		)
		if err != nil {
			return nil, err
		}

		out = append(out, topic)
	}

	return out, nil
}

func (app *application) Count() (int64, error) {
	return app.repository.Count()
}

func (app *application) RebuildTopics() error {
	cursor := uuid.Nil

	for {
		posts, err := app.posts.FindAfter(cursor, app.rebuildBatchSize)
		if err != nil {
			return err
		}

		if len(posts) == 0 {
			return nil
		}

		topics, err := app.builder.Build(posts)
		if err != nil {
			return err
		}

		for _, topic := range topics {
			if err := app.repository.Save(topic); err != nil {
				return err
			}
		}

		cursor = posts[len(posts)-1].Identifier()
	}
}
