package errs

// Ошибки приложения
const (
	ErrInvalidID           = "Некорректный формат ID"
	ErrNotFound            = "Ресурс не найден"
	ErrUserIsAuthorTasks   = "Не возможно удалить Юзера, который создал задачи"
	ErrUserIsExecutorTasks = "Не возможно удалить Юзера, у которого есть назначенные задачи"
	ErrEmptyBody           = "Нет полей для обновления"
	ErrInvalidBody         = "Некорректное тело запроса"
	ErrAuthorNotExist      = "Автор - несуществующий юзер"
	ErrExecutorNotExist    = "Исполнитель - несуществующий юзер"
	ErrInternal            = "Internal server error"
)
