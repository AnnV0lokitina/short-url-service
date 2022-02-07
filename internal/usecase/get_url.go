package usecase

func (u *Usecase) GetUrl(uuid string) (string, string, error) {
	collection := u.repo.GetInfo()
	url, err := collection.Get(uuid)
	if err != nil {
		return "", "", err
	}
	return url.Full, url.Short, nil
}
