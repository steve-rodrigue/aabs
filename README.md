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

Responsibilities:

- Process a single post
- Process multiple posts
- Rebuild semantic artifacts

The pipeline is the primary entry point for newly discovered content.

### Posts Application

The posts application manages content storage and retrieval.

Responsibilities:

- Save posts
- Find posts by identifier
- Find posts using index pagination
- Find posts using cursor pagination
- Count posts
- Find posts by user
- Find posts by community
- Find posts by platform

### Users Application

The users application manages platform users.

Responsibilities:

- Save users
- Find users by identifier
- Find users by platform-specific identifier
- Find users using index pagination
- Find users using cursor pagination
- Count users

### Communities Application

The communities application manages communities within platforms.

Responsibilities:

- Save communities
- Find communities by identifier
- Find communities by platform and handle
- Find communities using index pagination
- Find communities using cursor pagination
- Find communities by platform
- Count communities

### Platforms Application

The platforms application manages content platforms.

Responsibilities:

- Save platforms
- Find platforms by identifier
- Find platforms by handle
- Find platforms using index pagination
- Find platforms using cursor pagination
- Count platforms

### Searches Application

The searches application provides semantic search.

Responsibilities:

- Index searchable entities
- Execute semantic searches
- Return ranked search results

Any object implementing the searchable interface can be indexed.

Search results expose:

- Identifier
- Kind
- Optional title
- Text
- Similarity score

The search system is intentionally generic and independent from specific entity types.

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

Responsibilities:

- Access participation evidences
- Find participation by identifier
- Find participations by participant
- Find participations by target
- Find participation between two entities
- Rebuild participations

##### Participation Evidences Application

The participation evidences application exposes evidence supporting participation calculations.

Responsibilities:

- Find evidence by identifier
- Find evidence by participation
- Find evidence by post
- Find evidence by participant
- Find evidence by target

Evidence provides transparency and explainability for participation calculations.

### Relationships Application

The relationships application manages semantic relationships between entities.

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

Responsibilities:

- Compare two comparable entities
- Produce relationship results

Comparables isolate similarity calculations from relationship orchestration.

### Scores Application

The scores application calculates trust and risk metrics.

Responsibilities:

- Calculate scores for a target
- Retrieve latest scores
- Retrieve score history
- Recalculate all scores

Multiple score calculators can be combined to produce trust and risk measurements.

## Processing Pipeline

When content enters the system:

text Post  │  ▼ Save  │  ▼ Index Searchable Content  │  ▼ Rebuild Clusters  │  ▼ Rebuild Campaigns  │  ▼ Rebuild Topics  │  ▼ Rebuild Narratives  │  ▼ Rebuild Participations  │  ▼ Rebuild Relationships  │  ▼ Recalculate Scores 

This process continuously updates the semantic graph as new content is ingested.

## Development

### Build and start all services

bash docker compose up --build 

### Build and start all services in background

bash docker compose up -d --build 

### View logs

bash docker compose logs -f 

### Stop all services

bash docker compose down 

## License

MIT