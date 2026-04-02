package config

type Theme string

const (
	ThemePizza Theme = "pizza"
	ThemeSOC   Theme = "soc"
)

type ThemeTask struct {
	ID    string
	Label string
}

type ThemeMeta struct {
	Key   Theme
	Title string
	Tasks []ThemeTask
}

func AllThemes() []ThemeMeta {
	return []ThemeMeta{
		{
			Key:   ThemePizza,
			Title: "RocketSlice Dispatch",
			Tasks: []ThemeTask{
				{ID: "pizza-task-1", Label: "Фильтр заказов"},
				{ID: "pizza-task-2", Label: "Карта спроса"},
				{ID: "pizza-task-3", Label: "Генератор маршрутов"},
			},
		},
		{
			Key:   ThemeSOC,
			Title: "Enterprise SOC",
			Tasks: []ThemeTask{
				{ID: "soc-task-1", Label: "Threat Pipeline"},
				{ID: "soc-task-2", Label: "Threat Graph"},
				{ID: "soc-task-3", Label: "WAF Policies"},
			},
		},
	}
}
