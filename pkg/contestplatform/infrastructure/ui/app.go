package ui

import appmodel "contest-platform/pkg/contestplatform/app/model"

type StartupData struct {
	Title     string                `json:"title"`
	Languages []appmodel.UILanguage `json:"languages"`
	Tasks     []Task                `json:"tasks"`
}

type Task struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
}
