/**
 * Интерфейс API, который должен предоставить хост-приложение для инициализации.
 * Обязательно должен содержать метод `getStartupData`, возвращающий данные старта приложения.
 */
type StartupApi<STARTUP_DATA> = {
	getStartupData: () => Promise<STARTUP_DATA>,
}

export type {
	StartupApi,
}
