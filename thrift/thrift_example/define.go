package gamerec

import (
	"encoding/json"
	. "util/thrift/thrift_example/gen-go/gamerec"
)

//
var g_paramMapReset = map[string]string{
	"firstSession": "1",
}

//
type RecGame struct {
	PkgName string `json:"id"`
	Reason  string `json:"reason"`
	TraceId string `json:"traceId"`
}

//
type HttpRecGames struct {
	Status  int        `json:"status"`
	Count   int        `json:"count"`
	RecData []*RecGame `json:"data"`
}

//
type RecSubject struct {
	SubId    string   `json:"id"`
	PkgNames []string `json:"pkgNames"`
	TraceIds string   `json:"traceId"`
	Score    float64  `json:"score"`
}

//
type HttpRecSubjects struct {
	Status  int           `json:"status"`
	Count   int           `json:"count"`
	RecData []*RecSubject `json:"data"`
}

//
type HttpGameMateRec struct {
	Status  int            `json:"status"`
	Count   int            `json:"count"`
	RecData []*GameMateRec `json:"data"`
}

//
type GameMateRec struct {
	Id       string   `json:"id"`
	PkgNames []string `json:"pkgNames"`
	TraceId  string   `json:"traceId"`
	Score    float64  `json:"score"`
}

//
func Json2HttpRecGames(str string) (*HttpRecGames, error) {
	var rec *HttpRecGames
	if e := json.Unmarshal([]byte(str), &rec); e != nil {
		return nil, e
	}

	return rec, nil
}

//
func Json2HttpRecSubjects(str string) (*HttpRecSubjects, error) {
	var rec *HttpRecSubjects
	if e := json.Unmarshal([]byte(str), &rec); e != nil {
		return nil, e
	}

	return rec, nil
}

//
func Json2HttpGameMetaRec(str string) (*HttpGameMateRec, error) {
	var rec *HttpGameMateRec
	if e := json.Unmarshal([]byte(str), &rec); e != nil {
		return nil, e
	}

	return rec, nil
}

type MetaReqests []*MetaRequest

//
func NewMetaRequests() MetaReqests {
	return make([]*MetaRequest, 0)
}

//
func (this *MetaReqests) AppendReq(RecType int64, count int32) {
	var metaMetaReq = &MetaRequest{
		MetaType(RecType),
		count,
	}

	(*this) = append((*this), metaMetaReq)
}

//
type MetaRecGame map[MetaType]*HttpGameMateRec

//
func NewMetaRecGame() *MetaRecGame {
	return &MetaRecGame{}
}

//
func (this *MetaRecGame) SetHttpGameMateRec(field MetaType, data *HttpGameMateRec) {
	rec := map[MetaType]*HttpGameMateRec(*this)
	rec[field] = data
}

//
func (this *MetaRecGame) GetGameSubject() *HttpGameMateRec {
	rec := map[MetaType]*HttpGameMateRec(*this)
	if data, ok := rec[MetaType_GAME_SUBJECT]; ok {
		return data
	}
	return nil
}

//
func (this *MetaRecGame) GetGame() *HttpGameMateRec {
	rec := map[MetaType]*HttpGameMateRec(*this)
	if data, ok := rec[MetaType_GAME]; ok {
		return data
	}
	return nil
}

//
func (this *MetaRecGame) GetMiGame() *HttpGameMateRec {
	rec := map[MetaType]*HttpGameMateRec(*this)
	if data, ok := rec[MetaType_MI_GAME]; ok {
		return data
	}
	return nil
}
