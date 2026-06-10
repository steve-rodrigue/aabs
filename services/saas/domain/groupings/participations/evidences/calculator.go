package evidences

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/participatables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
)

type calculator struct {
	adapter     Adapter
	posts       posts.Repository
	comparables clusterables.ComparableRepository
	threshold   float64
}

func (calculator *calculator) Calculate(
	ctx context.Context,
	participation participations.Participation,
) ([]Evidence, error) {
	if participation == nil {
		return nil, ErrInvalidEvidenceCalculatorParticipation
	}

	participant := participation.Participant()
	if participant == nil {
		return nil, ErrInvalidEvidenceCalculatorParticipant
	}

	target := participation.Target()
	if target == nil {
		return nil, ErrInvalidEvidenceCalculatorTarget
	}

	targetComparable, err := calculator.comparables.FindByID(
		ctx,
		target.Identifier(),
	)
	if err != nil {
		return nil, err
	}

	if targetComparable == nil {
		return nil, ErrInvalidEvidenceCalculatorComparable
	}

	candidatePosts, err := calculator.participantPosts(
		ctx,
		participant,
	)
	if err != nil {
		return nil, err
	}

	evidences := []Evidence{}

	for _, post := range candidatePosts {
		if post == nil {
			continue
		}

		postComparable, err := calculator.comparables.FindByID(
			ctx,
			post.Identifier(),
		)
		if err != nil {
			return nil, err
		}

		if postComparable == nil {
			continue
		}

		score, err := comparableScore(
			postComparable.Vector(),
			targetComparable.Vector(),
		)
		if err != nil {
			return nil, err
		}

		if score < calculator.threshold {
			continue
		}

		evidence, err := calculator.adapter.ToDomain(
			EvidenceInput{
				Identifier:    uuid.New(),
				Participation: participation,
				Participant:   participant,
				Target:        target,
				Post:          post,
				Score:         score,
				DetectedOn:    time.Now().UTC(),
			},
		)
		if err != nil {
			return nil, err
		}

		evidences = append(evidences, evidence)
	}

	return evidences, nil
}

func (calculator *calculator) participantPosts(
	ctx context.Context,
	participant participatables.Participatable,
) ([]posts.Post, error) {
	switch participant.ParticipationKind() {
	case participatables.PostKind:
		post, err := calculator.posts.FindByID(
			ctx,
			participant.Identifier(),
		)
		if err != nil {
			return nil, err
		}

		if post == nil {
			return []posts.Post{}, nil
		}

		return []posts.Post{post}, nil

	case participatables.UserKind:
		user, ok := participant.(users.User)
		if !ok {
			return nil, ErrInvalidEvidenceCalculatorParticipant
		}

		return calculator.posts.FindByUser(ctx, user)

	case participatables.CommunityKind:
		community, ok := participant.(communities.Community)
		if !ok {
			return nil, ErrInvalidEvidenceCalculatorParticipant
		}

		return calculator.posts.FindByCommunity(ctx, community)

	case participatables.PlatformKind:
		platform, ok := participant.(platforms.Platform)
		if !ok {
			return nil, ErrInvalidEvidenceCalculatorParticipant
		}

		return calculator.posts.FindByPlatform(ctx, platform)

	default:
		return nil, ErrInvalidEvidenceCalculatorParticipant
	}
}

func comparableScore(
	source []float32,
	target []float32,
) (float64, error) {
	if len(source) == 0 ||
		len(target) == 0 ||
		len(source) != len(target) {
		return 0, ErrInvalidEvidenceCalculatorVector
	}

	score := cosineSimilarity(source, target)

	if score < 0 {
		return 0, nil
	}

	if score > 1 {
		return 1, nil
	}

	return score, nil
}

func cosineSimilarity(
	source []float32,
	target []float32,
) float64 {
	var dot float64
	var sourceMagnitude float64
	var targetMagnitude float64

	for index := range source {
		sourceValue := float64(source[index])
		targetValue := float64(target[index])

		dot += sourceValue * targetValue
		sourceMagnitude += sourceValue * sourceValue
		targetMagnitude += targetValue * targetValue
	}

	if sourceMagnitude == 0 ||
		targetMagnitude == 0 {
		return 0
	}

	return dot / (math.Sqrt(sourceMagnitude) * math.Sqrt(targetMagnitude))
}
