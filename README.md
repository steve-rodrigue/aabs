# AABS

> Project Status: Active Development

AABS (Anti-AI Bot Spam) is an open-source browser extension and platform designed to help users identify spam, bot networks, coordinated manipulation, and low-quality engagement on social media.

Rather than focusing solely on whether an account is a bot, AABS analyzes content similarity, posting behavior, semantic patterns, and network activity to detect coordinated campaigns, repetitive social-proof behavior, and other forms of artificial amplification.

AABS assigns trust scores to accounts, posts, comments, and conversations, helping users distinguish authentic discussion from coordinated or deceptive activity.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
  - [Layers](#layers)
  - [Application Graph](#application-graph)
- [Applications](#applications)
  - [Root Application](#root-application)
  - [Pipeline Application](#pipeline-application)
  - [Posts Application](#posts-application)
  - [Users Application](#users-application)
  - [Communities Application](#communities-application)
  - [Platforms Application](#platforms-application)
  - [Searches Application](#searches-application)
  - [Groupings Application](#groupings-application)
    - [Clusters Application](#clusters-application)
    - [Campaigns Application](#campaigns-application)
    - [Topics Application](#topics-application)
    - [Narratives Application](#narratives-application)
    - [Participations Application](#participations-application)
      - [Participation Evidences Application](#participation-evidences-application)
  - [Relationships Application](#relationships-application)
    - [Comparables Application](#comparables-application)
  - [Scores Application](#scores-application)
- [Processing Pipeline](#processing-pipeline)
- [Development](#development)
- [License](#license)

## Overview

AABS continuously ingests social media content and transforms it into a semantic graph.

The system stores content, indexes it for semantic search, groups related entities, calculates participations, builds relationships, and continuously updates trust scores.

The platform is built around a modular application layer where every feature is exposed through a dedicated application interface.

## Architecture

### Layers

The SaaS backend is organized into three primary layers.

- Applications
- Domain
- Infrastructure 

Applications orchestrate workflows.

The domain layer defines business rules and interfaces.

The infrastructure layer provides implementations such as databases, embedding engines, search indexes, and external integrations.

### Application Graph

- Application
  - Pipeline
  - Posts
  - Users
  - Communities
  - Platforms
  - Searches
  - Groupings
    - Clusters
    - Campaigns
    - Topics
    - Narratives
    - Participations
      - Evidences
  - Relationships
    - Comparables
  - Scores

Applications communicate through interfaces rather than concrete implementations.

This keeps the system testable, modular, and independent from storage technologies.

## Applications

### Root Application

The root application exposes the complete application graph.

Responsibilities:

- Provide access to all applications
- Centralize application wiring
- Act as the main entry point for consumers

Exposed applications:

- Pipeline
- Posts
- Users
- Communities
- Platforms
- Searches
- Groupings
- Relationships
- Scores

### Pipeline Application

The pipeline application orchestrates content ingestion and semantic graph rebuilding.

The pipeline serves as the primary entry point for new content entering the system. It coordinates multiple applications to transform raw posts into searchable entities, semantic groupings, relationships, and trust scores.

Responsibilities:

- Process a single post
- Process multiple posts
- Rebuild semantic artifacts

When a post is processed, the pipeline:

1. Stores the post
2. Indexes searchable content
3. Rebuilds clusters
4. Rebuilds campaigns
5. Rebuilds topics
6. Rebuilds narratives
7. Rebuilds participations
8. Rebuilds relationships
9. Recalculates scores

The pipeline ensures that all derived semantic data remains synchronized with the latest content available in the system.

### Posts Application

The posts application manages content storage and retrieval.

Posts are the primary source of information within AABS. Every semantic analysis, campaign detection, topic extraction, narrative discovery, participation calculation, relationship generation, and score calculation ultimately originates from posts.

A post may represent:

- A thread
- A reply
- A comment
- A social media publication
- Any other piece of user-generated content

The posts application provides a consistent interface for storing and querying content regardless of the platform from which it originated.

Responsibilities:

- Save posts
- Find posts by identifier
- Find posts using index pagination
- Find posts using cursor pagination
- Count posts
- Find posts by user
- Find posts by community
- Find posts by platform

Posts serve as the foundation of the semantic graph. Most higher-level entities such as campaigns, topics, narratives, participations, relationships, and trust scores are derived directly or indirectly from the content stored by this application.

### Users Application

The users application manages platform users.

Users represent accounts that create, share, or interact with content on supported platforms. They are one of the central entities within the AABS semantic graph and act as the primary actors behind posts, campaigns, narratives, relationships, participations, and trust scores.

A user may represent:

- A human account
- An automated account
- A bot network participant
- An organization account
- A brand account
- Any identifiable content publisher

The users application provides a unified interface for storing and querying accounts independently of the platform on which they exist.

Responsibilities:

- Save users
- Find users by identifier
- Find users by platform-specific identifier
- Find users using index pagination
- Find users using cursor pagination
- Count users

Users serve as the primary source of behavioral analysis within AABS. By analyzing the content they publish, the campaigns they participate in, the narratives they propagate, the relationships they form, and their historical activity patterns, the platform can identify coordinated behavior, detect suspicious amplification, and calculate trust and risk scores.

### Communities Application

The communities application manages communities within platforms.

Communities represent locations where users gather, interact, and publish content. They provide important context for understanding how information spreads across a platform and help identify where campaigns, narratives, and coordinated behavior originate or gain traction.

A community may represent:

- A subreddit on Reddit
- A Facebook group
- A Discord server
- A YouTube channel
- A Telegram group
- A forum section
- Any platform-specific content community

The communities application provides a unified interface for storing and querying communities independently of the platform on which they exist.

Responsibilities:

- Save communities
- Find communities by identifier
- Find communities by platform and handle
- Find communities using index pagination
- Find communities using cursor pagination
- Find communities by platform
- Count communities

Communities serve as an important layer of analysis within AABS. By examining the users, posts, campaigns, topics, narratives, participations, relationships, and trust scores associated with a community, the platform can identify information hubs, coordinated influence operations, echo chambers, highly trusted communities, and areas where manipulation is concentrated.

Community-level analysis also helps distinguish isolated behavior from platform-wide coordination by revealing how narratives and campaigns propagate between different groups over time.

### Platforms Application

The platforms application manages content platforms.

Platforms represent the highest-level source of information within AABS. Every community, user, and post originates from a platform, making platforms the foundation upon which the entire semantic graph is built.

A platform may represent:

- Reddit
- X (Twitter)
- Facebook
- YouTube
- TikTok
- Instagram
- LinkedIn
- Discord
- Telegram
- Forums, blogs, or custom social networks

The platforms application provides a unified interface for storing and querying content sources regardless of their underlying technology, APIs, or data structures.

Responsibilities:

- Save platforms
- Find platforms by identifier
- Find platforms by handle
- Find platforms using index pagination
- Find platforms using cursor pagination
- Count platforms

Platforms provide critical context for analysis because user behavior, content distribution patterns, moderation policies, and engagement mechanisms vary significantly between ecosystems.

By organizing data around platforms, AABS can:

- Compare activity across multiple social networks
- Detect campaigns operating on several platforms simultaneously
- Identify platform-specific manipulation techniques
- Measure how narratives spread between ecosystems
- Analyze the effectiveness of cross-platform influence operations
- Calculate trust and risk metrics at the platform level

Platforms act as the root of the content hierarchy:

```text
Platform
 ├── Communities
 │    ├── Users
 │    │    └── Posts
 │    └── Posts
 └── Platform-Level Analysis
      ├── Campaigns
      ├── Topics
      ├── Narratives
      ├── Relationships
      └── Scores
```

This structure allows AABS to analyze not only individual pieces of content, but also the broader ecosystems in which information is created, amplified, and propagated.

### Searches Application

The searches application provides semantic search across the entire AABS knowledge graph.

Unlike traditional keyword search, the search system understands meaning rather than exact wording. By leveraging embeddings and vector similarity, it can discover related content even when different words, phrases, languages, or writing styles are used.

The search engine acts as the discovery layer of AABS, allowing every searchable entity to be explored through semantic similarity rather than strict text matching.

Any object implementing the searchable interface can be indexed, making the search system completely independent from specific entity types.

Examples of searchable entities include:

- Posts
- Users
- Communities
- Campaigns
- Topics
- Narratives
- Relationships
- Future entity types

Responsibilities:

- Index searchable entities
- Execute semantic searches
- Return ranked search results

Search results expose:

- Identifier
- Kind
- Optional title
- Text
- Similarity score

The search system enables users and other applications to:

- Discover semantically similar posts
- Find related campaigns
- Explore topics and narratives
- Investigate coordinated messaging
- Identify similar users or communities
- Navigate the semantic graph through meaning rather than keywords

For example, a search for:

```text
Election fraud claims
```

may return results containing:

```text
Vote manipulation
Rigged election allegations
Ballot tampering accusations
Electoral corruption concerns
```

even if none of those results contain the exact words "election fraud claims".

Because all searchable entities share a common indexing mechanism, the search engine becomes a universal discovery system capable of connecting information across every layer of the platform.

This makes semantic search one of the core building blocks of AABS, providing the foundation for content exploration, investigation workflows, relationship discovery, campaign analysis, and future recommendation systems.

### Groupings Application

The groupings application provides access to semantic grouping subsystems.

Responsibilities:

- Access campaigns
- Access topics
- Access narratives
- Access participations
- Access clusters

#### Clusters Application

The clusters application groups semantically similar entities.

Clusters are the foundation of semantic analysis within AABS. They identify entities that are discussing similar subjects, promoting similar messages, exhibiting similar behaviors, or sharing similar semantic characteristics.

Examples:

- Posts discussing the same event
- Users repeatedly sharing similar content
- Communities promoting similar narratives
- Campaigns spreading related messages
- Topics covering closely related subjects

Clusters are used as the primary input for higher-level analysis such as campaign detection, topic extraction, narrative identification, participation calculations, relationship generation, and trust score computation.

Responsibilities:

- Build clusters for a target entity
- Find clusters by identifier
- Find clusters by target
- Find clusters by member
- Rebuild all clusters
- Rebuild clusters by entity type

Supported cluster rebuild operations:

- Posts
- Users
- Communities
- Platforms
- Campaigns
- Topics
- Narratives

Clusters are the foundation of higher-level semantic analysis.

#### Campaigns Application

The campaigns application identifies recurring semantic patterns.

Campaigns represent coordinated or recurring messaging discovered across clusters of semantically similar content. A campaign groups posts, users, communities, and other entities that appear to be promoting the same message, objective, or narrative.

Examples:

- Political influence operations
- Coordinated spam campaigns
- Product marketing campaigns
- Disinformation campaigns
- Reputation management campaigns

Campaigns provide a higher-level representation of coordinated activity within the semantic graph. They are derived from cluster analysis and are used to calculate participations, generate relationships, identify narratives, and contribute to trust score calculations.

Responsibilities:

- Find campaigns by identifier
- Find campaigns using index pagination
- Find campaigns using cursor pagination
- Find campaigns by user
- Find campaigns by community
- Find campaigns by platform
- Count campaigns
- Rebuild campaigns

#### Topics Application

The topics application identifies discussion subjects.

Topics represent the subjects being discussed within the platform, independent of any specific opinion, claim, or narrative. They organize content into semantic categories that help structure the overall conversation landscape.

Examples:

- Politics
- Elections
- Immigration
- Artificial Intelligence
- Cybersecurity
- Climate Change

Topics answer the question:

"What is being discussed?"

Unlike narratives, which describe specific claims or messages, topics describe the general subject matter of content.

Topics are derived from cluster analysis and are used to organize campaigns, narratives, participations, relationships, and trust score calculations within the semantic graph.

Responsibilities:

- Find topics by identifier
- Find topics using index pagination
- Find topics using cursor pagination
- Find topics by user
- Find topics by community
- Count topics
- Rebuild topics

#### Narratives Application

The narratives application identifies claims and stories being propagated.

Narratives represent the specific messages, claims, opinions, or stories being communicated across content. While topics describe the subject being discussed, narratives describe what is being said about that subject.

Examples:

- "Candidate X is corrupt"
- "Artificial intelligence will replace most jobs"
- "Product Y is dangerous"
- "Company Z is manipulating the market"
- "Government A is responsible for the crisis"

Narratives answer the question:

"What is being claimed?"

Multiple narratives may exist within the same topic. For example, a topic such as "Artificial Intelligence" may contain narratives that are supportive, critical, neutral, or contradictory.

Narratives are derived from cluster analysis and campaign detection. They provide a structured representation of the messages circulating through the platform and are used to calculate participations, generate relationships, identify coordinated activity, and contribute to trust score calculations.

Responsibilities:

- Find narratives by identifier
- Find narratives using index pagination
- Find narratives using cursor pagination
- Find narratives by user
- Find narratives by community
- Count narratives
- Rebuild narratives

#### Participations Application

The participations application measures how strongly entities contribute to other entities.

Participations quantify the relationship between a participant and a target within the semantic graph. They measure how much an entity contributes to, influences, or is associated with another entity based on observed content and behavior.

Examples:

- User → Campaign
- User → Topic
- User → Narrative
- Community → Campaign
- Community → Narrative
- Platform → Topic

Participations answer questions such as:

- "How strongly is this user associated with this campaign?"
- "How much does this community contribute to this narrative?"
- "Which platforms are most involved in this topic?"

Each participation is calculated from supporting evidence collected throughout the system and typically includes metrics such as:

- Matching content count
- Total analyzed content
- Participation percentage
- Evidence references

Participations serve as the connective layer of the semantic graph. They link entities together, provide explainability through evidence, support relationship generation, reveal patterns of coordinated activity, and contribute to trust score calculations.

Responsibilities:

- Access participation evidences
- Find participation by identifier
- Find participations by participant
- Find participations by target
- Find participation between two entities
- Rebuild participations

##### Participation Evidences Application

The participation evidences application exposes evidence supporting participation calculations.

Evidence provides transparency and explainability for participation measurements. Each evidence record represents a piece of content or observation that contributed to the calculation of a participation.

Examples:

- A post supporting a campaign
- A post contributing to a topic
- A post propagating a narrative
- A user interaction associated with a target entity
- Content linking a participant to a semantic grouping

Evidence answers questions such as:

- "Why is this user associated with this campaign?"
- "Which posts contributed to this narrative participation?"
- "What content supports this participation score?"

Evidence can be queried from multiple perspectives:

- By participation
- By post
- By participant
- By target

The evidence layer is a critical component of explainability within AABS. It allows analysts, researchers, and users to inspect the underlying data used to generate participations, verify conclusions, investigate coordinated activity, and understand how higher-level semantic relationships were derived.

Responsibilities:

- Find evidence by identifier
- Find evidence by participation
- Find evidence by post
- Find evidence by participant
- Find evidence by target

Evidence provides transparency and explainability for participation calculations.

### Relationships Application

The relationships application manages semantic relationships between entities.

Relationships represent semantic similarity and association between entities within the graph. They connect entities that exhibit related content, behavior, participation patterns, or semantic characteristics.

Examples:

- Campaign ↔ Campaign
- Topic ↔ Topic
- Narrative ↔ Narrative
- User ↔ User
- Community ↔ Community
- Campaign ↔ Narrative
- Topic ↔ Narrative

Relationships answer questions such as:

- "Which campaigns are promoting similar messages?"
- "Which narratives are closely related?"
- "Which users exhibit similar behavior?"
- "Which communities discuss similar subjects?"

Each relationship contains a similarity measurement derived from semantic comparison and may be supported by multiple observations throughout the graph.

Relationships form the connective tissue of the semantic graph. They allow the system to discover hidden associations, identify coordinated activity, reveal semantic proximity between entities, support graph exploration, and contribute to trust score calculations.

Responsibilities:

- Access comparables
- Build relationships
- Synchronize relationships
- Find relationships by identifier
- Find relationships using index pagination
- Find relationships using cursor pagination
- Count relationships
- Find relationships by source
- Find relationships by target
- Rebuild relationships

#### Comparables Application

The comparables application performs pairwise semantic comparison.

Comparables are responsible for measuring semantic similarity between entities. They analyze the characteristics, content, and semantic representations of two entities and determine how closely related they are.

Examples:

- Compare two campaigns
- Compare two topics
- Compare two narratives
- Compare two users
- Compare two communities
- Compare two platforms

Comparables answer questions such as:

- "How similar are these two campaigns?"
- "Do these narratives communicate the same message?"
- "Are these users participating in similar activities?"
- "Do these communities discuss the same subjects?"

The output of a comparison is a relationship containing a similarity score and any additional metadata required by the relationship graph.

Comparables serve as the semantic analysis engine behind relationships. By isolating similarity calculations into a dedicated application, comparison algorithms can evolve independently from relationship storage, querying, synchronization, and graph rebuilding processes.

Responsibilities:

- Compare two comparable entities
- Produce relationship results

Comparables isolate similarity calculations from relationship orchestration.

### Scores Application

The scores application calculates trust and risk metrics.

Scores provide quantitative assessments of entities within the semantic graph. They transform observations, relationships, participations, behaviors, and other signals into measurable indicators that help evaluate authenticity, influence, coordination, risk, and trustworthiness.

Examples:

- User trust scores
- Campaign risk scores
- Narrative credibility scores
- Community trust scores
- Relationship confidence scores
- Platform risk scores

Scores answer questions such as:

- "How trustworthy is this user?"
- "How likely is this campaign to represent coordinated activity?"
- "Which narratives present the highest risk?"
- "How reliable is this community?"

Scores are generated by one or more score calculators, each responsible for evaluating a specific aspect of an entity. Multiple score types can coexist and evolve independently as new detection techniques are introduced.

Examples of score inputs include:

- Participation patterns
- Relationship density
- Semantic repetition
- Campaign involvement
- Narrative propagation
- Community distribution
- Behavioral anomalies

The scores layer serves as the decision-support component of AABS. It converts complex graph structures into interpretable metrics that can be used for ranking, filtering, moderation, investigations, automated detection, and user-facing trust indicators.

Responsibilities:

- Calculate scores for a target
- Retrieve latest scores
- Retrieve score history
- Recalculate all scores

Multiple score calculators can be combined to produce trust and risk measurements.

## Processing Pipeline

When content enters the system:

```text
Post
 └─► Save
      └─► Index Searchable Content
            └─► Rebuild Clusters
                  └─► Rebuild Campaigns
                        └─► Rebuild Topics
                              └─► Rebuild Narratives
                                    └─► Rebuild Participations
                                          └─► Rebuild Relationships
                                                └─► Recalculate Scores
```

This process continuously updates the semantic graph as new content is ingested.

## Development

### Build and start all services

docker compose up --build 

### Build and start all services in background

docker compose up -d --build 

### View logs

docker compose logs -f 

### Stop all services

docker compose down 

## License

MIT