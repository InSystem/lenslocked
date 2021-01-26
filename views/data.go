package views

import (
	"log"

	"github.com/InSystem/lenslocked/models"
	// "log"
)

const (
	AlertLevelError   = "danger"
	AlertLevelWarning = "warning"
	AlertLevelInfo    = "info"
	AlertLevelSuccess = "success"

	//AlertMessageGeneric is display any user when error is encountered
	AlertMessageGeneric = "Something went wrong, please try again and contact us if the problem persits"
)

//Alert is used to render Bootstrap Alert messages in the bootstrap.gohtml templates
type Alert struct {
	Level   string
	Message string
}

//Data is the top level structure that views expect data to come in
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: pErr.Public(),
		}
	} else {
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLevelError,
			Message: AlertMessageGeneric,
		}
	}
}

// func (d *Data) AlertError(message string){
// 	d.Alert = &Alert{
// 		Level: AlertLevelError,
// 		Message: message,
// 	}
// }

type PublicError interface {
	error
	Public() string
}
