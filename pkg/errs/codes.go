package errs

// Ошибки приложения
const (
	ErrInvalidID                    = "Некорректный формат ID"
	ErrNotFound                     = "Ресурс не найден"
	ErrExecutorIdEmpty              = "Параметр executorId является обязательным"
	ErrUserIsAuthorTasks            = "Не возможно удалить Юзера, который создал задачи"
	ErrUserIsExecutorTasks          = "Не возможно удалить Юзера, у которого есть назначенные задачи"
	ErrEmptyBody                    = "Нет полей для обновления"
	ErrInvalidBody                  = "Некорректное тело запроса"
	ErrAuthorNotExist               = "Автор - несуществующий юзер"
	ErrTelegramUernameAlreadyExists = "Юзер с таким TelegramUername уже существует"
	ErrInvalidCreditails            = "Некоррекный логин или пароль"
	ErrExecutorNotExist             = "Исполнитель - несуществующий юзер"
	ErrTokenExpired                 = "Время жизни токена истекло"
	ErrUnauthorized                 = "Ошибка авторизации"
	ErrInternal                     = "Internal server error"
)
