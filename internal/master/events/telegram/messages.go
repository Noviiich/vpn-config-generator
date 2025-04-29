package telegram

const msgHelp = `Отправьте /wireguard чтобы получить конфигурацию WireGuard`

const msgWireGuardInstructions = `📱 Инструкция по установке WireGuard:

Шаг 1. Установите приложение WireGuard ⬇️
Шаг 2. В телеграм боте выберите сервер с протоколом WireGuard 🌎
Шаг 3. Скачайте файл на устройство или откройте его и нажмите "поделиться" с приложением WireGuard для iPhone
Шаг 4. Добавьте файл в приложение WireGuard и включите VPN ➕

📲 [Скачать Wireguard для iPhone (iOS)](https://apps.apple.com/ru/app/wireguard/id1441195209)
📲 [Скачать Wireguard для Android](https://play.google.com/store/apps/details?id=com.wireguard.android)
🖥 [Скачать Wireguard для MacOS](https://apps.apple.com/ru/app/wireguard/id1451685025?mt=12)
🖥 [Скачать Wireguard для Windows](https://www.wireguard.com/install/)`

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
	btnContact   = "📞 Связаться с админом"
	btnServers   = "🌎 Серверы"
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

	msgContact = `📞 Связь с администратором

Для связи с администратором напишите @alexnoviich.

Мы постараемся ответить вам как можно скорее!`
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
