from typing import List, Optional

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

import hdbscan
import numpy as np

app = FastAPI(title="AABS HDBSCAN Service")


class ClusterRequest(BaseModel):
    embeddings: List[List[float]] = Field(
        ...,
        description="List of embedding vectors",
    )

    min_cluster_size: int = 5
    min_samples: Optional[int] = None


class ClusterResult(BaseModel):
    index: int
    cluster_id: int
    probability: float
    is_noise: bool


class ClusterSummary(BaseModel):
    cluster_id: int
    size: int
    centroid: List[float]
    member_indexes: List[int]


class ClusterResponse(BaseModel):
    total_embeddings: int
    total_clusters: int
    noise_count: int
    results: List[ClusterResult]
    clusters: List[ClusterSummary]


@app.get("/health")
def health():
    return {
        "status": "ok",
        "service": "hdbscan",
    }


@app.post("/cluster", response_model=ClusterResponse)
def cluster(request: ClusterRequest):
    if len(request.embeddings) == 0:
        raise HTTPException(
            status_code=400,
            detail="embeddings cannot be empty",
        )

    dimensions = len(request.embeddings[0])

    for vector in request.embeddings:
        if len(vector) != dimensions:
            raise HTTPException(
                status_code=400,
                detail="all embeddings must have the same dimensions",
            )

    if len(request.embeddings) < request.min_cluster_size:
        raise HTTPException(
            status_code=400,
            detail="number of embeddings must be >= min_cluster_size",
        )

    data = np.array(
        request.embeddings,
        dtype=np.float32,
    )

    clusterer = hdbscan.HDBSCAN(
        min_cluster_size=request.min_cluster_size,
        min_samples=request.min_samples,
        metric="euclidean",
        prediction_data=True,
    )

    labels = clusterer.fit_predict(data)
    probabilities = clusterer.probabilities_

    results: List[ClusterResult] = []

    for index, label in enumerate(labels):
        results.append(
            ClusterResult(
                index=index,
                cluster_id=int(label),
                probability=float(probabilities[index]),
                is_noise=(label == -1),
            )
        )

    cluster_ids = set(labels)
    cluster_ids.discard(-1)

    clusters: List[ClusterSummary] = []

    for cluster_id in sorted(cluster_ids):
        member_indexes = np.where(
            labels == cluster_id
        )[0]

        cluster_vectors = data[member_indexes]

        centroid = np.mean(
            cluster_vectors,
            axis=0,
        )

        clusters.append(
            ClusterSummary(
                cluster_id=int(cluster_id),
                size=len(member_indexes),
                centroid=centroid.astype(np.float32).tolist(),
                member_indexes=member_indexes.tolist(),
            )
        )

    return ClusterResponse(
        total_embeddings=len(request.embeddings),
        total_clusters=len(cluster_ids),
        noise_count=int(np.sum(labels == -1)),
        results=results,
        clusters=clusters,
    )