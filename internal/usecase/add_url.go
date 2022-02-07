package usecase

func (u *Usecase) SetURL(full string) (string, string) {
	collection := u.repo.GetInfo()
	uuid, url := collection.Add(full)
	u.repo.SetInfo(collection)

	return uuid, url.Short
}
