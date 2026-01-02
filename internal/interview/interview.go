package interview

import (
	"context"
	"database/sql"
	"log/slog"

	"prepf/internal/db"
	"prepf/internal/event"
	"prepf/internal/pubsub"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusCompleted  TodoStatus = "completed"
)

type Todo struct {
	Content    string     `json:"content"`
	Status     TodoStatus `json:"status"`
	ActiveForm string     `json:"active_form"`
}

type Interview struct {
	ID               string
	Title            string
	MessageCount     int64
	PromptTokens     int64
	CompletionTokens int64
	SummaryMessageID string
	Cost             float64
	Todos            []Todo
	Difficulty       string
	Topic            string
	Status           string
	CreatedAt        int64
	UpdatedAt        int64
}

type InterviewService interface {
	pubsub.Subscriber[Interview]
	Create(ctx context.Context, title string) (Interview, error)
	Get(ctx context.Context, id string) (Interview, error)
	List(ctx context.Context) ([]Interview, error)
	Save(ctx context.Context, interview Interview) (Interview, error)
	UpdateTitleAndUsage(ctx context.Context, interviewID, title string, promptTokens, completionTokens int64, cost float64) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	*pubsub.Broker[Interview]
	q db.Querier
}

func (s *service) Create(ctx context.Context, title string) (Interview, error) {
	id := uuid.New().String()
	dbInterview, err := s.q.CreateInterview(ctx, db.CreateInterviewParams{
		ID:    id,
		Title: title,
	})
	if err != nil {
		return Interview{}, err
	}
	interview := s.fromDBItem(dbInterview)
	s.Publish(pubsub.CreatedEvent, interview)
	event.InterviewCreated()
	return interview, nil
}

func (s *service) Get(ctx context.Context, id string) (Interview, error) {
	dbInterview, err := s.q.GetInterviewByID(ctx, id)
	if err != nil {
		return Interview{}, err
	}
	return s.fromDBItem(dbInterview), nil
}

func (s *service) List(ctx context.Context) ([]Interview, error) {
	dbInterviews, err := s.q.ListInterviews(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Map(dbInterviews, func(item db.Interview, _ int) Interview {
		return s.fromDBItem(item)
	}), nil
}

func (s *service) Save(ctx context.Context, interview Interview) (Interview, error) {
	todosJSON, err := marshalTodos(interview.Todos)
	if err != nil {
		return Interview{}, err
	}

	dbInterview, err := s.q.UpdateInterview(ctx, db.UpdateInterviewParams{
		ID:               interview.ID,
		Title:            interview.Title,
		MessageCount:     interview.MessageCount,
		PromptTokens:     interview.PromptTokens,
		CompletionTokens: interview.CompletionTokens,
		Cost:             interview.Cost,
		SummaryMessageID: sql.NullString{
			String: interview.SummaryMessageID,
			Valid:  interview.SummaryMessageID != "",
		},
		Todos: sql.NullString{
			String: todosJSON,
			Valid:  todosJSON != "",
		},
		Difficulty: sql.NullString{
			String: interview.Difficulty,
			Valid:  interview.Difficulty != "",
		},
		Topic: sql.NullString{
			String: interview.Topic,
			Valid:  interview.Topic != "",
		},
		Status: sql.NullString{
			String: interview.Status,
			Valid:  interview.Status != "",
		},
	})
	if err != nil {
		return Interview{}, err
	}
	interview = s.fromDBItem(dbInterview)
	s.Publish(pubsub.UpdatedEvent, interview)
	return interview, nil
}

// Only update title and usage metrics atomically
func (s *service) UpdateTitleAndUsage(ctx context.Context, interviewID, title string, promptTokens, completionTokens int64, cost float64) error {
	return s.q.UpdateInterviewTitleAndUsage(ctx, db.UpdateInterviewTitleAndUsageParams{
		ID:               interviewID,
		Title:            title,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		Cost:             cost,
	})
}

func (s *service) Delete(ctx context.Context, id string) error {
	interview, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	err = s.q.DeleteInterview(ctx, interview.ID)
	if err != nil {
		return err
	}
	s.Publish(pubsub.DeletedEvent, interview)
	event.InterviewDeleted()
	return nil
}

// -- Helpers

func (s service) fromDBItem(item db.Interview) Interview {
	todos, err := unmarshalTodos(item.Todos.String)
	if err != nil {
		slog.Error("failed to unmarshal interview todos", "interview_id", item.ID, "error", err)
	}
	return Interview{
		ID:               item.ID,
		Title:            item.Title,
		MessageCount:     item.MessageCount,
		PromptTokens:     item.PromptTokens,
		CompletionTokens: item.CompletionTokens,
		SummaryMessageID: item.SummaryMessageID.String,
		Cost:             item.Cost,
		Todos:            todos,
		Difficulty:       item.Difficulty.String,
		Topic:            item.Topic.String,
		Status:           item.Status.String,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}
}

func marshalTodos(todos []Todo) (string, error) {
	if len(todos) == 0 {
		return "", nil
	}
	data, err := sonic.Marshal(todos)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshalTodos(data string) ([]Todo, error) {
	if data == "" {
		return []Todo{}, nil
	}
	var todos []Todo
	if err := sonic.Unmarshal([]byte(data), &todos); err != nil {
		return []Todo{}, err
	}
	return todos, nil
}

func NewService(q db.Querier) InterviewService {
	broker := pubsub.NewBroker[Interview]()
	return &service{
		Broker: broker,
		q:      q,
	}
}
