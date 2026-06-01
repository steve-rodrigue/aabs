# AABS

> Project Status: Active Development
>
> AABS is currently under active development. The initial goal is to deliver a functional Minimum Viable Product (MVP) by the end of June 1st, 2026.  I adjusted the timeline from the previous May 31st, 2026 to not work all night.  My wife is very happy that I made that decision.  Happy wife, happy life!

AABS (Anti-AI Bot Spam) is an open-source browser extension and platform designed to help users identify spam, bot networks, coordinated manipulation, and low-quality engagement on social media.

Rather than focusing solely on whether an account is a bot, AABS analyzes content similarity, posting behavior, semantic patterns, and network activity to detect coordinated campaigns, repetitive social-proof behavior, and other forms of artificial amplification.

AABS assigns trust scores to accounts, posts, comments, and conversations, helping users distinguish authentic discussion from coordinated or deceptive activity.

## SaaS Domain Model

The domain objects described below belong to the SaaS backend portion of AABS.

They define how the SaaS receives posts, analyzes content, stores semantic relationships, calculates trust scores, and exposes campaign detection results to the browser extension, dashboards, APIs, and other clients.

### Core Entities

#### Platform

A website or service where content is published.

Examples:

- Reddit
- X (Twitter)
- Facebook
- YouTube
- TikTok

#### Community

A group within a platform where users publish content.

Examples:

- r/politics
- r/worldnews
- Facebook groups
- Discord servers
- YouTube channels

#### User

A platform account that creates content.

Users belong to a platform and may participate in communities, topics, narratives, and campaigns.

#### Post

A piece of content published by a user.

Posts are the primary input processed by the SaaS pipeline.

### Semantic Groupings

#### Cluster

A cluster is a collection of semantically similar entities.

Clusters are generated automatically using embeddings and clustering algorithms.

Examples:

- Similar posts
- Similar users
- Similar campaigns
- Similar topics
- Similar narratives

#### Campaign

A campaign is a recurring pattern of semantically related content.

Campaigns are derived from clusters of posts and represent coordinated or repeated messaging.

Examples:

- Political messaging campaigns
- Product promotion campaigns
- Spam campaigns
- Influence operations

#### Topic

A topic represents a semantic subject.

Topics can form a hierarchical tree and may contain sub-topics.

Examples:

- Politics
  - Elections
  - Immigration
- Technology
  - Artificial Intelligence
  - Cybersecurity

#### Narrative

A narrative represents a specific story, claim, or message being propagated.

Examples:

- "Candidate X is corrupt"
- "Product Y is dangerous"
- "Technology Z will replace jobs"

Topics describe what is being discussed.

Narratives describe what is being claimed.

### Participations

Participations measure how strongly one entity contributes to another.

Examples:

- User → Campaign
- User → Topic
- User → Narrative
- Community → Campaign
- Platform → Narrative

Each participation contains:

- Number of matching posts
- Total posts analyzed
- Participation percentage
- Evidence supporting the participation

### Relationships

Relationships connect semantically related entities.

Examples:

- Campaign ↔ Campaign
- Topic ↔ Topic
- Narrative ↔ Narrative
- User ↔ User
- Campaign ↔ Narrative

Each relationship contains a similarity score derived from semantic comparison.

### Trust Scores

Trust scores estimate the likelihood that an entity represents authentic behavior.

Scores may be calculated for:

- Users
- Posts
- Campaigns
- Topics
- Narratives
- Communities
- Relationships
- Clusters

Trust scores are composed of weighted factors including:

- Semantic repetition
- Campaign participation
- User concentration
- Account age
- Posting velocity
- Relationship risk
- Community spread
- Content quality signals

### Processing Pipeline

When a post is received by the SaaS:

1. The post is stored.
2. An embedding is generated.
3. The embedding is indexed for semantic search.
4. Similar content is identified.
5. Clusters are updated.
6. Campaigns, topics, and narratives are rebuilt.
7. Participations are recalculated.
8. Relationships are generated.
9. Trust scores are updated.

This pipeline continuously updates the semantic graph as new content is ingested.

## Development

### Build and start all services
```bash
bash docker compose up --build 
```

### Build and start all services in background
```bash
bash docker compose up -d --build 
```

### View logs
```bash
bash docker compose logs -f 
```

### Stop all services
```bash
bash docker compose down 
```

## License

MIT