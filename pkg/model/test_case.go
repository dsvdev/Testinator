package model

import (
	"github.com/google/uuid"
)

type TestCaseID uuid.UUID

type TestStep string

func NewID(uuid uuid.UUID) TestCaseID {
	return TestCaseID(uuid)
}

type TestCase struct {
	ID    TestCaseID
	Steps []TestStep
}

func NewTestCase(steps ...TestStep) *TestCase {
	return &TestCase{
		ID:    NewID(uuid.New()),
		Steps: steps,
	}
}
