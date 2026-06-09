package evidences

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_communities "github.com/steve-rodrigue/aabs/services/saas/domain/entities/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/saas/domain/entities/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
	domain_users "github.com/steve-rodrigue/aabs/services/saas/domain/entities/users"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/participations/participatables"
)

var errTest = errors.New("test error")

func TestNewCalculator(t *testing.T) {
	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		domain_posts.NewMockPostRepository(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	if calculator == nil {
		t.Fatalf("expected calculator")
	}
}

func TestCalculatorCalculateForPostParticipant(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockEvidenceAdapter()
	postRepository := domain_posts.NewMockPostRepository()
	comparableRepository := clusterables.NewMockComparableRepository()

	post := domain_posts.NewMockPost("hello")
	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)
	participant := participatables.NewMockParticipatable(
		post.Identifier(),
		participatables.PostKind,
	)

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	postRepository.Items[post.Identifier()] = post

	comparableRepository.Items[target.Identifier()] =
		clusterables.NewMockComparableWithID(
			target.Identifier(),
			clusterables.TopicKind,
			[]float32{1, 0},
		)

	comparableRepository.Items[post.Identifier()] =
		clusterables.NewMockComparableWithID(
			post.Identifier(),
			clusterables.PostKind,
			[]float32{1, 0},
		)

	expected := NewMockEvidence()
	adapter.ToDomainValue = expected

	calculator := NewCalculator(
		adapter,
		postRepository,
		comparableRepository,
		0.7,
	)

	result, err := calculator.Calculate(
		ctx,
		participation,
	)

	if err != nil {
		t.Fatal(err)
	}

	if postRepository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if postRepository.LastContext != ctx {
		t.Fatalf("expected context to be passed to post repository")
	}

	if postRepository.LastID != post.Identifier() {
		t.Fatalf("expected post id to be passed")
	}

	if comparableRepository.FindByIDCalls != 2 {
		t.Fatalf(
			"expected 2 comparable lookups, got %d",
			comparableRepository.FindByIDCalls,
		)
	}

	if adapter.ToDomainCalls != 1 {
		t.Fatalf("expected 1 adapter call")
	}

	if adapter.LastInput.Participation != participation {
		t.Fatalf("expected participation")
	}

	if adapter.LastInput.Participant != participant {
		t.Fatalf("expected participant")
	}

	if adapter.LastInput.Target != target {
		t.Fatalf("expected target")
	}

	if adapter.LastInput.Post != post {
		t.Fatalf("expected post")
	}

	if adapter.LastInput.Score != 1 {
		t.Fatalf("expected score 1, got %f", adapter.LastInput.Score)
	}

	if adapter.LastInput.Identifier == uuid.Nil {
		t.Fatalf("expected generated identifier")
	}

	if adapter.LastInput.DetectedOn.IsZero() {
		t.Fatalf("expected detected on")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 evidence, got %d", len(result))
	}

	if result[0] != expected {
		t.Fatalf("expected evidence result")
	}
}

func TestCalculatorCalculateSkipsPostsBelowThreshold(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockEvidenceAdapter()
	postRepository := domain_posts.NewMockPostRepository()
	comparableRepository := clusterables.NewMockComparableRepository()

	post := domain_posts.NewMockPost("hello")
	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)
	participant := participatables.NewMockParticipatable(
		post.Identifier(),
		participatables.PostKind,
	)

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	postRepository.Items[post.Identifier()] = post

	comparableRepository.Items[target.Identifier()] =
		clusterables.NewMockComparableWithID(
			target.Identifier(),
			clusterables.TopicKind,
			[]float32{1, 0},
		)

	comparableRepository.Items[post.Identifier()] =
		clusterables.NewMockComparableWithID(
			post.Identifier(),
			clusterables.PostKind,
			[]float32{0, 1},
		)

	calculator := NewCalculator(
		adapter,
		postRepository,
		comparableRepository,
		0.7,
	)

	result, err := calculator.Calculate(
		ctx,
		participation,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected no evidence")
	}

	if adapter.ToDomainCalls != 0 {
		t.Fatalf("expected no adapter call")
	}
}

func TestCalculatorCalculateReturnsEmptyWhenPostParticipantNotFound(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockEvidenceAdapter()
	postRepository := domain_posts.NewMockPostRepository()
	comparableRepository := clusterables.NewMockComparableRepository()

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)

	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.PostKind,
	)

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	comparableRepository.Items[target.Identifier()] =
		clusterables.NewMockComparableWithID(
			target.Identifier(),
			clusterables.TopicKind,
			[]float32{1, 0},
		)

	calculator := NewCalculator(
		adapter,
		postRepository,
		comparableRepository,
		0.7,
	)

	result, err := calculator.Calculate(
		ctx,
		participation,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected empty result")
	}

	if adapter.ToDomainCalls != 0 {
		t.Fatalf("expected no adapter call")
	}
}

func TestCalculatorCalculateSkipsPostWithoutComparable(t *testing.T) {
	ctx := context.Background()

	adapter := NewMockEvidenceAdapter()
	postRepository := domain_posts.NewMockPostRepository()
	comparableRepository := clusterables.NewMockComparableRepository()

	post := domain_posts.NewMockPost("hello")
	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)
	participant := participatables.NewMockParticipatable(
		post.Identifier(),
		participatables.PostKind,
	)

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	postRepository.Items[post.Identifier()] = post

	comparableRepository.Items[target.Identifier()] =
		clusterables.NewMockComparableWithID(
			target.Identifier(),
			clusterables.TopicKind,
			[]float32{1, 0},
		)

	calculator := NewCalculator(
		adapter,
		postRepository,
		comparableRepository,
		0.7,
	)

	result, err := calculator.Calculate(
		ctx,
		participation,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected no evidence")
	}

	if adapter.ToDomainCalls != 0 {
		t.Fatalf("expected no adapter call")
	}
}

func TestCalculatorCalculateReturnsInvalidParticipationError(t *testing.T) {
	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		domain_posts.NewMockPostRepository(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		nil,
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorParticipation) {
		t.Fatalf("expected invalid participation error, got %v", err)
	}
}

func TestCalculatorCalculateReturnsInvalidParticipantError(t *testing.T) {
	participation := participations.NewMockParticipationWithParticipantAndTarget(
		nil,
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.TopicKind,
		),
	)

	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		domain_posts.NewMockPostRepository(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		participation,
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorParticipant) {
		t.Fatalf("expected invalid participant error, got %v", err)
	}
}

func TestCalculatorCalculateReturnsInvalidTargetError(t *testing.T) {
	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.PostKind,
		),
		nil,
	)

	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		domain_posts.NewMockPostRepository(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		participation,
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorTarget) {
		t.Fatalf("expected invalid target error, got %v", err)
	}
}

func TestCalculatorCalculateReturnsComparableRepositoryError(t *testing.T) {
	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.FindByIDErr = errTest

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.PostKind,
		),
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.TopicKind,
		),
	)

	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		domain_posts.NewMockPostRepository(),
		comparableRepository,
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		participation,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected comparable repository error, got %v", err)
	}
}

func TestCalculatorCalculateReturnsInvalidComparableError(t *testing.T) {
	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.PostKind,
		),
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.TopicKind,
		),
	)

	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		domain_posts.NewMockPostRepository(),
		clusterables.NewMockComparableRepository(),
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		participation,
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorComparable) {
		t.Fatalf("expected invalid comparable error, got %v", err)
	}
}

func TestCalculatorCalculateReturnsPostRepositoryError(t *testing.T) {
	postRepository := domain_posts.NewMockPostRepository()
	postRepository.FindByIDErr = errTest

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)
	participant := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.PostKind,
	)

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.Items[target.Identifier()] =
		clusterables.NewMockComparableWithID(
			target.Identifier(),
			clusterables.TopicKind,
			[]float32{1, 0},
		)

	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		postRepository,
		comparableRepository,
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		participation,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected post repository error, got %v", err)
	}
}

func TestCalculatorCalculateReturnsVectorError(t *testing.T) {
	post := domain_posts.NewMockPost("hello")

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)
	participant := participatables.NewMockParticipatable(
		post.Identifier(),
		participatables.PostKind,
	)

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	postRepository := domain_posts.NewMockPostRepository()
	postRepository.Items[post.Identifier()] = post

	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.Items[target.Identifier()] =
		clusterables.NewMockComparableWithID(
			target.Identifier(),
			clusterables.TopicKind,
			[]float32{1, 0},
		)

	comparableRepository.Items[post.Identifier()] =
		clusterables.NewMockComparableWithID(
			post.Identifier(),
			clusterables.PostKind,
			[]float32{1, 0, 0},
		)

	calculator := NewCalculator(
		NewMockEvidenceAdapter(),
		postRepository,
		comparableRepository,
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		participation,
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorVector) {
		t.Fatalf("expected vector error, got %v", err)
	}
}

func TestCalculatorCalculateReturnsAdapterError(t *testing.T) {
	adapter := NewMockEvidenceAdapter()
	adapter.ToDomainErr = errTest

	post := domain_posts.NewMockPost("hello")

	target := participatables.NewMockParticipatable(
		uuid.New(),
		participatables.TopicKind,
	)
	participant := participatables.NewMockParticipatable(
		post.Identifier(),
		participatables.PostKind,
	)

	participation := participations.NewMockParticipationWithParticipantAndTarget(
		participant,
		target,
	)

	postRepository := domain_posts.NewMockPostRepository()
	postRepository.Items[post.Identifier()] = post

	comparableRepository := clusterables.NewMockComparableRepository()
	comparableRepository.Items[target.Identifier()] =
		clusterables.NewMockComparableWithID(
			target.Identifier(),
			clusterables.TopicKind,
			[]float32{1, 0},
		)

	comparableRepository.Items[post.Identifier()] =
		clusterables.NewMockComparableWithID(
			post.Identifier(),
			clusterables.PostKind,
			[]float32{1, 0},
		)

	calculator := NewCalculator(
		adapter,
		postRepository,
		comparableRepository,
		0.7,
	)

	_, err := calculator.Calculate(
		context.Background(),
		participation,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected adapter error, got %v", err)
	}
}

func TestParticipantPostsReturnsInvalidParticipantForUserMock(t *testing.T) {
	calculator := &calculator{
		posts: domain_posts.NewMockPostRepository(),
	}

	_, err := calculator.participantPosts(
		context.Background(),
		participatables.NewMockParticipatable(
			uuid.New(),
			participatables.UserKind,
		),
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorParticipant) {
		t.Fatalf("expected invalid participant error, got %v", err)
	}
}

func TestComparableScore(t *testing.T) {
	score, err := comparableScore(
		[]float32{1, 0},
		[]float32{1, 0},
	)

	if err != nil {
		t.Fatal(err)
	}

	if score != 1 {
		t.Fatalf("expected score 1, got %f", score)
	}
}

func TestComparableScoreClampsNegativeToZero(t *testing.T) {
	score, err := comparableScore(
		[]float32{-1, 0},
		[]float32{1, 0},
	)

	if err != nil {
		t.Fatal(err)
	}

	if score != 0 {
		t.Fatalf("expected score 0, got %f", score)
	}
}

func TestComparableScoreReturnsVectorErrorWhenEmpty(t *testing.T) {
	_, err := comparableScore(
		[]float32{},
		[]float32{1, 0},
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorVector) {
		t.Fatalf("expected vector error, got %v", err)
	}
}

func TestComparableScoreReturnsVectorErrorWhenMismatch(t *testing.T) {
	_, err := comparableScore(
		[]float32{1, 0},
		[]float32{1, 0, 0},
	)

	if !errors.Is(err, ErrInvalidEvidenceCalculatorVector) {
		t.Fatalf("expected vector error, got %v", err)
	}
}

func TestCosineSimilarityReturnsZeroWhenSourceMagnitudeIsZero(t *testing.T) {
	score := cosineSimilarity(
		[]float32{0, 0},
		[]float32{1, 0},
	)

	if score != 0 {
		t.Fatalf("expected score 0, got %f", score)
	}
}

func TestCosineSimilarityReturnsZeroWhenTargetMagnitudeIsZero(t *testing.T) {
	score := cosineSimilarity(
		[]float32{1, 0},
		[]float32{0, 0},
	)

	if score != 0 {
		t.Fatalf("expected score 0, got %f", score)
	}
}

func TestParticipantPostsForUser(t *testing.T) {
	ctx := context.Background()

	user := domain_users.NewMockUser(
		"steve",
		"Steve",
	)

	postRepository := domain_posts.NewMockPostRepository()
	postRepository.FindByUserValue = []domain_posts.Post{
		domain_posts.NewMockPost("hello"),
	}

	calculator := &calculator{
		posts: postRepository,
	}

	result, err := calculator.participantPosts(
		ctx,
		user,
	)

	if err != nil {
		t.Fatal(err)
	}

	if postRepository.FindByUserCalls != 1 {
		t.Fatalf("expected 1 find by user call")
	}

	if postRepository.LastContext != ctx {
		t.Fatalf("expected context")
	}

	if postRepository.LastUser != user {
		t.Fatalf("expected user")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post")
	}
}

func TestParticipantPostsForCommunity(t *testing.T) {
	ctx := context.Background()

	community := domain_communities.NewMockCommunity(
		"community",
		"community text",
	)

	postRepository := domain_posts.NewMockPostRepository()
	postRepository.FindByCommunityValue = []domain_posts.Post{
		domain_posts.NewMockPost("hello"),
	}

	calculator := &calculator{
		posts: postRepository,
	}

	result, err := calculator.participantPosts(
		ctx,
		community,
	)

	if err != nil {
		t.Fatal(err)
	}

	if postRepository.FindByCommunityCalls != 1 {
		t.Fatalf("expected 1 find by community call")
	}

	if postRepository.LastCommunity != community {
		t.Fatalf("expected community")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post")
	}
}

func TestParticipantPostsForPlatform(t *testing.T) {
	ctx := context.Background()

	platform := domain_platforms.NewMockPlatform(
		"Reddit",
		"reddit",
	)

	postRepository := domain_posts.NewMockPostRepository()
	postRepository.FindByPlatformValue = []domain_posts.Post{
		domain_posts.NewMockPost("hello"),
	}

	calculator := &calculator{
		posts: postRepository,
	}

	result, err := calculator.participantPosts(
		ctx,
		platform,
	)

	if err != nil {
		t.Fatal(err)
	}

	if postRepository.FindByPlatformCalls != 1 {
		t.Fatalf("expected 1 find by platform call")
	}

	if postRepository.LastPlatform != platform {
		t.Fatalf("expected platform")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 post")
	}
}
