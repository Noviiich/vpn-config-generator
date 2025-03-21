package telegram

const msgHelp = `Отправьте /wireguard чтобы получить конфигурацию WireGuard`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
)

const (
	msgErrorCreateUser   = "Ошибка при создании пользователя"
	msgErrorGetStatus    = "Ошибка при получении статуса подписки"
	msgErrorSubscribe    = "Ошибка при подписке"
	msgSubscribe         = "Вы успешно подписались на рассылку конфигураций WireGuard"
	msgErrorGetConfig    = "Ошибка при получении конфигурации"
	msgErrorSendDocument = "Ошибка при отправке файла"
	msgNoSubscription    = "Вы не подписаны на рассылку конфигураций WireGuard"
)
