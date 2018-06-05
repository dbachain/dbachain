package project

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"time"
)

func TestNewSendMsg(t *testing.T) {}

func TestProjectMsgType(t *testing.T) {
	// Construct a ProjectMsg
	var msg = ProjectMsg{
		Address: nil,
		MsgType: "create",
		ID:      "project-id-1",
		Status:  "3",
		CTime:   time.Now().Unix(),
	}

	assert.Equal(t, msg.Type(), "project")
}

func TestValidateBasic(t *testing.T) {
	// Construct a ProjectMsg
	var msg = ProjectMsg{
		Address: nil,
		MsgType: "create",
		ID:      "project-id-1",
		Status:  "3",
		CTime:   time.Now().Unix(),
	}

	assert.Equal(t, msg.ValidateBasic(), nil)
}
