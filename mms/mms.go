/*
 * Copyright 2014 Canonical Ltd.
 *
 * Authors:
 * Sergio Schvezov: sergio.schvezov@cannical.com
 *
 * This file is part of mms.
 *
 * mms is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 3.
 *
 * mms is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package mms

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// MMS Field names from OMA-WAP-MMS section 7.3 Table 12
const (
	BCC                           = 0x01
	CC                            = 0x02
	X_MMS_CONTENT_LOCATION        = 0x03
	CONTENT_TYPE                  = 0x04
	DATE                          = 0x05
	X_MMS_DELIVERY_REPORT         = 0x06
	X_MMS_DELIVERY_TIME           = 0x07
	X_MMS_EXPIRY                  = 0x08
	FROM                          = 0x09
	X_MMS_MESSAGE_CLASS           = 0x0A
	MESSAGE_ID                    = 0x0B
	X_MMS_MESSAGE_TYPE            = 0x0C
	X_MMS_MMS_VERSION             = 0x0D
	X_MMS_MESSAGE_SIZE            = 0x0E
	X_MMS_PRIORITY                = 0x0F
	X_MMS_READ_REPORT             = 0x10
	X_MMS_REPORT_ALLOWED          = 0x11
	X_MMS_RESPONSE_STATUS         = 0x12
	X_MMS_RESPONSE_TEXT           = 0x13
	X_MMS_SENDER_VISIBILITY       = 0x14
	X_MMS_STATUS                  = 0x15
	SUBJECT                       = 0x16
	TO                            = 0x17
	X_MMS_TRANSACTION_ID          = 0x18
	X_MMS_RETRIEVE_STATUS         = 0x19
	X_MMS_RETRIEVE_TEXT           = 0x1A
	X_MMS_READ_STATUS             = 0x1B
	X_MMS_REPLY_CHARGING          = 0x1C
	X_MMS_REPLY_CHARGING_DEADLINE = 0x1D
	X_MMS_REPLY_CHARGING_ID       = 0x1E
	X_MMS_REPLY_CHARGING_SIZE     = 0x1F
	X_MMS_PREVIOUSLY_SENT_BY      = 0x20
	X_MMS_PREVIOUSLY_SENT_DATE    = 0x21
)

// MMS Content Type Assignments OMA-WAP-MMS section 7.3 Table 13
const (
	PUSH_APPLICATION_ID = 4
	VND_WAP_MMS_MESSAGE = "application/vnd.wap.mms-message"
)

const (
	TYPE_SEND_REQ         = 0x80
	TYPE_SEND_CONF        = 0x81
	TYPE_NOTIFICATION_IND = 0x82
	TYPE_NOTIFYRESP_IND   = 0x83
	TYPE_RETRIEVE_CONF    = 0x84
	TYPE_ACKNOWLEDGE_IND  = 0x85
	TYPE_DELIVERY_IND     = 0x86
)

const (
	MMS_MESSAGE_VERSION_1_0 = 0x90
	MMS_MESSAGE_VERSION_1_1 = 0x91
	MMS_MESSAGE_VERSION_1_2 = 0x92
	MMS_MESSAGE_VERSION_1_3 = 0x93
)

// Date tokens defined in OMA-WAP-MMS section 7.2.10
const (
	TOKEN_DATE_ABS = 0x80
	TOKEN_DATE_REL = 0x81
)

// From tokens defined in OMA-WAP-MMS section 7.2.11
const (
	TOKEN_ADDRESS_PRESENT = 0x80
	TOKEN_INSERT_ADDRESS  = 0x81
)

// Message classes defined in OMA-WAP-MMS section 7.2.14
const (
	CLASS_PERSONAL      = 0x80
	CLASS_ADVERTISEMENT = 0x81
	CLASS_INFORMATIONAL = 0x82
	CLASS_AUTO          = 0x83
)

// Report Allowed defined in OMA-WAP-MMS 7.2.19
const (
	REPORT_ALLOWED_YES = 128
	REPORT_ALLOWED_NO  = 129
)

// Response Status defined in OMA-WAP-MMS section 7.2.27
//
// An MMS Client MUST react the same to a value in range 196 to 223 as it
// does to the value 192 (Error-transient-failure).
//
// An MMS Client MUST react the same to a value in range 234 to 255 as it
// does to the value 224 (Error-permanent-failure).
//
// Any other values SHALL NOT be used. They are reserved for future use.
// An MMS Client that receives such a reserved value MUST react the same
// as it does to the value 224 (Error-permanent-failure).
const (
	ResponseStatusOk                            byte = 128
	ResponseStatusErrorUnspecified              byte = 129 // Obsolete
	ResponseStatusErrorServiceDenied            byte = 130 // Obsolete
	ResponseStatusErrorMessageFormatCorrupt     byte = 131 // Obsolete
	ResponseStatusErrorSendingAddressUnresolved byte = 132 // Obsolete
	ResponseStatusErrorMessageNotFound          byte = 133 // Obsolete
	ResponseStatusErrorNetworkProblem           byte = 134 // Obsolete
	ResponseStatusErrorContentNotAccepted       byte = 135 // Obsolete
	ResponseStatusErrorUnsupportedMessage       byte = 136

	ResponseStatusErrorTransientFailure           byte = 192
	ResponseStatusErrorTransientAddressUnresolved byte = 193
	ResponseStatusErrorTransientMessageNotFound   byte = 194
	ResponseStatusErrorTransientNetworkProblem    byte = 195

	ResponseStatusErrorTransientMaxReserved byte = 223

	ResponseStatusErrorPermanentFailure                         byte = 224
	ResponseStatusErrorPermanentServiceDenied                   byte = 225
	ResponseStatusErrorPermanentMessageFormatCorrupt            byte = 226
	ResponseStatusErrorPermanentAddressUnresolved               byte = 227
	ResponseStatusErrorPermanentMessageNotFound                 byte = 228
	ResponseStatusErrorPermanentContentNotAccepted              byte = 229
	ResponseStatusErrorPermanentReplyChargingLimitationsNotMet  byte = 230
	ResponseStatusErrorPermanentReplyChargingRequestNotAccepted byte = 231
	ResponseStatusErrorPermanentReplyChargingForwardingDenied   byte = 232
	ResponseStatusErrorPermanentReplyChargingNotSupported       byte = 233

	ResponseStatusErrorPermamentMaxReserved byte = 255
)

// Status defined in OMA-WAP-MMS section 7.2.23
const (
	STATUS_EXPIRED      = 128
	STATUS_RETRIEVED    = 129
	STATUS_REJECTED     = 130
	STATUS_DEFERRED     = 131
	STATUS_UNRECOGNIZED = 132
)

// MSendReq holds a m-send.req message defined in
// OMA-WAP-MMS-ENC-v1.1 section 6.1.1
type MSendReq struct {
	UUID             string `encode:"no"`
	Type             byte
	TransactionId    string
	Version          byte
	Date             uint64 `encode:"optional"`
	From             string
	To               string
	Cc               string `encode:"no"`
	Bcc              string `encode:"no"`
	Subject          string `encode:"optional"`
	Class            byte   `encode:"optional"`
	Expiry           uint64 `encode:"optional"`
	DeliveryTime     uint64 `encode:"optional"`
	Priority         byte   `encode:"optional"`
	SenderVisibility byte   `encode:"optional"`
	DeliveryReport   byte   `encode:"optional"`
	ReadReply        byte   `encode:"optional"`
	ContentType      string
	Attachments      []*Attachment `encode:"no"`
}

// MSendReq holds a m-send.conf message defined in
// OMA-WAP-MMS-ENC section 6.1.2
type MSendConf struct {
	Type           byte
	TransactionId  string
	Version        byte
	ResponseStatus byte
	ResponseText   string
	MessageId      string
}

// MNotificationInd holds a m-notification.ind message defined in
// OMA-WAP-MMS-ENC section 6.2
type MNotificationInd struct {
	MMSReader
	UUID                                 string
	Type, Version, Class, DeliveryReport byte
	ReplyCharging, ReplyChargingDeadline byte
	ReplyChargingId                      string
	TransactionId, ContentLocation       string
	From, Subject                        string
	Expiry, Size                         uint64
}

// MNotificationInd holds a m-notifyresp.ind message defined in
// OMA-WAP-MMS-ENC-v1.1 section 6.2
type MNotifyRespInd struct {
	UUID          string `encode:"no"`
	Type          byte
	TransactionId string
	Version       byte
	Status        byte
	ReportAllowed bool
}

// MRetrieveConf holds a m-retrieve.conf message defined in
// OMA-WAP-MMS-ENC-v1.1 section 6.3
type MRetrieveConf struct {
	MMSReader
	UUID                                       string
	Type, Version, Status, Class, Priority     byte
	ReplyCharging, ReplyChargingDeadline       byte
	ReplyChargingId                            string
	ReadReport, RetrieveStatus, DeliveryReport byte
	TransactionId, MessageId, RetrieveText     string
	From, Cc, Subject                          string
	To                                         string
	ReportAllowed                              bool
	Date                                       uint64
	Content                                    Attachment
	Attachments                                []Attachment
	Data                                       []byte
}

type MMSReader interface{}
type MMSWriter interface{}

func NewMSendReq(recipients []string, attachments []*Attachment) *MSendReq {
	for i := range recipients {
		recipients[i] += "/TYPE=PLMN"
	}
	uuid := genUUID()
	return &MSendReq{
		Type:          TYPE_SEND_REQ,
		To:            strings.Join(recipients, ","),
		TransactionId: uuid,
		Version:       MMS_MESSAGE_VERSION_1_3,
		UUID:          uuid,
		ContentType:   "application/vnd.wap.multipart.related",
		Attachments:   attachments,
	}
}

func NewMSendConf() *MSendConf {
	return &MSendConf{
		Type: TYPE_SEND_CONF,
	}
}

func NewMNotificationInd() *MNotificationInd {
	return &MNotificationInd{Type: TYPE_NOTIFICATION_IND, UUID: genUUID()}
}

func (mNotificationInd *MNotificationInd) NewMNotifyRespInd(status byte, deliveryReport bool) *MNotifyRespInd {
	return &MNotifyRespInd{
		Type:          TYPE_NOTIFYRESP_IND,
		UUID:          mNotificationInd.UUID,
		TransactionId: mNotificationInd.TransactionId,
		Version:       mNotificationInd.Version,
		Status:        status,
		ReportAllowed: deliveryReport,
	}
}

func (mRetrieveConf *MRetrieveConf) NewMNotifyRespInd(deliveryReport bool) *MNotifyRespInd {
	return &MNotifyRespInd{
		Type:          TYPE_NOTIFYRESP_IND,
		UUID:          mRetrieveConf.UUID,
		TransactionId: mRetrieveConf.TransactionId,
		Version:       mRetrieveConf.Version,
		Status:        STATUS_RETRIEVED,
		ReportAllowed: deliveryReport,
	}
}

func NewMNotifyRespInd() *MNotifyRespInd {
	return &MNotifyRespInd{Type: TYPE_NOTIFYRESP_IND}
}

func NewMRetrieveConf(uuid string) *MRetrieveConf {
	return &MRetrieveConf{Type: TYPE_RETRIEVE_CONF, UUID: uuid}
}

func genUUID() string {
	var id string
	random, err := os.Open("/dev/urandom")
	if err != nil {
		id = "1234567890ABCDEF"
	} else {
		defer random.Close()
		b := make([]byte, 16)
		random.Read(b)
		id = fmt.Sprintf("%x", b)
	}
	return id
}

var ErrTransient = errors.New("Error-transient-failure")
var ErrPermanent = errors.New("Error-permament-failure")

func (mSendConf *MSendConf) Status() error {
	s := mSendConf.ResponseStatus
	// these are case by case Response Status and we need to determine each one
	switch s {
	case ResponseStatusOk:
		return nil
	case ResponseStatusErrorUnspecified:
		return ErrTransient
	case ResponseStatusErrorServiceDenied:
		return ErrTransient
	case ResponseStatusErrorMessageFormatCorrupt:
		return ErrPermanent
	case ResponseStatusErrorSendingAddressUnresolved:
		return ErrPermanent
	case ResponseStatusErrorMessageNotFound:
		// this could be ErrTransient or ErrPermanent
		return ErrPermanent
	case ResponseStatusErrorNetworkProblem:
		return ErrTransient
	case ResponseStatusErrorContentNotAccepted:
		return ErrPermanent
	case ResponseStatusErrorUnsupportedMessage:
		return ErrPermanent
	}

	// these are the Response Status we can group
	if s >= ResponseStatusErrorTransientFailure && s <= ResponseStatusErrorTransientMaxReserved {
		return ErrTransient
	} else if s >= ResponseStatusErrorPermanentFailure && s <= ResponseStatusErrorPermamentMaxReserved {
		return ErrPermanent
	}

	// any case not handled is a permanent error
	return ErrPermanent
}
