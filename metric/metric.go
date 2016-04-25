package metric

type Updater interface {
	Update()
}

type Saver interface {
	Save()
}

type Metric interface {
	Updater
	Saver
}
