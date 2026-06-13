package posts

import (
	"context"
	"errors"
	"testing"
)

func TestNewService(t *testing.T) {
	service := NewService()

	if service == nil {
		t.Fatalf("expected service")
	}
}

func TestServiceSaveExecutesSubServicesInOrder(t *testing.T) {
	ctx := context.Background()
	post := NewMockPost("hello")

	first := NewMockPostService()
	second := NewMockPostService()
	third := NewMockPostService()

	service := NewService(
		first,
		second,
		third,
	)

	err := service.Save(ctx, post)
	if err != nil {
		t.Fatal(err)
	}

	if first.SaveCalls != 1 {
		t.Fatalf("expected first service to be called")
	}

	if second.SaveCalls != 1 {
		t.Fatalf("expected second service to be called")
	}

	if third.SaveCalls != 1 {
		t.Fatalf("expected third service to be called")
	}

	if first.LastContext != ctx ||
		second.LastContext != ctx ||
		third.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if first.LastPost != post ||
		second.LastPost != post ||
		third.LastPost != post {
		t.Fatalf("expected post to be passed")
	}
}

func TestServiceSaveSkipsNilServices(t *testing.T) {
	ctx := context.Background()
	post := NewMockPost("hello")

	first := NewMockPostService()
	second := NewMockPostService()

	service := NewService(
		first,
		nil,
		second,
	)

	err := service.Save(ctx, post)
	if err != nil {
		t.Fatal(err)
	}

	if first.SaveCalls != 1 {
		t.Fatalf("expected first service to be called")
	}

	if second.SaveCalls != 1 {
		t.Fatalf("expected second service to be called")
	}
}

func TestServiceSaveReturnsError(t *testing.T) {
	ctx := context.Background()
	post := NewMockPost("hello")

	first := NewMockPostService()
	second := NewMockPostService()
	third := NewMockPostService()

	second.SaveErr = errTest

	service := NewService(
		first,
		second,
		third,
	)

	err := service.Save(ctx, post)
	if !errors.Is(err, errTest) {
		t.Fatalf("expected error, got %v", err)
	}

	if first.SaveCalls != 1 {
		t.Fatalf("expected first service to be called")
	}

	if second.SaveCalls != 1 {
		t.Fatalf("expected second service to be called")
	}

	if third.SaveCalls != 0 {
		t.Fatalf("expected third service not to be called")
	}
}

func TestServiceSaveWithNoSubServices(t *testing.T) {
	ctx := context.Background()
	post := NewMockPost("hello")

	service := NewService()

	err := service.Save(ctx, post)
	if err != nil {
		t.Fatal(err)
	}
}
