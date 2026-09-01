INSERT INTO config_types (
    name,
    description,
    config_template,
    usage_instructions,
    client_links,
    is_active
)
SELECT
    'WireGuard',
    'Стандартный WireGuard клиент. Простой, быстрый, безопасный VPN протокол.',
    '[Interface]
PrivateKey = {{ private_key }}
Address = {{ client_ipv4 }}/32, {{ client_ipv6 }}/128
DNS = 1.1.1.1, 1.0.0.1, 2606:4700:4700::1111, 2606:4700:4700::1001
MTU = 1280

[Peer]
PublicKey = {{ peer_public_key }}
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = {{ endpoint }}:{{ port }}
PersistentKeepalive = 25
',
    '<div class="usage-instructions">
    <h4>Инструкция по использованию WireGuard</h4>

    <h5>1. Установите клиент WireGuard:</h5>
    <ul>
        <li><i class="fab fa-windows text-primary"></i><strong>Windows:</strong> Скачайте с <a href="https://www.wireguard.com/install/" target="_blank">официального сайта</a></li>
        <li><i class="fab fa-appley"></i><strong>MacOS:</strong> App Store</li>
        <li><i class="fab fa-apple"></i><strong>iOS:</strong> App Store</li>
        <li><i class="fab fa-android text-success"></i><strong>Android:</strong> Google Play</li>
        <li><i class="fab fa-linux"></i><strong>Linux:</strong> <code>sudo apt install wireguard</code> или <code>sudo yum install wireguard-tools</code></li>
    </ul>

    <h5>2. Импорт конфигурации:</h5>
    <ul>
        <li>Скачайте файл конфигурации (.conf)</li>
        <li>Откройте WireGuard клиент</li>
        <li>Нажмите +</li>
        <li>Нажмите "Import tunnel from file" или "Добавить туннель"</li>
        <li>Выберите скачанный файл</li>
        <li>Или добавьте конфигурацию через QR Code</li>
    </ul>

    <h5>3. Подключение:</h5>
    <ul>
        <li>Нажмите "Activate" или "Подключить"</li>
        <li>Проверьте подключение на <a href="https://1.1.1.1/help" target="_blank">1.1.1.1/help</a></li>
    </ul>

    <strong>Совет:</strong> Для автоматического подключения при запуске, включите опцию "Activate on boot" в настройках туннеля.
</div>',
    '{
        "windows": "https://www.wireguard.com/install/",
        "macos": "https://apps.apple.com/ru/app/wireguard/id1441195209",
        "ios": "https://apps.apple.com/ru/app/wireguard/id1441195209",
        "android": "https://play.google.com/store/apps/details?id=com.wireguard.android",
        "linux": "https://www.wireguard.com/install/"
    }',
    1
WHERE NOT EXISTS (
    SELECT 1 FROM config_types WHERE name = 'WireGuard'
);