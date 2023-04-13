package db

type Device struct {
	Model    string
	Customer string
	Vendor   string
	Version  string
}

func (d *Database) CreateDevice(device Device) error {
	_, err := d.devices.InsertOne(d.ctx, device, nil)
	return err
}

func (d *Database) RetrieveDevice() {

}

func (d *Database) DeleteDevice() {

}
