package model

// SetID implements dal.HasID interface
func (o *Notification) SetID(id string) {
	o.ID = id
}
