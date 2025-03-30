package telegram

const msgHelp = `Отправьте /wireguard чтобы получить конфигурацию WireGuard`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
	msgDeleteUser     = "Вы успешно удалили пользователя"
	msgSubscribe      = "Вы успешно подписались на рассылку конфигураций WireGuard"
)

const (
	msgErrorUnknown      = "Неизвестная ошибка"
	msgErrorGetConfig    = "Ошибка при получении конфигурации"
	msgErrorSendDocument = "Ошибка при отправке файла"
	msgErrorDeleteUser   = "Ошибка при удалении пользователя"
	msgUserNotFound      = "пользователя не существует"
	msgUsersNotFound     = "пользователей не существует"
	msgErrorGetUsers     = "Ошибка при получении списка пользователей"
	msgNoUsers           = "Список пользователей пуст"
)
