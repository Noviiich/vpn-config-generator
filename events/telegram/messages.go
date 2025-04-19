package telegram

const msgHelp = `Отправьте /wireguard чтобы получить конфигурацию WireGuard`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgAlreadyExists  = "You have already have this page in your list 🤗"
	msgDeleteUser     = "Вы успешно удалили пользователя"
	msgSubscribe      = "Вы успешно подписались на рассылку конфигураций WireGuard"
)

// Button labels
const (
	btnProfile   = "👤 Профиль"
	btnProtocols = "🔒 Протоколы VPN"
	btnHome      = "🏠 Главная"
	btnTariffs   = "💎 Тарифы"
)

// Page content
const (
	msgProfile = `👤 Ваш профиль:
• ID: %d
• Username: @%s
• Статус подписки: %s
• Осталось дней: %d`

	msgProtocols = `🔒 Доступные протоколы VPN:
• WireGuard (рекомендуется)
• OpenVPN

Выберите протокол для получения конфигурации:`

	msgHome = `🏠 Добро пожаловать в VPN сервис!

Выберите действие:
• Получить конфигурацию
• Проверить статус подписки
• Выбрать тариф
• Настройки`

	msgTariffs = `💎 Доступные тарифы:

1️⃣ Базовый
• 1 месяц
• До 3 устройств
• Безлимитный трафик
• Цена: 50₽/мес

2️⃣ Стандарт
• 3 месяца
• До 5 устройств
• Безлимитный трафик
• Цена: 130₽/3 мес

3️⃣ Премиум
• 6 месяцев
• До 10 устройств
• Безлимитный трафик
• Цена: 240₽/6 мес

Выберите тариф для оформления:`
)

const (
	msgErrorCreateUser   = "Ошибка при создании пользователя"
	msgErrorGetStatus    = "Ошибка при получении статуса подписки"
	msgErrorSubscribe    = "Ошибка при подписке"
	msgErrorGetConfig    = "Ошибка при получении конфигурации"
	msgErrorSendDocument = "Ошибка при отправке файла"
	msgNoSubscription    = `Вы не подписаны на рассылку конфигураций WireGuard
Не расстраивайтесь, выполните команду /subscribe`
	msgErrorDeleteUser = "Ошибка при удалении пользователя"
	msgErrorGetUsers   = "Ошибка при получении списка пользователей"
	msgNoUsers         = "Список пользователей пуст"
)
