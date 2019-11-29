package services

import (
	"efficient-api/domain"
	"efficient-api/utils/error_utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	tm = time.Now()
	getMessageDomain func(messageId int64) (*domain.Message, error_utils.MessageErr)
	createMessageDomain func(msg *domain.Message) (*domain.Message, error_utils.MessageErr)
)

type getDBMock struct {}

func (m *getDBMock) Get(messageId int64) (*domain.Message, error_utils.MessageErr){
	return getMessageDomain(messageId)
}
func (m *getDBMock) Create(msg *domain.Message) (*domain.Message, error_utils.MessageErr){
	return createMessageDomain(msg)
}
func (m *getDBMock) Initialize(string, string, string, string, string, string){}


///////////////////////////////////////////////////////////////
// Start of "GetMessage" test cases
///////////////////////////////////////////////////////////////
func TestMessagesService_GetMessage_Success(t *testing.T) {
	domain.MessageRepo = &getDBMock{} //this is where we swapped the functionality
	getMessageDomain = func(messageId int64) (*domain.Message, error_utils.MessageErr) {
		return &domain.Message{
			Id:        1,
			Title:     "the title",
			Body:      "the body",
			CreatedAt: tm,
		}, nil
	}
	msg, err := MessagesService.GetMessage(1)
	fmt.Println("this is the message: ", msg)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Id)
	assert.EqualValues(t, "the title", msg.Title)
	assert.EqualValues(t, "the body", msg.Body)
	assert.EqualValues(t, tm, msg.CreatedAt)
}

//Test the not found functionality
func TestMessagesService_GetMessage_NotFoundID(t *testing.T) {
	domain.MessageRepo = &getDBMock{}

	getMessageDomain = func(messageId int64) (*domain.Message, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("the id is not found")
	}
	msg, err := MessagesService.GetMessage(1)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "the id is not found", err.Message())
	assert.EqualValues(t, "not_found", err.Error())
}
///////////////////////////////////////////////////////////////
// End of "GetMessage" test cases
///////////////////////////////////////////////////////////////


///////////////////////////////////////////////////////////////
// Start of	"CreateMessage" test cases
///////////////////////////////////////////////////////////////

//Here we call the domain method, so we must mock it
func TestMessagesService_CreateMessage_Success(t *testing.T) {
	domain.MessageRepo = &getDBMock{}
	createMessageDomain  = func(msg *domain.Message) (*domain.Message, error_utils.MessageErr){
		return &domain.Message{
			Id:        1,
			Title:     "the title",
			Body:      "the body",
			CreatedAt: tm,
		}, nil
	}
	request := &domain.Message{
		Title:     "the title",
		Body:      "the body",
		CreatedAt: tm,
	}
	msg, err := MessagesService.CreateMessage(request)
	fmt.Println("this is the message: ", msg)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.Id)
	assert.EqualValues(t, "the title", msg.Title)
	assert.EqualValues(t, "the body", msg.Body)
	assert.EqualValues(t, tm, msg.CreatedAt)
}

//This is a table test that check both the title and the body
//Since this will never call the domain "Get" method, no need to mock that method here
func TestMessagesService_CreateMessage_Invalid_Request(t *testing.T) {
	tests := []struct {
		request *domain.Message
		statusCode int
		errMsg string
		errErr string
	}{
		{
			request: &domain.Message{
			  Title:     "",
			  Body:      "the body",
			  CreatedAt: tm,
		    },
		    statusCode: http.StatusUnprocessableEntity,
		    errMsg: "Please enter a valid title",
		    errErr: "invalid_request",
		},
		{
			request: &domain.Message{
				Title:     "the title",
				Body:      "",
				CreatedAt: tm,
			},
			statusCode: http.StatusUnprocessableEntity,
			errMsg: "Please enter a valid body",
			errErr: "invalid_request",
		},
	}
	for _, tt := range tests {
		msg, err := MessagesService.CreateMessage(tt.request)
		assert.Nil(t, msg)
		assert.NotNil(t, err)
		assert.EqualValues(t, tt.errMsg, err.Message())
		assert.EqualValues(t, tt.statusCode, err.Status())
		assert.EqualValues(t, tt.errErr, err.Error())
	}
}

//We mock the "Get" method in the domain here. What could go wrong?,
//Since the title of the message must be unique, an error must be thrown,
//Of course you can also mock when the sql query is wrong, etc(these where covered in the domain tests),
//For now, we have 100% coverage on the "CreateMessage" method in the service
func TestMessagesService_CreateMessage_Failure(t *testing.T) {
	domain.MessageRepo = &getDBMock{}
	createMessageDomain  = func(msg *domain.Message) (*domain.Message, error_utils.MessageErr){
		return nil, error_utils.NewInternalServerError("title already taken")
	}
	request := &domain.Message{
		Title:     "the title",
		Body:      "the body",
		CreatedAt: tm,
	}
	msg, err := MessagesService.CreateMessage(request)
	assert.Nil(t, msg)
	assert.NotNil(t, err)
	assert.EqualValues(t, "title already taken", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}

///////////////////////////////////////////////////////////////
// End of "CreateMessage" test cases
///////////////////////////////////////////////////////////////