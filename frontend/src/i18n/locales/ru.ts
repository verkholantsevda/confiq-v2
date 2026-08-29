export default {
    common: {
        login: "Войти",
        logout: "Выход",
        loading: "Загрузка...",
        theme: "Тема",
        language: "Язык",
        search: "Поиск",
        actions: "Действия",
        noData: "Нет данных",
        created: "Создана",
        cancel: "Отмена",
        create: "Создать"
    },

    auth: {
        username: "Логин",
        password: "Пароль",
        loginTitle: "Вход в систему",
        loginError: "Не удалось выполнить вход",
    },

    menu: {
        dashboard: "Панель управления",
        users: "Пользователи",
        endpoints: "Точки подключения",
        groups: "Группы",
        configtypes: "Типы конфигураций",
        configs: "Конфигурации",
        settings: "Настройки"
    },
    pages: {
        dashboard: {
            description: "Панель управления приложением",
            quickActions: "Быстрые действия",
            systemStatus: "Состояние системы",
            recentActivity: "Последняя активность",
        },
        users: {
            description: "Управление пользователями",
            create: "Создать пользователя",
            username: "Пользователь",
            group: "Группа",
            role: "Роль",
            configs: "Конфигурации",
            limit: "Лимит",
            created: "Создан",
            admin: "Администратор",
            user: "Пользователь",
            delete: "Удаление пользователя"
        },
        endpoints: {
            description: "Управление точками подключения",
            create: "Создать точку подключения",
            name: "Имя",
            address: "Адрес",
            port: "Порт",
            groups: "Группы",
        },
        groups: {
            description: "Управление группами",
            create: "Создать группу",
            name: "Название группы",
            groups_descriptions: "Описание",
            users: "Пользователи",
            endpoints: "Точки подключения",
        }, 
        configtypes: {
            name: "Название",
            description: "Управление типами конфигураций",
            description_type: "Описание",
            createbutton: "Создать тип конфигурации",
            edit: "Редактирование типа конфигурации",
            create: "Cоздание типа конфигурации",
            active: "Статус"
        },
        configtypes_create: {
            name: "Имя",
            description: "Описание",
            active: "Статус",
            template: "Шаблон конфигурации",
            instruction: "Инструкция по использованию",
            links_clients: "Ссылки на клиенты",
            create: "Cоздать тип конфигурации",
        },
        configs: {
            description: "Управление конфигурациями",
            create: "Создать конфигурацию",
            download: "Скачать",
            qrcode: "QR-код",
            view: "Посмотреть",
            total: "Всего",
            used: "Используется",
            available: "Доступно",
        },
        configs_edit: {
            endpoint: "Точка подключения",
            created: "Создан",
            configs: "Содержимое конфигурации",
            copy: "Копировать",
            download: "Скачать",
            detail: "Детали подключения",
            address: "адрес",
            instruction: "Инструкция по использованию",
            clients: "Ссылки на приложения",
            QR_Code: "Отсканируйте QR Код в приложении"
        },
    },
    dialog: {
        users: {
            username: {
                value: "Пользователь",
                description: " Минимум 3 символа, только латинские буквы, цифры и _ "
            },
            password: {
                value: "Пароль",
                description: "Минимум 3 символа",
                edit_description: "Оставьте пустым чтобы не менять пароль"
            },
            group: {
                value: "Группа",
                description: "Группа определяет доступные точки подключения для пользователей"
            },
            role: "Роль",
            configs: "Конфиги",
            limit: "Лимит",
            limit_configurations: "Лимит конфигураций",
            limit_configurations_hint: "Максимальное количество конфигураций для пользователя",
            created: "Создан",
            admin: {
                    value: "Администратор",
                    description: "Полный доступ к управлению системой"
                },
            user: {
                    value: "Пользователь",
                    description: "Создание и управление своими конфигурациями"
            },
            create: "Создать пользователя",
            info_title: "Информация",
            info_hint: "Роли пользователей",
            info2_title: "Лимиты конфигураций",
            info2_hint: "По умолчанию: 5 конфигураций",
            info2_hint2: "Рекомендуется: 10–20 для активных пользователей",
            info2_hint3: "Максимум: 100 для доверенных пользователей",
            message_delete: "Вы действительно уверены удалить пользователя?",
            cant_message_delete: "Это действие нельзя отменить"
        },
        endpoints: {
            name: "Имя",
            address: "Адрес",
            port: "Порт",
            configuration_types: "Типы конфигурации",        
            groups: "Группы",        
            name_dialog: "Создание точки подключения",
            examples_title: "Примеры адресов",
            config_types_title: "Типы конфигураций",
            config_types_description: "Выберите, какие типы конфигураций будут доступны для этой точки подключения.",
            address_hint: "IP-адрес или доменное имя",
            port_hint: "Стандартный порт WARP: 2408",
            config_types_hint: "Доступные типы конфигурации для этого endpoint (Ctrl+Click для множественного выбора)",
            groups_hint: "Группы пользователей, которым будет доступен эта точка подключения. Если не выбрано ни одной группы, endpoint будет доступен всем",
            name_placeholder: "Название точки подключения",
            address_placeholder: "Хост или IP-адрес",
            port_placeholder: "Порт",
            config_types_placeholder: "Выберите типы конфигурации",
            create: "Создать точку подключения"

        },
        groups: {
            name_dialog: "Создать группу",
            name: "Имя",
            name_placeholder: "Название группы",
            name_hint: "Уникальное имя для группы",
            description: "Описание",
            description_placeholder: "Описание группы",
            description_hint: "Необязательное описание группы",
            endpoints: "Доступные точки подключения",
            endpoints_hint: "Выберите точки подключения",
            users: "Пользователи",
        },
        config: {
            name_dialog: "Создание новой конфигурации",
            name: "Название конфигурации",
            endpoint: "Выберите точку подключения",
            configtype: "Выберите тип конфигурации",
            create: "Создать конфигурацию"
        }

    },
    donate: {
        title: "Поддержите развитие проекта",
        description: "Данный проект создан и поддерживается мной в свободное время на энтузиазме и остается полностью бесплатным.\nВаша поддержка поможет поддерживать проект в рабочем состоянии, а также уделять больше времени развитию, улучшению и внедрению новых функций.",
        buttonYoumoney: "Поддержать",
        buttonTips: "Поддержать через Tips",
    },
    user: {
        quickactions: {
            create: "Создать конфигурацию",
            allconf: "Все конфигурации",
            profile: "Профиль"
        },
        profile: {
            title: "Информация об аккаунте",
            user: "Пользователь",
            created: "Регистрация",
            conf: "Конфигурация",
            password: {
                title: "Смена пароля",
                current: "Текущий пароль",
                new: "Новый пароль",
                new_hint: "Минимум 3 символа",
                confirm: "Подтвердите пароль",
                change: "Изменить пароль"
            }
        }
    }
};