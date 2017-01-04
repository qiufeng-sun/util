package gamerec

////////////////////////////////////////////////////////////////////////////////
//
func (this *Wrap) NextGames(userId string,
	count int32, params map[string]string) (string, error) {
	//
	r, e := this.Check(func() (interface{}, error) {
		return this.Service().PersonalGamePortal(userId, count, params)
	})

	if e != nil || nil == r {
		return "", e
	}

	return r.(string), nil
}

func NextGames(userId string, count int32, reset bool) (*HttpRecGames, error) {
	s, e := GetPoolClient()
	if e != nil {
		return nil, e
	}
	defer s.Return()

	paramMap := make(map[string]string)
	if reset {
		paramMap = g_paramMapReset
	}

	str, e := s.NextGames(userId, count, paramMap)
	if e != nil {
		return nil, e
	}

	return Json2HttpRecGames(str)
}
