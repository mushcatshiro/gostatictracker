package workflow

func CreateNewReminder() {}

// recent is an input range, the default is End < + 14, status = INPROGRESS
func ReadRecentReminders() {}

// only update the repeat field
func UpdateReminderSettings() {}

// only moves one step ahead. e.g. weekly reminder is missed twice
// user had to clear it twice; this is intentionally designed as such
func CompleteCurrentTrigger() {}
