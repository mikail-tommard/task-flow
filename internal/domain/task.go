package domain

type Task struct {
	id          int
	title       string
	description string
	done        bool
	userID      int
}

func New(title string, description string, userID int) (*Task, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}

	return &Task{
		id:          0,
		title:       title,
		description: description,
		done:        false,
		userID:      userID,
	}, nil
}

func FromStorage(id int, title string, done bool, description string, userID int) *Task {
	return &Task{
		id:          id,
		title:       title,
		description: description,
		done:        done,
		userID:      userID,
	}
}

func (t *Task) Complete() {
	t.done = true
}

func (t *Task) ID() int {
	return t.id
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

func (t *Task) UserID() int {
	return t.userID
}

func (t *Task) Done() bool {
	return t.done
}

func (t *Task) Rename(title string) error {
	if title == "" {
		return ErrInvalidTitle
	}

	t.title = title

	return nil
}

func (t *Task) ChangeDescription(description string) {
	t.description = description
}
