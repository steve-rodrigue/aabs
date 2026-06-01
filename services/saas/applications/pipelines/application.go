package pipelines

import (
	"github.com/steve-rodrigue/aabs/services/saas/applications/groupings"
	"github.com/steve-rodrigue/aabs/services/saas/applications/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/applications/scores"
	"github.com/steve-rodrigue/aabs/services/saas/applications/searches"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
)

type application struct {
	postRepository posts.Repository

	searches searches.Application

	groupings     groupings.Application
	relationships relationships.Application
	scores        scores.Application
}

func createApplication(
	postRepository posts.Repository,
	searches searches.Application,
	groupings groupings.Application,
	relationships relationships.Application,
	scores scores.Application,
) Application {
	out := application{
		postRepository: postRepository,
		searches:       searches,
		groupings:      groupings,
		relationships:  relationships,
		scores:         scores,
	}

	return &out
}

// ProcessPost processes a single post
func (app *application) ProcessPost(post posts.Post) error {
	if err := app.postRepository.Save(post); err != nil {
		return err
	}

	text := post.Content().Text()

	if _, err := app.searches.SearchPosts(text, 50); err != nil {
		return err
	}

	return app.Rebuild()
}

// ProcessPosts processes multiple posts
func (app *application) ProcessPosts(posts []posts.Post) error {
	for _, post := range posts {
		if err := app.ProcessPost(post); err != nil {
			return err
		}
	}

	return nil
}

// Rebuild rebuilds all semantic analysis artifacts from stored posts
func (app *application) Rebuild() error {
	if err := app.groupings.Clusters().RebuildPostClusters(); err != nil {
		return err
	}

	if err := app.groupings.Campaigns().RebuildCampaigns(); err != nil {
		return err
	}

	if err := app.groupings.Topics().RebuildTopics(); err != nil {
		return err
	}

	if err := app.groupings.Narratives().RebuildNarratives(); err != nil {
		return err
	}

	if err := app.groupings.Participations().RebuildParticipations(); err != nil {
		return err
	}

	if err := app.relationships.RebuildRelationships(); err != nil {
		return err
	}

	if err := app.scores.RecalculateScores(); err != nil {
		return err
	}

	return nil
}
