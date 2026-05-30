package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	client "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	collectionName   = "comments"
	milvusAddress    = "localhost:19530"
	embeddingService = "http://localhost:8080/embed"
	hdbscanService   = "http://localhost:8090/cluster"
	llmService       = "http://localhost:8100/name-cluster"
	postgresDSN      = "postgres://aabs:aabs@localhost:5432/aabs?sslmode=disable"
	vectorDimension  = 1024
)

type EmbeddingRequest struct {
	Text string `json:"text"`
}

type EmbeddingResponse struct {
	Dimensions int       `json:"dimensions"`
	Embedding  []float32 `json:"embedding"`
}

type StoredPost struct {
	ID        int64
	Text      string
	Embedding []float32
}

type HDBSCANRequest struct {
	Embeddings     [][]float32 `json:"embeddings"`
	MinClusterSize int         `json:"min_cluster_size"`
	MinSamples     *int        `json:"min_samples"`
}

type HDBSCANResponse struct {
	TotalEmbeddings int              `json:"total_embeddings"`
	TotalClusters   int              `json:"total_clusters"`
	NoiseCount      int              `json:"noise_count"`
	Results         []HDBSCANCluster `json:"results"`
	Centroids       [][]float32      `json:"centroids"`
}

type HDBSCANCluster struct {
	Index       int     `json:"index"`
	ClusterID   int     `json:"cluster_id"`
	Probability float64 `json:"probability"`
	IsNoise     bool    `json:"is_noise"`
}

type LLMNameRequest struct {
	Posts        []string `json:"posts"`
	SystemPrompt string   `json:"system_prompt,omitempty"`
	UserPrompt   string   `json:"user_prompt,omitempty"`
	Temperature  float64  `json:"temperature,omitempty"`
	MaxTokens    int      `json:"max_tokens,omitempty"`
}

type LLMNameResponse struct {
	Name string `json:"name"`
	Raw  string `json:"raw"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	milvus, err := client.NewGrpcClient(ctx, milvusAddress)
	if err != nil {
		panic(err)
	}
	defer milvus.Close()

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	mustEnsurePostgres(ctx, db)

	command := os.Args[1]

	switch command {
	case "store":
		text := getTextArg()
		mustEnsureCollection(ctx, milvus)
		mustStoreEmbedding(ctx, milvus, text)

	case "verify":
		text := getTextArg()
		mustEnsureCollection(ctx, milvus)
		mustVerifyMeaning(ctx, milvus, text)

	case "cluster":
		mustEnsureCollection(ctx, milvus)
		mustClusterCampaigns(ctx, milvus, db)

	case "all":
		text := getTextArg()
		mustEnsureCollection(ctx, milvus)
		mustStoreEmbedding(ctx, milvus, text)
		mustVerifyMeaning(ctx, milvus, text)
		mustClusterCampaigns(ctx, milvus, db)

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run main.go store \"text to embed\"")
	fmt.Println("  go run main.go verify \"text to search\"")
	fmt.Println("  go run main.go cluster")
	fmt.Println("  go run main.go all \"text to embed, search, cluster, and name\"")
}

func getTextArg() string {
	if len(os.Args) >= 3 {
		return os.Args[2]
	}

	return "Trump is bad. Trump is terrible. Trump is awful."
}

func mustEnsurePostgres(ctx context.Context, db *sql.DB) {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS campaign_clusters (
			id SERIAL PRIMARY KEY,
			cluster_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			post_count INTEGER NOT NULL,
			avg_probability DOUBLE PRECISION NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS campaign_cluster_posts (
			id SERIAL PRIMARY KEY,
			cluster_id INTEGER NOT NULL,
			milvus_post_id BIGINT NOT NULL,
			text TEXT NOT NULL,
			probability DOUBLE PRECISION NOT NULL,
			is_noise BOOLEAN NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		`,
	}

	for _, query := range queries {
		if _, err := db.ExecContext(ctx, query); err != nil {
			panic(err)
		}
	}
}

func mustEnsureCollection(ctx context.Context, milvus client.Client) {
	exists, err := milvus.HasCollection(ctx, collectionName)
	if err != nil {
		panic(err)
	}

	if exists {
		mustEnsureIndex(ctx, milvus)
		return
	}

	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "AABS comment embeddings",
		AutoID:         true,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "text",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "65535",
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", vectorDimension),
				},
			},
		},
	}

	if err := milvus.CreateCollection(ctx, schema, 2); err != nil {
		panic(err)
	}

	fmt.Println("Collection created:", collectionName)

	mustEnsureIndex(ctx, milvus)
}

func mustEnsureIndex(ctx context.Context, milvus client.Client) {
	index, err := entity.NewIndexFlat(entity.COSINE)
	if err != nil {
		panic(err)
	}

	err = milvus.CreateIndex(
		ctx,
		collectionName,
		"embedding",
		index,
		false,
	)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "already") {
			return
		}

		fmt.Println("CreateIndex:", err)
		return
	}

	fmt.Println("Index created")
}

func mustFetchEmbedding(text string) EmbeddingResponse {
	body, err := json.Marshal(EmbeddingRequest{Text: text})
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{Timeout: 60 * time.Second}

	resp, err := httpClient.Post(
		embeddingService,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("unexpected status code %d: %s", resp.StatusCode, string(responseBody)))
	}

	var embedding EmbeddingResponse

	if err := json.NewDecoder(resp.Body).Decode(&embedding); err != nil {
		panic(err)
	}

	if embedding.Dimensions != vectorDimension {
		panic(fmt.Sprintf("unexpected embedding dimensions: got %d, expected %d", embedding.Dimensions, vectorDimension))
	}

	return embedding
}

func mustStoreEmbedding(ctx context.Context, milvus client.Client, text string) {
	embedding := mustFetchEmbedding(text)

	textColumn := entity.NewColumnVarChar(
		"text",
		[]string{text},
	)

	vectorColumn := entity.NewColumnFloatVector(
		"embedding",
		vectorDimension,
		[][]float32{embedding.Embedding},
	)

	_, err := milvus.Insert(
		ctx,
		collectionName,
		"",
		textColumn,
		vectorColumn,
	)
	if err != nil {
		panic(err)
	}

	if err := milvus.Flush(ctx, collectionName, false); err != nil {
		panic(err)
	}

	fmt.Println("Embedding stored in Milvus")
	fmt.Println("Text:", text)
	fmt.Println("Dimensions:", embedding.Dimensions)
}

func mustVerifyMeaning(ctx context.Context, milvus client.Client, text string) {
	if err := milvus.LoadCollection(ctx, collectionName, false); err != nil {
		panic(err)
	}

	embedding := mustFetchEmbedding(text)

	searchParam, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		panic(err)
	}

	results, err := milvus.Search(
		ctx,
		collectionName,
		[]string{},
		"",
		[]string{"text"},
		[]entity.Vector{
			entity.FloatVector(embedding.Embedding),
		},
		"embedding",
		entity.COSINE,
		10,
		searchParam,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Input text:")
	fmt.Println(text)
	fmt.Println()
	fmt.Println("Closest stored meanings:")

	for _, result := range results {
		for i, score := range result.Scores {
			fmt.Printf("Score: %.4f\n", score)

			for _, field := range result.Fields {
				if field.Name() == "text" {
					texts := field.(*entity.ColumnVarChar).Data()
					if i < len(texts) {
						fmt.Println("Text:", texts[i])
					}
				}
			}

			fmt.Println()
		}
	}
}

func mustClusterCampaigns(ctx context.Context, milvus client.Client, db *sql.DB) {
	posts := mustFetchStoredPosts(ctx, milvus)

	if len(posts) < 2 {
		fmt.Println("Not enough posts to cluster.")
		return
	}

	embeddings := make([][]float32, 0, len(posts))
	for _, post := range posts {
		embeddings = append(embeddings, post.Embedding)
	}

	minClusterSize := 2
	clusters := mustRunHDBSCAN(embeddings, minClusterSize)

	grouped := map[int][]HDBSCANCluster{}

	for _, result := range clusters.Results {
		if result.IsNoise || result.ClusterID < 0 {
			continue
		}

		grouped[result.ClusterID] = append(grouped[result.ClusterID], result)
	}

	if len(grouped) == 0 {
		fmt.Println("No campaign clusters found.")
		return
	}

	for clusterID, members := range grouped {
		clusterPosts := make([]StoredPost, 0, len(members))
		totalProbability := 0.0

		for _, member := range members {
			if member.Index >= len(posts) {
				continue
			}

			clusterPosts = append(clusterPosts, posts[member.Index])
			totalProbability += member.Probability
		}

		if len(clusterPosts) == 0 {
			continue
		}

		name := mustNameCluster(clusterPosts)
		avgProbability := totalProbability / float64(len(clusterPosts))

		mustStoreCampaignCluster(ctx, db, clusterID, name, clusterPosts, members, avgProbability)

		fmt.Println("Campaign cluster stored")
		fmt.Println("Cluster ID:", clusterID)
		fmt.Println("Name:", name)
		fmt.Println("Posts:", len(clusterPosts))
		fmt.Printf("Average probability: %.4f\n", avgProbability)
		fmt.Println()
	}
}

func mustFetchStoredPosts(ctx context.Context, milvus client.Client) []StoredPost {
	if err := milvus.LoadCollection(ctx, collectionName, false); err != nil {
		panic(err)
	}

	outputFields := []string{
		"id",
		"text",
		"embedding",
	}

	results, err := milvus.Query(
		ctx,
		collectionName,
		[]string{},
		"id >= 0",
		outputFields,
	)
	if err != nil {
		panic(err)
	}

	var ids []int64
	var texts []string
	var embeddings [][]float32

	for _, column := range results {
		switch c := column.(type) {
		case *entity.ColumnInt64:
			if c.Name() == "id" {
				ids = c.Data()
			}

		case *entity.ColumnVarChar:
			if c.Name() == "text" {
				texts = c.Data()
			}

		case *entity.ColumnFloatVector:
			if c.Name() == "embedding" {
				embeddings = c.Data()
			}
		}
	}

	posts := make([]StoredPost, 0, len(ids))

	for i := range ids {
		if i >= len(texts) || i >= len(embeddings) {
			continue
		}

		posts = append(posts, StoredPost{
			ID:        ids[i],
			Text:      texts[i],
			Embedding: embeddings[i],
		})
	}

	return posts
}

func mustRunHDBSCAN(embeddings [][]float32, minClusterSize int) HDBSCANResponse {
	body, err := json.Marshal(HDBSCANRequest{
		Embeddings:     embeddings,
		MinClusterSize: minClusterSize,
		MinSamples:     nil,
	})
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{Timeout: 120 * time.Second}

	resp, err := httpClient.Post(
		hdbscanService,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("hdbscan error %d: %s", resp.StatusCode, string(responseBody)))
	}

	var result HDBSCANResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	return result
}

func mustNameCluster(posts []StoredPost) string {
	postTexts := make([]string, 0, len(posts))

	for _, post := range posts {
		postTexts = append(postTexts, post.Text)
	}

	body, err := json.Marshal(
		LLMNameRequest{
			Posts: postTexts,

			UserPrompt: `
	Return only valid JSON:
	
	{"name":"cluster name"}
	
	Rules:
	- 2 to 5 words
	- neutral
	- no explanation
	- no markdown
	- no reasoning
	`,

			Temperature: 0.1,
			MaxTokens:   32,
		},
	)

	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{Timeout: 180 * time.Second}

	resp, err := httpClient.Post(
		llmService,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("llm error %d: %s", resp.StatusCode, string(responseBody)))
	}

	var result LLMNameResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	name := strings.TrimSpace(result.Name)
	if len(name) > 64 {
		name = name[:64]
	}

	if name == "" || name == "Unnamed Cluster" {
		name = fmt.Sprintf(
			"Cluster %d",
			time.Now().Unix(),
		)
	}

	if name == "" {
		return "Unnamed Campaign Cluster"
	}

	return name
}

func mustStoreCampaignCluster(
	ctx context.Context,
	db *sql.DB,
	clusterID int,
	name string,
	posts []StoredPost,
	members []HDBSCANCluster,
	avgProbability float64,
) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO campaign_clusters (
			cluster_id,
			name,
			post_count,
			avg_probability
		)
		VALUES ($1, $2, $3, $4)
		`,
		clusterID,
		name,
		len(posts),
		avgProbability,
	)
	if err != nil {
		panic(err)
	}

	for i, post := range posts {
		probability := 0.0
		isNoise := false

		if i < len(members) {
			probability = members[i].Probability
			isNoise = members[i].IsNoise
		}

		_, err = tx.ExecContext(
			ctx,
			`
			INSERT INTO campaign_cluster_posts (
				cluster_id,
				milvus_post_id,
				text,
				probability,
				is_noise
			)
			VALUES ($1, $2, $3, $4, $5)
			`,
			clusterID,
			post.ID,
			post.Text,
			probability,
			isNoise,
		)
		if err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
