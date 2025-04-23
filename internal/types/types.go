package types

/*
import (
	"database/sql/driver"
	"errors"
	"github.com/google/uuid"
)

type ChatID uuid.UUID

func (c ChatID) String() string {
	return uuid.UUID(c).String()
}
//TextMarshaler реализует интерфейс encoding.TextMarshaler
func (c ChatID) MarshalText() ([]byte, error) {
	return uuid.UUID(c).MarshalText()
}

//TextUnmarshaler реализует интерфейс encoding.TextUnmarshaler
func (c *ChatID) UnmarshalText(text []byte) error {
	return (*uuid.UUID)(c).UnmarshalText(text)
}

//ValueScanner реализует интерфейс entfield.ValueScanner
//из двух методов: Scan и Value
func (c *ChatID) Scan(src interface{}) error {
	return (*uuid.UUID)(c).Scan(src)
}

func (c ChatID) Value() (driver.Value, error) {
	return uuid.UUID(c).Value()
}

//Validator реализует интерфейс entfield.Validator
func (c ChatID) Validate() error {
	if c == ChatIDNil {
		return errors.New("chat id is nil")
	}
	return nil
}

//Matcher реализует интерфейс gomock.Matcher
func (c1 ChatID) Matches(x interface{}) bool {
	c2, ok := x.(ChatID)
	if !ok {
		return false
	}
	return c1 == c2
}


//NewChatID создает новый ChatID
func NewChatID() ChatID {
	return ChatID(uuid.New())
}

//ChatIDNil это nil ChatID
var ChatIDNil = ChatID(uuid.Nil)

//Parse парсит строку и возвращает ChatID
func Parse[T any](s string) (ChatID, error) {
	var zero ChatID
	uuid, err := uuid.Parse(s)
	if err != nil {
		return zero, err
	}
	return ChatID(uuid), nil
}

//MustParse парсит строку и возвращает ChatID
func MustParse[T any](s string) ChatID {
	uuid, err := uuid.Parse(s)
	if err != nil {
		panic(err)
	}
	return ChatID(uuid)
}

//IsZero проверяет, является ли ChatID нулевым
func (c ChatID) IsZero() bool {
	return c == ChatIDNil
}

//NewMessageID создает новый MessageID
func NewMessageID() ChatID {
	return ChatID(uuid.New())
}
*/
